//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/xdr"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LedgersProcessorTestSuiteLedger struct {
	suite.Suite
	processor     *LedgersProcessor
	mockQ         *history.MockQLedgers
	header        xdr.LedgerHeaderHistoryEntry
	successCount  int
	failedCount   int
	opCount       int
	ingestVersion int
	txs           []ingest.LedgerTransaction
	txSetOpCount  int
}

func TestLedgersProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(LedgersProcessorTestSuiteLedger))
}

func createTransaction(successful bool, numOps int) ingest.LedgerTransaction {
	code := xdr.TransactionResultCodeTxSuccess
	if !successful {
		code = xdr.TransactionResultCodeTxFailed
	}

	operations := []xdr.Operation{}
	op := xdr.BumpSequenceOp{BumpTo: 30000}
	for i := 0; i < numOps; i++ {
		operations = append(operations, xdr.Operation{
			Body: xdr.OperationBody{
				Type:           xdr.OperationTypeBumpSequence,
				BumpSequenceOp: &op,
			},
		})
	}
	sourceAID := xdr.MustAddress("GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY")
	return ingest.LedgerTransaction{
		Result: xdr.TransactionResultPair{
			TransactionHash: xdr.Hash{},
			Result: xdr.TransactionResult{
				Result: xdr.TransactionResultResult{
					Code:            code,
					InnerResultPair: &xdr.InnerTransactionResultPair{},
					Results:         &[]xdr.OperationResult{},
				},
			},
		},
		Envelope: xdr.TransactionEnvelope{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: sourceAID.ToMuxedAccount(),
					Operations:    operations,
				},
			},
		},
		UnsafeMeta: xdr.TransactionMeta{
			V: 2,
			V2: &xdr.TransactionMetaV2{
				Operations: make([]xdr.OperationMeta, numOps, numOps),
			},
		},
	}
}

func (s *LedgersProcessorTestSuiteLedger) SetupTest() {
	s.mockQ = &history.MockQLedgers{}
	s.ingestVersion = 100
	s.header = xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: xdr.Uint32(20),
		},
	}
	s.processor = NewLedgerProcessor(
		s.mockQ,
		s.header,
		s.ingestVersion,
	)

	s.txs = []ingest.LedgerTransaction{
		createTransaction(true, 1),
		createTransaction(false, 3),
		createTransaction(true, 4),
	}

	s.successCount = 2
	s.failedCount = 1
	s.opCount = 5
	s.txSetOpCount = 8
}

func (s *LedgersProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
}

func (s *LedgersProcessorTestSuiteLedger) TestInsertLedgerSucceeds() {
	ctx := context.Background()
	s.mockQ.On(
		"InsertLedger",
		ctx,
		s.header,
		s.successCount,
		s.failedCount,
		s.opCount,
		s.txSetOpCount,
		s.ingestVersion,
	).Return(int64(1), nil)

	for _, tx := range s.txs {
		err := s.processor.ProcessTransaction(ctx, tx)
		s.Assert().NoError(err)
	}

	err := s.processor.Commit(ctx)
	s.Assert().NoError(err)
}

func (s *LedgersProcessorTestSuiteLedger) TestInsertLedgerReturnsError() {
	ctx := context.Background()
	s.mockQ.On(
		"InsertLedger",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(int64(0), errors.New("transient error"))

	err := s.processor.Commit(ctx)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Could not insert ledger: transient error")
}

func (s *LedgersProcessorTestSuiteLedger) TestInsertLedgerNoRowsAffected() {
	ctx := context.Background()
	s.mockQ.On(
		"InsertLedger",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(int64(0), nil)

	err := s.processor.Commit(ctx)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "0 rows affected when ingesting new ledger: 20")
}
