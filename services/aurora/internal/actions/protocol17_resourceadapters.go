package actions

import (
	"context"
	"fmt"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/price"
	protocol "github.com/diamcircle/go/protocols/aurora"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ingest/processors"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/xdr"
)

// Populate fills out the details of a trade using a row from the history_trades
// table. This can be removed in release after P18 upgrade.
func PopulateTradeP17(
	ctx context.Context,
	dest *TradeP17,
	row history.Trade,
) {
	dest.ID = row.PagingToken()
	dest.PT = row.PagingToken()
	dest.BaseOfferID = ""
	if row.BaseOfferID.Valid {
		dest.BaseOfferID = fmt.Sprintf("%d", row.BaseOfferID.Int64)
	}
	dest.BaseAccount = row.BaseAccount.String // Always Valid in P17
	dest.BaseAssetType = row.BaseAssetType
	dest.BaseAssetCode = row.BaseAssetCode
	dest.BaseAssetIssuer = row.BaseAssetIssuer
	dest.BaseAmount = amount.String(xdr.Int64(row.BaseAmount))
	dest.CounterOfferID = ""
	if row.CounterOfferID.Valid {
		dest.CounterOfferID = fmt.Sprintf("%d", row.CounterOfferID.Int64)
	}
	dest.CounterAccount = row.CounterAccount.String // Always Valid in P17
	dest.CounterAssetType = row.CounterAssetType
	dest.CounterAssetCode = row.CounterAssetCode
	dest.CounterAssetIssuer = row.CounterAssetIssuer
	dest.CounterAmount = amount.String(xdr.Int64(row.CounterAmount))
	dest.LedgerCloseTime = row.LedgerCloseTime
	dest.BaseIsSeller = row.BaseIsSeller

	_, counterOfferIDType := processors.DecodeOfferID(row.CounterOfferID.Int64)
	_, baseOfferIDType := processors.DecodeOfferID(row.BaseOfferID.Int64)
	if counterOfferIDType == processors.CoreOfferIDType && baseOfferIDType == processors.CoreOfferIDType {
		// Both Core Offer ID
		// Confirm diff: https://aurora.diamcircle.org/trades?base_asset_code=USD&base_asset_issuer=GAEDZ7BHMDYEMU6IJT3CTTGDUSLZWS5CQWZHGP4XUOIDG5ISH3AFAEK2&base_asset_type=credit_alphanum4&counter_asset_type=native&cursor=142354492003246081-3&limit=200&order=asc
		if row.BaseIsSeller {
			dest.OfferID = fmt.Sprintf("%d", row.BaseOfferID.Int64)
		} else {
			dest.OfferID = fmt.Sprintf("%d", row.CounterOfferID.Int64)
		}
	} else if counterOfferIDType == processors.CoreOfferIDType {
		dest.OfferID = fmt.Sprintf("%d", row.CounterOfferID.Int64)
	} else {
		dest.OfferID = fmt.Sprintf("%d", row.BaseOfferID.Int64)
	}

	if row.HasPrice() {
		dest.Price = &protocol.Price{
			N: int32(row.PriceN.Int64),
			D: int32(row.PriceD.Int64),
		}
	}

	populateTradeLinks(ctx, dest, row.HistoryOperationID)
}

func populateTradeLinks(
	ctx context.Context,
	dest *TradeP17,
	opid int64,
) {
	lb := hal.LinkBuilder{auroraContext.BaseURL(ctx)}
	dest.Links.Base = lb.Link("/accounts", dest.BaseAccount)
	dest.Links.Counter = lb.Link("/accounts", dest.CounterAccount)
	dest.Links.Operation = lb.Link(
		"/operations",
		fmt.Sprintf("%d", opid),
	)
}

// PopulateTradeAggregationP17 fills out the details of a trade using a row from the history_trades
// table. This can be removed in release after P18 upgrade.
func PopulateTradeAggregationP17(
	ctx context.Context,
	dest *TradeAggregationP17,
	row history.TradeAggregation,
) error {
	var err error
	dest.Timestamp = row.Timestamp
	dest.TradeCount = row.TradeCount
	dest.BaseVolume, err = amount.IntStringToAmount(row.BaseVolume)
	if err != nil {
		return err
	}
	dest.CounterVolume, err = amount.IntStringToAmount(row.CounterVolume)
	if err != nil {
		return err
	}
	dest.Average = price.StringFromFloat64(row.Average)
	var (
		high = xdr.Price{
			N: xdr.Int32(row.HighN),
			D: xdr.Int32(row.HighD),
		}
		low = xdr.Price{
			N: xdr.Int32(row.LowN),
			D: xdr.Int32(row.LowD),
		}
		open = xdr.Price{
			N: xdr.Int32(row.OpenN),
			D: xdr.Int32(row.OpenD),
		}
		close = xdr.Price{
			N: xdr.Int32(row.CloseN),
			D: xdr.Int32(row.CloseD),
		}
	)
	dest.High = high.String()
	dest.HighR = high
	dest.Low = low.String()
	dest.LowR = low
	dest.Open = open.String()
	dest.OpenR = open
	dest.Close = close.String()
	dest.CloseR = close
	return nil
}
