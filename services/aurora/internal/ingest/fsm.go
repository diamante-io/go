package ingest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	logpkg "github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/xdr"
)

var (
	defaultSleep = time.Second
)

// ErrReingestRangeConflict indicates that the reingest range overlaps with
// aurora's most recently ingested ledger
type ErrReingestRangeConflict struct {
	maximumLedgerSequence uint32
}

func (e ErrReingestRangeConflict) Error() string {
	return fmt.Sprintf("reingest range overlaps with aurora ingestion, supplied range shouldn't contain ledger %d", e.maximumLedgerSequence)
}

type stateMachineNode interface {
	run(*system) (transition, error)
	String() string
}

type transition struct {
	node          stateMachineNode
	sleepDuration time.Duration
}

func stop() transition {
	return transition{node: stopState{}, sleepDuration: 0}
}

func start() transition {
	return transition{node: startState{}, sleepDuration: defaultSleep}
}

func rebuild(checkpointLedger uint32) transition {
	return transition{
		node: buildState{
			checkpointLedger: checkpointLedger,
		},
		sleepDuration: defaultSleep,
	}
}

func resume(latestSuccessfullyProcessedLedger uint32) transition {
	return transition{
		node: resumeState{
			latestSuccessfullyProcessedLedger: latestSuccessfullyProcessedLedger,
		},
		sleepDuration: defaultSleep,
	}
}

func resumeImmediately(latestSuccessfullyProcessedLedger uint32) transition {
	return transition{
		node: resumeState{
			latestSuccessfullyProcessedLedger: latestSuccessfullyProcessedLedger,
		},
		sleepDuration: 0,
	}
}

func retryResume(r resumeState) transition {
	return transition{
		node:          r,
		sleepDuration: defaultSleep,
	}
}

func historyRange(fromLedger, toLedger uint32) transition {
	return transition{
		node: historyRangeState{
			fromLedger: fromLedger,
			toLedger:   toLedger,
		},
		sleepDuration: defaultSleep,
	}
}

func waitForCheckPoint() transition {
	return transition{
		node:          waitForCheckpointState{},
		sleepDuration: 0,
	}
}

type stopState struct{}

func (stopState) String() string {
	return "stop"
}

func (stopState) run(s *system) (transition, error) {
	return stop(), errors.New("Cannot run terminal state")
}

type startState struct {
	suggestedCheckpoint uint32
}

func (startState) String() string {
	return "start"
}

