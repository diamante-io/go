//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package processors

import (
	"context"
	"testing"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/stretchr/testify/suite"
)

type TransactionsProcessorTestSuiteLedger struct {
	suite.Suite
	ctx                    context.Context
	processor              *TransactionProcessor
	mockQ                  *history.MockQTransactions
	mockBatchInsertBuilder *history.MockTransactionsBatchInsertBuilder
}

func TestTransactionsProcessorTestSuiteLedger(t *testing.T) {
	suite.Run(t, new(TransactionsProcessorTestSuiteLedger))
}

func (s *TransactionsProcessorTestSuiteLedger) SetupTest() {
	s.ctx = context.Background()
	s.mockQ = &history.MockQTransactions{}
	s.mockBatchInsertBuilder = &history.MockTransactionsBatchInsertBuilder{}

	s.mockQ.
		On("NewTransactionBatchInsertBuilder", maxBatchSize).
		Return(s.mockBatchInsertBuilder).Once()

	s.processor = NewTransactionProcessor(s.mockQ, 20)
}

func (s *TransactionsProcessorTestSuiteLedger) TearDownTest() {
	s.mockQ.AssertExpectations(s.T())
	s.mockBatchInsertBuilder.AssertExpectations(s.T())
}

func (s *TransactionsProcessorTestSuiteLedger) TestAddTransactionsSucceeds() {
	sequence := uint32(20)

	firstTx := createTransaction(true, 1)
	secondTx := createTransaction(false, 3)
	thirdTx := createTransaction(true, 4)

	s.mockBatchInsertBuilder.On("Add", s.ctx, firstTx, sequence).Return(nil).Once()
	s.mockBatchInsertBuilder.On("Add", s.ctx, secondTx, sequence).Return(nil).Once()
	s.mockBatchInsertBuilder.On("Add", s.ctx, thirdTx, sequence).Return(nil).Once()
	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(nil).Once()
	s.Assert().NoError(s.processor.Commit(s.ctx))

	err := s.processor.ProcessTransaction(s.ctx, firstTx)
	s.Assert().NoError(err)

	err = s.processor.ProcessTransaction(s.ctx, secondTx)
	s.Assert().NoError(err)

	err = s.processor.ProcessTransaction(s.ctx, thirdTx)
	s.Assert().NoError(err)
}

func (s *TransactionsProcessorTestSuiteLedger) TestAddTransactionsFails() {
	sequence := uint32(20)
	firstTx := createTransaction(true, 1)
	s.mockBatchInsertBuilder.On("Add", s.ctx, firstTx, sequence).
		Return(errors.New("transient error")).Once()

	err := s.processor.ProcessTransaction(s.ctx, firstTx)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error batch inserting transaction rows: transient error")
}

func (s *TransactionsProcessorTestSuiteLedger) TestExecFails() {
	sequence := uint32(20)
	firstTx := createTransaction(true, 1)

	s.mockBatchInsertBuilder.On("Add", s.ctx, firstTx, sequence).Return(nil).Once()
	s.mockBatchInsertBuilder.On("Exec", s.ctx).Return(errors.New("transient error")).Once()

	err := s.processor.ProcessTransaction(s.ctx, firstTx)
	s.Assert().NoError(err)

	err = s.processor.Commit(s.ctx)
	s.Assert().Error(err)
	s.Assert().EqualError(err, "Error flushing transaction batch: transient error")
}
