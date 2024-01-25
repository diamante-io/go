//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package ingest

import (
	"context"
	"fmt"
	"testing"

	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest/processors"
	"go/xdr"

	"github.com/stretchr/testify/suite"
)

type IngestionStatusTestSuite struct {
	suite.Suite
	ctx      context.Context
	historyQ *mockDBQ
	stream   *OrderBookStream
}

func TestIngestionStatus(t *testing.T) {
	suite.Run(t, new(IngestionStatusTestSuite))
}

func (t *IngestionStatusTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.historyQ = &mockDBQ{}
	t.stream = NewOrderBookStream(t.historyQ, &mockOrderBookGraph{})
}

func (t *IngestionStatusTestSuite) TearDownTest() {
	t.historyQ.AssertExpectations(t.T())
}

func (t *IngestionStatusTestSuite) TestGetExpStateInvalidError() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, fmt.Errorf("state invalid error")).
		Once()
	_, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().EqualError(err, "Error from GetExpStateInvalid: state invalid error")
}

func (t *IngestionStatusTestSuite) TestGetLatestLedgerError() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(0), fmt.Errorf("latest ledger error")).
		Once()
	_, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().EqualError(err, "Error from GetLatestHistoryLedger: latest ledger error")
}

func (t *IngestionStatusTestSuite) TestGetLastLedgerIngestNonBlockingError() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(0), fmt.Errorf("ingest error")).
		Once()

	_, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().EqualError(err, "Error from GetLastLedgerIngestNonBlocking: ingest error")
}

func (t *IngestionStatusTestSuite) TestGetOfferCompactionSequenceError() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetOfferCompactionSequence", t.ctx).
		Return(uint32(0), fmt.Errorf("compaction error")).
		Once()

	_, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().EqualError(err, "Error from GetOfferCompactionSequence: compaction error")
}

func (t *IngestionStatusTestSuite) TestLiquidityPoolCompactionSequenceError() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetOfferCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	t.historyQ.On("GetLiquidityPoolCompactionSequence", t.ctx).
		Return(uint32(0), fmt.Errorf("compaction error")).
		Once()

	_, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().EqualError(err, "Error from GetLiquidityPoolCompactionSequence: compaction error")
}

func (t *IngestionStatusTestSuite) TestStateInvalid() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(true, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetOfferCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	t.historyQ.On("GetLiquidityPoolCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	status, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().NoError(err)
	t.Assert().Equal(ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      true,
		LastIngestedLedger:                200,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}, status)
}

func (t *IngestionStatusTestSuite) TestHistoryInconsistentWithState() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(true, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(200), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(201), nil).
		Once()

	t.historyQ.On("GetOfferCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	t.historyQ.On("GetLiquidityPoolCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	status, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().NoError(err)
	t.Assert().Equal(ingestionStatus{
		HistoryConsistentWithState:        false,
		StateInvalid:                      true,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}, status)
}

func (t *IngestionStatusTestSuite) TestHistoryLatestLedgerZero() {
	t.historyQ.On("GetExpStateInvalid", t.ctx).
		Return(false, nil).
		Once()

	t.historyQ.On("GetLatestHistoryLedger", t.ctx).
		Return(uint32(0), nil).
		Once()

	t.historyQ.On("GetLastLedgerIngestNonBlocking", t.ctx).
		Return(uint32(201), nil).
		Once()

	t.historyQ.On("GetOfferCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	t.historyQ.On("GetLiquidityPoolCompactionSequence", t.ctx).
		Return(uint32(100), nil).
		Once()

	status, err := t.stream.getIngestionStatus(t.ctx)
	t.Assert().NoError(err)
	t.Assert().Equal(ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}, status)
}

type UpdateOrderBookStreamTestSuite struct {
	suite.Suite
	ctx      context.Context
	historyQ *mockDBQ
	graph    *mockOrderBookGraph
	stream   *OrderBookStream
}

func TestUpdateOrderBookStream(t *testing.T) {
	suite.Run(t, new(UpdateOrderBookStreamTestSuite))
}

func (t *UpdateOrderBookStreamTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.historyQ = &mockDBQ{}
	t.graph = &mockOrderBookGraph{}
	t.stream = NewOrderBookStream(t.historyQ, t.graph)
}

func (t *UpdateOrderBookStreamTestSuite) TearDownTest() {
	t.historyQ.AssertExpectations(t.T())
	t.graph.AssertExpectations(t.T())
}

func (t *UpdateOrderBookStreamTestSuite) TestGetAllOffersError() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Clear").Return().Once()
	t.graph.On("Discard").Return().Once()
	t.historyQ.On("GetAllOffers", t.ctx).
		Return([]history.Offer{}, fmt.Errorf("offers error")).
		Once()

	t.stream.lastLedger = 300
	_, err := t.stream.update(t.ctx, status)
	t.Assert().EqualError(err, "Error from GetAllOffers: offers error")
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
}

func (t *UpdateOrderBookStreamTestSuite) TestResetApplyError() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Clear").Return().Once()
	t.graph.On("Discard").Return().Once()

	sellerID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	offer := history.Offer{OfferID: 1, SellerID: sellerID}
	offerEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  1,
	}}
	otherOffer := history.Offer{OfferID: 20, SellerID: sellerID}
	otherOfferEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  20,
	}}
	t.historyQ.On("GetAllOffers", t.ctx).
		Return([]history.Offer{offer, otherOffer}, nil).
		Once()

	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return([]history.LiquidityPool{}, nil).
		Once()

	t.graph.On("AddOffers", offerEntry).Return().Once()
	t.graph.On("AddOffers", otherOfferEntry).Return().Once()

	t.graph.On("Apply", status.LastIngestedLedger).
		Return(fmt.Errorf("apply error")).
		Once()

	t.stream.lastLedger = 300
	_, err := t.stream.update(t.ctx, status)
	t.Assert().EqualError(err, "Error applying changes to order book: apply error")
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
}

func (t *UpdateOrderBookStreamTestSuite) mockReset(status ingestionStatus) {
	t.graph.On("Clear").Return().Once()
	t.graph.On("Discard").Return().Once()

	sellerID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	offer := history.Offer{OfferID: 1, SellerID: sellerID}
	offerEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  1,
	}}
	otherOffer := history.Offer{OfferID: 20, SellerID: sellerID}
	otherOfferEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  20,
	}}
	offers := []history.Offer{offer, otherOffer}
	t.historyQ.On("GetAllOffers", t.ctx).
		Return(offers, nil).
		Once()

	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return([]history.LiquidityPool{}, nil).
		Once()

	t.graph.On("AddOffers", offerEntry).Return().Once()
	t.graph.On("AddOffers", otherOfferEntry).Return().Once()

	t.graph.On("Apply", status.LastIngestedLedger).
		Return(nil).
		Once()
}

func (t *UpdateOrderBookStreamTestSuite) TestFirstUpdateSucceeds() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.mockReset(status)

	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(201), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestInvalidState() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      true,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Clear").Return().Once()

	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
	t.Assert().True(reset)

	t.stream.lastLedger = 123

	t.graph.On("Clear").Return().Once()

	reset, err = t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestHistoryInconsistentWithState() {
	status := ingestionStatus{
		HistoryConsistentWithState:        false,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Clear").Return().Once()

	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
	t.Assert().True(reset)

	t.stream.lastLedger = 123

	t.graph.On("Clear").Return().Once()

	reset, err = t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(0), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestOfferCompactionDoesNotMatchLiquidityPoolCompaction() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 110,
	}
	t.mockReset(status)

	t.stream.lastLedger = 201
	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(201), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestLastIngestedLedgerBehindStream() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.mockReset(status)

	t.stream.lastLedger = 300
	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(201), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestStreamBehindLastCompactionLedger() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.mockReset(status)

	t.stream.lastLedger = 99
	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(201), t.stream.lastLedger)
	t.Assert().True(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestStreamLedgerEqualsLastIngestedLedger() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}

	t.stream.lastLedger = 201
	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(uint32(201), t.stream.lastLedger)
	t.Assert().False(reset)
}

