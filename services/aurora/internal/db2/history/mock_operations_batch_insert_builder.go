package history

import (
	"context"

	"go/xdr"

	"github.com/guregu/null"
	"github.com/stretchr/testify/mock"
)

// MockOperationsBatchInsertBuilder OperationsBatchInsertBuilder mock
type MockOperationsBatchInsertBuilder struct {
	mock.Mock
}

// Add mock
func (m *MockOperationsBatchInsertBuilder) Add(ctx context.Context,
	id int64,
	transactionID int64,
	applicationOrder uint32,
	operationType xdr.OperationType,
	details []byte,
	sourceAccount string,
	sourceAccountMuxed null.String,
) error {
	a := m.Called(ctx,
		id,
		transactionID,
		applicationOrder,
		operationType,
		details,
		sourceAccount,
		sourceAccountMuxed,
	)
	return a.Error(0)
}

// Exec mock
func (m *MockOperationsBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
