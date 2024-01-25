package history

import (
	"database/sql"
	"math"
	"testing"

	"go/services/aurora/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderBookSummaryRequiresTransaction(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	_, err := q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, 10)
	assert.EqualError(t, err, "cannot be called outside of a transaction")

	assert.NoError(t, q.Begin())
	defer q.Rollback()

	_, err = q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, 10)
	assert.EqualError(t, err, "should only be called in a repeatable read transaction")
}

func TestGetOrderBookSummary(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	asksButNoBids := []Offer{twoEurOffer}
	asksButNoBidsResponse := OrderBookSummary{
		Asks: []PriceLevel{
			{
				Pricen: int32(twoEurOffer.Pricen),
				Priced: int32(twoEurOffer.Priced),
				Pricef: "2.0000000",
				Amount: "0.0000500",
			},
		},
	}

	sellEurOffer := Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(15),

		BuyingAsset:  nativeAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	bidsButNoAsks := []Offer{sellEurOffer}
	bidsButNoAsksResponse := OrderBookSummary{
		Bids: []PriceLevel{
			{
				Pricen: int32(sellEurOffer.Priced),
				Priced: int32(sellEurOffer.Pricen),
				Pricef: "0.5000000",
				Amount: "0.0000500",
			},
		},
	}

	otherEurOffer := Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(6),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(math.MaxInt64),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	nonCanonicalPriceTwoEurOffer := Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(30),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(2 * 15),
		Priced:             int32(1 * 15),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	otherSellEurOffer := Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(17),

		BuyingAsset:  nativeAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(9),
		Priced:             int32(5),
		Price:              float64(9) / float64(5),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	fullOffers := []Offer{
		twoEurOffer,
		otherEurOffer,
		nonCanonicalPriceTwoEurOffer,
		threeEurOffer,
		sellEurOffer,
		otherSellEurOffer,
	}

	fullResponse := OrderBookSummary{
		Asks: []PriceLevel{
			{
				Pricen: int32(twoEurOffer.Pricen),
				Priced: int32(twoEurOffer.Priced),
				Pricef: "2.0000000",
				Amount: "922337203685.4776807",
			},
			{
				Pricen: int32(threeEurOffer.Pricen),
				Priced: int32(threeEurOffer.Priced),
				Pricef: "3.0000000",
				Amount: "0.0000500",
			},
		},
		Bids: []PriceLevel{
			{
				Pricen: int32(otherSellEurOffer.Priced),
				Priced: int32(otherSellEurOffer.Pricen),
				Pricef: "0.5555556",
				Amount: "0.0000500",
			},
			{
				Pricen: int32(sellEurOffer.Priced),
				Priced: int32(sellEurOffer.Pricen),
				Pricef: "0.5000000",
				Amount: "0.0000500",
			},
		},
	}

	limitResponse := OrderBookSummary{
		Asks: []PriceLevel{
			{
				Pricen: int32(twoEurOffer.Pricen),
				Priced: int32(twoEurOffer.Priced),
				Pricef: "2.0000000",
				Amount: "922337203685.4776807",
			},
		},
		Bids: []PriceLevel{
			{
				Pricen: int32(otherSellEurOffer.Priced),
				Priced: int32(otherSellEurOffer.Pricen),
				Pricef: "0.5555556",
				Amount: "0.0000500",
			},
		},
	}

	for _, testCase := range []struct {
		name     string
		offers   []Offer
		limit    int
		expected OrderBookSummary
	}{
		{
			"empty orderbook",
			[]Offer{},
			10,
			OrderBookSummary{},
		},
		{
			"orderbook with asks but no bids",
			asksButNoBids,
			10,
			asksButNoBidsResponse,
		},
		{
			"orderbook with bids but no asks",
			bidsButNoAsks,
			10,
			bidsButNoAsksResponse,
		},
		{
			"full orderbook",
			fullOffers,
			10,
			fullResponse,
		},
		{
			"limit summaries",
			fullOffers,
			1,
			limitResponse,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			assert.NoError(t, q.TruncateTables(tt.Ctx, []string{"offers"}))
			assert.NoError(t, q.UpsertOffers(tt.Ctx, testCase.offers))

			assert.NoError(t, q.BeginTx(&sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
				ReadOnly:  true,
			}))
			defer q.Rollback()

			result, err := q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, testCase.limit)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestGetOrderBookSummaryExcludesRemovedOffers(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	sellEurOffer := Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(15),

		BuyingAsset:  nativeAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}

	offers := []Offer{
		twoEurOffer,
		threeEurOffer,
		sellEurOffer,
	}

	assert.NoError(t, q.UpsertOffers(tt.Ctx, offers))

	assert.NoError(t, q.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}))

	result, err := q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, 100)
	assert.NoError(t, err)
	assert.Len(t, result.Asks, 2)
	assert.Len(t, result.Bids, 1)

	assert.NoError(t, q.Rollback())

	var offersToDelete []Offer
	for i, offer := range offers {
		toDelete := offer
		toDelete.Deleted = true
		toDelete.LastModifiedLedger = uint32(i + 2)
		offersToDelete = append(offersToDelete, toDelete)
	}
	assert.NoError(t, q.UpsertOffers(tt.Ctx, offersToDelete))

	assert.NoError(t, q.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}))

	result, err = q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, 100)
	assert.NoError(t, err)
	assert.Len(t, result.Asks, 0)
	assert.Len(t, result.Bids, 0)

	assert.NoError(t, q.Rollback())

	count, err := q.CompactOffers(tt.Ctx, 1000)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(offers)), count)

	assert.NoError(t, q.BeginTx(&sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	}))

	result, err = q.GetOrderBookSummary(tt.Ctx, nativeAsset, eurAsset, 100)
	assert.NoError(t, err)
	assert.Len(t, result.Asks, 0)
	assert.Len(t, result.Bids, 0)

	assert.NoError(t, q.Rollback())
}