func (state startState) run(s *system) (transition, error) {
	if err := s.historyQ.Begin(); err != nil {
		return start(), errors.Wrap(err, "Error starting a transaction")
	}
	defer s.historyQ.Rollback()

	// This will get the value `FOR UPDATE`, blocking it for other nodes.
	lastIngestedLedger, err := s.historyQ.GetLastLedgerIngest(s.ctx)
	if err != nil {
		return start(), errors.Wrap(err, getLastIngestedErrMsg)
	}

	ingestVersion, err := s.historyQ.GetIngestVersion(s.ctx)
	if err != nil {
		return start(), errors.Wrap(err, getIngestVersionErrMsg)
	}

	if ingestVersion > CurrentVersion {
		log.WithFields(logpkg.F{
			"ingestVersion":  ingestVersion,
			"currentVersion": CurrentVersion,
		}).Info("ingestion version in db is greater than current version, going to terminate")
		return stop(), nil
	}

	lastHistoryLedger, err := s.historyQ.GetLatestHistoryLedger(s.ctx)
	if err != nil {
		return start(), errors.Wrap(err, "Error getting last history ledger sequence")
	}

	if ingestVersion != CurrentVersion || lastIngestedLedger == 0 {
		// This block is either starting from empty state or ingestion
		// version upgrade.
		// This will always run on a single instance due to the fact that
		// `LastLedgerIngest` value is blocked for update and will always
		// be updated when leading instance finishes processing state.
		// In case of errors it will start `Init` from the beginning.
		var lastCheckpoint uint32
		if state.suggestedCheckpoint != 0 {
			lastCheckpoint = state.suggestedCheckpoint
		} else {
			lastCheckpoint, err = s.historyAdapter.GetLatestLedgerSequence()
			if err != nil {
				return start(), errors.Wrap(err, "Error getting last checkpoint")
			}
		}

		if lastHistoryLedger != 0 {
			// There are ledgers in history_ledgers table. This means that the
			// old or new ingest system was running prior the upgrade. In both
			// cases we need to:
			// * Wait for the checkpoint ledger if the latest history ledger is
			//   greater that the latest checkpoint ledger.
			// * Catchup history data if the latest history ledger is less than
			//   the latest checkpoint ledger.
			// * Build state from the last checkpoint if the latest history ledger
			//   is equal to the latest checkpoint ledger.
			switch {
			case lastHistoryLedger > lastCheckpoint:
				return waitForCheckPoint(), nil
			case lastHistoryLedger < lastCheckpoint:
				return historyRange(lastHistoryLedger+1, lastCheckpoint), nil
			default: // lastHistoryLedger == lastCheckpoint
				// Build state but make sure it's using `lastCheckpoint`. It's possible
				// that the new checkpoint will be created during state transition.
				return rebuild(lastCheckpoint), nil
			}
		}

		return rebuild(lastCheckpoint), nil
	}

	switch {
	case lastHistoryLedger > lastIngestedLedger:
		// Ingestion was running at some point the past but was turned off.
		// Now it's on by default but the latest history ledger is greater
		// than the latest ingest ledger. We reset the exp ledger sequence
		// so init state will rebuild the state correctly.
		err = s.historyQ.UpdateLastLedgerIngest(s.ctx, 0)
		if err != nil {
			return start(), errors.Wrap(err, updateLastLedgerIngestErrMsg)
		}
		err = s.historyQ.Commit()
		if err != nil {
			return start(), errors.Wrap(err, commitErrMsg)
		}
		return start(), nil
	// lastHistoryLedger != 0 check is here to check the case when one node ingested
	// the state (so latest ingestion is > 0) but no history has been ingested yet.
	// In such case we execute default case and resume from the last ingested
	// ledger.
	case lastHistoryLedger != 0 && lastHistoryLedger < lastIngestedLedger:
		// Ingestion was running at some point the past but was turned off.
		// Now it's on by default but the latest history ledger is less
		// than the latest ingest ledger. We catchup history.
		return historyRange(lastHistoryLedger+1, lastIngestedLedger), nil
	default: // lastHistoryLedger == lastIngestedLedger
		// The other node already ingested a state (just now or in the past)
		// so we need to get offers from a DB, then resume session normally.
		// State pipeline is NOT processed.
		log.WithField("last_ledger", lastIngestedLedger).
			Info("Resuming ingestion system from last processed ledger...")

		return resume(lastIngestedLedger), nil
	}
}

type buildState struct {
	checkpointLedger uint32
	stop             bool
}

func (b buildState) String() string {
	return fmt.Sprintf("buildFromCheckpoint(checkpointLedger=%d)", b.checkpointLedger)
}

