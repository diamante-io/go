package actions

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

var (
	issuer = xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	seller = xdr.MustAddress("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2")

	nativeAsset = xdr.MustNewNativeAsset()
	usdAsset    = xdr.MustNewCreditAsset("USD", issuer.Address())
	eurAsset    = xdr.MustNewCreditAsset("EUR", issuer.Address())

	eurOffer = history.Offer{
		SellerID: issuer.Address(),
		OfferID:  int64(4),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(1),
		Priced:             int32(1),
		Price:              float64(1),
		Flags:              1,
		LastModifiedLedger: uint32(3),
	}
	twoEurOffer = history.Offer{
		SellerID: seller.Address(),
		OfferID:  int64(5),

		BuyingAsset:  eurAsset,
		SellingAsset: nativeAsset,

		Amount:             int64(500),
		Pricen:             int32(2),
		Priced:             int32(1),
		Price:              float64(2),
		Flags:              2,
		LastModifiedLedger: uint32(4),
		Sponsor:            null.StringFrom(sponsor),
	}
	usdOffer = history.Offer{
		SellerID: issuer.Address(),
		OfferID:  int64(6),

		BuyingAsset:  usdAsset,
		SellingAsset: eurAsset,

		Amount:             int64(500),
		Pricen:             int32(1),
		Priced:             int32(1),
		Price:              float64(1),
		Flags:              1,
		LastModifiedLedger: uint32(4),
	}
)

