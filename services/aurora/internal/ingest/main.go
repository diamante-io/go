// Package ingest contains the new ingestion system for aurora.
// It currently runs completely independent of the old one, that means
// that the new system can be ledgers behind/ahead the old system.
package ingest

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"go/clients/diamcirclecore"
	"go/historyarchive"
	"go/ingest"
	"go/ingest/ledgerbackend"
	"go/services/aurora/internal/db2/history"
	"go/support/db"
	"go/support/errors"
	logpkg "go/support/log"
)

const (
	// MaxSupportedProtocolVersion defines the maximum supported version of
	// the diamcircle protocol.
	MaxSupportedProtocolVersion uint32 = 18

	// CurrentVersion reflects the latest version of the ingestion
	// algorithm. This value is stored in KV store and is used to decide
	// if there's a need to reprocess the ledger state or reingest data.
	//
	// Version history:
	// - 1: Initial version
	// - 2: Added the orderbook, offers processors and distributed ingestion.
	// - 3: Fixed a bug that could potentialy result in invalid state
	//      (#1722). Update the version to clear the state.
	// - 4: Fixed a bug in AccountSignersChanged method.
	// - 5: Added trust lines.
	// - 6: Added accounts and accounts data.
	// - 7: Fixes a bug in AccountSignersChanged method.
	// - 8: Fixes AccountSigners processor to remove preauth tx signer
	//      when preauth tx is failed.
	// - 9: Fixes a bug in asset stats processor that counted unauthorized
	//      trustlines.
	// - 10: Fixes a bug in meta processing (fees are now processed before
	//      everything else).
	// - 11: Protocol 14: CAP-23 and CAP-33.
	// - 12: Trigger state rebuild due to `absTime` -> `abs_time` rename
	//       in ClaimableBalances predicates.
	// - 13: Trigger state rebuild to include more than just authorized assets.
	// - 14: Trigger state rebuild to include claimable balances in the asset stats processor.
	// - 15: Fixed bug in asset stat ingestion where clawback is enabled (#3846).
	CurrentVersion = 15

	// MaxDBConnections is the size of the postgres connection pool dedicated to aurora ingestion:
	//  * Ledger ingestion,
	//  * State verifications,
	//  * Metrics updates.
	MaxDBConnections = 3

	defaultCoreCursorName           = "aurora"
	stateVerificationErrorThreshold = 3
)

var log = logpkg.DefaultLogger.WithField("service", "ingest")

type Config struct {
	CoreSession            db.SessionInterface
	diamcircleCoreURL      string
	diamcircleCoreCursor   string
	EnableCaptiveCore      bool
	CaptiveCoreBinaryPath  string
	CaptiveCoreStoragePath string
	CaptiveCoreToml        *ledgerbackend.CaptiveCoreToml
	RemoteCaptiveCoreURL   string
	NetworkPassphrase      string

	HistorySession    db.SessionInterface
	HistoryArchiveURL string

	DisableStateVerification     bool
	EnableExtendedLogLedgerStats bool

	MaxReingestRetries          int
	ReingestRetryBackoffSeconds int

	// The checkpoint frequency will be 64 unless you are using an exotic test setup.
	CheckpointFrequency uint32
}

const (
	getLastIngestedErrMsg        string = "Error getting last ingested ledger"
	getIngestVersionErrMsg       string = "Error getting ingestion version"
	updateLastLedgerIngestErrMsg string = "Error updating last ingested ledger"
	commitErrMsg                 string = "Error committing db transaction"
	updateExpStateInvalidErrMsg  string = "Error updating state invalid value"
)

type diamcircleCoreClient interface {
	SetCursor(ctx context.Context, id string, cursor int32) error
}

