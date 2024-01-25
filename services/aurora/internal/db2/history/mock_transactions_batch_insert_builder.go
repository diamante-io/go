package history

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/ingest"
)

type MockTransactionsBatchInsertBuilder struct {
	mock.Mock
}

func (m *MockTransactionsBatchInsertBuilder) Add(ctx context.Context, transaction ingest.LedgerTransaction, sequence uint32) error {
	a := m.Called(ctx, transaction, sequence)
	return a.Error(0)
}

func (m *MockTransactionsBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
