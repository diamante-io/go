package resourceadapter

import (
	"context"
	"testing"

	"go/protocols/aurora"
	protocol "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestPopulateExpAssetStat(t *testing.T) {
	row := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   "XIM",
		AssetIssuer: "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      429,
			AuthorizedToMaintainLiabilities: 214,
			Unauthorized:                    107,
			ClaimableBalances:               12,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "100000000000000000000",
			AuthorizedToMaintainLiabilities: "50000000000000000000",
			Unauthorized:                    "2500000000000000000",
			ClaimableBalances:               "1200000000000000000",
			LiquidityPools:                  "7700000000000000000",
		},
		Amount:      "100000000000000000000", // 10T
		NumAccounts: 429,
	}
	issuer := history.AccountEntry{
		AccountID:  "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM",
		Flags:      0,
		HomeDomain: "xim.com",
	}

	var res protocol.AssetStat
	err := PopulateAssetStat(context.Background(), &res, row, issuer)
	assert.NoError(t, err)

	assert.Equal(t, "credit_alphanum4", res.Type)
	assert.Equal(t, "XIM", res.Code)
	assert.Equal(t, "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM", res.Issuer)
	assert.Equal(t, int32(429), res.Accounts.Authorized)
	assert.Equal(t, int32(214), res.Accounts.AuthorizedToMaintainLiabilities)
	assert.Equal(t, int32(107), res.Accounts.Unauthorized)
	assert.Equal(t, int32(12), res.NumClaimableBalances)
	assert.Equal(t, "10000000000000.0000000", res.Balances.Authorized)
	assert.Equal(t, "5000000000000.0000000", res.Balances.AuthorizedToMaintainLiabilities)
	assert.Equal(t, "250000000000.0000000", res.Balances.Unauthorized)
	assert.Equal(t, "120000000000.0000000", res.ClaimableBalancesAmount)
	assert.Equal(t, "770000000000.0000000", res.LiquidityPoolsAmount)
	assert.Equal(t, "10000000000000.0000000", res.Amount)
	assert.Equal(t, int32(429), res.NumAccounts)
	assert.Equal(t, aurora.AccountFlags{}, res.Flags)
	assert.Equal(t, "https://xim.com/.well-known/diamcircle.toml", res.Links.Toml.Href)
	assert.Equal(t, row.PagingToken(), res.PagingToken())

	issuer.HomeDomain = ""
	issuer.Flags = uint32(xdr.AccountFlagsAuthRequiredFlag) |
		uint32(xdr.AccountFlagsAuthImmutableFlag) |
		uint32(xdr.AccountFlagsAuthClawbackEnabledFlag)

	err = PopulateAssetStat(context.Background(), &res, row, issuer)
	assert.NoError(t, err)

	assert.Equal(t, "credit_alphanum4", res.Type)
	assert.Equal(t, "XIM", res.Code)
	assert.Equal(t, "GBZ35ZJRIKJGYH5PBKLKOZ5L6EXCNTO7BKIL7DAVVDFQ2ODJEEHHJXIM", res.Issuer)
	assert.Equal(t, "10000000000000.0000000", res.Amount)
	assert.Equal(t, int32(429), res.NumAccounts)
	assert.Equal(
		t,
		aurora.AccountFlags{
			AuthRequired:        true,
			AuthImmutable:       true,
			AuthClawbackEnabled: true,
		},
		res.Flags,
	)
	assert.Equal(t, "", res.Links.Toml.Href)
	assert.Equal(t, row.PagingToken(), res.PagingToken())
}