func (b buildState) run(s *system) (transition, error) {
	var nextFailState = start()
	if b.stop {
		nextFailState = stop()
	}

	if b.checkpointLedger == 0 {
		return nextFailState, errors.New("unexpected checkpointLedger value")
	}

	// We don't need to prepare range for genesis checkpoint because we don't
	// perform protocol version and bucket list hash checks.
	// In the long term we should probably create artificial xdr.LedgerCloseMeta
	// for ledger #1 instead of using `ingest.GenesisChange` reader in
	// ProcessorRunner.RunHistoryArchiveIngestion().
	var ledgerCloseMeta xdr.LedgerCloseMeta
	if b.checkpointLedger != 1 {
		err := s.maybePrepareRange(s.ctx, b.checkpointLedger)
		if err != nil {
			return nextFailState, err
		}

		log.WithField("sequence", b.checkpointLedger).Info("Waiting for ledger to be available in the backend...")
		startTime := time.Now()
		ledgerCloseMeta, err = s.ledgerBackend.GetLedger(s.ctx, b.checkpointLedger)
		if err != nil {
			return nextFailState, errors.Wrap(err, "error getting ledger blocking")
		}
		log.WithFields(logpkg.F{
			"sequence": b.checkpointLedger,
			"duration": time.Since(startTime).Seconds(),
		}).Info("Ledger returned from the backend")
	}

	if err := s.historyQ.Begin(); err != nil {
		return nextFailState, errors.Wrap(err, "Error starting a transaction")
	}
	defer s.historyQ.Rollback()

	// We need to get this value `FOR UPDATE` so all other instances
	// are blocked.
	lastIngestedLedger, err := s.historyQ.GetLastLedgerIngest(s.ctx)
	if err != nil {
		return nextFailState, errors.Wrap(err, getLastIngestedErrMsg)
	}

	ingestVersion, err := s.historyQ.GetIngestVersion(s.ctx)
	if err != nil {
		return nextFailState, errors.Wrap(err, getIngestVersionErrMsg)
	}

	// Double check if we should proceed with state ingestion. It's possible that
	// another ingesting instance will be redirected to this state from `init`
	// but it's first to complete the task.
	if ingestVersion == CurrentVersion && lastIngestedLedger > 0 {
		log.Info("Another instance completed `buildState`. Skipping...")
		return nextFailState, nil
	}

	if err = s.updateCursor(b.checkpointLedger - 1); err != nil {
		// Don't return updateCursor error.
		log.WithError(err).Warn("error updating diamcircle-core cursor")
	}

	log.Info("Starting ingestion system from empty state...")

	// Clear last_ingested_ledger in key value store
	err = s.historyQ.UpdateLastLedgerIngest(s.ctx, 0)
	if err != nil {
		return nextFailState, errors.Wrap(err, updateLastLedgerIngestErrMsg)
	}

	// Clear invalid state in key value store. It's possible that upgraded
	// ingestion is fixing it.
	err = s.historyQ.UpdateExpStateInvalid(s.ctx, false)
	if err != nil {
		return nextFailState, errors.Wrap(err, updateExpStateInvalidErrMsg)
	}

	// State tables should be empty.
	err = s.historyQ.TruncateIngestStateTables(s.ctx)
	if err != nil {
		return nextFailState, errors.Wrap(err, "Error clearing ingest tables")
	}

	log.WithFields(logpkg.F{
		"sequence": b.checkpointLedger,
	}).Info("Processing state")
	startTime := time.Now()

	var stats ingest.StatsChangeProcessorResults
	if b.checkpointLedger == 1 {
		stats, err = s.runner.RunGenesisStateIngestion()
	} else {
		stats, err = s.runner.RunHistoryArchiveIngestion(
			ledgerCloseMeta.LedgerSequence(),
			ledgerCloseMeta.ProtocolVersion(),
			ledgerCloseMeta.BucketListHash(),
		)
	}

	if err != nil {
		return nextFailState, errors.Wrap(err, "Error ingesting history archive")
	}

	if err = s.historyQ.UpdateIngestVersion(s.ctx, CurrentVersion); err != nil {
		return nextFailState, errors.Wrap(err, "Error updating ingestion version")
	}

	if err = s.completeIngestion(s.ctx, b.checkpointLedger); err != nil {
		return nextFailState, err
	}

	log.
		WithFields(stats.Map()).
		WithFields(logpkg.F{
			"sequence": b.checkpointLedger,
			"duration": time.Since(startTime).Seconds(),
		}).
		Info("Processed state")

	if b.stop {
		return stop(), nil
	}
	// If successful, continue from the next ledger
	return resume(b.checkpointLedger), nil
}

type resumeState struct {
	latestSuccessfullyProcessedLedger uint32
}

func (r resumeState) String() string {
	return fmt.Sprintf("resume(latestSuccessfullyProcessedLedger=%d)", r.latestSuccessfullyProcessedLedger)
}

