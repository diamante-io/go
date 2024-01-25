package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQClaimableBalances is a mock implementation of the QAccounts interface
type MockQClaimableBalances struct {
	mock.Mock
}

func (m *MockQClaimableBalances) CountClaimableBalances(ctx context.Context) (int, error) {
	a := m.Called(ctx)
	return a.Get(0).(int), a.Error(1)
}

func (m *MockQClaimableBalances) GetClaimableBalancesByID(ctx context.Context, ids []string) ([]ClaimableBalance, error) {
	a := m.Called(ctx, ids)
	return a.Get(0).([]ClaimableBalance), a.Error(1)
}

func (m *MockQClaimableBalances) UpsertClaimableBalances(ctx context.Context, cbs []ClaimableBalance) error {
	a := m.Called(ctx, cbs)
	return a.Error(0)
}

func (m *MockQClaimableBalances) RemoveClaimableBalances(ctx context.Context, ids []string) (int64, error) {
	a := m.Called(ctx, ids)
	return a.Get(0).(int64), a.Error(1)
}
