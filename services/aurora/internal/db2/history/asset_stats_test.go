package history

import (
	"database/sql"
	"testing"

	"go/services/aurora/internal/db2"
	"go/services/aurora/internal/test"
	"go/xdr"
)

func TestInsertAssetStats(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}
	tt.Assert.NoError(q.InsertAssetStats(tt.Ctx, []ExpAssetStat{}, 1))

	assetStats := []ExpAssetStat{
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			AssetCode:   "USD",
			Accounts: ExpAssetStatAccounts{
				Authorized:                      2,
				AuthorizedToMaintainLiabilities: 3,
				Unauthorized:                    4,
			},
			Balances: ExpAssetStatBalances{
				Authorized:                      "1",
				AuthorizedToMaintainLiabilities: "2",
				Unauthorized:                    "3",
				ClaimableBalances:               "4",
				LiquidityPools:                  "5",
			},
			Amount:      "1",
			NumAccounts: 2,
		},
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
			AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			AssetCode:   "ETHER",
			Accounts: ExpAssetStatAccounts{
				Authorized:                      1,
				AuthorizedToMaintainLiabilities: 3,
				Unauthorized:                    4,
			},
			Balances: ExpAssetStatBalances{
				Authorized:                      "23",
				AuthorizedToMaintainLiabilities: "2",
				Unauthorized:                    "3",
				ClaimableBalances:               "4",
				LiquidityPools:                  "5",
			},
			Amount:      "23",
			NumAccounts: 1,
		},
	}
	tt.Assert.NoError(q.InsertAssetStats(tt.Ctx, assetStats, 1))

	for _, assetStat := range assetStats {
		got, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
		tt.Assert.NoError(err)
		tt.Assert.Equal(got, assetStat)
	}
}

func TestInsertAssetStat(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	assetStats := []ExpAssetStat{
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			AssetCode:   "USD",
			Accounts: ExpAssetStatAccounts{
				Authorized:                      2,
				AuthorizedToMaintainLiabilities: 3,
				Unauthorized:                    4,
			},
			Balances: ExpAssetStatBalances{
				Authorized:                      "1",
				AuthorizedToMaintainLiabilities: "2",
				Unauthorized:                    "3",
				ClaimableBalances:               "4",
				LiquidityPools:                  "5",
			},
			Amount:      "1",
			NumAccounts: 2,
		},
		{
			AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
			AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			AssetCode:   "ETHER",
			Accounts: ExpAssetStatAccounts{
				Authorized:                      1,
				AuthorizedToMaintainLiabilities: 3,
				Unauthorized:                    4,
			},
			Balances: ExpAssetStatBalances{
				Authorized:                      "23",
				AuthorizedToMaintainLiabilities: "2",
				Unauthorized:                    "3",
				ClaimableBalances:               "4",
				LiquidityPools:                  "5",
			},
			Amount:      "23",
			NumAccounts: 1,
		},
	}

	for _, assetStat := range assetStats {
		numChanged, err := q.InsertAssetStat(tt.Ctx, assetStat)
		tt.Assert.NoError(err)
		tt.Assert.Equal(numChanged, int64(1))

		got, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
		tt.Assert.NoError(err)
		tt.Assert.Equal(got, assetStat)
	}
}

func TestInsertAssetStatAlreadyExistsError(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	assetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "4",
			LiquidityPools:                  "5",
		},
		Amount:      "1",
		NumAccounts: 2,
	}

	numChanged, err := q.InsertAssetStat(tt.Ctx, assetStat)
	tt.Assert.NoError(err)
	tt.Assert.Equal(numChanged, int64(1))

	numChanged, err = q.InsertAssetStat(tt.Ctx, assetStat)
	tt.Assert.Error(err)
	tt.Assert.Equal(numChanged, int64(0))

	assetStat.NumAccounts = 4
	assetStat.Amount = "3"
	numChanged, err = q.InsertAssetStat(tt.Ctx, assetStat)
	tt.Assert.Error(err)
	tt.Assert.Equal(numChanged, int64(0))

	assetStat.NumAccounts = 2
	assetStat.Amount = "1"
	got, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.NoError(err)
	tt.Assert.Equal(got, assetStat)
}