func (t *UpdateOrderBookStreamTestSuite) TestGetUpdatedOffersError() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Discard").Return().Once()

	t.stream.lastLedger = 100
	t.historyQ.MockQOffers.On("GetUpdatedOffers", t.ctx, uint32(100)).
		Return([]history.Offer{}, fmt.Errorf("updated offers error")).
		Once()

	_, err := t.stream.update(t.ctx, status)
	t.Assert().EqualError(err, "Error from GetUpdatedOffers: updated offers error")
	t.Assert().Equal(uint32(100), t.stream.lastLedger)
}

func (t *UpdateOrderBookStreamTestSuite) TestGetUpdatedLiquidityPoolsError() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}
	t.graph.On("Discard").Return().Once()

	t.stream.lastLedger = 100
	t.historyQ.MockQOffers.On("GetUpdatedOffers", t.ctx, uint32(100)).
		Return([]history.Offer{}, nil).
		Once()

	t.historyQ.MockQLiquidityPools.On("GetUpdatedLiquidityPools", t.ctx, t.stream.lastLedger).
		Return([]history.LiquidityPool{}, fmt.Errorf("updated liquidity pools error")).
		Once()

	_, err := t.stream.update(t.ctx, status)
	t.Assert().EqualError(err, "Error from GetUpdatedLiquidityPools: updated liquidity pools error")
	t.Assert().Equal(uint32(100), t.stream.lastLedger)
}

func (t *UpdateOrderBookStreamTestSuite) mockUpdate() {
	t.stream.lastLedger = 100

	t.graph.On("Discard").Return().Once()
	sellerID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	offer := history.Offer{OfferID: 1, SellerID: sellerID, LastModifiedLedger: 101}
	offerEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  1,
	}}
	otherOffer := history.Offer{OfferID: 20, SellerID: sellerID, LastModifiedLedger: 102}
	otherOfferEntry := []xdr.OfferEntry{{
		SellerId: xdr.MustAddress(sellerID),
		OfferId:  20,
	}}
	deletedOffer := history.Offer{OfferID: 30, SellerID: sellerID, LastModifiedLedger: 103, Deleted: true}
	offers := []history.Offer{offer, otherOffer, deletedOffer}
	t.historyQ.MockQOffers.On("GetUpdatedOffers", t.ctx, t.stream.lastLedger).
		Return(offers, nil).
		Once()

	t.historyQ.MockQLiquidityPools.On("GetUpdatedLiquidityPools", t.ctx, t.stream.lastLedger).
		Return([]history.LiquidityPool{}, nil).
		Once()

	t.graph.On("AddOffers", offerEntry).Return().Once()
	t.graph.On("AddOffers", otherOfferEntry).Return().Once()
	t.graph.On("RemoveOffer", xdr.Int64(deletedOffer.OfferID)).Return(t.graph).Once()
}

func (t *UpdateOrderBookStreamTestSuite) TestApplyUpdatesError() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}

	t.mockUpdate()

	t.graph.On("Apply", status.LastIngestedLedger).
		Return(fmt.Errorf("apply error")).
		Once()

	_, err := t.stream.update(t.ctx, status)
	t.Assert().EqualError(err, "Error applying changes to order book: apply error")
	t.Assert().Equal(uint32(100), t.stream.lastLedger)
}

func (t *UpdateOrderBookStreamTestSuite) TestApplyUpdatesSucceeds() {
	status := ingestionStatus{
		HistoryConsistentWithState:        true,
		StateInvalid:                      false,
		LastIngestedLedger:                201,
		LastOfferCompactionLedger:         100,
		LastLiquidityPoolCompactionLedger: 100,
	}

	t.mockUpdate()

	t.graph.On("Apply", status.LastIngestedLedger).
		Return(nil).
		Once()

	reset, err := t.stream.update(t.ctx, status)
	t.Assert().NoError(err)
	t.Assert().Equal(status.LastIngestedLedger, t.stream.lastLedger)
	t.Assert().False(reset)
}

