package ingest

import (
	"github.com/diamcircle/go/exp/orderbook"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/mock"
)

var _ orderbook.OBGraph = (*mockOrderBookGraph)(nil)

type mockOrderBookGraph struct {
	mock.Mock
}

func (m *mockOrderBookGraph) AddOffers(offer ...xdr.OfferEntry) {
	m.Called(offer)
}

func (m *mockOrderBookGraph) AddLiquidityPools(liquidityPool ...xdr.LiquidityPoolEntry) {
	m.Called(liquidityPool)
}

func (m *mockOrderBookGraph) Apply(ledger uint32) error {
	args := m.Called(ledger)
	return args.Error(0)

}

func (m *mockOrderBookGraph) Discard() {
	m.Called()
}

func (m *mockOrderBookGraph) Offers() []xdr.OfferEntry {
	args := m.Called()
	return args.Get(0).([]xdr.OfferEntry)
}

func (m *mockOrderBookGraph) LiquidityPools() []xdr.LiquidityPoolEntry {
	args := m.Called()
	return args.Get(0).([]xdr.LiquidityPoolEntry)
}

func (m *mockOrderBookGraph) RemoveOffer(id xdr.Int64) orderbook.OBGraph {
	args := m.Called(id)
	return args.Get(0).(orderbook.OBGraph)
}

func (m *mockOrderBookGraph) RemoveLiquidityPool(pool xdr.LiquidityPoolEntry) orderbook.OBGraph {
	args := m.Called(pool)
	return args.Get(0).(orderbook.OBGraph)
}

func (m *mockOrderBookGraph) Clear() {
	m.Called()
}

func (m *mockOrderBookGraph) Verify() ([]xdr.OfferEntry, []xdr.LiquidityPoolEntry, error) {
	args := m.Called()
	return args.Get(0).([]xdr.OfferEntry), args.Get(1).([]xdr.LiquidityPoolEntry), args.Error(2)
}
