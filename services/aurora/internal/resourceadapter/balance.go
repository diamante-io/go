package resourceadapter

import (
	"github.com/diamcircle/go/amount"
	protocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/assets"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

func PopulatePoolShareBalance(dest *protocol.Balance, row history.TrustLine) (err error) {
	if row.AssetType == xdr.AssetTypeAssetTypePoolShare {
		dest.Type = "liquidity_pool_shares"
	} else {
		dest.Type, err = assets.String(row.AssetType)
		if err != nil {
			return err
		}

		if dest.Type != "liquidity_pool_shares" {
			return PopulateAssetBalance(dest, row)
		}
	}

	dest.LiquidityPoolId = row.LiquidityPoolID
	dest.Balance = amount.StringFromInt64(row.Balance)
	dest.Limit = amount.StringFromInt64(row.Limit)
	dest.LastModifiedLedger = row.LastModifiedLedger
	fillAuthorizationFlags(dest, row)

	return
}

func PopulateAssetBalance(dest *protocol.Balance, row history.TrustLine) (err error) {
	dest.Type, err = assets.String(row.AssetType)
	if err != nil {
		return err
	}

	dest.Balance = amount.StringFromInt64(row.Balance)
	dest.BuyingLiabilities = amount.StringFromInt64(row.BuyingLiabilities)
	dest.SellingLiabilities = amount.StringFromInt64(row.SellingLiabilities)
	dest.Limit = amount.StringFromInt64(row.Limit)
	dest.Issuer = row.AssetIssuer
	dest.Code = row.AssetCode
	dest.LastModifiedLedger = row.LastModifiedLedger
	fillAuthorizationFlags(dest, row)
	if row.Sponsor.Valid {
		dest.Sponsor = row.Sponsor.String
	}

	return
}

func PopulateNativeBalance(dest *protocol.Balance, stroops, buyingLiabilities, sellingLiabilities xdr.Int64) (err error) {
	dest.Type, err = assets.String(xdr.AssetTypeAssetTypeNative)
	if err != nil {
		return errors.Wrap(err, "getting the string representation from the provided xdr asset type")
	}

	dest.Balance = amount.String(stroops)
	dest.BuyingLiabilities = amount.String(buyingLiabilities)
	dest.SellingLiabilities = amount.String(sellingLiabilities)
	dest.LastModifiedLedger = 0
	dest.Limit = ""
	dest.Issuer = ""
	dest.Code = ""
	dest.IsAuthorized = nil
	dest.IsAuthorizedToMaintainLiabilities = nil
	dest.IsClawbackEnabled = nil
	return
}

func fillAuthorizationFlags(dest *protocol.Balance, row history.TrustLine) {
	isAuthorized := row.IsAuthorized()
	dest.IsAuthorized = &isAuthorized

	// After CAP-18, isAuth => isAuthToMaintain, so the following code does this
	// in a backwards compatible manner.
	dest.IsAuthorizedToMaintainLiabilities = &isAuthorized
	isAuthorizedToMaintainLiabilities := row.IsAuthorizedToMaintainLiabilities()
	if isAuthorizedToMaintainLiabilities {
		dest.IsAuthorizedToMaintainLiabilities = &isAuthorizedToMaintainLiabilities
	}

	isClawbackEnabled := row.IsClawbackEnabled()
	if isClawbackEnabled {
		dest.IsClawbackEnabled = &isClawbackEnabled
	}
}
