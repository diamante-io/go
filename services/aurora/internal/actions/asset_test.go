package actions

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go/protocols/aurora"
	"go/protocols/aurora/base"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/test"
	"go/support/render/hal"
	"go/support/render/problem"
	"go/xdr"
)

func TestAssetStatsValidation(t *testing.T) {
	handler := AssetStatsHandler{}

	for _, testCase := range []struct {
		name               string
		queryParams        map[string]string
		expectedErrorField string
		expectedError      string
	}{
		{
			"invalid asset code",
			map[string]string{
				"asset_code": "tooooooooolong",
			},
			"asset_code",
			"not a valid asset code",
		},
		{
			"invalid asset issuer",
			map[string]string{
				"asset_issuer": "invalid",
			},
			"asset_issuer",
			"not a valid asset issuer",
		},
		{
			"cursor has too many underscores",
			map[string]string{
				"cursor": "ABC_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum4_",
			},
			"cursor",
			"credit_alphanum4_ is not a valid asset type",
		},
		{
			"invalid cursor code",
			map[string]string{
				"cursor": "tooooooooolong_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum12",
			},
			"cursor",
			"not a valid asset code",
		},
		{
			"invalid cursor issuer",
			map[string]string{
				"cursor": "ABC_invalidissuer_credit_alphanum4",
			},
			"cursor",
			"not a valid asset issuer",
		},
		{
			"invalid cursor type",
			map[string]string{
				"cursor": "ABC_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum123",
			},
			"cursor",
			"credit_alphanum123 is not a valid asset type",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			r := makeRequest(t, testCase.queryParams, map[string]string{}, nil)
			_, err := handler.GetResourcePage(httptest.NewRecorder(), r)
			if err == nil {
				t.Fatalf("expected error %v but got %v", testCase.expectedError, err)
			}

			problem := err.(*problem.P)
			if field := problem.Extras["invalid_field"]; field != testCase.expectedErrorField {
				t.Fatalf(
					"expected error field %v but got %v",
					testCase.expectedErrorField,
					field,
				)
			}

			reason := problem.Extras["reason"]
			if !strings.Contains(reason.(string), testCase.expectedError) {
				t.Fatalf("expected reason %v but got %v", testCase.expectedError, reason)
			}
		})
	}
}

