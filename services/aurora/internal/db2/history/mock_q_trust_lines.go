package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQTrustLines is a mock implementation of the QOffers interface
type MockQTrustLines struct {
	mock.Mock
}

func (m *MockQTrustLines) GetTrustLinesByKeys(ctx context.Context, keys []string) ([]TrustLine, error) {
	a := m.Called(ctx, keys)
	return a.Get(0).([]TrustLine), a.Error(1)
}

func (m *MockQTrustLines) UpsertTrustLines(ctx context.Context, trustLines []TrustLine) error {
	a := m.Called(ctx, trustLines)
	return a.Error(0)
}

func (m *MockQTrustLines) RemoveTrustLines(ctx context.Context, ledgerKeys []string) (int64, error) {
	a := m.Called(ctx, ledgerKeys)
	return a.Get(0).(int64), a.Error(1)
}
