package txsub

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/txsub/sequence"
	"go/support/log"
	"go/xdr"

	"github.com/prometheus/client_golang/prometheus"
)

type auroraDB interface {
	GetLatestHistoryLedger(ctx context.Context) (uint32, error)
	TransactionByHash(ctx context.Context, dest interface{}, hash string) error
	TransactionsByHashesSinceLedger(ctx context.Context, hashes []string, sinceLedgerSeq uint32) ([]history.Transaction, error)
	GetSequenceNumbers(ctx context.Context, addresses []string) (map[string]uint64, error)
	BeginTx(*sql.TxOptions) error
	Rollback() error
	NoRows(error) bool
}

// System represents a completely configured transaction submission system.
// Its methods tie together the various pieces used to reliably submit transactions
// to a diamcircle-core instance.
type System struct {
	initializer sync.Once

	tickMutex      sync.Mutex
	tickInProgress bool

	accountSeqPollInterval time.Duration

	DB                func(context.Context) auroraDB
	Pending           OpenSubmissionList
	Submitter         Submitter
	SubmissionQueue   *sequence.Manager
	SubmissionTimeout time.Duration
	Log               *log.Entry

	Metrics struct {
		// SubmissionDuration exposes timing metrics about the rate and latency of
		// submissions to diamcircle-core
		SubmissionDuration prometheus.Summary

		// BufferedSubmissionGauge tracks the count of submissions buffered
		// behind this system's SubmissionQueue
		BufferedSubmissionsGauge prometheus.Gauge

		// OpenSubmissionsGauge tracks the count of "open" submissions (i.e.
		// submissions whose transactions haven't been confirmed successful or failed
		OpenSubmissionsGauge prometheus.Gauge

		// FailedSubmissionsCounter tracks the rate of failed transactions that have
		// been submitted to this process
		FailedSubmissionsCounter prometheus.Counter

		// SuccessfulSubmissionsCounter tracks the rate of successful transactions that
		// have been submitted to this process
		SuccessfulSubmissionsCounter prometheus.Counter

		// V0TransactionsCounter tracks the rate of v0 transaction envelopes that
		// have been submitted to this process
		V0TransactionsCounter prometheus.Counter

		// V1TransactionsCounter tracks the rate of v1 transaction envelopes that
		// have been submitted to this process
		V1TransactionsCounter prometheus.Counter

		// FeeBumpTransactionsCounter tracks the rate of fee bump transaction envelopes that
		// have been submitted to this process
		FeeBumpTransactionsCounter prometheus.Counter
	}
}

