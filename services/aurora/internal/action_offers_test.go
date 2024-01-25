package aurora

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest"
	"go/xdr"
)

func TestOfferActions_Show(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	q := &history.Q{ht.auroraSession()}
	ctx := context.Background()

	err := q.UpdateLastLedgerIngest(ctx, 100)
	ht.Assert.NoError(err)
	err = q.UpdateIngestVersion(ctx, ingest.CurrentVersion)
	ht.Assert.NoError(err)

	ledgerCloseTime := time.Now().Unix()
	_, err = q.InsertLedger(ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 100,
			ScpValue: xdr.diamcircleValue{
				CloseTime: xdr.TimePoint(ledgerCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	ht.Assert.NoError(err)

	issuer := xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	nativeAsset := xdr.MustNewNativeAsset()
	usdAsset := xdr.MustNewCreditAsset("USD", issuer.Address())
	eurAsset := xdr.MustNewCreditAsset("EUR", issuer.Address())

	eurOffer := history.Offer{
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
	usdOffer := history.Offer{
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

	err = q.UpsertOffers(ctx, []history.Offer{eurOffer, usdOffer})
	ht.Assert.NoError(err)

	w := ht.Get("/offers")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(2, w.Body)
	}

	w = ht.Get("/offers/4")
	if ht.Assert.Equal(200, w.Code) {
		var response aurora.Offer
		err = json.Unmarshal(w.Body.Bytes(), &response)
		ht.Assert.NoError(err)
		ht.Assert.Equal(int64(4), response.ID)
	}
}
