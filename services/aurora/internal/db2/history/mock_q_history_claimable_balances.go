package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQHistoryClaimableBalances is a mock implementation of the QClaimableBalances interface
type MockQHistoryClaimableBalances struct {
	mock.Mock
}

func (m *MockQHistoryClaimableBalances) CreateHistoryClaimableBalances(ctx context.Context, ids []string, maxBatchSize int) (map[string]int64, error) {
	a := m.Called(ctx, ids, maxBatchSize)
	return a.Get(0).(map[string]int64), a.Error(1)
}

func (m *MockQHistoryClaimableBalances) NewTransactionClaimableBalanceBatchInsertBuilder(maxBatchSize int) TransactionClaimableBalanceBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(TransactionClaimableBalanceBatchInsertBuilder)
}

// MockTransactionClaimableBalanceBatchInsertBuilder is a mock implementation of the
// TransactionClaimableBalanceBatchInsertBuilder interface
type MockTransactionClaimableBalanceBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockTransactionClaimableBalanceBatchInsertBuilder) Add(ctx context.Context, transactionID, accountID int64) error {
	a := m.Called(ctx, transactionID, accountID)
	return a.Error(0)
}

func (m *MockTransactionClaimableBalanceBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}

// NewOperationClaimableBalanceBatchInsertBuilder mock
func (m *MockQHistoryClaimableBalances) NewOperationClaimableBalanceBatchInsertBuilder(maxBatchSize int) OperationClaimableBalanceBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(OperationClaimableBalanceBatchInsertBuilder)
}

// MockOperationClaimableBalanceBatchInsertBuilder is a mock implementation of the
// OperationClaimableBalanceBatchInsertBuilder interface
type MockOperationClaimableBalanceBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockOperationClaimableBalanceBatchInsertBuilder) Add(ctx context.Context, transactionID, accountID int64) error {
	a := m.Called(ctx, transactionID, accountID)
	return a.Error(0)
}

func (m *MockOperationClaimableBalanceBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
