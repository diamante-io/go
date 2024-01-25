package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQAccounts is a mock implementation of the QAccounts interface
type MockQAccounts struct {
	mock.Mock
}

func (m *MockQAccounts) GetAccountsByIDs(ctx context.Context, ids []string) ([]AccountEntry, error) {
	a := m.Called(ctx, ids)
	return a.Get(0).([]AccountEntry), a.Error(1)
}

func (m *MockQAccounts) UpsertAccounts(ctx context.Context, accounts []AccountEntry) error {
	a := m.Called(ctx, accounts)
	return a.Error(0)
}

func (m *MockQAccounts) RemoveAccounts(ctx context.Context, accountIDs []string) (int64, error) {
	a := m.Called(ctx, accountIDs)
	return a.Get(0).(int64), a.Error(1)
}
