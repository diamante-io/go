package ingest

import (
	"bytes"
	"context"
	"fmt"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest/processors"
	"go/support/errors"
	"go/xdr"
)

type ingestionSource int

const (
	_                    = iota
	historyArchiveSource = ingestionSource(iota)
	ledgerSource         = ingestionSource(iota)
	logFrequency         = 50000
)

type auroraChangeProcessor interface {
	processors.ChangeProcessor
	Commit(context.Context) error
}

type auroraTransactionProcessor interface {
	processors.LedgerTransactionProcessor
	Commit(context.Context) error
}

type statsChangeProcessor struct {
	*ingest.StatsChangeProcessor
}

func (statsChangeProcessor) Commit(ctx context.Context) error {
	return nil
}

type statsLedgerTransactionProcessor struct {
	*processors.StatsLedgerTransactionProcessor
}

func (statsLedgerTransactionProcessor) Commit(ctx context.Context) error {
	return nil
}

type ledgerStats struct {
	changeStats          ingest.StatsChangeProcessorResults
	changeDurations      processorsRunDurations
	transactionStats     processors.StatsLedgerTransactionProcessorResults
	transactionDurations processorsRunDurations
	tradeStats           processors.TradeStats
}

type ProcessorRunnerInterface interface {
	SetHistoryAdapter(historyAdapter historyArchiveAdapterInterface)
	EnableMemoryStatsLogging()
	DisableMemoryStatsLogging()
	RunGenesisStateIngestion() (ingest.StatsChangeProcessorResults, error)
	RunHistoryArchiveIngestion(
		checkpointLedger uint32,
		ledgerProtocolVersion uint32,
		bucketListHash xdr.Hash,
	) (ingest.StatsChangeProcessorResults, error)
	RunTransactionProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
		transactionStats processors.StatsLedgerTransactionProcessorResults,
		transactionDurations processorsRunDurations,
		tradeStats processors.TradeStats,
		err error,
	)
	RunAllProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
		stats ledgerStats,
		err error,
	)
}

var _ ProcessorRunnerInterface = (*ProcessorRunner)(nil)

type ProcessorRunner struct {
	config Config

	ctx            context.Context
	historyQ       history.IngestionQ
	historyAdapter historyArchiveAdapterInterface
	logMemoryStats bool
}

func (s *ProcessorRunner) SetHistoryAdapter(historyAdapter historyArchiveAdapterInterface) {
	s.historyAdapter = historyAdapter
}

func (s *ProcessorRunner) EnableMemoryStatsLogging() {
	s.logMemoryStats = true
}

func (s *ProcessorRunner) DisableMemoryStatsLogging() {
	s.logMemoryStats = false
}

func buildChangeProcessor(
	historyQ history.IngestionQ,
	changeStats *ingest.StatsChangeProcessor,
	source ingestionSource,
	ledgerSequence uint32,
) *groupChangeProcessors {
	statsChangeProcessor := &statsChangeProcessor{
		StatsChangeProcessor: changeStats,
	}

	useLedgerCache := source == ledgerSource
	return newGroupChangeProcessors([]auroraChangeProcessor{
		statsChangeProcessor,
		processors.NewAccountDataProcessor(historyQ),
		processors.NewAccountsProcessor(historyQ),
		processors.NewOffersProcessor(historyQ, ledgerSequence),
		processors.NewAssetStatsProcessor(historyQ, useLedgerCache),
		processors.NewSignersProcessor(historyQ, useLedgerCache),
		processors.NewTrustLinesProcessor(historyQ),
		processors.NewClaimableBalancesChangeProcessor(historyQ),
		processors.NewLiquidityPoolsChangeProcessor(historyQ, ledgerSequence),
	})
}