func (r resumeState) run(s *system) (transition, error) {
	if r.latestSuccessfullyProcessedLedger == 0 {
		return start(), errors.New("unexpected latestSuccessfullyProcessedLedger value")
	}

	s.metrics.LocalLatestLedger.Set(float64(r.latestSuccessfullyProcessedLedger))

	ingestLedger := r.latestSuccessfullyProcessedLedger + 1

	err := s.maybePrepareRange(s.ctx, ingestLedger)
	if err != nil {
		return start(), err
	}

	log.WithField("sequence", ingestLedger).Info("Waiting for ledger to be available in the backend...")
	startTime := time.Now()
	ledgerCloseMeta, err := s.ledgerBackend.GetLedger(s.ctx, ingestLedger)
	if err != nil {
		return start(), errors.Wrap(err, "error getting ledger blocking")
	}
	duration := time.Since(startTime).Seconds()
	log.WithFields(logpkg.F{
		"sequence": ingestLedger,
		"duration": duration,
	}).Info("Ledger returned from the backend")

	s.Metrics().LedgerFetchDurationSummary.Observe(float64(duration))

	if err = s.historyQ.Begin(); err != nil {
		return retryResume(r),
			errors.Wrap(err, "Error starting a transaction")
	}
	defer s.historyQ.Rollback()

	// This will get the value `FOR UPDATE`, blocking it for other nodes.
	lastIngestedLedger, err := s.historyQ.GetLastLedgerIngest(s.ctx)
	if err != nil {
		return retryResume(r), errors.Wrap(err, getLastIngestedErrMsg)
	}

	if ingestLedger > lastIngestedLedger+1 {
		return start(), errors.New("expected ingest ledger to be at most one greater " +
			"than last ingested ledger in db")
	} else if ingestLedger <= lastIngestedLedger {
		log.WithField("ingestLedger", ingestLedger).
			WithField("lastIngestedLedger", lastIngestedLedger).
			Info("bumping ingest ledger to next ledger after ingested ledger in db")
		return retryResume(resumeState{
			latestSuccessfullyProcessedLedger: lastIngestedLedger,
		}), nil
	}

	ingestVersion, err := s.historyQ.GetIngestVersion(s.ctx)
	if err != nil {
		return retryResume(r), errors.Wrap(err, getIngestVersionErrMsg)
	}

	if ingestVersion != CurrentVersion {
		log.WithFields(logpkg.F{
			"ingestVersion":  ingestVersion,
			"currentVersion": CurrentVersion,
		}).Info("ingestion version in db is not current, going back to start state")
		return start(), nil
	}

	lastHistoryLedger, err := s.historyQ.GetLatestHistoryLedger(s.ctx)
	if err != nil {
		return retryResume(r), errors.Wrap(err, "could not get latest history ledger")
	}

	if lastHistoryLedger != 0 && lastHistoryLedger != lastIngestedLedger {
		log.WithFields(logpkg.F{
			"lastHistoryLedger":  lastHistoryLedger,
			"lastIngestedLedger": lastIngestedLedger,
		}).Info(
			"last history ledger does not match last ingested ledger, " +
				"going back to start state",
		)
		return start(), nil
	}

	startTime = time.Now()

	log.WithFields(logpkg.F{
		"sequence": ingestLedger,
		"state":    true,
		"ledger":   true,
		"commit":   true,
	}).Info("Processing ledger")

	stats, err :=
		s.runner.RunAllProcessorsOnLedger(ledgerCloseMeta)
	if err != nil {
		return retryResume(r), errors.Wrap(err, "Error running processors on ledger")
	}

	rebuildStart := time.Now()
	err = s.historyQ.RebuildTradeAggregationBuckets(s.ctx, ingestLedger, ingestLedger)
	if err != nil {
		return stop(), errors.Wrap(err, "error rebuilding trade aggregations")
	}
	rebuildDuration := time.Since(rebuildStart).Seconds()
	s.Metrics().LedgerIngestionTradeAggregationDuration.Observe(float64(rebuildDuration))

	if err = s.completeIngestion(s.ctx, ingestLedger); err != nil {
		return retryResume(r), err
	}

	if err = s.updateCursor(ingestLedger); err != nil {
		// Don't return updateCursor error.
		log.WithError(err).Warn("error updating diamcircle-core cursor")
	}

	duration = time.Since(startTime).Seconds()
	s.Metrics().LedgerIngestionDuration.Observe(float64(duration))

	// Update stats metrics
	changeStatsMap := stats.changeStats.Map()
	r.addLedgerStatsMetricFromMap(s, "change", changeStatsMap)
	r.addProcessorDurationsMetricFromMap(s, stats.changeDurations)

	transactionStatsMap := stats.transactionStats.Map()
	r.addLedgerStatsMetricFromMap(s, "ledger", transactionStatsMap)
	tradeStatsMap := stats.tradeStats.Map()
	r.addLedgerStatsMetricFromMap(s, "trades", tradeStatsMap)
	r.addProcessorDurationsMetricFromMap(s, stats.transactionDurations)

	localLog := log.WithFields(logpkg.F{
		"sequence": ingestLedger,
		"duration": duration,
		"state":    true,
		"ledger":   true,
		"commit":   true,
	})

	if s.config.EnableExtendedLogLedgerStats {
		localLog = localLog.
			WithFields(changeStatsMap).
			WithFields(transactionStatsMap).
			WithFields(tradeStatsMap)
	}

	localLog.Info("Processed ledger")

	s.maybeVerifyState(ingestLedger)

	return resumeImmediately(ingestLedger), nil
}

func (r resumeState) addLedgerStatsMetricFromMap(s *system, prefix string, m map[string]interface{}) {
	for stat, value := range m {
		stat = strings.Replace(stat, "stats_", prefix+"_", 1)
		s.Metrics().LedgerStatsCounter.
			With(prometheus.Labels{"type": stat}).Add(float64(value.(int64)))
	}
}

