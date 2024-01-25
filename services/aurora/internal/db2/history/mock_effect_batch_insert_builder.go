package history

import (
	"context"

	"github.com/guregu/null"
	"github.com/stretchr/testify/mock"
)

// MockEffectBatchInsertBuilder mock EffectBatchInsertBuilder
type MockEffectBatchInsertBuilder struct {
	mock.Mock
}

// Add mock
func (m *MockEffectBatchInsertBuilder) Add(ctx context.Context,
	accountID int64,
	muxedAccount null.String,
	operationID int64,
	order uint32,
	effectType EffectType,
	details []byte,
) error {
	a := m.Called(ctx,
		accountID,
		muxedAccount,
		operationID,
		order,
		effectType,
		details,
	)
	return a.Error(0)
}

// Exec mock
func (m *MockEffectBatchInsertBuilder) Exec(ctx context.Context) error {
	a := m.Called(ctx)
	return a.Error(0)
}
