//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package txsub

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/services/aurora/internal/txsub/sequence"
	"github.com/diamcircle/go/xdr"
)

type SystemTestSuite struct {
	suite.Suite
	ctx           context.Context
	submitter     *MockSubmitter
	db            *mockDBQ
	system        *System
	noResults     Result
	successTx     Result
	successXDR    xdr.TransactionEnvelope
	badSeq        SubmissionResult
	unmuxedSource xdr.AccountId
}

func (suite *SystemTestSuite) SetupTest() {
	suite.ctx = test.Context()
	suite.submitter = &MockSubmitter{}
	suite.db = &mockDBQ{}

	suite.system = &System{
		Pending:         NewDefaultSubmissionList(),
		Submitter:       suite.submitter,
		SubmissionQueue: sequence.NewManager(),
		DB: func(ctx context.Context) AuroraDB {
			return suite.db
		},
	}

	suite.unmuxedSource = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	source := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xcafebabe,
			Ed25519: *suite.unmuxedSource.Ed25519,
		},
	}

	tx := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTx,
		V1: &xdr.TransactionV1Envelope{
			Tx: xdr.Transaction{
				SourceAccount: source,
				Fee:           100,
				SeqNum:        1,
				Operations: []xdr.Operation{
					{
						Body: xdr.OperationBody{
							CreateAccountOp: &xdr.CreateAccountOp{
								Destination:     xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"),
								StartingBalance: 1000000000,
							},
						},
					},
				},
			},
			Signatures: []xdr.DecoratedSignature{
				{
					Hint:      xdr.SignatureHint{86, 252, 5, 247},
					Signature: xdr.Signature{131, 206, 171, 228, 64, 20, 40, 52, 2, 98, 124, 244, 87, 14, 130, 225, 190, 220, 156, 79, 121, 69, 60, 36, 57, 214, 9, 29, 176, 81, 218, 4, 213, 176, 211, 148, 191, 86, 21, 180, 94, 9, 43, 208, 32, 79, 19, 131, 90, 21, 93, 138, 153, 203, 55, 103, 2, 230, 137, 190, 19, 70, 179, 11},
				},
			},
		},
	}

	result := xdr.TransactionResult{
		FeeCharged: 123,
		Result: xdr.TransactionResultResult{
			Code:    xdr.TransactionResultCodeTxSuccess,
			Results: &[]xdr.OperationResult{},
		},
	}
	resultXDR, err := xdr.MarshalBase64(result)
	suite.Assert().NoError(err)

	suite.noResults = Result{Err: ErrNoResults}
	envelopeBase64, _ := xdr.MarshalBase64(tx)
	suite.successTx = Result{
		Transaction: history.Transaction{
			TransactionWithoutLedger: history.TransactionWithoutLedger{
				TransactionHash: "2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d",
				LedgerSequence:  2,
				TxEnvelope:      envelopeBase64,
				TxResult:        resultXDR,
			},
		},
	}
	assert.NoError(suite.T(), xdr.SafeUnmarshalBase64(suite.successTx.Transaction.TxEnvelope, &suite.successXDR))

	suite.badSeq = SubmissionResult{
		Err: ErrBadSequence,
	}

	suite.db.On("GetLatestHistoryLedger").Return(uint32(1000), nil).Maybe()
}

func (suite *SystemTestSuite) TearDownTest() {
	t := suite.T()
	suite.db.AssertExpectations(t)
}

// Returns the result provided by the ResultProvider.
func (suite *SystemTestSuite) TestSubmit_Basic() {
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Run(func(args mock.Arguments) {
			ptr := args.Get(1).(*history.Transaction)
			*ptr = suite.successTx.Transaction
		}).
		Return(nil).Once()

	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.Nil(suite.T(), r.Err)
	assert.Equal(suite.T(), suite.successTx, r)
	assert.False(suite.T(), suite.submitter.WasSubmittedTo)
}

func (suite *SystemTestSuite) TestTimeoutDuringSequnceLoop() {
	var cancel context.CancelFunc
	suite.ctx, cancel = context.WithTimeout(suite.ctx, time.Duration(0))
	defer cancel()

	suite.submitter.R = suite.badSeq
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil)

	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.NotNil(suite.T(), r.Err)
	assert.Equal(suite.T(), ErrTimeout, r.Err)
}

func (suite *SystemTestSuite) TestClientDisconnectedDuringSequnceLoop() {
	var cancel context.CancelFunc
	suite.ctx, cancel = context.WithCancel(suite.ctx)

	suite.submitter.R = suite.badSeq
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil).
		Run(func(args mock.Arguments) {
			// simulate client disconnecting while looping on sequnce number check
			cancel()
			suite.ctx.Deadline()
		}).
		Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil)

	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.NotNil(suite.T(), r.Err)
	assert.Equal(suite.T(), ErrCanceled, r.Err)
}

