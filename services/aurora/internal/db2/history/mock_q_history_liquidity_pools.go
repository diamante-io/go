package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQHistoryLiquidityPools is a mock implementation of the QLiquidityPools interface
type MockQHistoryLiquidityPools struct {
	mock.Mock
}

func (m *MockQHistoryLiquidityPools) CreateHistoryLiquidityPools(ctx context.Context, poolIDs []string, maxBatchSize int) (map[string]int64, error) {
	a := m.Called(ctx, poolIDs, maxBatchSize)
	return a.Get(0).(map[string]int64), a.Error(1)
}

func (m *MockQHistoryLiquidityPools) NewTransactionLiquidityPoolBatchInsertBuilder(maxBatchSize int) TransactionLiquidityPoolBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(TransactionLiquidityPoolBatchInsertBuilder)
}

// MockTransactionLiquidityPoolBatchInsertBuilder is a mock implementation of the
// TransactionLiquidityPoolBatchInsertBuilder interface
type MockTransactionLiquidityPoolBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockTransactionLiquidityPoolBatchInsertBuilder) Add(ctx context.Context, transactionID, accountID int64) error {
	a := m.Called(ctx, transactionID, accountID)
	return a.Error(0)
}

func (m *MockTransactionLiquidityPoolBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}

// NewOperationLiquidityPoolBatchInsertBuilder mock
func (m *MockQHistoryLiquidityPools) NewOperationLiquidityPoolBatchInsertBuilder(maxBatchSize int) OperationLiquidityPoolBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(OperationLiquidityPoolBatchInsertBuilder)
}

// MockOperationLiquidityPoolBatchInsertBuilder is a mock implementation of the
// OperationLiquidityPoolBatchInsertBuilder interface
type MockOperationLiquidityPoolBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockOperationLiquidityPoolBatchInsertBuilder) Add(ctx context.Context, transactionID, accountID int64) error {
	a := m.Called(ctx, transactionID, accountID)
	return a.Error(0)
}

func (m *MockOperationLiquidityPoolBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