func (r resumeState) addProcessorDurationsMetricFromMap(s *system, m map[string]time.Duration) {
	for processorName, value := range m {
		// * is not accepted in Prometheus labels
		processorName = strings.Replace(processorName, "*", "", -1)
		s.Metrics().ProcessorsRunDuration.
			With(prometheus.Labels{"name": processorName}).Add(value.Seconds())
		s.Metrics().ProcessorsRunDurationSummary.
			With(prometheus.Labels{"name": processorName}).Observe(value.Seconds())
	}
}

type historyRangeState struct {
	fromLedger uint32
	toLedger   uint32
}

func (h historyRangeState) String() string {
	return fmt.Sprintf(
		"historyRange(fromLedger=%d, toLedger=%d)",
		h.fromLedger,
		h.toLedger,
	)
}

// historyRangeState is used when catching up history data
func (h historyRangeState) run(s *system) (transition, error) {
	if h.fromLedger == 0 || h.toLedger == 0 ||
		h.fromLedger > h.toLedger {
		return start(), errors.Errorf("invalid range: [%d, %d]", h.fromLedger, h.toLedger)
	}

	err := s.maybePrepareRange(s.ctx, h.fromLedger)
	if err != nil {
		return start(), err
	}

	if err = s.historyQ.Begin(); err != nil {
		return start(), errors.Wrap(err, "Error starting a transaction")
	}
	defer s.historyQ.Rollback()

	// acquire distributed lock so no one else can perform ingestion operations.
	if _, err = s.historyQ.GetLastLedgerIngest(s.ctx); err != nil {
		return start(), errors.Wrap(err, getLastIngestedErrMsg)
	}

	lastHistoryLedger, err := s.historyQ.GetLatestHistoryLedger(s.ctx)
	if err != nil {
		return start(), errors.Wrap(err, "could not get latest history ledger")
	}

	// We should be ingesting the ledger which occurs after
	// lastHistoryLedger. Otherwise, some other aurora node has
	// already completed the ingest history range operation and
	// we should go back to the init state
	if lastHistoryLedger != h.fromLedger-1 {
		return start(), nil
	}

	for cur := h.fromLedger; cur <= h.toLedger; cur++ {
		var ledgerCloseMeta xdr.LedgerCloseMeta

		log.WithField("sequence", cur).Info("Waiting for ledger to be available in the backend...")
		startTime := time.Now()

		ledgerCloseMeta, err = s.ledgerBackend.GetLedger(s.ctx, cur)
		if err != nil {
			// Commit finished work in case of ledger backend error.
			commitErr := s.historyQ.Commit()
			if commitErr != nil {
				log.WithError(commitErr).Error("Error commiting partial range results")
			} else {
				log.Info("Commited partial range results")
			}
			return start(), errors.Wrap(err, "error getting ledger")
		}

		log.WithFields(logpkg.F{
			"sequence": cur,
			"duration": time.Since(startTime).Seconds(),
		}).Info("Ledger returned from the backend")

		if err = runTransactionProcessorsOnLedger(s, ledgerCloseMeta); err != nil {
			return start(), err
		}
	}

	if err = s.historyQ.Commit(); err != nil {
		return start(), errors.Wrap(err, commitErrMsg)
	}

	return start(), nil
}

func runTransactionProcessorsOnLedger(s *system, ledger xdr.LedgerCloseMeta) error {
	log.WithFields(logpkg.F{
		"sequence": ledger.LedgerSequence(),
		"state":    false,
		"ledger":   true,
		"commit":   false,
	}).Info("Processing ledger")
	startTime := time.Now()

	ledgerTransactionStats, _, tradeStats, err := s.runner.RunTransactionProcessorsOnLedger(ledger)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error processing ledger sequence=%d", ledger.LedgerSequence()))
	}

	log.
		WithFields(ledgerTransactionStats.Map()).
		WithFields(tradeStats.Map()).
		WithFields(logpkg.F{
			"sequence": ledger.LedgerSequence(),
			"duration": time.Since(startTime).Seconds(),
			"state":    false,
			"ledger":   true,
			"commit":   false,
		}).
		Info("Processed ledger")
	return nil
}

type reingestHistoryRangeState struct {
	fromLedger uint32
	toLedger   uint32
	force      bool
}

