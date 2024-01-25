package history

import (
	"strconv"
	"testing"

	"go/services/aurora/internal/db2"
	"go/services/aurora/internal/test"
	"go/xdr"

	"github.com/guregu/null"
)

var (
	issuer            = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	twoEurOfferSeller = xdr.MustAddress("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2")

	nativeAsset = xdr.MustNewNativeAsset()
	eurAsset    = xdr.MustNewCreditAsset("EUR", issuer.Address())
	usdAsset    = xdr.MustNewCreditAsset("USD", issuer.Address())

	eurOffer = Offer{
		SellerID: issuer.Address(),
		OfferID:  int64(4),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(1),
		Priced:             int32(1),
		Price:              float64(1),
		Flags:              1,
		LastModifiedLedger: uint32(1234),
		Sponsor:            null.StringFrom(sponsor),
	}
	twoEurOffer = Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(5),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}
	threeEurOffer = Offer{
		SellerID: twoEurOfferSeller.Address(),
		OfferID:  int64(50),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(3),
		Priced:             int32(1),
		Price:              float64(3),
		Flags:              2,
		LastModifiedLedger: uint32(1234),
	}
)

func insertOffer(tt *test.T, q *Q, offer Offer) error {
	return q.UpsertOffers(tt.Ctx, []Offer{offer})
}

func TestGetOfferByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := insertOffer(tt, q, eurOffer)
	tt.Assert.NoError(err)
	offer, err := q.GetOfferByID(tt.Ctx, eurOffer.OfferID)
	tt.Assert.NoError(err)
	tt.Assert.Equal(offer, eurOffer)
}

func TestGetNonExistentOfferByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	_, err := q.GetOfferByID(tt.Ctx, 12345)
	tt.Assert.True(q.NoRows(err))
}

func TestQueryEmptyOffers(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	offers, err := q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 0)

	updated, err := q.GetUpdatedOffers(tt.Ctx, 0)
	tt.Assert.NoError(err)
	tt.Assert.Len(updated, 0)

	count, err := q.CountOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(0, count)

	numRemoved, err := q.CompactOffers(tt.Ctx, 100)
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(0), numRemoved)
	seq, err := q.GetOfferCompactionSequence(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(uint32(100), seq)
}

func TestInsertOffers(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := insertOffer(tt, q, eurOffer)
	tt.Assert.NoError(err)
	err = insertOffer(tt, q, twoEurOffer)
	tt.Assert.NoError(err)

	offers, err := q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 2)

	offersByID := map[int64]Offer{
		offers[0].OfferID: offers[0],
		offers[1].OfferID: offers[1],
	}

	tt.Assert.Equal(offersByID[eurOffer.OfferID], eurOffer)
	tt.Assert.Equal(offersByID[twoEurOffer.OfferID], twoEurOffer)

	count, err := q.CountOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(2, count)

	numRemoved, err := q.CompactOffers(tt.Ctx, 12350)
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(0), numRemoved)
	seq, err := q.GetOfferCompactionSequence(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(uint32(12350), seq)

	afterCompactionCount, err := q.CountOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(2, afterCompactionCount)

	afterCompactionOffers, err := q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(afterCompactionOffers, 2)
}

func TestUpdateOffer(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := insertOffer(tt, q, eurOffer)
	tt.Assert.NoError(err)

	offers, err := q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 1)

	updatedOffers, err := q.GetUpdatedOffers(tt.Ctx, 1233)
	tt.Assert.NoError(err)
	tt.Assert.Equal(offers, updatedOffers)

	updatedOffers, err = q.GetUpdatedOffers(tt.Ctx, 100)
	tt.Assert.NoError(err)
	tt.Assert.Equal(offers, updatedOffers)

	updatedOffers, err = q.GetUpdatedOffers(tt.Ctx, 1234)
	tt.Assert.NoError(err)
	tt.Assert.Len(updatedOffers, 0)

	tt.Assert.Equal(offers[0], eurOffer)

	modifiedEurOffer := eurOffer
	modifiedEurOffer.Amount -= 10

	err = q.UpsertOffers(tt.Ctx, []Offer{modifiedEurOffer})
	tt.Assert.NoError(err)

	offers, err = q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 1)

	updatedOffers, err = q.GetUpdatedOffers(tt.Ctx, 1233)
	tt.Assert.NoError(err)
	tt.Assert.Equal(offers, updatedOffers)

	updatedOffers, err = q.GetUpdatedOffers(tt.Ctx, 1235)
	tt.Assert.NoError(err)
	tt.Assert.Len(updatedOffers, 0)

	tt.Assert.Equal(offers[0], modifiedEurOffer)
}

