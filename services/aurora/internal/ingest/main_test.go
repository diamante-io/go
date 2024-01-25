package ingest

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/historyarchive"
	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ingest/processors"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	logpkg "github.com/diamcircle/go/support/log"
	strtime "github.com/diamcircle/go/support/time"
	"github.com/diamcircle/go/xdr"
)

var (
	issuer   = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	usdAsset = xdr.Asset{
		Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
		AlphaNum4: &xdr.AlphaNum4{
			AssetCode: [4]byte{'u', 's', 'd', 0},
			Issuer:    issuer,
		},
	}

	nativeAsset = xdr.Asset{
		Type: xdr.AssetTypeAssetTypeNative,
	}

	eurAsset = xdr.Asset{
		Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
		AlphaNum4: &xdr.AlphaNum4{
			AssetCode: [4]byte{'e', 'u', 'r', 0},
			Issuer:    issuer,
		},
	}
	eurOffer = xdr.OfferEntry{
		SellerId: issuer,
		OfferId:  xdr.Int64(4),
		Buying:   eurAsset,
		Selling:  nativeAsset,
		Price: xdr.Price{
			N: 1,
			D: 1,
		},
		Flags:  1,
		Amount: xdr.Int64(500),
	}
	twoEurOffer = xdr.OfferEntry{
		SellerId: issuer,
		OfferId:  xdr.Int64(5),
		Buying:   eurAsset,
		Selling:  nativeAsset,
		Price: xdr.Price{
			N: 2,
			D: 1,
		},
		Flags:  2,
		Amount: xdr.Int64(500),
	}
)

func TestCheckVerifyStateVersion(t *testing.T) {
	assert.Equal(
		t,
		CurrentVersion,
		stateVerifierExpectedIngestionVersion,
		"State verifier is outdated, update it, then update stateVerifierExpectedIngestionVersion value",
	)
}

func TestNewSystem(t *testing.T) {
	config := Config{
		CoreSession:              &db.Session{DB: &sqlx.DB{}},
		HistorySession:           &db.Session{DB: &sqlx.DB{}},
		DisableStateVerification: true,
		HistoryArchiveURL:        "https://history.diamcircle.org/prd/core-live/core_live_001",
		CheckpointFrequency:      64,
	}

	sIface, err := NewSystem(config)
	assert.NoError(t, err)
	system := sIface.(*system)

	assert.Equal(t, config, system.config)
	assert.Equal(t, config.DisableStateVerification, system.disableStateVerification)

	assert.Equal(t, config, system.runner.(*ProcessorRunner).config)
	assert.Equal(t, system.ctx, system.runner.(*ProcessorRunner).ctx)
}

func TestStateMachineRunReturnsUnexpectedTransaction(t *testing.T) {
	historyQ := &mockDBQ{}
	system := &system{
		historyQ: historyQ,
		ctx:      context.Background(),
	}

	historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()

	assert.PanicsWithValue(t, "unexpected transaction", func() {
		system.Run()
	})
}

func TestStateMachineTransition(t *testing.T) {
	historyQ := &mockDBQ{}
	system := &system{
		historyQ: historyQ,
		ctx:      context.Background(),
	}

	historyQ.On("GetTx").Return(nil).Once()
	historyQ.On("Begin").Return(errors.New("my error")).Once()
	historyQ.On("GetTx").Return(&sqlx.Tx{}).Once()

	assert.PanicsWithValue(t, "unexpected transaction", func() {
		system.Run()
	})
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	historyQ := &mockDBQ{}
	system := &system{
		historyQ: historyQ,
		ctx:      ctx,
	}

	historyQ.On("GetTx").Return(nil).Once()
	historyQ.On("Begin").Return(errors.New("my error")).Once()

	cancel()
	assert.NoError(t, system.runStateMachine(startState{}))
}

// TestStateMachineRunReturnsErrorWhenNextStateIsShutdownWithError checks if the
// state that goes to shutdownState and returns an error will make `run` function
// return that error. This is essential because some commands rely on this to return
// non-zero exit code.
func TestStateMachineRunReturnsErrorWhenNextStateIsShutdownWithError(t *testing.T) {
	historyQ := &mockDBQ{}
	system := &system{
		ctx:      context.Background(),
		historyQ: historyQ,
	}

	historyQ.On("GetTx").Return(nil).Once()

	err := system.runStateMachine(verifyRangeState{})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid range: [0, 0]")
}

