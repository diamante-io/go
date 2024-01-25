package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQLiquidityPools is a mock implementation of the QAccounts interface
type MockQLiquidityPools struct {
	mock.Mock
}

func (m *MockQLiquidityPools) UpsertLiquidityPools(ctx context.Context, lps []LiquidityPool) error {
	a := m.Called(ctx, lps)
	return a.Error(0)
}

func (m *MockQLiquidityPools) CountLiquidityPools(ctx context.Context) (int, error) {
	a := m.Called(ctx)
	return a.Get(0).(int), a.Error(1)
}

func (m *MockQLiquidityPools) GetLiquidityPoolsByID(ctx context.Context, poolIDs []string) ([]LiquidityPool, error) {
	a := m.Called(ctx, poolIDs)
	return a.Get(0).([]LiquidityPool), a.Error(1)
}

func (m *MockQLiquidityPools) FindLiquidityPoolByID(ctx context.Context, liquidityPoolID string) (LiquidityPool, error) {
	a := m.Called(ctx, liquidityPoolID)
	return a.Get(0).(LiquidityPool), a.Error(1)
}

func (m *MockQLiquidityPools) GetAllLiquidityPools(ctx context.Context) ([]LiquidityPool, error) {
	a := m.Called(ctx)
	return a.Get(0).([]LiquidityPool), a.Error(1)
}

func (m *MockQLiquidityPools) GetUpdatedLiquidityPools(ctx context.Context, sequence uint32) ([]LiquidityPool, error) {
	a := m.Called(ctx, sequence)
	return a.Get(0).([]LiquidityPool), a.Error(1)
}

func (m *MockQLiquidityPools) CompactLiquidityPools(ctx context.Context, cutOffSequence uint32) (int64, error) {
	a := m.Called(ctx, cutOffSequence)
	return a.Get(0).(int64), a.Error(1)
}