type VerifyOffersStreamTestSuite struct {
	suite.Suite
	ctx      context.Context
	historyQ *mockDBQ
	graph    *mockOrderBookGraph
	stream   *OrderBookStream
}

func TestVerifyOffersStreamTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyOffersStreamTestSuite))
}

func (t *VerifyOffersStreamTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.historyQ = &mockDBQ{}
	t.graph = &mockOrderBookGraph{}
	t.stream = NewOrderBookStream(t.historyQ, t.graph)

	sellerID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	otherSellerID := "GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK"
	t.graph.On("Offers").Return([]xdr.OfferEntry{
		{
			SellerId: xdr.MustAddress(sellerID),
			OfferId:  1,
			Selling:  xdr.MustNewNativeAsset(),
			Buying:   xdr.MustNewCreditAsset("USD", sellerID),
			Amount:   123,
			Price: xdr.Price{
				N: 1,
				D: 2,
			},
			Flags: 1,
			Ext:   xdr.OfferEntryExt{},
		},
		{
			SellerId: xdr.MustAddress(otherSellerID),
			OfferId:  3,
			Selling:  xdr.MustNewCreditAsset("EUR", sellerID),
			Buying:   xdr.MustNewCreditAsset("CHF", sellerID),
			Amount:   9,
			Price: xdr.Price{
				N: 3,
				D: 1,
			},
			Flags: 0,
			Ext:   xdr.OfferEntryExt{},
		},
	}).Once()
}

func (t *VerifyOffersStreamTestSuite) TearDownTest() {
	t.historyQ.AssertExpectations(t.T())
	t.graph.AssertExpectations(t.T())
}

func (t *VerifyOffersStreamTestSuite) TestGetAllOffersError() {
	t.historyQ.On("GetAllOffers", t.ctx).
		Return([]history.Offer{}, fmt.Errorf("offers error")).
		Once()

	offersOk, err := t.stream.verifyAllOffers(t.ctx, t.graph.Offers())
	t.Assert().EqualError(err, "Error from GetAllOffers: offers error")
	t.Assert().False(offersOk)
}

func (t *VerifyOffersStreamTestSuite) TestEmptyDBOffers() {
	var offers []history.Offer
	t.historyQ.On("GetAllOffers", t.ctx).Return(offers, nil).Once()

	offersOk, err := t.stream.verifyAllOffers(t.ctx, t.graph.Offers())
	t.Assert().NoError(err)
	t.Assert().False(offersOk)
}

func (t *VerifyOffersStreamTestSuite) TestLengthMismatch() {
	offers := []history.Offer{
		{
			OfferID:            1,
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			SellingAsset:       xdr.MustNewNativeAsset(),
			BuyingAsset:        xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Amount:             123,
			Pricen:             1,
			Priced:             2,
			Price:              0.5,
			Flags:              1,
			Deleted:            false,
			LastModifiedLedger: 1,
		},
	}
	t.historyQ.On("GetAllOffers", t.ctx).Return(offers, nil).Once()

	offersOk, err := t.stream.verifyAllOffers(t.ctx, t.graph.Offers())
	t.Assert().NoError(err)
	t.Assert().False(offersOk)
}

func (t *VerifyOffersStreamTestSuite) TestContentMismatch() {
	offers := []history.Offer{
		{
			OfferID:            1,
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			SellingAsset:       xdr.MustNewNativeAsset(),
			BuyingAsset:        xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Amount:             123,
			Pricen:             1,
			Priced:             2,
			Price:              0.5,
			Flags:              1,
			Deleted:            false,
			LastModifiedLedger: 1,
		},
		{
			OfferID:            3,
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			SellingAsset:       xdr.MustNewNativeAsset(),
			BuyingAsset:        xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Amount:             123,
			Pricen:             1,
			Priced:             2,
			Price:              0.5,
			Flags:              1,
			Deleted:            false,
			LastModifiedLedger: 1,
		},
	}
	t.historyQ.On("GetAllOffers", t.ctx).Return(offers, nil).Once()

	t.stream.lastLedger = 300
	offersOk, err := t.stream.verifyAllOffers(t.ctx, t.graph.Offers())
	t.Assert().NoError(err)
	t.Assert().False(offersOk)
}