func TestMaybeVerifyStateGetExpStateInvalidDBErrCancelOrContextCanceled(t *testing.T) {
	historyQ := &mockDBQ{}
	system := &system{
		historyQ:          historyQ,
		ctx:               context.Background(),
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	var out bytes.Buffer
	logger := logpkg.New()
	logger.SetOutput(&out)
	done := logger.StartTest(logpkg.InfoLevel)

	oldLogger := log
	log = logger
	defer func() { log = oldLogger }()

	historyQ.On("GetExpStateInvalid", system.ctx).Return(false, db.ErrCancelled).Once()
	system.maybeVerifyState(0)

	historyQ.On("GetExpStateInvalid", system.ctx).Return(false, context.Canceled).Once()
	system.maybeVerifyState(0)

	logged := done()
	assert.Len(t, logged, 0)
	historyQ.AssertExpectations(t)
}
func TestMaybeVerifyInternalDBErrCancelOrContextCanceled(t *testing.T) {
	historyQ := &mockDBQ{}
	system := &system{
		historyQ:          historyQ,
		ctx:               context.Background(),
		checkpointManager: historyarchive.NewCheckpointManager(64),
	}

	var out bytes.Buffer
	logger := logpkg.New()
	logger.SetOutput(&out)
	done := logger.StartTest(logpkg.InfoLevel)

	oldLogger := log
	log = logger
	defer func() { log = oldLogger }()

	historyQ.On("GetExpStateInvalid", system.ctx, mock.Anything).Return(false, nil).Twice()
	historyQ.On("Rollback").Return(nil).Twice()
	historyQ.On("CloneIngestionQ").Return(historyQ).Twice()

	historyQ.On("BeginTx", mock.Anything).Return(db.ErrCancelled).Once()
	system.maybeVerifyState(63)
	system.wg.Wait()

	historyQ.On("BeginTx", mock.Anything).Return(context.Canceled).Once()
	system.maybeVerifyState(63)
	system.wg.Wait()

	logged := done()

	// it logs "State verification finished" twice, but no errors
	assert.Len(t, logged, 0)

	historyQ.AssertExpectations(t)
}

type mockDBQ struct {
	mock.Mock

	history.MockQAccounts
	history.MockQClaimableBalances
	history.MockQHistoryClaimableBalances
	history.MockQLiquidityPools
	history.MockQHistoryLiquidityPools
	history.MockQAssetStats
	history.MockQData
	history.MockQEffects
	history.MockQLedgers
	history.MockQOffers
	history.MockQOperations
	history.MockQSigners
	history.MockQTransactions
	history.MockQTrustLines
}

func (m *mockDBQ) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDBQ) BeginTx(txOpts *sql.TxOptions) error {
	args := m.Called(txOpts)
	return args.Error(0)
}

func (m *mockDBQ) CloneIngestionQ() history.IngestionQ {
	args := m.Called()
	return args.Get(0).(history.IngestionQ)
}

func (m *mockDBQ) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDBQ) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockDBQ) GetTx() *sqlx.Tx {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sqlx.Tx)
}

func (m *mockDBQ) GetLastLedgerIngest(ctx context.Context) (uint32, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) GetOfferCompactionSequence(ctx context.Context) (uint32, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) GetLiquidityPoolCompactionSequence(ctx context.Context) (uint32, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) GetLastLedgerIngestNonBlocking(ctx context.Context) (uint32, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) GetIngestVersion(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *mockDBQ) UpdateLastLedgerIngest(ctx context.Context, sequence uint32) error {
	args := m.Called(ctx, sequence)
	return args.Error(0)
}

func (m *mockDBQ) UpdateExpStateInvalid(ctx context.Context, invalid bool) error {
	args := m.Called(ctx, invalid)
	return args.Error(0)
}

func (m *mockDBQ) UpdateIngestVersion(ctx context.Context, version int) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *mockDBQ) GetExpStateInvalid(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockDBQ) GetAllOffers(ctx context.Context) ([]history.Offer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]history.Offer), args.Error(1)
}

func (m *mockDBQ) GetLatestHistoryLedger(ctx context.Context) (uint32, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockDBQ) TruncateIngestStateTables(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockDBQ) DeleteRangeAll(ctx context.Context, start, end int64) error {
	args := m.Called(ctx, start, end)
	return args.Error(0)
}

// Methods from interfaces duplicating methods:

func (m *mockDBQ) NewTransactionParticipantsBatchInsertBuilder(maxBatchSize int) history.TransactionParticipantsBatchInsertBuilder {
	args := m.Called(maxBatchSize)
	return args.Get(0).(history.TransactionParticipantsBatchInsertBuilder)
}

func (m *mockDBQ) NewOperationParticipantBatchInsertBuilder(maxBatchSize int) history.OperationParticipantBatchInsertBuilder {
	args := m.Called(maxBatchSize)
	return args.Get(0).(history.TransactionParticipantsBatchInsertBuilder)
}

func (m *mockDBQ) NewTradeBatchInsertBuilder(maxBatchSize int) history.TradeBatchInsertBuilder {
	args := m.Called(maxBatchSize)
	return args.Get(0).(history.TradeBatchInsertBuilder)
}

func (m *mockDBQ) RebuildTradeAggregationTimes(ctx context.Context, from, to strtime.Millis) error {
	args := m.Called(ctx, from, to)
	return args.Error(0)
}

