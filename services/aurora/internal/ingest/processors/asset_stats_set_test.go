package processors

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"go/ingest"
	"go/services/aurora/internal/db2/history"
	"go/xdr"
)

func TestEmptyAssetStatSet(t *testing.T) {
	set := AssetStatSet{}
	if all := set.All(); len(all) != 0 {
		t.Fatalf("expected empty list but got %v", all)
	}

	_, ok := set.Remove(
		xdr.AssetTypeAssetTypeCreditAlphanum4,
		"USD",
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
	)
	if ok {
		t.Fatal("expected remove to return false")
	}
}

func assertAllEquals(t *testing.T, set AssetStatSet, expected []history.ExpAssetStat) {
	all := set.All()
	assert.Len(t, all, len(expected))
	sort.Slice(all, func(i, j int) bool {
		return all[i].AssetCode < all[j].AssetCode
	})
	for i, got := range all {
		assert.Equal(t, expected[i], got)
	}
}

func TestAddNativeClaimableBalance(t *testing.T) {
	set := AssetStatSet{}
	claimableBalance := xdr.ClaimableBalanceEntry{
		BalanceId: xdr.ClaimableBalanceId{},
		Claimants: nil,
		Asset:     xdr.MustNewNativeAsset(),
		Amount:    100,
	}
	assert.NoError(t, set.AddClaimableBalance(
		ingest.Change{
			Post: &xdr.LedgerEntry{
				Data: xdr.LedgerEntryData{
					ClaimableBalance: &claimableBalance,
				},
			},
		},
	))
	assert.Empty(t, set.All())
}

func trustlineChange(pre, post *xdr.TrustLineEntry) ingest.Change {
	c := ingest.Change{}
	if pre != nil {
		c.Pre = &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				TrustLine: pre,
			},
		}
	}
	if post != nil {
		c.Post = &xdr.LedgerEntry{
			Data: xdr.LedgerEntryData{
				TrustLine: post,
			},
		}
	}
	return c
}

func TestAddPoolShareTrustline(t *testing.T) {
	set := AssetStatSet{}
	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset: xdr.TrustLineAsset{
				Type:            xdr.AssetTypeAssetTypePoolShare,
				LiquidityPoolId: &xdr.PoolId{1, 2, 3},
			},
			Balance: 1,
			Flags:   xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
		},
		)),
	)
	assert.Empty(t, set.All())
}