type Metrics struct {
	// MaxSupportedProtocolVersion exposes the maximum protocol version
	// supported by this version.
	MaxSupportedProtocolVersion prometheus.Gauge

	// LocalLedger exposes the last ingested ledger by this ingesting instance.
	LocalLatestLedger prometheus.Gauge

	// LedgerIngestionDuration exposes timing metrics about the rate and
	// duration of ledger ingestion (including updating DB and graph).
	LedgerIngestionDuration prometheus.Summary

	// LedgerIngestionTradeAggregationDuration exposes timing metrics about the rate and
	// duration of rebuilding trade aggregation buckets.
	LedgerIngestionTradeAggregationDuration prometheus.Summary

	// StateVerifyDuration exposes timing metrics about the rate and
	// duration of state verification.
	StateVerifyDuration prometheus.Summary

	// StateInvalidGauge exposes state invalid metric. 1 if state is invalid,
	// 0 otherwise.
	StateInvalidGauge prometheus.GaugeFunc

	// StateVerifyLedgerEntriesCount exposes total number of ledger entries
	// checked by the state verifier by type.
	StateVerifyLedgerEntriesCount *prometheus.GaugeVec

	// LedgerStatsCounter exposes ledger stats counters (like number of ops/changes).
	LedgerStatsCounter *prometheus.CounterVec

	// ProcessorsRunDuration exposes processors run durations.
	// Deprecated in favour of: ProcessorsRunDurationSummary.
	ProcessorsRunDuration *prometheus.CounterVec

	// ProcessorsRunDurationSummary exposes processors run durations.
	ProcessorsRunDurationSummary *prometheus.SummaryVec

	// LedgerFetchDurationSummary exposes a summary of durations required to
	// fetch data from ledger backend.
	LedgerFetchDurationSummary prometheus.Summary

	// CaptivediamcircleCoreSynced exposes synced status of Captive diamcircle-Core.
	// 1 if sync, 0 if not synced, -1 if unable to connect or HTTP server disabled.
	CaptivediamcircleCoreSynced prometheus.GaugeFunc

	// CaptiveCoreSupportedProtocolVersion exposes the maximum protocol version
	// supported by the running Captive-Core.
	CaptiveCoreSupportedProtocolVersion prometheus.GaugeFunc
}

type System interface {
	Run()
	RegisterMetrics(*prometheus.Registry)
	Metrics() Metrics
	StressTest(numTransactions, changesPerTransaction int) error
	VerifyRange(fromLedger, toLedger uint32, verifyState bool) error
	ReingestRange(ledgerRanges []history.LedgerRange, force bool) error
	BuildGenesisState() error
	Shutdown()
}

type system struct {
	metrics Metrics
	ctx     context.Context
	cancel  context.CancelFunc

	config Config

	historyQ history.IngestionQ
	runner   ProcessorRunnerInterface

	ledgerBackend  ledgerbackend.LedgerBackend
	historyAdapter historyArchiveAdapterInterface

	diamcircleCoreClient diamcircleCoreClient

	maxReingestRetries          int
	reingestRetryBackoffSeconds int
	wg                          sync.WaitGroup

	// stateVerificationRunning is true when verification routine is currently
	// running.
	stateVerificationMutex sync.Mutex
	// number of consecutive state verification runs which encountered errors
	stateVerificationErrors  int
	stateVerificationRunning bool
	disableStateVerification bool

	checkpointManager historyarchive.CheckpointManager
}