func (h reingestHistoryRangeState) String() string {
	return fmt.Sprintf(
		"reingestHistoryRange(fromLedger=%d, toLedger=%d, force=%t)",
		h.fromLedger,
		h.toLedger,
		h.force,
	)
}

func (h reingestHistoryRangeState) ingestRange(s *system, fromLedger, toLedger uint32) error {
	if s.historyQ.GetTx() == nil {
		return errors.New("expected transaction to be present")
	}

	// Clear history data before ingesting - used in `reingest range` command.
	start, end, err := toid.LedgerRangeInclusive(
		int32(fromLedger),
		int32(toLedger),
	)
	if err != nil {
		return errors.Wrap(err, "Invalid range")
	}

	err = s.historyQ.DeleteRangeAll(s.ctx, start, end)
	if err != nil {
		return errors.Wrap(err, "error in DeleteRangeAll")
	}

	for cur := fromLedger; cur <= toLedger; cur++ {
		var ledgerCloseMeta xdr.LedgerCloseMeta
		ledgerCloseMeta, err = s.ledgerBackend.GetLedger(s.ctx, cur)
		if err != nil {
			return errors.Wrap(err, "error getting ledger")
		}

		if err = runTransactionProcessorsOnLedger(s, ledgerCloseMeta); err != nil {
			return err
		}
	}

	return nil
}

func (h reingestHistoryRangeState) prepareRange(s *system) (transition, error) {
	log.WithFields(logpkg.F{
		"from": h.fromLedger,
		"to":   h.toLedger,
	}).Info("Preparing ledger backend to retrieve range")
	startTime := time.Now()

	err := s.ledgerBackend.PrepareRange(s.ctx, ledgerbackend.BoundedRange(h.fromLedger, h.toLedger))
	if err != nil {
		return stop(), errors.Wrap(err, "error preparing range")
	}

	log.WithFields(logpkg.F{
		"from":     h.fromLedger,
		"to":       h.toLedger,
		"duration": time.Since(startTime).Seconds(),
	}).Info("Range ready")

	return transition{}, nil
}

// reingestHistoryRangeState is used as a command to reingest historical data
func (h reingestHistoryRangeState) run(s *system) (transition, error) {
	if h.fromLedger == 0 || h.toLedger == 0 ||
		h.fromLedger > h.toLedger {
		return stop(), errors.Errorf("invalid range: [%d, %d]", h.fromLedger, h.toLedger)
	}

	if h.fromLedger == 1 {
		log.Warn("Ledger 1 is pregenerated and not available, starting from ledger 2.")
		h.fromLedger = 2
	}

	var startTime time.Time

	if h.force {
		if t, err := h.prepareRange(s); err != nil {
			return t, err
		}
		startTime = time.Now()

		if err := s.historyQ.Begin(); err != nil {
			return stop(), errors.Wrap(err, "Error starting a transaction")
		}
		defer s.historyQ.Rollback()

		// acquire distributed lock so no one else can perform ingestion operations.
		if _, err := s.historyQ.GetLastLedgerIngest(s.ctx); err != nil {
			return stop(), errors.Wrap(err, getLastIngestedErrMsg)
		}

		if err := h.ingestRange(s, h.fromLedger, h.toLedger); err != nil {
			return stop(), err
		}

		if err := s.historyQ.Commit(); err != nil {
			return stop(), errors.Wrap(err, commitErrMsg)
		}
	} else {
		lastIngestedLedger, err := s.historyQ.GetLastLedgerIngestNonBlocking(s.ctx)
		if err != nil {
			return stop(), errors.Wrap(err, getLastIngestedErrMsg)
		}

		if lastIngestedLedger > 0 && h.toLedger >= lastIngestedLedger {
			return stop(), ErrReingestRangeConflict{lastIngestedLedger}
		}

		// Only prepare the range after checking the bounds to enable an early error return
		var t transition
		if t, err = h.prepareRange(s); err != nil {
			return t, err
		}
		startTime = time.Now()

		for cur := h.fromLedger; cur <= h.toLedger; cur++ {
			err = func(ledger uint32) error {
				if e := s.historyQ.Begin(); e != nil {
					return errors.Wrap(e, "Error starting a transaction")
				}
				defer s.historyQ.Rollback()

				// ingest each ledger in a separate transaction to prevent deadlocks
				// when acquiring ShareLocks from multiple parallel reingest range processes
				if e := h.ingestRange(s, ledger, ledger); e != nil {
					return e
				}

				if e := s.historyQ.Commit(); e != nil {
					return errors.Wrap(e, commitErrMsg)
				}

				return nil
			}(cur)
			if err != nil {
				return stop(), err
			}
		}
	}

	err := s.historyQ.RebuildTradeAggregationBuckets(s.ctx, h.fromLedger, h.toLedger)
	if err != nil {
		return stop(), errors.Wrap(err, "Error rebuilding trade aggregations")
	}

	log.WithFields(logpkg.F{
		"from":     h.fromLedger,
		"to":       h.toLedger,
		"duration": time.Since(startTime).Seconds(),
	}).Info("Reingestion done")

	return stop(), nil
}