func TestAddAndRemoveAssetStats(t *testing.T) {
	set := AssetStatSet{}
	eur := "EUR"
	eurAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   eur,
		AssetIssuer: trustLineIssuer.Address(),
		Accounts: history.ExpAssetStatAccounts{
			Authorized: 1,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "0",
			Unauthorized:                    "0",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "1",
		NumAccounts: 1,
	}

	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset(eur, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   1,
			Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
		},
		)),
	)
	assertAllEquals(t, set, []history.ExpAssetStat{eurAssetStat})

	eurAssetStat.Accounts.ClaimableBalances++
	eurAssetStat.Balances.ClaimableBalances = "23"
	eurAsset := xdr.MustNewCreditAsset(eur, trustLineIssuer.Address())
	assert.NoError(
		t,
		set.addDelta(
			eurAsset,
			delta{ClaimableBalances: 23},
			delta{ClaimableBalances: 1},
		),
	)

	assertAllEquals(t, set, []history.ExpAssetStat{eurAssetStat})

	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset(eur, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   24,
			Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag | xdr.TrustLineFlagsTrustlineClawbackEnabledFlag),
		})),
	)

	eurAssetStat.Balances.Authorized = "25"
	eurAssetStat.Amount = "25"
	eurAssetStat.Accounts.Authorized++
	eurAssetStat.NumAccounts++
	assertAllEquals(t, set, []history.ExpAssetStat{eurAssetStat})

	usd := "USD"
	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Asset:     xdr.MustNewCreditAsset(usd, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   10,
			Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
		})),
	)

	ether := "ETHER"
	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
			Asset:     xdr.MustNewCreditAsset(ether, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   3,
			Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
		})),
	)

	// AddTrustline an authorized_to_maintain_liabilities trust line
	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset(ether, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   4,
			Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag),
		})),
	)

	// AddTrustline an unauthorized trust line
	assert.NoError(
		t,
		set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
			AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
			Asset:     xdr.MustNewCreditAsset(ether, trustLineIssuer.Address()).ToTrustLineAsset(),
			Balance:   5,
		})),
	)
	expected := []history.ExpAssetStat{
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
			AssetCode:   ether,
			AssetIssuer: trustLineIssuer.Address(),
			Accounts: history.ExpAssetStatAccounts{
				Authorized:                      1,
				AuthorizedToMaintainLiabilities: 1,
				Unauthorized:                    1,
			},
			Balances: history.ExpAssetStatBalances{
				Authorized:                      "3",
				AuthorizedToMaintainLiabilities: "4",
				Unauthorized:                    "5",
				ClaimableBalances:               "0",
				LiquidityPools:                  "0",
			},
			Amount:      "3",
			NumAccounts: 1,
		},
		eurAssetStat,
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetCode:   usd,
			AssetIssuer: trustLineIssuer.Address(),
			Accounts: history.ExpAssetStatAccounts{
				Authorized: 1,
			},
			Balances: history.ExpAssetStatBalances{
				Authorized:                      "10",
				AuthorizedToMaintainLiabilities: "0",
				Unauthorized:                    "0",
				ClaimableBalances:               "0",
				LiquidityPools:                  "0",
			},
			Amount:      "10",
			NumAccounts: 1,
		},
	}
	assertAllEquals(t, set, expected)

	for i, assetStat := range expected {
		removed, ok := set.Remove(assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
		if !ok {
			t.Fatal("expected remove to return true")
		}
		if removed != assetStat {
			t.Fatalf("expected removed asset stat to be %v but got %v", assetStat, removed)
		}

		assertAllEquals(t, set, expected[i+1:])
	}
}

func TestOverflowAssetStatSet(t *testing.T) {
	set := AssetStatSet{}
	eur := "EUR"
	err := set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset(eur, trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   math.MaxInt64,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	all := set.All()
	if len(all) != 1 {
		t.Fatalf("expected list of 1 asset stat but got %v", all)
	}

	eurAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   eur,
		AssetIssuer: trustLineIssuer.Address(),
		Accounts: history.ExpAssetStatAccounts{
			Authorized: 1,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "9223372036854775807",
			AuthorizedToMaintainLiabilities: "0",
			Unauthorized:                    "0",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "9223372036854775807",
		NumAccounts: 1,
	}
	if all[0] != eurAssetStat {
		t.Fatalf("expected asset stat to be %v but got %v", eurAssetStat, all[0])
	}

	err = set.AddTrustline(trustlineChange(nil, &xdr.TrustLineEntry{
		AccountId: xdr.MustAddress("GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB"),
		Asset:     xdr.MustNewCreditAsset(eur, trustLineIssuer.Address()).ToTrustLineAsset(),
		Balance:   math.MaxInt64,
		Flags:     xdr.Uint32(xdr.TrustLineFlagsAuthorizedFlag),
	}))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	all = set.All()
	if len(all) != 1 {
		t.Fatalf("expected list of 1 asset stat but got %v", all)
	}

	eurAssetStat = history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetCode:   eur,
		AssetIssuer: trustLineIssuer.Address(),
		Accounts: history.ExpAssetStatAccounts{
			Authorized: 2,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "18446744073709551614",
			AuthorizedToMaintainLiabilities: "0",
			Unauthorized:                    "0",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "18446744073709551614",
		NumAccounts: 2,
	}
	if all[0] != eurAssetStat {
		t.Fatalf("expected asset stat to be %v but got %v", eurAssetStat, all[0])
	}
}