func NewSystem(config Config) (System, error) {
	ctx, cancel := context.WithCancel(context.Background())

	archive, err := historyarchive.Connect(
		config.HistoryArchiveURL,
		historyarchive.ConnectOptions{
			Context:             ctx,
			NetworkPassphrase:   config.NetworkPassphrase,
			CheckpointFrequency: config.CheckpointFrequency,
		},
	)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "error creating history archive")
	}

	var ledgerBackend ledgerbackend.LedgerBackend
	if config.EnableCaptiveCore {
		if len(config.RemoteCaptiveCoreURL) > 0 {
			ledgerBackend, err = ledgerbackend.NewRemoteCaptive(config.RemoteCaptiveCoreURL)
			if err != nil {
				cancel()
				return nil, errors.Wrap(err, "error creating captive core backend")
			}
		} else {
			logger := log.WithField("subservice", "diamcircle-core")
			ledgerBackend, err = ledgerbackend.NewCaptive(
				ledgerbackend.CaptiveCoreConfig{
					BinaryPath:          config.CaptiveCoreBinaryPath,
					StoragePath:         config.CaptiveCoreStoragePath,
					Toml:                config.CaptiveCoreToml,
					NetworkPassphrase:   config.NetworkPassphrase,
					HistoryArchiveURLs:  []string{config.HistoryArchiveURL},
					CheckpointFrequency: config.CheckpointFrequency,
					LedgerHashStore:     ledgerbackend.NewauroraDBLedgerHashStore(config.HistorySession),
					Log:                 logger,
					Context:             ctx,
				},
			)
			if err != nil {
				cancel()
				return nil, errors.Wrap(err, "error creating captive core backend")
			}
		}
	} else {
		coreSession := config.CoreSession.Clone()
		ledgerBackend, err = ledgerbackend.NewDatabaseBackendFromSession(coreSession, config.NetworkPassphrase)
		if err != nil {
			cancel()
			return nil, errors.Wrap(err, "error creating ledger backend")
		}
	}

	historyQ := &history.Q{config.HistorySession.Clone()}

	historyAdapter := newHistoryArchiveAdapter(archive)

	system := &system{
		cancel:                      cancel,
		config:                      config,
		ctx:                         ctx,
		disableStateVerification:    config.DisableStateVerification,
		historyAdapter:              historyAdapter,
		historyQ:                    historyQ,
		ledgerBackend:               ledgerBackend,
		maxReingestRetries:          config.MaxReingestRetries,
		reingestRetryBackoffSeconds: config.ReingestRetryBackoffSeconds,
		diamcircleCoreClient: &diamcirclecore.Client{
			URL: config.diamcircleCoreURL,
		},
		runner: &ProcessorRunner{
			ctx:            ctx,
			config:         config,
			historyQ:       historyQ,
			historyAdapter: historyAdapter,
		},
		checkpointManager: historyarchive.NewCheckpointManager(config.CheckpointFrequency),
	}

	system.initMetrics()
	return system, nil
}