func TestAssetStats(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}
	handler := AssetStatsHandler{}

	issuer := history.AccountEntry{
		AccountID: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		Flags: uint32(xdr.AccountFlagsAuthRequiredFlag) |
			uint32(xdr.AccountFlagsAuthImmutableFlag) |
			uint32(xdr.AccountFlagsAuthClawbackEnabledFlag),
	}
	issuerFlags := aurora.AccountFlags{
		AuthRequired:        true,
		AuthImmutable:       true,
		AuthClawbackEnabled: true,
	}
	otherIssuer := history.AccountEntry{
		AccountID:  "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
		HomeDomain: "xim.com",
	}

	usdAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: issuer.AccountID,
		AssetCode:   "USD",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
			ClaimableBalances:               1,
			LiquidityPools:                  5,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "10",
			LiquidityPools:                  "20",
		},
		Amount:      "1",
		NumAccounts: 2,
	}
	usdAssetStatResponse := aurora.AssetStat{
		Accounts: aurora.AssetStatAccounts{
			Authorized:                      usdAssetStat.Accounts.Authorized,
			AuthorizedToMaintainLiabilities: usdAssetStat.Accounts.AuthorizedToMaintainLiabilities,
			Unauthorized:                    usdAssetStat.Accounts.Unauthorized,
		},
		NumClaimableBalances: usdAssetStat.Accounts.ClaimableBalances,
		NumLiquidityPools:    usdAssetStat.Accounts.LiquidityPools,
		Balances: aurora.AssetStatBalances{
			Authorized:                      "0.0000001",
			AuthorizedToMaintainLiabilities: "0.0000002",
			Unauthorized:                    "0.0000003",
		},
		ClaimableBalancesAmount: "0.0000010",
		LiquidityPoolsAmount:    "0.0000020",
		Amount:                  "0.0000001",
		NumAccounts:             usdAssetStat.NumAccounts,
		Asset: base.Asset{
			Type:   "credit_alphanum4",
			Code:   usdAssetStat.AssetCode,
			Issuer: usdAssetStat.AssetIssuer,
		},
		PT:    usdAssetStat.PagingToken(),
		Flags: issuerFlags,
	}

	etherAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: issuer.AccountID,
		AssetCode:   "ETHER",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      1,
			AuthorizedToMaintainLiabilities: 2,
			Unauthorized:                    3,
			ClaimableBalances:               0,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "23",
			AuthorizedToMaintainLiabilities: "46",
			Unauthorized:                    "92",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "23",
		NumAccounts: 1,
	}
	etherAssetStatResponse := aurora.AssetStat{
		Accounts: aurora.AssetStatAccounts{
			Authorized:                      etherAssetStat.Accounts.Authorized,
			AuthorizedToMaintainLiabilities: etherAssetStat.Accounts.AuthorizedToMaintainLiabilities,
			Unauthorized:                    etherAssetStat.Accounts.Unauthorized,
		},
		NumClaimableBalances: etherAssetStat.Accounts.ClaimableBalances,
		Balances: aurora.AssetStatBalances{
			Authorized:                      "0.0000023",
			AuthorizedToMaintainLiabilities: "0.0000046",
			Unauthorized:                    "0.0000092",
		},
		ClaimableBalancesAmount: "0.0000000",
		LiquidityPoolsAmount:    "0.0000000",
		Amount:                  "0.0000023",
		NumAccounts:             etherAssetStat.NumAccounts,
		Asset: base.Asset{
			Type:   "credit_alphanum4",
			Code:   etherAssetStat.AssetCode,
			Issuer: etherAssetStat.AssetIssuer,
		},
		PT:    etherAssetStat.PagingToken(),
		Flags: issuerFlags,
	}

	otherUSDAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: otherIssuer.AccountID,
		AssetCode:   "USD",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
			ClaimableBalances:               0,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "1",
		NumAccounts: 2,
	}
	otherUSDAssetStatResponse := aurora.AssetStat{
		Accounts: aurora.AssetStatAccounts{
			Authorized:                      otherUSDAssetStat.Accounts.Authorized,
			AuthorizedToMaintainLiabilities: otherUSDAssetStat.Accounts.AuthorizedToMaintainLiabilities,
			Unauthorized:                    otherUSDAssetStat.Accounts.Unauthorized,
		},
		NumClaimableBalances: otherUSDAssetStat.Accounts.ClaimableBalances,
		Balances: aurora.AssetStatBalances{
			Authorized:                      "0.0000001",
			AuthorizedToMaintainLiabilities: "0.0000002",
			Unauthorized:                    "0.0000003",
		},
		ClaimableBalancesAmount: "0.0000000",
		LiquidityPoolsAmount:    "0.0000000",
		Amount:                  "0.0000001",
		NumAccounts:             otherUSDAssetStat.NumAccounts,
		Asset: base.Asset{
			Type:   "credit_alphanum4",
			Code:   otherUSDAssetStat.AssetCode,
			Issuer: otherUSDAssetStat.AssetIssuer,
		},
		PT: otherUSDAssetStat.PagingToken(),
	}
	otherUSDAssetStatResponse.Links.Toml = hal.NewLink(
		"https://" + otherIssuer.HomeDomain + "/.well-known/diamcircle.toml",
	)

	eurAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: otherIssuer.AccountID,
		AssetCode:   "EUR",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      3,
			AuthorizedToMaintainLiabilities: 4,
			Unauthorized:                    5,
			ClaimableBalances:               0,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "111",
			AuthorizedToMaintainLiabilities: "222",
			Unauthorized:                    "333",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "111",
		NumAccounts: 3,
	}
	eurAssetStatResponse := aurora.AssetStat{
		Accounts: aurora.AssetStatAccounts{
			Authorized:                      eurAssetStat.Accounts.Authorized,
			AuthorizedToMaintainLiabilities: eurAssetStat.Accounts.AuthorizedToMaintainLiabilities,
			Unauthorized:                    eurAssetStat.Accounts.Unauthorized,
		},
		NumClaimableBalances: eurAssetStat.Accounts.ClaimableBalances,
		Balances: aurora.AssetStatBalances{
			Authorized:                      "0.0000111",
			AuthorizedToMaintainLiabilities: "0.0000222",
			Unauthorized:                    "0.0000333",
		},
		ClaimableBalancesAmount: "0.0000000",
		LiquidityPoolsAmount:    "0.0000000",
		Amount:                  "0.0000111",
		NumAccounts:             eurAssetStat.NumAccounts,
		Asset: base.Asset{
			Type:   "credit_alphanum4",
			Code:   eurAssetStat.AssetCode,
			Issuer: eurAssetStat.AssetIssuer,
		},
		PT: eurAssetStat.PagingToken(),
	}
	eurAssetStatResponse.Links.Toml = hal.NewLink(
		"https://" + otherIssuer.HomeDomain + "/.well-known/diamcircle.toml",
	)

	for _, assetStat := range []history.ExpAssetStat{
		etherAssetStat,
		eurAssetStat,
		otherUSDAssetStat,
		usdAssetStat,
	} {
		numChanged, err := q.InsertAssetStat(tt.Ctx, assetStat)
		tt.Assert.NoError(err)
		tt.Assert.Equal(numChanged, int64(1))
	}

	for _, account := range []history.AccountEntry{
		issuer,
		otherIssuer,
	} {
		accountEntry := history.AccountEntry{
			LastModifiedLedger: 100,
			AccountID:          account.AccountID,
			Flags:              account.Flags,
			HomeDomain:         account.HomeDomain,
		}

		err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{accountEntry})
		tt.Assert.NoError(err)
	}

	for _, testCase := range []struct {
		name        string
		queryParams map[string]string
		expected    []aurora.AssetStat
	}{
		{
			"default parameters",
			map[string]string{},
			[]aurora.AssetStat{
				etherAssetStatResponse,
				eurAssetStatResponse,
				otherUSDAssetStatResponse,
				usdAssetStatResponse,
			},
		},
		{
			"with cursor",
			map[string]string{
				"cursor": etherAssetStatResponse.PagingToken(),
			},
			[]aurora.AssetStat{
				eurAssetStatResponse,
				otherUSDAssetStatResponse,
				usdAssetStatResponse,
			},
		},
		{
			"descending order",
			map[string]string{"order": "desc"},
			[]aurora.AssetStat{
				usdAssetStatResponse,
				otherUSDAssetStatResponse,
				eurAssetStatResponse,
				etherAssetStatResponse,
			},
		},
		{
			"filter by asset code",
			map[string]string{
				"asset_code": "USD",
			},
			[]aurora.AssetStat{
				otherUSDAssetStatResponse,
				usdAssetStatResponse,
			},
		},
		{
			"filter by asset issuer",
			map[string]string{
				"asset_issuer": issuer.AccountID,
			},
			[]aurora.AssetStat{
				etherAssetStatResponse,
				usdAssetStatResponse,
			},
		},
		{
			"filter by both asset code and asset issuer",
			map[string]string{
				"asset_code":   "USD",
				"asset_issuer": issuer.AccountID,
			},
			[]aurora.AssetStat{
				usdAssetStatResponse,
			},
		},
		{
			"filter produces empty set",
			map[string]string{
				"asset_code":   "XYZ",
				"asset_issuer": issuer.AccountID,
			},
			[]aurora.AssetStat{},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			r := makeRequest(t, testCase.queryParams, map[string]string{}, q)
			results, err := handler.GetResourcePage(httptest.NewRecorder(), r)
			assert.NoError(t, err)

			assert.Len(t, results, len(testCase.expected))
			for i, item := range results {
				assetStat := item.(aurora.AssetStat)
				assert.Equal(t, testCase.expected[i], assetStat)
			}
		})
	}
}

