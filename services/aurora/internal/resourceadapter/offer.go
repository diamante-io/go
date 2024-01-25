package resourceadapter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/diamcircle/go/amount"
	protocol "github.com/diamcircle/go/protocols/aurora"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/xdr"
)

// PopulateOffer constructs an offer response struct from an offer row extracted from the
// the aurora offers table.
func PopulateOffer(ctx context.Context, dest *protocol.Offer, row history.Offer, ledger *history.Ledger) {
	dest.ID = int64(row.OfferID)
	dest.PT = fmt.Sprintf("%d", row.OfferID)
	dest.Seller = row.SellerID
	dest.Amount = amount.String(xdr.Int64(row.Amount))
	dest.PriceR.N = row.Pricen
	dest.PriceR.D = row.Priced
	dest.Price = big.NewRat(int64(row.Pricen), int64(row.Priced)).FloatString(7)
	if row.Sponsor.Valid {
		dest.Sponsor = row.Sponsor.String
	}

	row.SellingAsset.MustExtract(&dest.Selling.Type, &dest.Selling.Code, &dest.Selling.Issuer)
	row.BuyingAsset.MustExtract(&dest.Buying.Type, &dest.Buying.Code, &dest.Buying.Issuer)

	dest.LastModifiedLedger = int32(row.LastModifiedLedger)
	if ledger != nil {
		dest.LastModifiedTime = &ledger.ClosedAt
	}
	lb := hal.LinkBuilder{auroraContext.BaseURL(ctx)}
	dest.Links.Self = lb.Linkf("/offers/%d", row.OfferID)
	dest.Links.OfferMaker = lb.Linkf("/accounts/%s", row.SellerID)
}