// RegisterMetrics registers the prometheus metrics
func (sys *System) RegisterMetrics(registry *prometheus.Registry) {
	registry.MustRegister(sys.Metrics.SubmissionDuration)
	registry.MustRegister(sys.Metrics.BufferedSubmissionsGauge)
	registry.MustRegister(sys.Metrics.OpenSubmissionsGauge)
	registry.MustRegister(sys.Metrics.FailedSubmissionsCounter)
	registry.MustRegister(sys.Metrics.SuccessfulSubmissionsCounter)
	registry.MustRegister(sys.Metrics.V0TransactionsCounter)
	registry.MustRegister(sys.Metrics.V1TransactionsCounter)
	registry.MustRegister(sys.Metrics.FeeBumpTransactionsCounter)
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (sys *System) Submit(
	ctx context.Context,
	rawTx string,
	envelope xdr.TransactionEnvelope,
	hash string,
) (result <-chan Result) {
	sys.Init()
	response := make(chan Result, 1)
	result = response

	db := sys.DB(ctx)
	// The database doesn't (yet) store muxed accounts, so we query
	// the corresponding AccountId
	sourceAccount := envelope.SourceAccount().ToAccountId()
	sourceAddress := sourceAccount.Address()

	sys.Log.Ctx(ctx).WithFields(log.F{
		"hash":    hash,
		"tx_type": envelope.Type.String(),
		"tx":      rawTx,
	}).Info("Processing transaction")

	if envelope.SeqNum() < 0 {
		sys.finish(ctx, hash, response, Result{Err: ErrBadSequence})
		return
	}

	tx, sequenceNumber, err := checkTxAlreadyExists(ctx, db, hash, sourceAddress)
	if err == nil {
		sys.Log.Ctx(ctx).WithField("hash", hash).Info("Found submission result in a DB")
		sys.finish(ctx, hash, response, Result{Transaction: tx})
		return
	}
	if err != ErrNoResults {
		sys.Log.Ctx(ctx).WithField("hash", hash).Info("Error getting submission result from a DB")
		sys.finish(ctx, hash, response, Result{Transaction: tx, Err: err})
		return
	}

	// queue the submission and get the channel that will emit when
	// submission is valid
	seq := sys.SubmissionQueue.Push(sourceAddress, uint64(envelope.SeqNum()))

	// update the submission queue with the source accounts current sequence value
	// which will cause the channel returned by Push() to emit if possible.
	sys.SubmissionQueue.Update(map[string]uint64{
		sourceAddress: sequenceNumber,
	})

	select {
	case err := <-seq:
		if err == sequence.ErrBadSequence {
			// convert the internal only ErrBadSequence into the FailedTransactionError
			err = ErrBadSequence
		}

		if err != nil {
			sys.finish(ctx, hash, response, Result{Err: err})
			return
		}

		sr := sys.submitOnce(ctx, rawTx)
		sys.updateTransactionTypeMetrics(envelope)

		// if submission succeeded
		if sr.Err == nil {
			// add transactions to open list
			sys.Pending.Add(ctx, hash, response)
			// update the submission queue, allowing the next submission to proceed
			sys.SubmissionQueue.Update(map[string]uint64{
				sourceAddress: uint64(envelope.SeqNum()),
			})
			return
		}

		// any error other than "txBAD_SEQ" is a failure
		isBad, err := sr.IsBadSeq()
		if err != nil {
			sys.finish(ctx, hash, response, Result{Err: err})
			return
		}

		if !isBad {
			sys.finish(ctx, hash, response, Result{Err: sr.Err})
			return
		}
		if err = sys.waitUntilAccountSequence(ctx, db, sourceAddress, uint64(envelope.SeqNum())); err != nil {
			sys.finish(ctx, hash, response, Result{Err: err})
			return
		}

		// If error is txBAD_SEQ, check for the result again
		tx, err = txResultByHash(ctx, db, hash)
		if err == nil {
			// If the found use it as the result
			sys.finish(ctx, hash, response, Result{Transaction: tx})
		} else {
			// finally, return the bad_seq error if no result was found on 2nd attempt
			sys.finish(ctx, hash, response, Result{Err: sr.Err})
		}

	case <-ctx.Done():
		sys.finish(ctx, hash, response, Result{Err: sys.deriveTxSubError(ctx)})
	}

	return
}

// waitUntilAccountSequence blocks until either the context times out or the sequence number of the
// given source account is greater than or equal to `seq`
func (sys *System) waitUntilAccountSequence(ctx context.Context, db auroraDB, sourceAddress string, seq uint64) error {
	timer := time.NewTimer(sys.accountSeqPollInterval)
	defer timer.Stop()

	for {
		sequenceNumbers, err := db.GetSequenceNumbers(ctx, []string{sourceAddress})
		if err != nil {
			sys.Log.Ctx(ctx).
				WithError(err).
				WithField("sourceAddress", sourceAddress).
				Warn("cannot fetch sequence number")
		} else {
			num, ok := sequenceNumbers[sourceAddress]
			if !ok {
				sys.Log.Ctx(ctx).
					WithField("sequenceNumbers", sequenceNumbers).
					WithField("sourceAddress", sourceAddress).
					Warn("missing sequence number for account")
			}
			if num >= seq {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return sys.deriveTxSubError(ctx)
		case <-timer.C:
			timer.Reset(sys.accountSeqPollInterval)
		}
	}
}

func (sys *System) deriveTxSubError(ctx context.Context) error {
	if ctx.Err() == context.Canceled {
		return ErrCanceled
	}
	return ErrTimeout
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (sys *System) submitOnce(ctx context.Context, env string) SubmissionResult {
	// submit to diamcircle-core
	sr := sys.Submitter.Submit(ctx, env)
	sys.Metrics.SubmissionDuration.Observe(float64(sr.Duration.Seconds()))

	// if received or duplicate, add to the open submissions list
	if sr.Err == nil {
		sys.Metrics.SuccessfulSubmissionsCounter.Inc()
	} else {
		sys.Metrics.FailedSubmissionsCounter.Inc()
	}

	return sr
}

func (sys *System) updateTransactionTypeMetrics(envelope xdr.TransactionEnvelope) {
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		sys.Metrics.V0TransactionsCounter.Inc()
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		sys.Metrics.V1TransactionsCounter.Inc()
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		sys.Metrics.FeeBumpTransactionsCounter.Inc()
	}
}

// setTickInProgress sets `tickInProgress` to `true` if it's
// `false`. Returns `true` if `tickInProgress` has been switched
// to `true` inside this method and `Tick()` should continue.
func (sys *System) setTickInProgress(ctx context.Context) bool {
	sys.tickMutex.Lock()
	defer sys.tickMutex.Unlock()

	if sys.tickInProgress {
		logger := log.Ctx(ctx)
		logger.Info("ticking in progress")
		return false
	}

	sys.tickInProgress = true
	return true
}

func (sys *System) unsetTickInProgress() {
	sys.tickMutex.Lock()
	defer sys.tickMutex.Unlock()
	sys.tickInProgress = false
}

// Tick triggers the system to update itself with any new data available.
func (sys *System) Tick(ctx context.Context) {
	sys.Init()
	logger := log.Ctx(ctx)

	// Make sure Tick is not run concurrently
	if !sys.setTickInProgress(ctx) {
		return
	}

	defer sys.unsetTickInProgress()

	logger.
		WithField("queued", sys.SubmissionQueue.String()).
		Debug("ticking txsub system")

	db := sys.DB(ctx)
	options := &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}
	if err := db.BeginTx(options); err != nil {
		logger.WithError(err).Error("could not start repeatable read transaction for txsub tick")
		return
	}
	defer db.Rollback()

	addys := sys.SubmissionQueue.Addresses()
	if len(addys) > 0 {
		curSeq, err := db.GetSequenceNumbers(ctx, addys)
		if err != nil {
			logger.WithStack(err).Error(err)
			return
		} else {
			sys.SubmissionQueue.Update(curSeq)
		}
	}

	pending := sys.Pending.Pending(ctx)

	if len(pending) > 0 {
		latestLedger, err := db.GetLatestHistoryLedger(ctx)
		if err != nil {
			logger.WithError(err).Error("error getting latest history ledger")
			return
		}

		// In Tick we only check txs in a queue so those which did not have results before Tick
		// so we check for them in the last 5 mins of ledgers: 60.
		var sinceLedgerSeq int32 = int32(latestLedger) - 60
		if sinceLedgerSeq < 0 {
			sinceLedgerSeq = 0
		}

		txs, err := db.TransactionsByHashesSinceLedger(ctx, pending, uint32(sinceLedgerSeq))
		if err != nil && !db.NoRows(err) {
			logger.WithError(err).Error("error getting transactions by hashes")
			return
		}

		txMap := make(map[string]history.Transaction, len(txs))
		for _, tx := range txs {
			txMap[tx.TransactionHash] = tx
			if tx.InnerTransactionHash.Valid {
				txMap[tx.InnerTransactionHash.String] = tx
			}
		}

		for _, hash := range pending {
			tx, found := txMap[hash]
			if !found {
				continue
			}
			_, err := txResultFromHistory(tx)

			if err == nil {
				logger.WithField("hash", hash).Debug("finishing open submission")
				sys.Pending.Finish(ctx, hash, Result{Transaction: tx})
				continue
			}

			if _, ok := err.(*FailedTransactionError); ok {
				logger.WithField("hash", hash).Debug("finishing open submission")
				sys.Pending.Finish(ctx, hash, Result{Transaction: tx, Err: err})
				continue
			}

			if err != nil {
				logger.WithStack(err).Error(err)
			}
		}
	}

	stillOpen, err := sys.Pending.Clean(ctx, sys.SubmissionTimeout)
	if err != nil {
		logger.WithStack(err).Error(err)
		return
	}

	sys.Metrics.OpenSubmissionsGauge.Set(float64(stillOpen))
	sys.Metrics.BufferedSubmissionsGauge.Set(float64(sys.SubmissionQueue.Size()))
}

// Init initializes `sys`
func (sys *System) Init() {
	sys.initializer.Do(func() {
		sys.Log = log.DefaultLogger.WithField("service", "txsub.System")

		sys.Metrics.SubmissionDuration = prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "submission_duration_seconds",
			Help: "submission durations to diamcircle-Core, sliding window = 10m",
		})
		sys.Metrics.FailedSubmissionsCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "failed",
		})
		sys.Metrics.SuccessfulSubmissionsCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "succeeded",
		})
		sys.Metrics.OpenSubmissionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "open",
		})
		sys.Metrics.BufferedSubmissionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "buffered",
		})
		sys.Metrics.V0TransactionsCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "v0",
		})
		sys.Metrics.V1TransactionsCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "v1",
		})
		sys.Metrics.FeeBumpTransactionsCounter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "aurora", Subsystem: "txsub", Name: "feebump",
		})

		sys.accountSeqPollInterval = time.Second

		if sys.SubmissionTimeout == 0 {
			// HTTP clients in SDKs usually timeout in 60 seconds. We want SubmissionTimeout
			// to be lower than that to make sure that they read the response before the client
			// timeout.
			// 30 seconds is 6 ledgers (with avg. close time = 5 sec), enough for diamcircle-core
			// to drop the transaction if not added to the ledger and ask client to try again
			// by sending a Timeout response.
			sys.SubmissionTimeout = 30 * time.Second
		}
	})
}

func (sys *System) finish(ctx context.Context, hash string, response chan<- Result, r Result) {
	sys.Log.Ctx(ctx).
		WithField("result", fmt.Sprintf("%+v", r)).
		WithField("hash", hash).
		Info("Submission system result")
	response <- r
	close(response)
}