func (s *system) initMetrics() {
	s.metrics.MaxSupportedProtocolVersion = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "aurora", Subsystem: "ingest", Name: "max_supported_protocol_version",
		Help: "the maximum protocol version supported by this version.",
	})

	s.metrics.MaxSupportedProtocolVersion.Set(float64(MaxSupportedProtocolVersion))

	s.metrics.LocalLatestLedger = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "aurora", Subsystem: "ingest", Name: "local_latest_ledger",
		Help: "sequence number of the latest ledger ingested by this ingesting instance",
	})

	s.metrics.LedgerIngestionDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "aurora", Subsystem: "ingest", Name: "ledger_ingestion_duration_seconds",
		Help: "ledger ingestion durations, sliding window = 10m",
	})

	s.metrics.LedgerIngestionTradeAggregationDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "aurora", Subsystem: "ingest", Name: "ledger_ingestion_trade_aggregation_duration_seconds",
		Help: "ledger ingestion trade aggregation rebuild durations, sliding window = 10m",
	})

	s.metrics.StateVerifyDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "aurora", Subsystem: "ingest", Name: "state_verify_duration_seconds",
		Help: "state verification durations, sliding window = 10m",
	})

	s.metrics.StateInvalidGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "state_invalid",
			Help: "equals 1 if state invalid, 0 otherwise",
		},
		func() float64 {
			invalid, err := s.historyQ.CloneIngestionQ().GetExpStateInvalid(s.ctx)
			if err != nil {
				log.WithError(err).Error("Error in initMetrics/GetExpStateInvalid")
				return 0
			}
			invalidFloat := float64(0)
			if invalid {
				invalidFloat = 1
			}
			return invalidFloat
		},
	)

	s.metrics.StateVerifyLedgerEntriesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "state_verify_ledger_entries",
			Help: "number of ledger entries downloaded from buckets in a single state verifier run",
		},
		[]string{"type"},
	)

	s.metrics.LedgerStatsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "ledger_stats_total",
			Help: "counters of different ledger stats",
		},
		[]string{"type"},
	)

	s.metrics.ProcessorsRunDuration = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "processor_run_duration_seconds_total",
			Help: "run durations of ingestion processors",
		},
		[]string{"name"},
	)

	s.metrics.ProcessorsRunDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "processor_run_duration_seconds",
			Help: "run durations of ingestion processors, sliding window = 10m",
		},
		[]string{"name"},
	)

	s.metrics.LedgerFetchDurationSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "ledger_fetch_duration_seconds",
			Help: "duration of fetching ledgers from ledger backend, sliding window = 10m",
		},
	)

	s.metrics.CaptivediamcircleCoreSynced = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "captive_diamcircle_core_synced",
			Help: "1 if sync, 0 if not synced, -1 if unable to connect or HTTP server disabled.",
		},
		func() float64 {
			if !s.config.EnableCaptiveCore || (s.config.CaptiveCoreToml.HTTPPort == 0) {
				return -1
			}

			client := diamcirclecore.Client{
				HTTP: &http.Client{
					Timeout: 2 * time.Second,
				},
				URL: fmt.Sprintf("http://localhost:%d", s.config.CaptiveCoreToml.HTTPPort),
			}

			info, err := client.Info(s.ctx)
			if err != nil {
				log.WithError(err).Error("Cannot connect to Captive diamcircle-Core HTTP server")
				return -1
			}

			if info.IsSynced() {
				return 1
			} else {
				return 0
			}
		},
	)

	s.metrics.CaptiveCoreSupportedProtocolVersion = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "ingest", Name: "captive_diamcircle_core_supported_protocol_version",
			Help: "determines the supported version of the protocol by Captive-Core",
		},
		func() float64 {
			if !s.config.EnableCaptiveCore || (s.config.CaptiveCoreToml.HTTPPort == 0) {
				return -1
			}

			client := diamcirclecore.Client{
				HTTP: &http.Client{
					Timeout: 2 * time.Second,
				},
				URL: fmt.Sprintf("http://localhost:%d", s.config.CaptiveCoreToml.HTTPPort),
			}

			info, err := client.Info(s.ctx)
			if err != nil {
				log.WithError(err).Error("Cannot connect to Captive diamcircle-Core HTTP server")
				return -1
			}

			return float64(info.Info.ProtocolVersion)
		},
	)
}

func (s *system) Metrics() Metrics {
	return s.metrics
}

// RegisterMetrics registers the prometheus metrics
func (s *system) RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(s.metrics.MaxSupportedProtocolVersion)
	registry.MustRegister(s.metrics.LocalLatestLedger)
	registry.MustRegister(s.metrics.LedgerIngestionDuration)
	registry.MustRegister(s.metrics.LedgerIngestionTradeAggregationDuration)
	registry.MustRegister(s.metrics.StateVerifyDuration)
	registry.MustRegister(s.metrics.StateInvalidGauge)
	registry.MustRegister(s.metrics.LedgerStatsCounter)
	registry.MustRegister(s.metrics.ProcessorsRunDuration)
	registry.MustRegister(s.metrics.ProcessorsRunDurationSummary)
	registry.MustRegister(s.metrics.CaptivediamcircleCoreSynced)
	registry.MustRegister(s.metrics.CaptiveCoreSupportedProtocolVersion)
	registry.MustRegister(s.metrics.LedgerFetchDurationSummary)
	registry.MustRegister(s.metrics.StateVerifyLedgerEntriesCount)
}