func TestAssetStatsIssuerDoesNotExist(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}
	handler := AssetStatsHandler{}

	usdAssetStat := history.ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: history.ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
			ClaimableBalances:               0,
		},
		Balances: history.ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "0",
		},
		Amount:      "1",
		NumAccounts: 2,
	}
	numChanged, err := q.InsertAssetStat(tt.Ctx, usdAssetStat)
	tt.Assert.NoError(err)
	tt.Assert.Equal(numChanged, int64(1))

	r := makeRequest(t, map[string]string{}, map[string]string{}, q)
	results, err := handler.GetResourcePage(httptest.NewRecorder(), r)
	tt.Assert.NoError(err)

	expectedAssetStatResponse := aurora.AssetStat{
		Accounts: aurora.AssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		NumClaimableBalances: 0,
		Balances: aurora.AssetStatBalances{
			Authorized:                      "0.0000001",
			AuthorizedToMaintainLiabilities: "0.0000002",
			Unauthorized:                    "0.0000003",
		},
		ClaimableBalancesAmount: "0.0000000",
		LiquidityPoolsAmount:    "0.0000000",
		Amount:                  "0.0000001",
		NumAccounts:             usdAssetStat.NumAccounts,
		Asset: base.Asset{
			Type:   "credit_alphanum4",
			Code:   usdAssetStat.AssetCode,
			Issuer: usdAssetStat.AssetIssuer,
		},
		PT: usdAssetStat.PagingToken(),
	}

	tt.Assert.Len(results, 1)
	assetStat := results[0].(aurora.AssetStat)
	tt.Assert.Equal(assetStat, expectedAssetStatResponse)
}