func TestUpdateAssetStatDoesNotExistsError(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	assetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "4",
			LiquidityPools:                  "5",
		},
		Amount:      "1",
		NumAccounts: 2,
	}

	numChanged, err := q.UpdateAssetStat(tt.Ctx, assetStat)
	tt.Assert.Nil(err)
	tt.Assert.Equal(numChanged, int64(0))

	_, err = q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.Equal(err, sql.ErrNoRows)
}

func TestUpdateStat(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	assetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "4",
			LiquidityPools:                  "5",
		},
		Amount:      "1",
		NumAccounts: 2,
	}

	numChanged, err := q.InsertAssetStat(tt.Ctx, assetStat)
	tt.Assert.NoError(err)
	tt.Assert.Equal(numChanged, int64(1))

	got, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.NoError(err)
	tt.Assert.Equal(got, assetStat)

	assetStat.NumAccounts = 50
	assetStat.Amount = "23"

	numChanged, err = q.UpdateAssetStat(tt.Ctx, assetStat)
	tt.Assert.Nil(err)
	tt.Assert.Equal(numChanged, int64(1))

	got, err = q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.NoError(err)
	tt.Assert.Equal(got, assetStat)
}

func TestGetAssetStatDoesNotExist(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	assetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "4",
			LiquidityPools:                  "5",
		},
		Amount:      "1",
		NumAccounts: 2,
	}

	_, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.Equal(err, sql.ErrNoRows)
}

func TestRemoveAssetStat(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	assetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "4",
			LiquidityPools:                  "5",
		},
		Amount:      "1",
		NumAccounts: 2,
	}

	numChanged, err := q.RemoveAssetStat(tt.Ctx,
		assetStat.AssetType,
		assetStat.AssetCode,
		assetStat.AssetIssuer,
	)
	tt.Assert.Nil(err)
	tt.Assert.Equal(numChanged, int64(0))

	numChanged, err = q.InsertAssetStat(tt.Ctx, assetStat)
	tt.Assert.NoError(err)
	tt.Assert.Equal(numChanged, int64(1))

	got, err := q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.NoError(err)
	tt.Assert.Equal(got, assetStat)

	numChanged, err = q.RemoveAssetStat(tt.Ctx,
		assetStat.AssetType,
		assetStat.AssetCode,
		assetStat.AssetIssuer,
	)
	tt.Assert.Nil(err)
	tt.Assert.Equal(numChanged, int64(1))

	_, err = q.GetAssetStat(tt.Ctx, assetStat.AssetType, assetStat.AssetCode, assetStat.AssetIssuer)
	tt.Assert.Equal(err, sql.ErrNoRows)
}

func TestGetAssetStatsCursorValidation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	for _, testCase := range []struct {
		name          string
		cursor        string
		expectedError string
	}{
		{
			"cursor does not use underscore as serpator",
			"usdc-GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			"invalid asset stats cursor",
		},
		{
			"cursor has no underscore",
			"usdc",
			"invalid asset stats cursor",
		},
		{
			"cursor has too many underscores",
			"usdc_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum4_",
			"invalid asset type in asset stats cursor",
		},
		{
			"issuer in cursor is invalid",
			"usd_abcdefghijklmnopqrstuv_credit_alphanum4",
			"invalid issuer in asset stats cursor",
		},
		{
			"asset type in cursor is invalid",
			"usd_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum",
			"invalid asset type in asset stats cursor",
		},
		{
			"asset code in cursor is too long",
			"abcdefghijklmnopqrstuv_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum12",
			"invalid asset stats cursor",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			page := db2.PageQuery{
				Cursor: testCase.cursor,
				Order:  "asc",
				Limit:  5,
			}
			results, err := q.GetAssetStats(tt.Ctx, "", "", page)
			tt.Assert.Empty(results)
			tt.Assert.NotNil(err)
			tt.Assert.Contains(err.Error(), testCase.expectedError)
		})
	}
}