func getMetricValue(metric prometheus.Metric) *dto.Metric {
	value := &dto.Metric{}
	err := metric.Write(value)
	if err != nil {
		panic(err)
	}
	return value
}

// Returns the error from submission if no result is found by hash and the suite.submitter returns an error.
func (suite *SystemTestSuite) TestSubmit_NotFoundError() {
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil).
		Once()

	suite.submitter.R.Err = errors.New("busted for some reason")
	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.NotNil(suite.T(), r.Err)
	assert.True(suite.T(), suite.submitter.WasSubmittedTo)
	assert.Equal(suite.T(), float64(0), getMetricValue(suite.system.Metrics.SuccessfulSubmissionsCounter).GetCounter().GetValue())
	assert.Equal(suite.T(), float64(1), getMetricValue(suite.system.Metrics.FailedSubmissionsCounter).GetCounter().GetValue())
	assert.Equal(suite.T(), uint64(1), getMetricValue(suite.system.Metrics.SubmissionDuration).GetSummary().GetSampleCount())
}

// If the error is bad_seq and the result at the transaction's sequence number is for the same hash, return result.
func (suite *SystemTestSuite) TestSubmit_BadSeq() {
	suite.submitter.R = suite.badSeq
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil).
		Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 1}, nil).
		Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Run(func(args mock.Arguments) {
			ptr := args.Get(1).(*history.Transaction)
			*ptr = suite.successTx.Transaction
		}).
		Return(nil).Once()

	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.Nil(suite.T(), r.Err)
	assert.Equal(suite.T(), suite.successTx, r)
	assert.True(suite.T(), suite.submitter.WasSubmittedTo)
}

// If error is bad_seq and no result is found, return error.
func (suite *SystemTestSuite) TestSubmit_BadSeqNotFound() {
	suite.submitter.R = suite.badSeq
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Twice()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Twice()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil).
		Times(3)
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 1}, nil).
		Once()

	// set poll interval to 1ms so we don't need to wait 3 seconds for the test to complete
	suite.system.Init()
	suite.system.accountSeqPollInterval = time.Millisecond

	r := <-suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)

	assert.NotNil(suite.T(), r.Err)
	assert.True(suite.T(), suite.submitter.WasSubmittedTo)
}

// If no result found and no error submitting, add to open transaction list.
func (suite *SystemTestSuite) TestSubmit_OpenTransactionList() {
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, suite.successTx.Transaction.TransactionHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{suite.unmuxedSource.Address()}).
		Return(map[string]uint64{suite.unmuxedSource.Address(): 0}, nil).
		Once()

	suite.system.Submit(
		suite.ctx,
		suite.successTx.Transaction.TxEnvelope,
		suite.successXDR,
		suite.successTx.Transaction.TransactionHash,
	)
	pending := suite.system.Pending.Pending(suite.ctx)
	assert.Equal(suite.T(), 1, len(pending))
	assert.Equal(suite.T(), suite.successTx.Transaction.TransactionHash, pending[0])
	assert.Equal(suite.T(), float64(1), getMetricValue(suite.system.Metrics.SuccessfulSubmissionsCounter).GetCounter().GetValue())
	assert.Equal(suite.T(), float64(0), getMetricValue(suite.system.Metrics.FailedSubmissionsCounter).GetCounter().GetValue())
	assert.Equal(suite.T(), uint64(1), getMetricValue(suite.system.Metrics.SubmissionDuration).GetSummary().GetSampleCount())
}

// Tick should be a no-op if there are no open submissions.
func (suite *SystemTestSuite) TestTick_Noop() {
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()

	suite.system.Tick(suite.ctx)
}

// TestTick_Deadlock is a regression test for Tick() deadlock: if for any reason
// call to Tick() takes more time and another Tick() is called.
// This test starts two go routines: both calling Tick() but the call to
// `sys.Sequences.Get(addys)` is delayed by 1 second. It allows to simulate two
// calls to `Tick()` executed at the same time.
func (suite *SystemTestSuite) TestTick_Deadlock() {
	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()

	// Start first Tick
	suite.system.SubmissionQueue.Push("address", 0)
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{"address"}).
		Return(map[string]uint64{}, nil).
		Run(func(args mock.Arguments) {
			// Start second tick
			suite.system.Tick(suite.ctx)
		}).
		Once()

	suite.system.Tick(suite.ctx)
}