type waitForCheckpointState struct{}

func (waitForCheckpointState) String() string {
	return "waitForCheckpoint"
}

func (waitForCheckpointState) run(*system) (transition, error) {
	log.Info("Waiting for the next checkpoint...")
	time.Sleep(10 * time.Second)
	return start(), nil
}

type verifyRangeState struct {
	fromLedger  uint32
	toLedger    uint32
	verifyState bool
}

func (v verifyRangeState) String() string {
	return fmt.Sprintf(
		"verifyRange(fromLedger=%d, toLedger=%d, verifyState=%t)",
		v.fromLedger,
		v.toLedger,
		v.verifyState,
	)
}

func (v verifyRangeState) run(s *system) (transition, error) {
	if v.fromLedger == 0 || v.toLedger == 0 ||
		v.fromLedger > v.toLedger {
		return stop(), errors.Errorf("invalid range: [%d, %d]", v.fromLedger, v.toLedger)
	}

	if err := s.historyQ.Begin(); err != nil {
		err = errors.Wrap(err, "Error starting a transaction")
		return stop(), err
	}
	defer s.historyQ.Rollback()

	// Simple check if DB clean
	lastIngestedLedger, err := s.historyQ.GetLastLedgerIngest(s.ctx)
	if err != nil {
		err = errors.Wrap(err, getLastIngestedErrMsg)
		return stop(), err
	}

	if lastIngestedLedger != 0 {
		err = errors.New("Database not empty")
		return stop(), err
	}

	log.WithField("sequence", v.fromLedger).Info("Preparing range")
	startTime := time.Now()

	err = s.ledgerBackend.PrepareRange(s.ctx, ledgerbackend.BoundedRange(v.fromLedger, v.toLedger))
	if err != nil {
		return stop(), errors.Wrap(err, "Error preparing range")
	}

	log.WithFields(logpkg.F{
		"sequence": v.fromLedger,
		"duration": time.Since(startTime).Seconds(),
	}).Info("Range prepared")

	log.WithField("sequence", v.fromLedger).Info("Processing state")
	startTime = time.Now()

	ledgerCloseMeta, err := s.ledgerBackend.GetLedger(s.ctx, v.fromLedger)
	if err != nil {
		return stop(), errors.Wrap(err, "error getting ledger")
	}

	stats, err := s.runner.RunHistoryArchiveIngestion(
		ledgerCloseMeta.LedgerSequence(),
		ledgerCloseMeta.ProtocolVersion(),
		ledgerCloseMeta.BucketListHash(),
	)
	if err != nil {
		err = errors.Wrap(err, "Error ingesting history archive")
		return stop(), err
	}

	if err = s.completeIngestion(s.ctx, v.fromLedger); err != nil {
		return stop(), err
	}

	log.
		WithFields(stats.Map()).
		WithFields(logpkg.F{
			"sequence": v.fromLedger,
			"duration": time.Since(startTime).Seconds(),
		}).
		Info("Processed state")

	for sequence := v.fromLedger + 1; sequence <= v.toLedger; sequence++ {
		log.WithFields(logpkg.F{
			"sequence": sequence,
			"state":    true,
			"ledger":   true,
			"commit":   true,
		}).Info("Processing ledger")
		startTime := time.Now()

		if err = s.historyQ.Begin(); err != nil {
			err = errors.Wrap(err, "Error starting a transaction")
			return stop(), err
		}

		var ledgerCloseMeta xdr.LedgerCloseMeta
		ledgerCloseMeta, err = s.ledgerBackend.GetLedger(s.ctx, sequence)
		if err != nil {
			return stop(), errors.Wrap(err, "error getting ledger")
		}

		var ledgerStats ledgerStats
		ledgerStats, err = s.runner.RunAllProcessorsOnLedger(ledgerCloseMeta)
		if err != nil {
			err = errors.Wrap(err, "Error running processors on ledger")
			return stop(), err
		}

		if err = s.completeIngestion(s.ctx, sequence); err != nil {
			return stop(), err
		}

		log.
			WithFields(ledgerStats.changeStats.Map()).
			WithFields(ledgerStats.transactionStats.Map()).
			WithFields(ledgerStats.tradeStats.Map()).
			WithFields(logpkg.F{
				"sequence": sequence,
				"duration": time.Since(startTime).Seconds(),
				"state":    true,
				"ledger":   true,
				"commit":   true,
			}).
			Info("Processed ledger")
	}

	err = s.historyQ.RebuildTradeAggregationBuckets(s.ctx, v.fromLedger, v.toLedger)
	if err != nil {
		return stop(), errors.Wrap(err, "error rebuilding trade aggregations")
	}

	if v.verifyState {
		err = s.verifyState(false)
	}

	return stop(), err
}

