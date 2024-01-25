package resourceadapter

import (
	"context"
	"fmt"

	"go/xdr"

	"go/amount"
	protocol "go/protocols/aurora"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/support/render/hal"
)

// Populate fills out the details of a trade using a row from the history_trades
// table.
func PopulateTrade(
	ctx context.Context,
	dest *protocol.Trade,
	row history.Trade,
) {
	dest.ID = row.PagingToken()
	dest.PT = row.PagingToken()
	dest.BaseOfferID = ""
	if row.BaseOfferID.Valid {
		dest.BaseOfferID = fmt.Sprintf("%d", row.BaseOfferID.Int64)
	}
	if row.BaseAccount.Valid {
		dest.BaseAccount = row.BaseAccount.String
	}
	if row.BaseLiquidityPoolID.Valid {
		dest.BaseLiquidityPoolID = row.BaseLiquidityPoolID.String
	}
	dest.BaseAssetType = row.BaseAssetType
	dest.BaseAssetCode = row.BaseAssetCode
	dest.BaseAssetIssuer = row.BaseAssetIssuer
	dest.BaseAmount = amount.String(xdr.Int64(row.BaseAmount))
	dest.CounterOfferID = ""
	if row.CounterOfferID.Valid {
		dest.CounterOfferID = fmt.Sprintf("%d", row.CounterOfferID.Int64)
	}
	if row.CounterAccount.Valid {
		dest.CounterAccount = row.CounterAccount.String
	}
	if row.CounterLiquidityPoolID.Valid {
		dest.CounterLiquidityPoolID = row.CounterLiquidityPoolID.String
	}
	dest.CounterAssetType = row.CounterAssetType
	dest.CounterAssetCode = row.CounterAssetCode
	dest.CounterAssetIssuer = row.CounterAssetIssuer
	dest.CounterAmount = amount.String(xdr.Int64(row.CounterAmount))
	dest.LedgerCloseTime = row.LedgerCloseTime
	dest.BaseIsSeller = row.BaseIsSeller

	if row.LiquidityPoolFee.Valid {
		dest.LiquidityPoolFeeBP = uint32(row.LiquidityPoolFee.Int64)
	}

	switch row.Type {
	case history.OrderbookTradeType:
		dest.TradeType = history.OrderbookTrades
	case history.LiquidityPoolTradeType:
		dest.TradeType = history.LiquidityPoolTrades
	}

	if row.HasPrice() {
		dest.Price = protocol.TradePrice{
			N: row.PriceN.Int64,
			D: row.PriceD.Int64,
		}
	}

	populateTradeLinks(ctx, dest, row.HistoryOperationID)
}

func populateTradeLinks(
	ctx context.Context,
	dest *protocol.Trade,
	opid int64,
) {
	lb := hal.LinkBuilder{auroraContext.BaseURL(ctx)}
	switch {
	case dest.BaseOfferID != "":
		dest.Links.Base = lb.Link("/accounts", dest.BaseAccount)
	case dest.BaseLiquidityPoolID != "":
		dest.Links.Base = lb.Link("/liquidity_pools", dest.BaseLiquidityPoolID)
	}
	switch {
	case dest.CounterOfferID != "":
		dest.Links.Counter = lb.Link("/accounts", dest.CounterAccount)
	case dest.CounterLiquidityPoolID != "":
		dest.Links.Counter = lb.Link("/liquidity_pools", dest.CounterLiquidityPoolID)
	}
	dest.Links.Operation = lb.Link(
		"/operations",
		fmt.Sprintf("%d", opid),
	)
}
