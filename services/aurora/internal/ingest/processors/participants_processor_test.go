//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

type ParticipantsProcessorTestSuiteLedger struct {
	suite.Suite
	ctx                              context.Context
	processor                        *ParticipantsProcessor
	mockQ                            *history.MockQParticipants
	mockBatchInsertBuilder           *history.MockTransactionParticipantsBatchInsertBuilder
	mockOperationsBatchInsertBuilder *history.MockOperationParticipantBatchInsertBuilder

	firstTx     ingest.LedgerTransaction
	secondTx    ingest.LedgerTransaction
	thirdTx     ingest.LedgerTransaction
	firstTxID   int64
	secondTxID  int64
	thirdTxID   int64
	addresses   []string
	addressToID map[string]int64
	txs         []ingest.LedgerTransaction
}

func TestParticipantsProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(ParticipantsProcessorTestSuiteLedger))
}

func (s *ParticipantsProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQParticipants{}
	s.mockBatchInsertBuilder = &history.MockTransactionParticipantsBatchInsertBuilder{}
	s.mockOperationsBatchInsertBuilder = &history.MockOperationParticipantBatchInsertBuilder{}
	sequence := uint32(20)

	s.addresses = []string{
		"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
		"GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	}

	s.firstTx = createTransaction(true, 1)
	s.firstTx.Index = 1
	aid := xdr.MustAddress(s.addresses[0])
	s.firstTx.Envelope.V1.Tx.SourceAccount = aid.ToMuxedAccount()
	s.firstTxID = toid.New(int32(sequence), 1, 0).ToInt64()

	s.secondTx = createTransaction(true, 1)
	s.secondTx.Index = 2
	s.secondTx.Envelope.Operations()[0].Body = xdr.OperationBody{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccountOp: &xdr.CreateAccountOp{
			Destination: xdr.MustAddress(s.addresses[1]),
		},
	}
	aid = xdr.MustAddress(s.addresses[2])
	s.secondTx.Envelope.V1.Tx.SourceAccount = aid.ToMuxedAccount()
	s.secondTxID = toid.New(int32(sequence), 2, 0).ToInt64()

	s.thirdTx = createTransaction(true, 1)
	s.thirdTx.Index = 3
	aid = xdr.MustAddress(s.addresses[0])
	s.thirdTx.Envelope.V1.Tx.SourceAccount = aid.ToMuxedAccount()
	s.thirdTxID = toid.New(int32(sequence), 3, 0).ToInt64()

	s.addressToID = map[string]int64{
		s.addresses[0]: 2,
		s.addresses[1]: 20,
		s.addresses[2]: 200,
	}

	s.processor = NewParticipantsProcessor(
		s.mockQ,
		sequence,
	)

	s.txs = []ingest.LedgerTransaction{
		s.firstTx,
		s.secondTx,
		s.thirdTx,
	}
}

func (s *ParticipantsProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
	s.mockBatchInsertBuilder.AssertExpectations(s.T())
	s.mockOperationsBatchInsertBuilder.AssertExpectations(s.T())
}

func (s *ParticipantsProcessorTestSuiteLedger) mockSuccessfulTransactionBatchAdds() {
	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.firstTxID, s.addressToID[s.addresses[0]],
	).Return(nil).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID, s.addressToID[s.addresses[1]],
	).Return(nil).Once()
	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID, s.addressToID[s.addresses[2]],
	).Return(nil).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.thirdTxID, s.addressToID[s.addresses[0]],
	).Return(nil).Once()
}

func (s *ParticipantsProcessorTestSuiteLedger) mockSuccessfulOperationBatchAdds() {
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.firstTxID+1, s.addressToID[s.addresses[0]],
	).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID+1, s.addressToID[s.addresses[1]],
	).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID+1, s.addressToID[s.addresses[2]],
	).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.thirdTxID+1, s.addressToID[s.addresses[0]],
	).Return(nil).Once()
}
func (s *ParticipantsProcessorTestSuiteLedger) TestEmptyParticipants() {
	err := s.processor.Commit(s.ctx)
	s.Assert().NoError(err)
}