func TestGetAssetStatsOrderValidation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	page := db2.PageQuery{
		Order: "invalid",
		Limit: 5,
	}
	results, err := q.GetAssetStats(tt.Ctx, "", "", page)
	tt.Assert.Empty(results)
	tt.Assert.NotNil(err)
	tt.Assert.Contains(err.Error(), "invalid page order")
}

func reverseAssetStats(a []ExpAssetStat) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func TestGetAssetStatsFiltersAndCursor(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	q := &Q{tt.auroraSession()}

	usdAssetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "1",
		NumAccounts: 2,
	}
	etherAssetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
		AssetIssuer: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		AssetCode:   "ETHER",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      1,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "23",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "23",
		NumAccounts: 1,
	}
	otherUSDAssetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
		AssetCode:   "USD",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      2,
			AuthorizedToMaintainLiabilities: 3,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "1",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "0",
			LiquidityPools:                  "0",
		},
		Amount:      "1",
		NumAccounts: 2,
	}
	eurAssetStat := ExpAssetStat{
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum4,
		AssetIssuer: "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
		AssetCode:   "EUR",
		Accounts: ExpAssetStatAccounts{
			Authorized:                      3,
			AuthorizedToMaintainLiabilities: 2,
			Unauthorized:                    4,
		},
		Balances: ExpAssetStatBalances{
			Authorized:                      "111",
			AuthorizedToMaintainLiabilities: "2",
			Unauthorized:                    "3",
			ClaimableBalances:               "1",
			LiquidityPools:                  "2",
		},
		Amount:      "111",
		NumAccounts: 3,
	}
	assetStats := []ExpAssetStat{
		etherAssetStat,
		eurAssetStat,
		otherUSDAssetStat,
		usdAssetStat,
	}
	for _, assetStat := range assetStats {
		numChanged, err := q.InsertAssetStat(tt.Ctx, assetStat)
		tt.Assert.NoError(err)
		tt.Assert.Equal(numChanged, int64(1))
	}

	for _, testCase := range []struct {
		name        string
		assetCode   string
		assetIssuer string
		cursor      string
		order       string
		expected    []ExpAssetStat
	}{
		{
			"no filter without cursor",
			"",
			"",
			"",
			"asc",
			[]ExpAssetStat{
				etherAssetStat,
				eurAssetStat,
				otherUSDAssetStat,
				usdAssetStat,
			},
		},
		{
			"no filter with cursor",
			"",
			"",
			"ABC_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			[]ExpAssetStat{
				etherAssetStat,
				eurAssetStat,
				otherUSDAssetStat,
				usdAssetStat,
			},
		},
		{
			"no filter with cursor descending",
			"",
			"",
			"ZZZ_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				usdAssetStat,
				otherUSDAssetStat,
				eurAssetStat,
				etherAssetStat,
			},
		},
		{
			"no filter with cursor and offset",
			"",
			"",
			"ETHER_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum12",
			"asc",
			[]ExpAssetStat{
				eurAssetStat,
				otherUSDAssetStat,
				usdAssetStat,
			},
		},
		{
			"no filter with cursor and offset descending",
			"",
			"",
			"EUR_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				etherAssetStat,
			},
		},
		{
			"no filter with cursor and offset descending including eur",
			"",
			"",
			"EUR_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				eurAssetStat,
				etherAssetStat,
			},
		},
		{
			"filter on code without cursor",
			"USD",
			"",
			"",
			"asc",
			[]ExpAssetStat{
				otherUSDAssetStat,
				usdAssetStat,
			},
		},
		{
			"filter on code with cursor",
			"USD",
			"",
			"USD_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			[]ExpAssetStat{
				usdAssetStat,
			},
		},
		{
			"filter on code with cursor descending",
			"USD",
			"",
			"USD_GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				otherUSDAssetStat,
			},
		},
		{
			"filter on issuer without cursor",
			"",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"",
			"asc",
			[]ExpAssetStat{
				eurAssetStat,
				otherUSDAssetStat,
			},
		},
		{
			"filter on issuer with cursor",
			"",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"EUR_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			[]ExpAssetStat{
				otherUSDAssetStat,
			},
		},
		{
			"filter on issuer with cursor descending",
			"",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"USD_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				eurAssetStat,
			},
		},
		{
			"filter on non existant code without cursor",
			"BTC",
			"",
			"",
			"asc",
			nil,
		},
		{
			"filter on non existant code with cursor",
			"BTC",
			"",
			"BTC_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			nil,
		},
		{
			"filter on non existant issuer without cursor",
			"",
			"GAEIHD6U4WSBHJGA2HPWOQ3OQEFQ3Y7QZE2DR76YKZNKPW5YDLYW4UGF",
			"",
			"asc",
			nil,
		},
		{
			"filter on non existant issuer with cursor",
			"",
			"GAEIHD6U4WSBHJGA2HPWOQ3OQEFQ3Y7QZE2DR76YKZNKPW5YDLYW4UGF",
			"AAA_GAEIHD6U4WSBHJGA2HPWOQ3OQEFQ3Y7QZE2DR76YKZNKPW5YDLYW4UGF_credit_alphanum4",
			"asc",
			nil,
		},
		{
			"filter on non existant code and non existant issuer without cursor",
			"BTC",
			"GAEIHD6U4WSBHJGA2HPWOQ3OQEFQ3Y7QZE2DR76YKZNKPW5YDLYW4UGF",
			"",
			"asc",
			nil,
		},
		{
			"filter on non existant code and non existant issuer with cursor",
			"BTC",
			"GAEIHD6U4WSBHJGA2HPWOQ3OQEFQ3Y7QZE2DR76YKZNKPW5YDLYW4UGF",
			"AAA_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			nil,
		},
		{
			"filter on both code and issuer without cursor",
			"USD",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"",
			"asc",
			[]ExpAssetStat{
				otherUSDAssetStat,
			},
		},
		{
			"filter on both code and issuer with cursor",
			"USD",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"USC_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			[]ExpAssetStat{
				otherUSDAssetStat,
			},
		},
		{
			"filter on both code and issuer with cursor descending",
			"USD",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"USE_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"desc",
			[]ExpAssetStat{
				otherUSDAssetStat,
			},
		},
		{
			"cursor negates filter",
			"USD",
			"GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			"USD_GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2_credit_alphanum4",
			"asc",
			nil,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			page := db2.PageQuery{
				Order:  testCase.order,
				Cursor: testCase.cursor,
				Limit:  5,
			}
			results, err := q.GetAssetStats(tt.Ctx, testCase.assetCode, testCase.assetIssuer, page)
			tt.Assert.NoError(err)
			tt.Assert.Equal(testCase.expected, results)

			page.Limit = 1
			results, err = q.GetAssetStats(tt.Ctx, testCase.assetCode, testCase.assetIssuer, page)
			tt.Assert.NoError(err)
			if len(testCase.expected) == 0 {
				tt.Assert.Equal(testCase.expected, results)
			} else {
				tt.Assert.Equal(testCase.expected[:1], results)
			}

			if page.Cursor == "" {
				page = page.Invert()
				page.Limit = 5

				results, err = q.GetAssetStats(tt.Ctx, testCase.assetCode, testCase.assetIssuer, page)
				tt.Assert.NoError(err)
				reverseAssetStats(results)
				tt.Assert.Equal(testCase.expected, results)
			}
		})
	}
}
