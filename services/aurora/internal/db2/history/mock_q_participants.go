package history

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockQParticipants is a mock implementation of the QParticipants interface
type MockQParticipants struct {
	mock.Mock
}

func (m *MockQParticipants) CreateAccounts(ctx context.Context, addresses []string, maxBatchSize int) (map[string]int64, error) {
	a := m.Called(ctx, addresses, maxBatchSize)
	return a.Get(0).(map[string]int64), a.Error(1)
}

func (m *MockQParticipants) NewTransactionParticipantsBatchInsertBuilder(maxBatchSize int) TransactionParticipantsBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(TransactionParticipantsBatchInsertBuilder)
}

// MockTransactionParticipantsBatchInsertBuilder is a mock implementation of the
// TransactionParticipantsBatchInsertBuilder interface
type MockTransactionParticipantsBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockTransactionParticipantsBatchInsertBuilder) Add(ctx context.Context, transactionID, accountID int64) error {
	a := m.Called(ctx, transactionID, accountID)
	return a.Error(0)
}

func (m *MockTransactionParticipantsBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}

// NewOperationParticipantBatchInsertBuilder mock
func (m *MockQParticipants) NewOperationParticipantBatchInsertBuilder(maxBatchSize int) OperationParticipantBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(OperationParticipantBatchInsertBuilder)
}