type stressTestState struct{}

func (stressTestState) String() string {
	return "stressTest"
}

func (stressTestState) run(s *system) (transition, error) {
	if err := s.historyQ.Begin(); err != nil {
		err = errors.Wrap(err, "Error starting a transaction")
		return stop(), err
	}
	defer s.historyQ.Rollback()

	// Simple check if DB clean
	lastIngestedLedger, err := s.historyQ.GetLastLedgerIngest(s.ctx)
	if err != nil {
		err = errors.Wrap(err, getLastIngestedErrMsg)
		return stop(), err
	}

	if lastIngestedLedger != 0 {
		err = errors.New("Database not empty")
		return stop(), err
	}

	curHeap, sysHeap := getMemStats()
	sequence := lastIngestedLedger + 1
	log.WithFields(logpkg.F{
		"currentHeapSizeMB": curHeap,
		"systemHeapSizeMB":  sysHeap,
		"sequence":          sequence,
		"state":             true,
		"ledger":            true,
		"commit":            true,
	}).Info("Processing ledger")
	startTime := time.Now()

	ledgerCloseMeta, err := s.ledgerBackend.GetLedger(s.ctx, sequence)
	if err != nil {
		return stop(), errors.Wrap(err, "error getting ledger")
	}

	stats, err := s.runner.RunAllProcessorsOnLedger(ledgerCloseMeta)
	if err != nil {
		err = errors.Wrap(err, "Error running processors on ledger")
		return stop(), err
	}

	if err = s.completeIngestion(s.ctx, sequence); err != nil {
		return stop(), err
	}

	curHeap, sysHeap = getMemStats()
	log.
		WithFields(stats.changeStats.Map()).
		WithFields(stats.transactionStats.Map()).
		WithFields(stats.tradeStats.Map()).
		WithFields(logpkg.F{
			"currentHeapSizeMB": curHeap,
			"systemHeapSizeMB":  sysHeap,
			"sequence":          sequence,
			"duration":          time.Since(startTime).Seconds(),
			"state":             true,
			"ledger":            true,
			"commit":            true,
		}).
		Info("Processed ledger")

	return stop(), nil
}

func (s *system) completeIngestion(ctx context.Context, ledger uint32) error {
	if ledger == 0 {
		return errors.New("ledger must be positive")
	}

	if err := s.historyQ.UpdateLastLedgerIngest(ctx, ledger); err != nil {
		err = errors.Wrap(err, updateLastLedgerIngestErrMsg)
		return err
	}

	if err := s.historyQ.Commit(); err != nil {
		return errors.Wrap(err, commitErrMsg)
	}

	return nil
}

// maybePrepareRange checks if the range is prepared and, if not, prepares it.
func (s *system) maybePrepareRange(ctx context.Context, from uint32) error {
	ledgerRange := ledgerbackend.UnboundedRange(from)

	prepared, err := s.ledgerBackend.IsPrepared(ctx, ledgerRange)
	if err != nil {
		return errors.Wrap(err, "error checking prepared range")
	}

	if !prepared {
		log.WithFields(logpkg.F{"from": from}).Info("Preparing range")
		startTime := time.Now()

		err = s.ledgerBackend.PrepareRange(ctx, ledgerRange)
		if err != nil {
			return errors.Wrap(err, "error preparing range")
		}

		log.WithFields(logpkg.F{
			"from":     from,
			"duration": time.Since(startTime).Seconds(),
		}).Info("Range prepared")

		return nil
	}

	return nil
}
