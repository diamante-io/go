package history

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/xdr"
)

// MockQAssetStats is a mock implementation of the QAssetStats interface
type MockQAssetStats struct {
	mock.Mock
}

func (m *MockQAssetStats) InsertAssetStats(ctx context.Context, assetStats []ExpAssetStat, batchSize int) error {
	a := m.Called(ctx, assetStats, batchSize)
	return a.Error(0)
}

func (m *MockQAssetStats) InsertAssetStat(ctx context.Context, assetStat ExpAssetStat) (int64, error) {
	a := m.Called(ctx, assetStat)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockQAssetStats) UpdateAssetStat(ctx context.Context, assetStat ExpAssetStat) (int64, error) {
	a := m.Called(ctx, assetStat)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockQAssetStats) GetAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (ExpAssetStat, error) {
	a := m.Called(ctx, assetType, assetCode, assetIssuer)
	return a.Get(0).(ExpAssetStat), a.Error(1)
}

func (m *MockQAssetStats) RemoveAssetStat(ctx context.Context, assetType xdr.AssetType, assetCode, assetIssuer string) (int64, error) {
	a := m.Called(ctx, assetType, assetCode, assetIssuer)
	return a.Get(0).(int64), a.Error(1)
}

func (m *MockQAssetStats) GetAssetStats(ctx context.Context, assetCode, assetIssuer string, page db2.PageQuery) ([]ExpAssetStat, error) {
	a := m.Called(ctx, assetCode, assetIssuer, page)
	return a.Get(0).([]ExpAssetStat), a.Error(1)
}

func (m *MockQAssetStats) CountTrustLines(ctx context.Context) (int, error) {
	a := m.Called(ctx)
	return a.Get(0).(int), a.Error(1)
}
