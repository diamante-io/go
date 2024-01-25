package history

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockQOffers is a mock implementation of the QOffers interface
type MockQOffers struct {
	mock.Mock
}

func (m *MockQOffers) GetAllOffers(ctx context.Context) ([]Offer, error) {
	a := m.Called(ctx)
	return a.Get(0).([]Offer), a.Error(1)
}

func (m *MockQOffers) GetOffersByIDs(ctx context.Context, ids []int64) ([]Offer, error) {
	a := m.Called(ctx, ids)
	return a.Get(0).([]Offer), a.Error(1)
}

func (m *MockQOffers) GetUpdatedOffers(ctx context.Context, newerThanSequence uint32) ([]Offer, error) {
	a := m.Called(ctx, newerThanSequence)
	return a.Get(0).([]Offer), a.Error(1)
}

func (m *MockQOffers) CountOffers(ctx context.Context) (int, error) {
	a := m.Called(ctx)
	return a.Get(0).(int), a.Error(1)
}

func (m *MockQOffers) UpsertOffers(ctx context.Context, rows []Offer) error {
	a := m.Called(ctx, rows)
	return a.Error(0)
}

func (m *MockQOffers) CompactOffers(ctx context.Context, cutOffSequence uint32) (int64, error) {
	a := m.Called(ctx, cutOffSequence)
	return a.Get(0).(int64), a.Error(1)
}