// Run starts ingestion system. Ingestion system supports distributed ingestion
// that means that aurora ingestion can be running on multiple machines and
// only one, random node will lead the ingestion.
//
// It needs to support cartesian product of the following run scenarios cases:
// - Init from empty state (1a) and resuming from existing state (1b).
// - Ingestion system version has been upgraded (2a) or not (2b).
// - Current node is leading ingestion (3a) or not (3b).
//
// We always clear state when ingestion system is upgraded so 2a and 2b are
// included in 1a.
//
// We ensure that only one instance is a leader because in each round instances
// try to acquire a lock on `LastLedgerIngest value in key value store and only
// one instance will be able to acquire it. This happens in both initial processing
// and ledger processing. So this solves 3a and 3b in both 1a and 1b.
//
// Finally, 1a and 1b are tricky because we need to keep the latest version
// of order book graph in memory of each aurora instance. To solve this:
// * For state init:
//   * If instance is a leader, we update the order book graph by running state
//     pipeline normally.
//   * If instance is NOT a leader, we build a graph from offers present in a
//     database. We completely omit state pipeline in this case.
// * For resuming:
//   * If instances is a leader, it runs full ledger pipeline, including updating
//     a database.
//   * If instances is a NOT leader, it runs ledger pipeline without updating a
//     a database so order book graph is updated but database is not overwritten.
func (s *system) Run() {
	s.runStateMachine(startState{})
}

func (s *system) StressTest(numTransactions, changesPerTransaction int) error {
	if numTransactions <= 0 {
		return errors.New("transactions must be positive")
	}
	if changesPerTransaction <= 0 {
		return errors.New("changes per transaction must be positive")
	}

	s.runner.EnableMemoryStatsLogging()
	s.ledgerBackend = &fakeLedgerBackend{
		numTransactions:       numTransactions,
		changesPerTransaction: changesPerTransaction,
	}
	return s.runStateMachine(stressTestState{})
}

// VerifyRange runs the ingestion pipeline on the range of ledgers. When
// verifyState is true it verifies the state when ingestion is complete.
func (s *system) VerifyRange(fromLedger, toLedger uint32, verifyState bool) error {
	return s.runStateMachine(verifyRangeState{
		fromLedger:  fromLedger,
		toLedger:    toLedger,
		verifyState: verifyState,
	})
}

func validateRanges(ledgerRanges []history.LedgerRange) error {
	for i, cur := range ledgerRanges {
		if cur.StartSequence > cur.EndSequence {
			return errors.Errorf("Invalid range: %v from > to", cur)
		}
		if cur.StartSequence == 0 {
			return errors.Errorf("Invalid range: %v genesis ledger starts at 1", cur)
		}
		if i == 0 {
			continue
		}
		prev := ledgerRanges[i-1]
		if prev.EndSequence >= cur.StartSequence {
			return errors.Errorf("ranges are not sorted prevRange %v curRange %v", prev, cur)
		}
	}
	return nil
}

