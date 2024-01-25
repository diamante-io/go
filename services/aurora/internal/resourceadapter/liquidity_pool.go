package resourceadapter

import (
	"context"
	"fmt"

	"github.com/diamcircle/go/amount"
	protocol "github.com/diamcircle/go/protocols/aurora"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/xdr"
)

// PopulateLiquidityPool fills out the resource's fields
func PopulateLiquidityPool(
	ctx context.Context,
	dest *protocol.LiquidityPool,
	liquidityPool history.LiquidityPool,
	ledger *history.Ledger,
) error {
	dest.ID = liquidityPool.PoolID
	dest.FeeBP = liquidityPool.Fee
	typ, ok := xdr.LiquidityPoolTypeToString[liquidityPool.Type]
	if !ok {
		return errors.Errorf("unknown liquidity pool type: %d", liquidityPool.Type)
	}
	dest.Type = typ
	dest.TotalTrustlines = liquidityPool.TrustlineCount
	dest.TotalShares = amount.StringFromInt64(int64(liquidityPool.ShareCount))
	for _, reserve := range liquidityPool.AssetReserves {
		dest.Reserves = append(dest.Reserves, protocol.LiquidityPoolReserve{
			Asset:  reserve.Asset.StringCanonical(),
			Amount: amount.StringFromInt64(int64(reserve.Reserve)),
		})
	}

	dest.LastModifiedLedger = liquidityPool.LastModifiedLedger

	if ledger != nil {
		dest.LastModifiedTime = &ledger.ClosedAt
	}

	lb := hal.LinkBuilder{Base: auroraContext.BaseURL(ctx)}
	self := fmt.Sprintf("/liquidity_pools/%s", dest.ID)
	dest.Links.Self = lb.Link(self)
	dest.PT = dest.ID
	dest.Links.Transactions = lb.PagedLink(self, "transactions")
	dest.Links.Operations = lb.PagedLink(self, "operations")
	return nil
}