func TestRemoveOffer(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := insertOffer(tt, q, eurOffer)
	tt.Assert.NoError(err)
	offers, err := q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 1)
	tt.Assert.Equal(offers[0], eurOffer)

	deletedOffer := eurOffer
	deletedOffer.Deleted = true
	deletedOffer.LastModifiedLedger = 1236
	err = q.UpsertOffers(tt.Ctx, []Offer{deletedOffer})
	tt.Assert.NoError(err)
	expectedUpdates := offers
	expectedUpdates[0].LastModifiedLedger = 1236
	expectedUpdates[0].Deleted = true

	offers, err = q.GetAllOffers(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 0)

	offers, err = q.GetOffersByIDs(tt.Ctx, []int64{expectedUpdates[0].OfferID})
	tt.Assert.NoError(err)
	tt.Assert.Len(offers, 0)

	_, err = q.GetOfferByID(tt.Ctx, int64(expectedUpdates[0].OfferID))
	tt.Assert.True(q.NoRows(err))

	updated, err := q.GetUpdatedOffers(tt.Ctx, 1234)
	tt.Assert.NoError(err)
	tt.Assert.Equal(expectedUpdates, updated)

	count, err := q.CompactOffers(tt.Ctx, 1235)
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(0), count)

	updated, err = q.GetUpdatedOffers(tt.Ctx, 1234)
	tt.Assert.NoError(err)
	tt.Assert.Equal(expectedUpdates, updated)

	count, err = q.CompactOffers(tt.Ctx, 1236)
	tt.Assert.NoError(err)
	tt.Assert.Equal(int64(1), count)
	seq, err := q.GetOfferCompactionSequence(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(uint32(1236), seq)

	updated, err = q.GetUpdatedOffers(tt.Ctx, 1234)
	tt.Assert.NoError(err)
	tt.Assert.Len(updated, 0)
}

func TestGetOffers(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	err := insertOffer(tt, q, eurOffer)
	tt.Assert.NoError(err)
	err = insertOffer(tt, q, twoEurOffer)
	tt.Assert.NoError(err)

	// check removed offers aren't included in GetOffer queries
	err = insertOffer(tt, q, threeEurOffer)
	tt.Assert.NoError(err)
	deletedOffer := threeEurOffer
	deletedOffer.Deleted = true
	deletedOffer.LastModifiedLedger = 1235
	err = q.UpsertOffers(tt.Ctx, []Offer{deletedOffer})
	tt.Assert.NoError(err)

	pageQuery, err := db2.NewPageQuery("", false, "", 10)
	tt.Assert.NoError(err)

	t.Run("Filter by selling asset", func(t *testing.T) {
		query := OffersQuery{
			PageQuery: pageQuery,
			Selling:   &usdAsset,
		}

		offers, err := q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 0)

		query = OffersQuery{
			PageQuery: pageQuery,
			Selling:   &nativeAsset,
		}

		offers, err = q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 2)

		for _, offer := range offers {
			tt.Assert.Equal(nativeAsset, offer.SellingAsset)
		}
	})

	t.Run("Filter by buying asset", func(t *testing.T) {
		query := OffersQuery{
			PageQuery: pageQuery,
			Buying:    &eurAsset,
		}

		offers, err := q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 2)

		for _, offer := range offers {
			tt.Assert.Equal(eurAsset, offer.BuyingAsset)
		}

		query = OffersQuery{
			PageQuery: pageQuery,
			Buying:    &usdAsset,
		}

		offers, err = q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 0)
	})

	t.Run("Filter by seller", func(t *testing.T) {
		sellerID := issuer.Address()
		query := OffersQuery{
			PageQuery: pageQuery,
			SellerID:  sellerID,
		}

		offers, err := q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 1)

		tt.Assert.Equal(offers[0], eurOffer)
	})

	t.Run("Filter by sponsor", func(t *testing.T) {
		query := OffersQuery{
			PageQuery: pageQuery,
			Sponsor:   sponsor,
		}

		offers, err := q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 1)

		tt.Assert.Equal(offers[0], eurOffer)
	})

	t.Run("PageQuery", func(t *testing.T) {
		pageQuery, err := db2.NewPageQuery("", false, "", 10)
		tt.Assert.NoError(err)

		query := OffersQuery{
			PageQuery: pageQuery,
		}

		offers, err := q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 2)

		offersByID := map[int64]Offer{
			offers[0].OfferID: offers[0],
			offers[1].OfferID: offers[1],
		}

		tt.Assert.Equal(offersByID[eurOffer.OfferID], eurOffer)
		tt.Assert.Equal(offersByID[twoEurOffer.OfferID], twoEurOffer)

		pageQuery, err = db2.NewPageQuery("", false, "asc", 1)
		tt.Assert.NoError(err)
		query = OffersQuery{
			PageQuery: pageQuery,
		}

		offers, err = q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 1)

		tt.Assert.Equal(offers[0], eurOffer)

		pageQuery, err = db2.NewPageQuery("", false, "desc", 1)
		tt.Assert.NoError(err)
		query = OffersQuery{
			PageQuery: pageQuery,
		}

		offers, err = q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 1)

		tt.Assert.Equal(offers[0], twoEurOffer)

		pageQuery, err = db2.NewPageQuery(
			strconv.FormatInt(int64(eurOffer.OfferID), 10),
			false,
			"",
			10,
		)
		tt.Assert.NoError(err)
		query = OffersQuery{
			PageQuery: pageQuery,
		}

		offers, err = q.GetOffers(tt.Ctx, query)
		tt.Assert.NoError(err)
		tt.Assert.Len(offers, 1)

		tt.Assert.Equal(offers[0], twoEurOffer)
	})
}
