package history

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockOperationParticipantBatchInsertBuilder OperationParticipantBatchInsertBuilder mock
type MockOperationParticipantBatchInsertBuilder struct {
	mock.Mock
}

// Add mock
func (m *MockOperationParticipantBatchInsertBuilder) Add(ctx context.Context, operationID int64, accountID int64) error {
	a := m.Called(ctx, operationID, accountID)
	return a.Error(0)
}

// Exec mock
func (m *MockOperationParticipantBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
