package history

import "github.com/stretchr/testify/mock"

// MockQTransactions is a mock implementation of the QTransactions interface
type MockQTransactions struct {
	mock.Mock
}

func (m *MockQTransactions) NewTransactionBatchInsertBuilder(maxBatchSize int) TransactionBatchInsertBuilder {
	a := m.Called(maxBatchSize)
	return a.Get(0).(TransactionBatchInsertBuilder)
}