func (m *mockDBQ) RebuildTradeAggregationBuckets(ctx context.Context, fromLedger, toLedger uint32) error {
	args := m.Called(ctx, fromLedger, toLedger)
	return args.Error(0)
}

func (m *mockDBQ) CreateAssets(ctx context.Context, assets []xdr.Asset, batchSize int) (map[string]history.Asset, error) {
	args := m.Called(ctx, assets)
	return args.Get(0).(map[string]history.Asset), args.Error(1)
}

type mockLedgerBackend struct {
	mock.Mock
}

func (m *mockLedgerBackend) GetLatestLedgerSequence(ctx context.Context) (sequence uint32, err error) {
	args := m.Called(ctx)
	return args.Get(0).(uint32), args.Error(1)
}

func (m *mockLedgerBackend) GetLedger(ctx context.Context, sequence uint32) (xdr.LedgerCloseMeta, error) {
	args := m.Called(ctx, sequence)
	return args.Get(0).(xdr.LedgerCloseMeta), args.Error(1)
}

func (m *mockLedgerBackend) PrepareRange(ctx context.Context, ledgerRange ledgerbackend.Range) error {
	args := m.Called(ctx, ledgerRange)
	return args.Error(0)
}

func (m *mockLedgerBackend) IsPrepared(ctx context.Context, ledgerRange ledgerbackend.Range) (bool, error) {
	args := m.Called(ctx, ledgerRange)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockLedgerBackend) Close() error {
	args := m.Called()
	return args.Error(0)
}

type mockProcessorsRunner struct {
	mock.Mock
}

func (m *mockProcessorsRunner) SetHistoryAdapter(historyAdapter historyArchiveAdapterInterface) {
	m.Called(historyAdapter)
}

func (m *mockProcessorsRunner) EnableMemoryStatsLogging() {
	m.Called()
}

func (m *mockProcessorsRunner) DisableMemoryStatsLogging() {
	m.Called()
}

func (m *mockProcessorsRunner) RunGenesisStateIngestion() (ingest.StatsChangeProcessorResults, error) {
	args := m.Called()
	return args.Get(0).(ingest.StatsChangeProcessorResults), args.Error(1)
}

func (m *mockProcessorsRunner) RunHistoryArchiveIngestion(
	checkpointLedger uint32,
	ledgerProtocolVersion uint32,
	bucketListHash xdr.Hash,
) (ingest.StatsChangeProcessorResults, error) {
	args := m.Called(checkpointLedger, ledgerProtocolVersion, bucketListHash)
	return args.Get(0).(ingest.StatsChangeProcessorResults), args.Error(1)
}

func (m *mockProcessorsRunner) RunAllProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
	ledgerStats,
	error,
) {
	args := m.Called(ledger)
	return args.Get(0).(ledgerStats),
		args.Error(1)
}

func (m *mockProcessorsRunner) RunTransactionProcessorsOnLedger(ledger xdr.LedgerCloseMeta) (
	processors.StatsLedgerTransactionProcessorResults,
	processorsRunDurations,
	processors.TradeStats,
	error,
) {
	args := m.Called(ledger)
	return args.Get(0).(processors.StatsLedgerTransactionProcessorResults),
		args.Get(1).(processorsRunDurations),
		args.Get(2).(processors.TradeStats),
		args.Error(3)
}

var _ ProcessorRunnerInterface = (*mockProcessorsRunner)(nil)

type mockDiamcircleCoreClient struct {
	mock.Mock
}

func (m *mockDiamcircleCoreClient) SetCursor(ctx context.Context, id string, cursor int32) error {
	args := m.Called(ctx, id, cursor)
	return args.Error(0)
}

var _ diamcircleCoreClient = (*mockDiamcircleCoreClient)(nil)

type mockSystem struct {
	mock.Mock
}

func (m *mockSystem) Run() {
	m.Called()
}

func (m *mockSystem) Metrics() Metrics {
	args := m.Called()
	return args.Get(0).(Metrics)
}

func (m *mockSystem) RegisterMetrics(registry *prometheus.Registry) {
	m.Called(registry)
}

func (m *mockSystem) StressTest(numTransactions, changesPerTransaction int) error {
	args := m.Called(numTransactions, changesPerTransaction)
	return args.Error(0)
}

func (m *mockSystem) VerifyRange(fromLedger, toLedger uint32, verifyState bool) error {
	args := m.Called(fromLedger, toLedger, verifyState)
	return args.Error(0)
}

func (m *mockSystem) ReingestRange(ledgerRanges []history.LedgerRange, force bool) error {
	args := m.Called(ledgerRanges, force)
	return args.Error(0)
}

func (m *mockSystem) BuildGenesisState() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSystem) Shutdown() {
	m.Called()
}

var _ System = (*mockSystem)(nil)
