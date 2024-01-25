package resourceadapter

import (
	"context"
	"strings"

	"go/amount"
	protocol "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/support/errors"
	"go/support/render/hal"
	"go/xdr"
)

// PopulateAssetStat populates an AssetStat using asset stats and account entries
// generated from the ingestion system.
func PopulateAssetStat(
	ctx context.Context,
	res *protocol.AssetStat,
	row history.ExpAssetStat,
	issuer history.AccountEntry,
) (err error) {
	res.Asset.Type = xdr.AssetTypeToString[row.AssetType]
	res.Asset.Code = row.AssetCode
	res.Asset.Issuer = row.AssetIssuer
	res.Accounts = protocol.AssetStatAccounts{
		Authorized:                      row.Accounts.Authorized,
		AuthorizedToMaintainLiabilities: row.Accounts.AuthorizedToMaintainLiabilities,
		Unauthorized:                    row.Accounts.Unauthorized,
	}
	res.NumClaimableBalances = row.Accounts.ClaimableBalances
	res.NumLiquidityPools = row.Accounts.LiquidityPools
	res.NumAccounts = row.NumAccounts
	err = populateAssetStatBalances(res, row.Balances)
	if err != nil {
		return err
	}
	flags := int8(issuer.Flags)
	res.Flags = protocol.AccountFlags{
		(flags & int8(xdr.AccountFlagsAuthRequiredFlag)) != 0,
		(flags & int8(xdr.AccountFlagsAuthRevocableFlag)) != 0,
		(flags & int8(xdr.AccountFlagsAuthImmutableFlag)) != 0,
		(flags & int8(xdr.AccountFlagsAuthClawbackEnabledFlag)) != 0,
	}
	res.PT = row.PagingToken()

	trimmed := strings.TrimSpace(issuer.HomeDomain)
	var toml string
	if trimmed != "" {
		toml = "https://" + issuer.HomeDomain + "/.well-known/diamcircle.toml"
	}
	res.Links.Toml = hal.NewLink(toml)
	return
}

func populateAssetStatBalances(res *protocol.AssetStat, row history.ExpAssetStatBalances) (err error) {
	res.Amount, err = amount.IntStringToAmount(row.Authorized)
	if err != nil {
		return errors.Wrap(err, "Invalid amount in PopulateAssetStat")
	}

	res.Balances.Authorized, err = amount.IntStringToAmount(row.Authorized)
	if err != nil {
		return errors.Wrapf(err, "Invalid amount in PopulateAssetStatBalances: %q", row.Authorized)
	}

	res.Balances.AuthorizedToMaintainLiabilities, err = amount.IntStringToAmount(row.AuthorizedToMaintainLiabilities)
	if err != nil {
		return errors.Wrapf(err, "Invalid amount in PopulateAssetStatBalances: %q", row.AuthorizedToMaintainLiabilities)
	}

	res.Balances.Unauthorized, err = amount.IntStringToAmount(row.Unauthorized)
	if err != nil {
		return errors.Wrapf(err, "Invalid amount in PopulateAssetStatBalances: %q", row.Unauthorized)
	}

	res.ClaimableBalancesAmount, err = amount.IntStringToAmount(row.ClaimableBalances)
	if err != nil {
		return errors.Wrapf(err, "Invalid amount in PopulateAssetStatBalances: %q", row.ClaimableBalances)
	}

	res.LiquidityPoolsAmount, err = amount.IntStringToAmount(row.LiquidityPools)
	if err != nil {
		return errors.Wrapf(err, "Invalid amount in PopulateAssetStatBalances: %q", row.LiquidityPools)
	}

	return nil
}