func (t *VerifyOffersStreamTestSuite) TestSuccess() {
	offers := []history.Offer{
		{
			OfferID:            1,
			SellerID:           "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			SellingAsset:       xdr.MustNewNativeAsset(),
			BuyingAsset:        xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Amount:             123,
			Pricen:             1,
			Priced:             2,
			Price:              0.5,
			Flags:              1,
			Deleted:            false,
			LastModifiedLedger: 1,
		},
		{
			OfferID:            3,
			SellerID:           "GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK",
			SellingAsset:       xdr.MustNewCreditAsset("EUR", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			BuyingAsset:        xdr.MustNewCreditAsset("CHF", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Amount:             9,
			Pricen:             3,
			Priced:             1,
			Price:              3,
			Flags:              0,
			Deleted:            false,
			LastModifiedLedger: 1,
		},
	}
	t.historyQ.On("GetAllOffers", t.ctx).Return(offers, nil).Once()

	offersOk, err := t.stream.verifyAllOffers(t.ctx, t.graph.Offers())
	t.Assert().NoError(err)
	t.Assert().True(offersOk)
}

type VerifyLiquidityPoolsStreamTestSuite struct {
	suite.Suite
	ctx      context.Context
	historyQ *mockDBQ
	graph    *mockOrderBookGraph
	stream   *OrderBookStream
}

func TestVerifyLiquidityPoolsStreamTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyLiquidityPoolsStreamTestSuite))
}

func (t *VerifyLiquidityPoolsStreamTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.historyQ = &mockDBQ{}
	t.graph = &mockOrderBookGraph{}
	t.stream = NewOrderBookStream(t.historyQ, t.graph)

	t.graph.On("LiquidityPools").Return([]xdr.LiquidityPoolEntry{
		{
			LiquidityPoolId: xdr.PoolId{1, 2, 3},
			Body: xdr.LiquidityPoolEntryBody{
				Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
				ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
					Params: xdr.LiquidityPoolConstantProductParameters{
						AssetA: xdr.MustNewNativeAsset(),
						AssetB: xdr.MustNewCreditAsset("USD", issuer.Address()),
						Fee:    xdr.LiquidityPoolFeeV18,
					},
					ReserveA:                 789,
					ReserveB:                 456,
					TotalPoolShares:          11,
					PoolSharesTrustLineCount: 13,
				},
			},
		},
		{
			LiquidityPoolId: xdr.PoolId{4, 5, 6},
			Body: xdr.LiquidityPoolEntryBody{
				Type: xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
				ConstantProduct: &xdr.LiquidityPoolEntryConstantProduct{
					Params: xdr.LiquidityPoolConstantProductParameters{
						AssetA: xdr.MustNewNativeAsset(),
						AssetB: xdr.MustNewCreditAsset("EUR", issuer.Address()),
						Fee:    xdr.LiquidityPoolFeeV18,
					},
					ReserveA:                 19,
					ReserveB:                 1234,
					TotalPoolShares:          456,
					PoolSharesTrustLineCount: 90,
				},
			},
		},
	}).Once()
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TearDownTest() {
	t.historyQ.AssertExpectations(t.T())
	t.graph.AssertExpectations(t.T())
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TestGetAllLiquidityPoolsError() {
	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return([]history.LiquidityPool{}, fmt.Errorf("liquidity pools error")).
		Once()

	liquidityPoolsOk, err := t.stream.verifyAllLiquidityPools(t.ctx, t.graph.LiquidityPools())
	t.Assert().EqualError(err, "Error from GetAllLiquidityPools: liquidity pools error")
	t.Assert().False(liquidityPoolsOk)
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TestEmptyDBOffers() {
	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return([]history.LiquidityPool{}, nil).
		Once()

	liquidityPoolsOk, err := t.stream.verifyAllLiquidityPools(t.ctx, t.graph.LiquidityPools())
	t.Assert().NoError(err)
	t.Assert().False(liquidityPoolsOk)
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TestLengthMismatch() {
	liquidityPools := []history.LiquidityPool{
		{
			PoolID:         processors.PoolIDToString(xdr.PoolId{1, 2, 3}),
			Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			Fee:            xdr.LiquidityPoolFeeV18,
			TrustlineCount: 13,
			ShareCount:     11,
			AssetReserves: history.LiquidityPoolAssetReserves{
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewNativeAsset(),
					Reserve: 789,
				},
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewCreditAsset("USD", issuer.Address()),
					Reserve: 456,
				},
			},
			LastModifiedLedger: 100,
			Deleted:            false,
		},
	}

	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return(liquidityPools, nil).
		Once()

	liquidityPoolsOk, err := t.stream.verifyAllLiquidityPools(t.ctx, t.graph.LiquidityPools())
	t.Assert().NoError(err)
	t.Assert().False(liquidityPoolsOk)
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TestContentMismatch() {
	liquidityPools := []history.LiquidityPool{
		{
			PoolID:         processors.PoolIDToString(xdr.PoolId{1, 2, 3}),
			Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			Fee:            xdr.LiquidityPoolFeeV18,
			TrustlineCount: 0,
			ShareCount:     11,
			AssetReserves: history.LiquidityPoolAssetReserves{
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewNativeAsset(),
					Reserve: 789,
				},
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewCreditAsset("USD", issuer.Address()),
					Reserve: 456,
				},
			},
			LastModifiedLedger: 100,
			Deleted:            false,
		},
		{
			PoolID:         processors.PoolIDToString(xdr.PoolId{4, 5, 6}),
			Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			Fee:            xdr.LiquidityPoolFeeV18,
			TrustlineCount: 90,
			ShareCount:     456,
			AssetReserves: history.LiquidityPoolAssetReserves{
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewNativeAsset(),
					Reserve: 19,
				},
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewCreditAsset("EUR", issuer.Address()),
					Reserve: 1234,
				},
			},
			LastModifiedLedger: 50,
			Deleted:            false,
		},
	}
	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return(liquidityPools, nil).
		Once()

	liquidityPoolsOk, err := t.stream.verifyAllLiquidityPools(t.ctx, t.graph.LiquidityPools())
	t.Assert().NoError(err)
	t.Assert().False(liquidityPoolsOk)
}

func (t *VerifyLiquidityPoolsStreamTestSuite) TestSuccess() {
	liquidityPools := []history.LiquidityPool{
		{
			PoolID:         processors.PoolIDToString(xdr.PoolId{1, 2, 3}),
			Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			Fee:            xdr.LiquidityPoolFeeV18,
			TrustlineCount: 13,
			ShareCount:     11,
			AssetReserves: history.LiquidityPoolAssetReserves{
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewNativeAsset(),
					Reserve: 789,
				},
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewCreditAsset("USD", issuer.Address()),
					Reserve: 456,
				},
			},
			LastModifiedLedger: 100,
			Deleted:            false,
		},
		{
			PoolID:         processors.PoolIDToString(xdr.PoolId{4, 5, 6}),
			Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
			Fee:            xdr.LiquidityPoolFeeV18,
			TrustlineCount: 90,
			ShareCount:     456,
			AssetReserves: history.LiquidityPoolAssetReserves{
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewNativeAsset(),
					Reserve: 19,
				},
				history.LiquidityPoolAssetReserve{
					Asset:   xdr.MustNewCreditAsset("EUR", issuer.Address()),
					Reserve: 1234,
				},
			},
			LastModifiedLedger: 50,
			Deleted:            false,
		},
	}
	t.historyQ.MockQLiquidityPools.On("GetAllLiquidityPools", t.ctx).
		Return(liquidityPools, nil).
		Once()

	offersOk, err := t.stream.verifyAllLiquidityPools(t.ctx, t.graph.LiquidityPools())
	t.Assert().NoError(err)
	t.Assert().True(offersOk)
}
