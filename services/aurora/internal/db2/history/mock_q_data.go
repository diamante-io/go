package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQData is a mock implementation of the QAccounts interface
type MockQData struct {
	mock.Mock
}

func (m *MockQData) GetAccountDataByKeys(ctx context.Context, keys []AccountDataKey) ([]Data, error) {
	a := m.Called(ctx)
	return a.Get(0).([]Data), a.Error(1)
}

func (m *MockQData) CountAccountsData(ctx context.Context) (int, error) {
	a := m.Called(ctx)
	return a.Get(0).(int), a.Error(1)
}

func (m *MockQData) UpsertAccountData(ctx context.Context, data []Data) error {
	a := m.Called(ctx, data)
	return a.Error(0)
}

func (m *MockQData) RemoveAccountData(ctx context.Context, keys []AccountDataKey) (int64, error) {
	a := m.Called(ctx, keys)
	return a.Get(0).(int64), a.Error(1)
}