// Test that Tick finishes any available transactions,
func (suite *SystemTestSuite) TestTick_FinishesTransactions() {
	l := make(chan Result, 1)
	suite.system.Pending.Add(suite.ctx, suite.successTx.Transaction.TransactionHash, l)

	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionsByHashesSinceLedger", suite.ctx, []string{suite.successTx.Transaction.TransactionHash}, uint32(940)).
		Return(nil, sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()

	suite.system.Tick(suite.ctx)

	assert.Equal(suite.T(), 0, len(l))
	assert.Equal(suite.T(), 1, len(suite.system.Pending.Pending(suite.ctx)))

	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionsByHashesSinceLedger", suite.ctx, []string{suite.successTx.Transaction.TransactionHash}, uint32(940)).
		Return([]history.Transaction{suite.successTx.Transaction}, nil).Once()

	suite.system.Tick(suite.ctx)

	assert.Equal(suite.T(), 1, len(l))
	assert.Equal(suite.T(), 0, len(suite.system.Pending.Pending(suite.ctx)))
}

func (suite *SystemTestSuite) TestTickFinishFeeBumpTransaction() {
	innerTxEnvelope := "AAAAAAMDAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYwAAAAAAAABhAAAAAQAAAAAAAAACAAAAAAAAAAQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAGIAAAAAAAAAAQICAgIAAAADFBQUAA=="
	innerHash := "e98869bba8bce08c10b78406202127f3888c25454cd37b02600862452751f526"
	var parsedInnerTx xdr.TransactionEnvelope
	assert.NoError(suite.T(), xdr.SafeUnmarshalBase64(innerTxEnvelope, &parsedInnerTx))

	feeBumpTx := Result{
		Transaction: history.Transaction{
			TransactionWithoutLedger: history.TransactionWithoutLedger{
				Account:              "GABQGAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB2MX",
				TransactionHash:      "3dfef7d7226995b504f2827cc63d45ad41e9687bb0a8abcf08ba755fedca0352",
				LedgerSequence:       123,
				TxEnvelope:           "AAAABQAAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMIAAAAAgAAAAADAwMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGMAAAAAAAAAYQAAAAEAAAAAAAAAAgAAAAAAAAAEAAAAAAAAAAEAAAAAAAAACwAAAAAAAABiAAAAAAAAAAECAgICAAAAAxQUFAAAAAAAAAAAAQMDAwMAAAADHh4eAA==",
				TxResult:             "AAAAAAAAAHsAAAAB6Yhpu6i84IwQt4QGICEn84iMJUVM03sCYAhiRSdR9SYAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAsAAAAAAAAAAAAAAAA=",
				TxMeta:               "AAAAAQAAAAAAAAAA",
				InnerTransactionHash: null.StringFrom("e98869bba8bce08c10b78406202127f3888c25454cd37b02600862452751f526"),
			},
		},
	}

	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionByHash", suite.ctx, mock.Anything, innerHash).
		Return(sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()
	suite.db.On("GetSequenceNumbers", suite.ctx, []string{"GABQGAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB2MX"}).
		Return(map[string]uint64{"GABQGAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB2MX": 96}, nil).
		Once()

	l := suite.system.Submit(suite.ctx, innerTxEnvelope, parsedInnerTx, innerHash)
	assert.Equal(suite.T(), 0, len(l))
	assert.Equal(suite.T(), 1, len(suite.system.Pending.Pending(suite.ctx)))

	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionsByHashesSinceLedger", suite.ctx, []string{innerHash}, uint32(940)).
		Return([]history.Transaction{feeBumpTx.Transaction}, nil).Once()

	suite.system.Tick(suite.ctx)

	assert.Equal(suite.T(), 1, len(l))
	assert.Equal(suite.T(), 0, len(suite.system.Pending.Pending(suite.ctx)))
	r := <-l
	assert.NoError(suite.T(), r.Err)
	assert.Equal(suite.T(), feeBumpTx, r)
}

// Test that Tick removes old submissions that have timed out.
func (suite *SystemTestSuite) TestTick_RemovesStaleSubmissions() {
	l := make(chan Result, 1)
	suite.system.SubmissionTimeout = 100 * time.Millisecond
	suite.system.Pending.Add(suite.ctx, suite.successTx.Transaction.TransactionHash, l)
	<-time.After(101 * time.Millisecond)

	suite.db.On("BeginTx", &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}).Return(nil).Once()
	suite.db.On("Rollback").Return(nil).Once()
	suite.db.On("TransactionsByHashesSinceLedger", suite.ctx, []string{suite.successTx.Transaction.TransactionHash}, uint32(940)).
		Return(nil, sql.ErrNoRows).Once()
	suite.db.On("NoRows", sql.ErrNoRows).Return(true).Once()

	suite.system.Tick(suite.ctx)

	assert.Equal(suite.T(), 0, len(suite.system.Pending.Pending(suite.ctx)))
	assert.Equal(suite.T(), 1, len(l))
	<-l
	select {
	case _, stillOpen := <-l:
		assert.False(suite.T(), stillOpen)
	default:
		panic("could not read from listener")
	}
}

func TestSystemTestSuite(t *testing.T) {
	suite.Run(t, new(SystemTestSuite))
}