func (s *ParticipantsProcessorTestSuiteLedger) TestFeeBumptransaction() {
	feeBumpTx := createTransaction(true, 0)
	feeBumpTx.Index = 1
	aid := xdr.MustAddress(s.addresses[0])
	feeBumpTx.Envelope.V1.Tx.SourceAccount = aid.ToMuxedAccount()
	aid = xdr.MustAddress(s.addresses[1])
	feeBumpTx.Envelope.FeeBump = &xdr.FeeBumpTransactionEnvelope{
		Tx: xdr.FeeBumpTransaction{
			FeeSource: aid.ToMuxedAccount(),
			Fee:       100,
			InnerTx: xdr.FeeBumpTransactionInnerTx{
				Type: xdr.EnvelopeTypeEnvelopeTypeTx,
				V1:   feeBumpTx.Envelope.V1,
			},
		},
	}
	feeBumpTx.Envelope.V1 = nil
	feeBumpTx.Envelope.Type = xdr.EnvelopeTypeEnvelopeTypeTxFeeBump
	feeBumpTx.Result.Result.Result.Code = xdr.TransactionResultCodeTxFeeBumpInnerSuccess
	feeBumpTx.Result.Result.Result.InnerResultPair = &xdr.InnerTransactionResultPair{
		Result: xdr.InnerTransactionResult{
			Result: xdr.InnerTransactionResultResult{
				Code:    xdr.TransactionResultCodeTxSuccess,
				Results: &[]xdr.OperationResult{},
			},
		},
	}
	feeBumpTx.Result.Result.Result.Results = nil
	feeBumpTxID := toid.New(20, 1, 0).ToInt64()

	addresses := s.addresses[:2]
	addressToID := map[string]int64{
		addresses[0]: s.addressToID[addresses[0]],
		addresses[1]: s.addressToID[addresses[1]],
	}
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				addresses,
				arg,
			)
		}).Return(addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()
	s.mockQ.On("NewOperationParticipantBatchInsertBuilder", maxBatchSize).
		Return(s.mockOperationsBatchInsertBuilder).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, feeBumpTxID, addressToID[addresses[0]],
	).Return(nil).Once()
	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, feeBumpTxID, addressToID[addresses[1]],
	).Return(nil).Once()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	s.Assert().NoError(s.processor.ProcessTransaction(s.ctx, feeBumpTx))
	s.Assert().NoError(s.processor.Commit(s.ctx))
}

func (s *ParticipantsProcessorTestSuiteLedger) TestIngestParticipantsSucceeds() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				s.addresses,
				arg,
			)
		}).Return(s.addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()
	s.mockQ.On("NewOperationParticipantBatchInsertBuilder", maxBatchSize).
		Return(s.mockOperationsBatchInsertBuilder).Once()

	s.mockSuccessfulTransactionBatchAdds()
	s.mockSuccessfulOperationBatchAdds()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().NoError(err)
}

func (s *ParticipantsProcessorTestSuiteLedger) TestCreateAccountsFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Return(s.addressToID, errors.New("transient error")).Once()
	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "Could not create account ids: transient error")
}

func (s *ParticipantsProcessorTestSuiteLedger) TestBatchAddFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				s.addresses,
				arg,
			)
		}).Return(s.addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.firstTxID, s.addressToID[s.addresses[0]],
	).Return(errors.New("transient error")).Once()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID, s.addressToID[s.addresses[1]],
	).Return(nil).Maybe()
	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID, s.addressToID[s.addresses[2]],
	).Return(nil).Maybe()

	s.mockBatchInsertBuilder.On(
		"Add", s.ctx, s.thirdTxID, s.addressToID[s.addresses[0]],
	).Return(nil).Maybe()
	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "Could not insert transaction participant in db: transient error")
}

func (s *ParticipantsProcessorTestSuiteLedger) TestOperationParticipantsBatchAddFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				s.addresses,
				arg,
			)
		}).Return(s.addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()
	s.mockQ.On("NewOperationParticipantBatchInsertBuilder", maxBatchSize).
		Return(s.mockOperationsBatchInsertBuilder).Once()

	s.mockSuccessfulTransactionBatchAdds()

	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.firstTxID+1, s.addressToID[s.addresses[0]],
	).Return(errors.New("transient error")).Once()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID+1, s.addressToID[s.addresses[1]],
	).Return(nil).Maybe()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.secondTxID+1, s.addressToID[s.addresses[2]],
	).Return(nil).Maybe()
	s.mockOperationsBatchInsertBuilder.On(
		"Add", s.ctx, s.thirdTxID+1, s.addressToID[s.addresses[0]],
	).Return(nil).Maybe()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "could not insert operation participant in db: transient error")
}

func (s *ParticipantsProcessorTestSuiteLedger) TestBatchAddExecFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				s.addresses,
				arg,
			)
		}).Return(s.addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	s.mockSuccessfulTransactionBatchAdds()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(errors.New("transient error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "Could not flush transaction participants to db: transient error")
}

func (s *ParticipantsProcessorTestSuiteLedger) TestOpeartionBatchAddExecFails() {
	s.mockQ.On("CreateAccounts", s.ctx, mock.AnythingOfType("[]string"), maxBatchSize).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).([]string)
			s.Assert().ElementsMatch(
				s.addresses,
				arg,
			)
		}).Return(s.addressToID, nil).Once()
	s.mockQ.On("NewTransactionParticipantsBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()
	s.mockQ.On("NewOperationParticipantBatchInsertBuilder", maxBatchSize).
		Return(s.mockOperationsBatchInsertBuilder).Once()

	s.mockSuccessfulTransactionBatchAdds()
	s.mockSuccessfulOperationBatchAdds()

	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()
	s.mockOperationsBatchInsertBuilder.On("Exec", s.ctx).Return(errors.New("transient error")).Once()

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(s.ctx, tx)
		s.Assert().NoError(err)
	}
	err := s.processor.Commit(s.ctx)
	s.Assert().EqualError(err, "could not flush operation participants to db: transient error")
}