func TestGetOfferByIDHandler(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := GetOfferByID{}

	ledgerCloseTime := time.Now().Unix()
	_, err := q.InsertLedger(tt.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 3,
			ScpValue: xdr.DiamcircleValue{
				CloseTime: xdr.TimePoint(ledgerCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	tt.Assert.NoError(err)

	err = q.UpsertOffers(tt.Ctx, []history.Offer{eurOffer, usdOffer})
	tt.Assert.NoError(err)

	for _, testCase := range []struct {
		name          string
		request       *http.Request
		expectedError func(error)
		expectedOffer func(interface{})
	}{
		{
			"offer id is invalid",
			makeRequest(
				t, map[string]string{}, map[string]string{"offer_id": "invalid"}, q,
			),
			func(err error) {
				tt.Assert.Error(err)
				p := err.(*problem.P)
				tt.Assert.Equal("bad_request", p.Type)
				tt.Assert.Equal("offer_id", p.Extras["invalid_field"])
				tt.Assert.Equal("Offer ID must be an integer higher than 0", p.Extras["reason"])
			},
			func(response interface{}) {
				tt.Assert.Nil(response)
			},
		},
		{
			"offer does not exist",
			makeRequest(
				t, map[string]string{}, map[string]string{"offer_id": "1234567"}, q,
			),
			func(err error) {
				tt.Assert.Equal(err, sql.ErrNoRows)
			},
			func(response interface{}) {
				tt.Assert.Nil(response)
			},
		},
		{
			"offer with ledger close time",
			makeRequest(
				t, map[string]string{}, map[string]string{"offer_id": "4"}, q,
			),
			func(err error) {
				tt.Assert.NoError(err)
			},
			func(response interface{}) {
				offer := response.(aurora.Offer)
				tt.Assert.Equal(int64(eurOffer.OfferID), offer.ID)
				tt.Assert.Equal("native", offer.Selling.Type)
				tt.Assert.Equal("credit_alphanum4", offer.Buying.Type)
				tt.Assert.Equal("EUR", offer.Buying.Code)
				tt.Assert.Equal(issuer.Address(), offer.Seller)
				tt.Assert.Equal(issuer.Address(), offer.Buying.Issuer)
				tt.Assert.Equal(int32(3), offer.LastModifiedLedger)
				tt.Assert.Equal(ledgerCloseTime, offer.LastModifiedTime.Unix())
			},
		},
		{
			"offer without ledger close time",
			makeRequest(
				t, map[string]string{}, map[string]string{"offer_id": "6"}, q,
			),
			func(err error) {
				tt.Assert.NoError(err)
			},
			func(response interface{}) {
				offer := response.(aurora.Offer)
				tt.Assert.Equal(int64(usdOffer.OfferID), offer.ID)
				tt.Assert.Equal("credit_alphanum4", offer.Selling.Type)
				tt.Assert.Equal("EUR", offer.Selling.Code)
				tt.Assert.Equal("credit_alphanum4", offer.Buying.Type)
				tt.Assert.Equal("USD", offer.Buying.Code)
				tt.Assert.Equal(issuer.Address(), offer.Seller)
				tt.Assert.Equal(issuer.Address(), offer.Selling.Issuer)
				tt.Assert.Equal(issuer.Address(), offer.Buying.Issuer)
				tt.Assert.Equal(int32(4), offer.LastModifiedLedger)
				tt.Assert.Nil(offer.LastModifiedTime)
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			offer, err := handler.GetResource(httptest.NewRecorder(), testCase.request)
			testCase.expectedError(err)
			testCase.expectedOffer(offer)
		})
	}
}

func TestGetOffersHandler(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := GetOffersHandler{}

	ledgerCloseTime := time.Now().Unix()
	_, err := q.InsertLedger(tt.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 3,
			ScpValue: xdr.DiamcircleValue{
				CloseTime: xdr.TimePoint(ledgerCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	tt.Assert.NoError(err)

	err = q.UpsertOffers(tt.Ctx, []history.Offer{eurOffer, twoEurOffer, usdOffer})
	tt.Assert.NoError(err)

	t.Run("No filter", func(t *testing.T) {
		records, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t, map[string]string{}, map[string]string{}, q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 3)

		offers := pageableToOffers(t, records)

		tt.Assert.Equal(int64(eurOffer.OfferID), offers[0].ID)
		tt.Assert.Equal("native", offers[0].Selling.Type)
		tt.Assert.Equal("credit_alphanum4", offers[0].Buying.Type)
		tt.Assert.Equal(issuer.Address(), offers[0].Seller)
		tt.Assert.Equal(issuer.Address(), offers[0].Buying.Issuer)
		tt.Assert.Equal(int32(3), offers[0].LastModifiedLedger)
		tt.Assert.Equal(ledgerCloseTime, offers[0].LastModifiedTime.Unix())
	})

	t.Run("Filter by seller", func(t *testing.T) {
		records, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"seller": issuer.Address(),
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 2)

		offers := pageableToOffers(t, records)
		for _, offer := range offers {
			tt.Assert.Equal(issuer.Address(), offer.Seller)
		}

		_, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"seller": "GCXEWJ6U4KPGTNTBY5HX4WQ2EEVPWV2QKXEYIQ32IDYIX",
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.Error(err)
		tt.Assert.IsType(&problem.P{}, err)
		p := err.(*problem.P)
		tt.Assert.Equal("bad_request", p.Type)
		tt.Assert.Equal("seller", p.Extras["invalid_field"])
		tt.Assert.Equal(
			"Account ID must start with `G` and contain 56 alphanum characters",
			p.Extras["reason"],
		)
	})

	t.Run("Filter by sponsor", func(t *testing.T) {
		records, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"sponsor": sponsor,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 1)

		offers := pageableToOffers(t, records)
		tt.Assert.Equal(int64(twoEurOffer.OfferID), offers[0].ID)

		_, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"sponsor": "GCXEWJ6U4KPGTNTBY5HX4WQ2EEVPWV2QKXEYIQ32IDYIX",
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.Error(err)
		tt.Assert.IsType(&problem.P{}, err)
		p := err.(*problem.P)
		tt.Assert.Equal("bad_request", p.Type)
		tt.Assert.Equal("sponsor", p.Extras["invalid_field"])
		tt.Assert.Equal(
			"Account ID must start with `G` and contain 56 alphanum characters",
			p.Extras["reason"],
		)
	})

	t.Run("Filter by selling asset", func(t *testing.T) {
		asset := aurora.Asset{}
		nativeAsset.Extract(&asset.Type, &asset.Code, &asset.Issuer)
		records, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"selling_asset_type": asset.Type,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 2)

		offers := pageableToOffers(t, records)
		for _, offer := range offers {
			tt.Assert.Equal(asset, offer.Selling)
		}

		asset = aurora.Asset{}
		eurAsset.Extract(&asset.Type, &asset.Code, &asset.Issuer)

		records, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"selling_asset_type":   asset.Type,
					"selling_asset_code":   asset.Code,
					"selling_asset_issuer": asset.Issuer,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 1)

		offers = pageableToOffers(t, records)
		tt.Assert.Equal(asset, offers[0].Selling)

		records, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"selling": asset.Code + ":" + asset.Issuer,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 1)

		offers = pageableToOffers(t, records)
		tt.Assert.Equal(asset, offers[0].Selling)
	})

	t.Run("Filter by buying asset", func(t *testing.T) {
		asset := aurora.Asset{}
		eurAsset.Extract(&asset.Type, &asset.Code, &asset.Issuer)

		records, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"buying_asset_type":   asset.Type,
					"buying_asset_code":   asset.Code,
					"buying_asset_issuer": asset.Issuer,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 2)

		offers := pageableToOffers(t, records)
		for _, offer := range offers {
			tt.Assert.Equal(asset, offer.Buying)
		}

		asset = aurora.Asset{}
		usdAsset.Extract(&asset.Type, &asset.Code, &asset.Issuer)

		records, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"buying_asset_type":   asset.Type,
					"buying_asset_code":   asset.Code,
					"buying_asset_issuer": asset.Issuer,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 1)

		offers = pageableToOffers(t, records)
		for _, offer := range offers {
			tt.Assert.Equal(asset, offer.Buying)
		}

		records, err = handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"buying": asset.Code + ":" + asset.Issuer,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.NoError(err)
		tt.Assert.Len(records, 1)

		offers = pageableToOffers(t, records)
		for _, offer := range offers {
			tt.Assert.Equal(asset, offer.Buying)
		}
	})

	t.Run("Wrong buying query parameter", func(t *testing.T) {
		asset := aurora.Asset{}
		eurAsset.Extract(&asset.Type, &asset.Code, &asset.Issuer)

		_, err := handler.GetResourcePage(
			httptest.NewRecorder(),
			makeRequest(
				t,
				map[string]string{
					"buying": `native\\u0026cursor=\\u0026limit=10\\u0026order=asc\\u0026selling=BTC:GAEDZ7BHMDYEMU6IJT3CTTGDUSLZWS5CQWZHGP4XUOIDG5ISH3AFAEK2`,
				},
				map[string]string{},
				q,
			),
		)
		tt.Assert.Error(err)
		p, ok := err.(*problem.P)
		if tt.Assert.True(ok) {
			tt.Assert.Equal(400, p.Status)
			tt.Assert.NotNil(p.Extras)
			tt.Assert.Equal(p.Extras["invalid_field"], "buying")
			tt.Assert.Equal(p.Extras["reason"], "Asset code length is invalid")
		}
	})
}

func TestGetAccountOffersHandler(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()

	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &history.Q{tt.AuroraSession()}
	handler := GetAccountOffersHandler{}

	err := q.UpsertOffers(tt.Ctx, []history.Offer{eurOffer, twoEurOffer, usdOffer})
	tt.Assert.NoError(err)

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{},
			map[string]string{"account_id": issuer.Address()},
			q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 2)

	offers := pageableToOffers(t, records)

	for _, offer := range offers {
		tt.Assert.Equal(issuer.Address(), offer.Seller)
	}

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{},
			map[string]string{},
			q,
		),
	)
	tt.Assert.Error(err)
}

func pageableToOffers(t *testing.T, page []hal.Pageable) []aurora.Offer {
	var offers []aurora.Offer
	for _, entry := range page {
		offers = append(offers, entry.(aurora.Offer))
	}
	return offers
}

func TestOffersQueryURLTemplate(t *testing.T) {
	tt := assert.New(t)
	expected := "/offers{?selling,buying,seller,sponsor,cursor,limit,order}"
	offersQuery := OffersQuery{}
	tt.Equal(expected, offersQuery.URITemplate())
}