func (s *ProcessorRunner) buildTransactionProcessor(
	ledgerTransactionStats *processors.StatsLedgerTransactionProcessor,
	tradeProcessor *processors.TradeProcessor,
	ledger xdr.LedgerHeaderHistoryEntry,
) *groupTransactionProcessors {
	statsLedgerTransactionProcessor := &statsLedgerTransactionProcessor{
		StatsLedgerTransactionProcessor: ledgerTransactionStats,
	}
	*tradeProcessor = *processors.NewTradeProcessor(s.historyQ, ledger)
	sequence := uint32(ledger.Header.LedgerSeq)
	return newGroupTransactionProcessors([]auroraTransactionProcessor{
		statsLedgerTransactionProcessor,
		processors.NewEffectProcessor(s.historyQ, sequence),
		processors.NewLedgerProcessor(s.historyQ, ledger, CurrentVersion),
		processors.NewOperationProcessor(s.historyQ, sequence),
		tradeProcessor,
		processors.NewParticipantsProcessor(s.historyQ, sequence),
		processors.NewTransactionProcessor(s.historyQ, sequence),
		processors.NewClaimableBalancesTransactionProcessor(s.historyQ, sequence),
		processors.NewLiquidityPoolsTransactionProcessor(s.historyQ, sequence),
	})
}

// checkIfProtocolVersionSupported checks if this aurora version supports the
// protocol version of a ledger with the given sequence number.
func (s *ProcessorRunner) checkIfProtocolVersionSupported(ledgerProtocolVersion uint32) error {
	if ledgerProtocolVersion > MaxSupportedProtocolVersion {
		return fmt.Errorf(
			"This aurora version does not support protocol version %d. "+
				"The latest supported protocol version is %d. Please upgrade to the latest aurora version.",
			ledgerProtocolVersion,
			MaxSupportedProtocolVersion,
		)
	}

	return nil
}

// validateBucketList validates if the bucket list hash in history archive
// matches the one in corresponding ledger header in diamcircle-core backend.
// This gives you full security if data in diamcircle-core backend can be trusted
// (ex. you run it in your infrastructure).
// The hashes of actual buckets of this HAS file are checked using
// historyarchive.XdrStream.SetExpectedHash (this is done in
// CheckpointChangeReader).
func (s *ProcessorRunner) validateBucketList(ledgerSequence uint32, ledgerBucketHashList xdr.Hash) error {
	historyBucketListHash, err := s.historyAdapter.BucketListHash(ledgerSequence)
	if err != nil {
		return errors.Wrap(err, "Error getting bucket list hash")
	}

	if !bytes.Equal(historyBucketListHash[:], ledgerBucketHashList[:]) {
		return fmt.Errorf(
			"Bucket list hash of history archive and ledger header does not match: %#x %#x",
			historyBucketListHash,
			ledgerBucketHashList,
		)
	}

	return nil
}

func (s *ProcessorRunner) RunGenesisStateIngestion() (ingest.StatsChangeProcessorResults, error) {
	return s.RunHistoryArchiveIngestion(1, 0, xdr.Hash{})
}

func (s *ProcessorRunner) RunHistoryArchiveIngestion(
	checkpointLedger uint32,
	ledgerProtocolVersion uint32,
	bucketListHash xdr.Hash,
) (ingest.StatsChangeProcessorResults, error) {
	changeStats := ingest.StatsChangeProcessor{}
	changeProcessor := buildChangeProcessor(s.historyQ, &changeStats, historyArchiveSource, checkpointLedger)

	if checkpointLedger == 1 {
		if err := changeProcessor.ProcessChange(s.ctx, ingest.GenesisChange(s.config.NetworkPassphrase)); err != nil {
			return changeStats.GetResults(), errors.Wrap(err, "Error ingesting genesis ledger")
		}
	} else {
		if err := s.checkIfProtocolVersionSupported(ledgerProtocolVersion); err != nil {
			return changeStats.GetResults(), errors.Wrap(err, "Error while checking for supported protocol version")
		}

		if err := s.validateBucketList(checkpointLedger, bucketListHash); err != nil {
			return changeStats.GetResults(), errors.Wrap(err, "Error validating bucket list from HAS")
		}

		changeReader, err := s.historyAdapter.GetState(s.ctx, checkpointLedger)
		if err != nil {
			return changeStats.GetResults(), errors.Wrap(err, "Error creating HAS reader")
		}

		defer changeReader.Close()

		log.WithField("sequence", checkpointLedger).
			Info("Processing entries from History Archive Snapshot")

		err = processors.StreamChanges(s.ctx, changeProcessor, newloggingChangeReader(
			changeReader,
			"historyArchive",
			checkpointLedger,
			logFrequency,
			s.logMemoryStats,
		))
		if err != nil {
			return changeStats.GetResults(), errors.Wrap(err, "Error streaming changes from HAS")
		}
	}

	if err := changeProcessor.Commit(s.ctx); err != nil {
		return changeStats.GetResults(), errors.Wrap(err, "Error commiting changes from processor")
	}

	return changeStats.GetResults(), nil
}