// ReingestRange runs the ingestion pipeline on the range of ledgers ingesting
// history data only.
func (s *system) ReingestRange(ledgerRanges []history.LedgerRange, force bool) error {
	if err := validateRanges(ledgerRanges); err != nil {
		return err
	}
	for _, cur := range ledgerRanges {
		run := func() error {
			return s.runStateMachine(reingestHistoryRangeState{
				fromLedger: cur.StartSequence,
				toLedger:   cur.EndSequence,
				force:      force,
			})
		}
		err := run()
		for retry := 0; err != nil && retry < s.maxReingestRetries; retry++ {
			log.Warnf("reingest range [%d, %d] failed (%s), retrying", cur.StartSequence, cur.EndSequence, err.Error())
			time.Sleep(time.Second * time.Duration(s.reingestRetryBackoffSeconds))
			err = run()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildGenesisState runs the ingestion pipeline on genesis ledger. Transitions
// to stopState when done.
func (s *system) BuildGenesisState() error {
	return s.runStateMachine(buildState{
		checkpointLedger: 1,
		stop:             true,
	})
}

func (s *system) runStateMachine(cur stateMachineNode) error {
	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		s.wg.Wait()
	}()

	log.WithFields(logpkg.F{"current_state": cur}).Info("Ingestion system initial state")

	for {
		// Every node in the state machine is responsible for
		// creating and disposing its own transaction.
		// We should never enter a new state with the transaction
		// from the previous state.
		if s.historyQ.GetTx() != nil {
			panic("unexpected transaction")
		}

		next, err := cur.run(s)
		if err != nil {
			logger := log.WithFields(logpkg.F{
				"error":         err,
				"current_state": cur,
				"next_state":    next.node,
			})
			if isCancelledError(err) {
				// We only expect context.Canceled errors to occur when aurora is shutting down
				// so we log these errors using the info log level
				logger.Info("Error in ingestion state machine")
			} else {
				logger.Error("Error in ingestion state machine")
			}
		}

		// Exit after processing shutdownState
		if next.node == (stopState{}) {
			log.Info("Shut down")
			return err
		}

		select {
		case <-s.ctx.Done():
			log.Info("Received shut down signal...")
			return nil
		case <-time.After(next.sleepDuration):
		}

		log.WithFields(logpkg.F{
			"current_state": cur,
			"next_state":    next.node,
		}).Info("Ingestion system state machine transition")
		cur = next.node
	}
}

func (s *system) maybeVerifyState(lastIngestedLedger uint32) {
	stateInvalid, err := s.historyQ.GetExpStateInvalid(s.ctx)
	if err != nil && !isCancelledError(err) {
		log.WithField("err", err).Error("Error getting state invalid value")
	}

	// Run verification routine only when...
	if !stateInvalid && // state has not been proved to be invalid...
		!s.disableStateVerification && // state verification is not disabled...
		s.checkpointManager.IsCheckpoint(lastIngestedLedger) { // it's a checkpoint ledger.
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			err := s.verifyState(true)
			if err != nil {
				if isCancelledError(err) {
					return
				}

				errorCount := s.incrementStateVerificationErrors()
				switch errors.Cause(err).(type) {
				case ingest.StateError:
					markStateInvalid(s.ctx, s.historyQ, err)
				default:
					logger := log.WithField("err", err).Warn
					if errorCount >= stateVerificationErrorThreshold {
						logger = log.WithField("err", err).Error
					}
					logger("State verification errored")
				}
			} else {
				s.resetStateVerificationErrors()
			}
		}()
	}
}

func (s *system) incrementStateVerificationErrors() int {
	s.stateVerificationMutex.Lock()
	defer s.stateVerificationMutex.Unlock()
	s.stateVerificationErrors++
	return s.stateVerificationErrors
}

func (s *system) resetStateVerificationErrors() {
	s.stateVerificationMutex.Lock()
	defer s.stateVerificationMutex.Unlock()
	s.stateVerificationErrors = 0
}

func (s *system) updateCursor(ledgerSequence uint32) error {
	if s.diamcircleCoreClient == nil || s.config.EnableCaptiveCore {
		return nil
	}

	cursor := defaultCoreCursorName
	if s.config.diamcircleCoreCursor != "" {
		cursor = s.config.diamcircleCoreCursor
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()
	err := s.diamcircleCoreClient.SetCursor(ctx, cursor, int32(ledgerSequence))
	if err != nil {
		return errors.Wrap(err, "Setting diamcircle-core cursor failed")
	}

	return nil
}

func (s *system) Shutdown() {
	log.Info("Shutting down ingestion system...")
	s.stateVerificationMutex.Lock()
	if s.stateVerificationRunning {
		log.Info("Shutting down state verifier...")
	}
	s.stateVerificationMutex.Unlock()
	s.cancel()
	// wait for ingestion state machine to terminate
	s.wg.Wait()
	if err := s.ledgerBackend.Close(); err != nil {
		log.WithError(err).Info("could not close ledger backend")
	}
}

func markStateInvalid(ctx context.Context, historyQ history.IngestionQ, err error) {
	log.WithField("err", err).Error("STATE IS INVALID!")
	q := historyQ.CloneIngestionQ()
	if err := q.UpdateExpStateInvalid(ctx, true); err != nil {
		log.WithField("err", err).Error(updateExpStateInvalidErrMsg)
	}
}

func isCancelledError(err error) bool {
	cause := errors.Cause(err)
	return cause == context.Canceled || cause == db.ErrCancelled
}