func (s *ProcessorRunner) runChangeProcessorOnLedger(
	changeProcessor auroraChangeProcessor, ledger xdr.LedgerCloseMeta,
) error {
	var changeReader ingest.ChangeReader
	var err error
	changeReader, err = ingest.NewLedgerChangeReaderFromLedgerCloseMeta(s.config.NetworkPassphrase, ledger)
	if err != nil {
		return errors.Wrap(err, "Error creating ledger change reader")
	}
	changeReader = newloggingChangeReader(
		changeReader,
		"ledger",
		ledger.LedgerSequence(),
		logFrequency,
		s.logMemoryStats,
	)
	if err = processors.StreamChanges(s.ctx, changeProcessor, changeReader); err != nil {
		return errors.Wrap(err, "Error streaming changes from ledger")
	}

	err = changeProcessor.Commit(s.ctx)
	if err != nil {
		return errors.Wrap(err, "Error commiting changes from processor")
	}

	return nil
}

func (s *ProcessorRunner) RunTransactionProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
	transactionStats processors.StatsLedgerTransactionProcessorResults,
	transactionDurations processorsRunDurations,
	tradeStats processors.TradeStats,
	err error,
) {
	var (
		ledgerTransactionStats processors.StatsLedgerTransactionProcessor
		tradeProcessor         processors.TradeProcessor
		transactionReader      *ingest.LedgerTransactionReader
	)

	transactionReader, err = ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(s.config.NetworkPassphrase, ledger)
	if err != nil {
		err = errors.Wrap(err, "Error creating ledger reader")
		return
	}

	if err = s.checkIfProtocolVersionSupported(ledger.ProtocolVersion()); err != nil {
		err = errors.Wrap(err, "Error while checking for supported protocol version")
		return
	}

	groupTransactionProcessors := s.buildTransactionProcessor(
		&ledgerTransactionStats, &tradeProcessor, transactionReader.GetHeader())
	err = processors.StreamLedgerTransactions(s.ctx, groupTransactionProcessors, transactionReader)
	if err != nil {
		err = errors.Wrap(err, "Error streaming changes from ledger")
		return
	}

	err = groupTransactionProcessors.Commit(s.ctx)
	if err != nil {
		err = errors.Wrap(err, "Error committing changes from processor")
		return
	}

	transactionStats = ledgerTransactionStats.GetResults()
	transactionDurations = groupTransactionProcessors.processorsRunDurations
	tradeStats = tradeProcessor.GetStats()
	return
}

func (s *ProcessorRunner) RunAllProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
	stats ledgerStats,
	err error,
) {
	changeStatsProcessor := ingest.StatsChangeProcessor{}

	if err = s.checkIfProtocolVersionSupported(ledger.ProtocolVersion()); err != nil {
		err = errors.Wrap(err, "Error while checking for supported protocol version")
		return
	}

	groupChangeProcessors := buildChangeProcessor(s.historyQ, &changeStatsProcessor, ledgerSource, ledger.LedgerSequence())
	err = s.runChangeProcessorOnLedger(groupChangeProcessors, ledger)
	if err != nil {
		return
	}

	stats.changeStats = changeStatsProcessor.GetResults()
	stats.changeDurations = groupChangeProcessors.processorsRunDurations

	stats.transactionStats, stats.transactionDurations, stats.tradeStats, err =
		s.RunTransactionProcessorsOnLedger(ledger)
	if err != nil {
		return
	}

	return
}
