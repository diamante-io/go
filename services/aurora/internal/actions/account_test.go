package actions

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/guregu/null"

	"github.com/stretchr/testify/assert"

	protocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/protocols/aurora/base"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/xdr"
)

var (
	trustLineIssuer = "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"
	accountOne      = "GABGMPEKKDWR2WFH5AJOZV5PDKLJEHGCR3Q24ALETWR5H3A7GI3YTS7V"
	accountTwo      = "GADTXHUTHIAESMMQ2ZWSTIIGBZRLHUCBLCHPLLUEIAWDEFRDC4SYDKOZ"
	signer          = "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"
	sponsor         = "GCO26ZSBD63TKYX45H2C7D2WOFWOUSG5BMTNC3BG4QMXM3PAYI6WHKVZ"
	usd             = xdr.MustNewCreditAsset("USD", trustLineIssuer)
	euro            = xdr.MustNewCreditAsset("EUR", trustLineIssuer)

	account1 = history.AccountEntry{
		LastModifiedLedger: 1234,
		AccountID:          accountOne,
		Balance:            20000,
		SequenceNumber:     223456789,
		NumSubEntries:      10,
		Flags:              1,
		HomeDomain:         "diamcircle.org",
		MasterWeight:       1,
		ThresholdLow:       2,
		ThresholdMedium:    3,
		ThresholdHigh:      4,
		BuyingLiabilities:  3,
		SellingLiabilities: 3,
	}

	account2 = history.AccountEntry{
		LastModifiedLedger: 1234,
		AccountID:          accountTwo,
		Balance:            50000,
		SequenceNumber:     648736,
		NumSubEntries:      10,
		Flags:              2,
		HomeDomain:         "meridian.diamcircle.org",
		MasterWeight:       5,
		ThresholdLow:       6,
		ThresholdMedium:    7,
		ThresholdHigh:      8,
		BuyingLiabilities:  30,
		SellingLiabilities: 40,
	}

	account3 = history.AccountEntry{
		LastModifiedLedger: 1234,
		AccountID:          signer,
		Balance:            50000,
		SequenceNumber:     648736,
		NumSubEntries:      10,
		Flags:              2,
		MasterWeight:       5,
		ThresholdLow:       6,
		ThresholdMedium:    7,
		ThresholdHigh:      8,
		BuyingLiabilities:  30,
		SellingLiabilities: 40,
		NumSponsored:       1,
		NumSponsoring:      2,
		Sponsor:            null.StringFrom(sponsor),
	}

	eurTrustLine = history.TrustLine{
		AccountID:          accountOne,
		AssetType:          euro.Type,
		AssetIssuer:        trustLineIssuer,
		AssetCode:          "EUR",
		Balance:            20000,
		LedgerKey:          "eur-trustline1",
		Limit:              223456789,
		LiquidityPoolID:    "",
		BuyingLiabilities:  3,
		SellingLiabilities: 4,
		Flags:              1,
		LastModifiedLedger: 1234,
		Sponsor:            null.String{},
	}

	usdTrustLine = history.TrustLine{
		AccountID:          accountTwo,
		AssetType:          usd.Type,
		AssetIssuer:        trustLineIssuer,
		AssetCode:          "USD",
		Balance:            10000,
		LedgerKey:          "usd-trustline1",
		Limit:              123456789,
		LiquidityPoolID:    "",
		BuyingLiabilities:  1,
		SellingLiabilities: 2,
		Flags:              0,
		LastModifiedLedger: 1234,
		Sponsor:            null.String{},
	}

	data1 = history.Data{
		AccountID: accountOne,
		Name:      "test data",
		// This also tests if base64 encoding is working as 0 is invalid UTF-8 byte
		Value:              []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		LastModifiedLedger: 1234,
	}

	data2 = history.Data{
		AccountID:          accountTwo,
		Name:               "test data2",
		Value:              []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		LastModifiedLedger: 1234,
		Sponsor:            null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
	}

	accountSigners = []history.AccountSigner{
		{
			Account: accountOne,
			Signer:  accountOne,
			Weight:  1,
		},
		{
			Account: accountTwo,
			Signer:  accountTwo,
			Weight:  1,
		},
		{
			Account: accountOne,
			Signer:  signer,
			Weight:  1,
		},
		{
			Account: accountTwo,
			Signer:  signer,
			Weight:  2,
		},
		{
			Account: signer,
			Signer:  signer,
			Weight:  3,
		},
	}
)

func TestAccountInfo(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &history.Q{tt.AuroraSession()}

	var thresholds xdr.Thresholds
	tt.Assert.NoError(
		xdr.SafeUnmarshalBase64("AQAAAA==", &thresholds),
	)
	accountID := "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"
	accountEntry := history.AccountEntry{
		LastModifiedLedger:   4,
		AccountID:            accountID,
		Balance:              9999999900,
		SequenceNumber:       8589934593,
		NumSubEntries:        1,
		InflationDestination: "",
		HomeDomain:           "",
		MasterWeight:         thresholds[0],
		ThresholdLow:         thresholds[1],
		ThresholdMedium:      thresholds[2],
		ThresholdHigh:        thresholds[3],
		Flags:                0,
	}
	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{accountEntry})
	assert.NoError(t, err)

	tt.Assert.NoError(err)

	err = q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
		{
			AccountID:          accountID,
			AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetIssuer:        "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4",
			AssetCode:          "USD",
			Balance:            0,
			LedgerKey:          "test-usd-tl-1",
			Limit:              9223372036854775807,
			LiquidityPoolID:    "",
			BuyingLiabilities:  0,
			SellingLiabilities: 0,
			Flags:              1,
			LastModifiedLedger: 6,
			Sponsor:            null.String{},
		},
		{
			AccountID:          accountID,
			AssetType:          xdr.AssetTypeAssetTypeCreditAlphanum4,
			AssetIssuer:        "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4",
			AssetCode:          "EUR",
			Balance:            0,
			LedgerKey:          "test-eur-tl-1",
			Limit:              9223372036854775807,
			LiquidityPoolID:    "",
			BuyingLiabilities:  0,
			SellingLiabilities: 0,
			Flags:              1,
			LastModifiedLedger: 1234,
			Sponsor:            null.String{},
		},
	})
	assert.NoError(t, err)

	ledgerFourCloseTime := time.Now().Unix()
	_, err = q.InsertLedger(tt.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 4,
			ScpValue: xdr.DiamcircleValue{
				CloseTime: xdr.TimePoint(ledgerFourCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	assert.NoError(t, err)

	account, err := AccountInfo(tt.Ctx, &history.Q{tt.AuroraSession()}, accountID)
	tt.Assert.NoError(err)

	tt.Assert.Equal("8589934593", account.Sequence)
	tt.Assert.Equal(uint32(4), account.LastModifiedLedger)
	tt.Assert.NotNil(account.LastModifiedTime)
	tt.Assert.Equal(ledgerFourCloseTime, account.LastModifiedTime.Unix())
	tt.Assert.Len(account.Balances, 3)

	tt.Assert.Equal(account.Balances[0].Code, "EUR")
	tt.Assert.Equal(account.Balances[1].Code, "USD")
	tt.Assert.Equal(
		"GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4",
		account.Balances[1].Issuer,
	)
	tt.Assert.NotEqual(uint32(0), account.Balances[1].LastModifiedLedger)
	tt.Assert.Equal(account.Balances[2].Type, "native")
	tt.Assert.Equal(uint32(0), account.Balances[2].LastModifiedLedger)
	tt.Assert.Len(account.Signers, 1)

	// Regression: no trades link
	tt.Assert.Contains(account.Links.Trades.Href, "/trades")
	// Regression: no data link
	tt.Assert.Contains(account.Links.Data.Href, "/data/{key}")
	tt.Assert.True(account.Links.Data.Templated)

	// try to fetch account which does not exist
	_, err = AccountInfo(tt.Ctx, &history.Q{tt.AuroraSession()}, "GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ")
	tt.Assert.True(q.NoRows(errors.Cause(err)))
}

func TestGetAccountsHandlerPageNoResults(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := &GetAccountsHandler{}
	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{
				"signer": signer,
			},
			map[string]string{},
			q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 0)
}

func TestGetAccountsHandlerPageResultsBySigner(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := &GetAccountsHandler{}

	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{account1, account2, account3})
	assert.NoError(t, err)

	for _, row := range accountSigners {
		q.CreateAccountSigner(tt.Ctx, row.Account, row.Signer, row.Weight, nil)
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{
				"signer": signer,
			},
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(records))

	want := map[string]bool{
		accountOne: true,
		accountTwo: true,
		signer:     true,
	}

	for _, row := range records {
		result := row.(protocol.Account)
		tt.Assert.True(want[result.AccountID])
		delete(want, result.AccountID)
	}

	tt.Assert.Empty(want)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{
				"signer": signer,
				"cursor": accountOne,
			},
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(2, len(records))

	want = map[string]bool{
		accountTwo: true,
		signer:     true,
	}

	for _, row := range records {
		result := row.(protocol.Account)
		tt.Assert.True(want[result.AccountID])
		delete(want, result.AccountID)
	}

	tt.Assert.Empty(want)
}

func TestGetAccountsHandlerPageResultsBySponsor(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := &GetAccountsHandler{}

	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{account1, account2, account3})
	assert.NoError(t, err)

	for _, row := range accountSigners {
		q.CreateAccountSigner(tt.Ctx, row.Account, row.Signer, row.Weight, nil)
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			map[string]string{
				"sponsor": sponsor,
			},
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(records))
	tt.Assert.Equal(signer, records[0].(protocol.Account).ID)
}

func TestGetAccountsHandlerPageResultsByAsset(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := &GetAccountsHandler{}

	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{account1, account2})
	assert.NoError(t, err)
	ledgerCloseTime := time.Now().Unix()
	_, err = q.InsertLedger(tt.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 1234,
			ScpValue: xdr.DiamcircleValue{
				CloseTime: xdr.TimePoint(ledgerCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	assert.NoError(t, err)

	for _, row := range accountSigners {
		_, err = q.CreateAccountSigner(tt.Ctx, row.Account, row.Signer, row.Weight, nil)
		tt.Assert.NoError(err)
	}

	err = q.UpsertAccountData(tt.Ctx, []history.Data{data1, data2})
	assert.NoError(t, err)

	var assetType, code, issuer string
	usd.MustExtract(&assetType, &code, &issuer)
	params := map[string]string{
		"asset": code + ":" + issuer,
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			params,
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(0, len(records))

	err = q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
		eurTrustLine,
		usdTrustLine,
	})
	assert.NoError(t, err)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			params,
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(records))
	result := records[0].(protocol.Account)
	tt.Assert.Equal(accountTwo, result.AccountID)
	tt.Assert.NotNil(result.LastModifiedTime)
	tt.Assert.Equal(ledgerCloseTime, result.LastModifiedTime.Unix())
	tt.Assert.Len(result.Balances, 2)
	tt.Assert.Len(result.Signers, 2)

	_, ok := result.Data[data2.Name]
	tt.Assert.True(ok)
}

func createLP(tt *test.T, q *history.Q) history.LiquidityPool {
	lp := history.LiquidityPool{
		PoolID:         "cafebabedeadbeef000000000000000000000000000000000000000000000000",
		Type:           xdr.LiquidityPoolTypeLiquidityPoolConstantProduct,
		Fee:            34,
		TrustlineCount: 52115,
		ShareCount:     412241,
		AssetReserves: []history.LiquidityPoolAssetReserve{
			{
				Asset:   xdr.MustNewCreditAsset("USD", "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"),
				Reserve: 450,
			},
			{
				Asset:   xdr.MustNewNativeAsset(),
				Reserve: 450,
			},
		},
		LastModifiedLedger: 123,
	}

	err := q.UpsertLiquidityPools(tt.Ctx, []history.LiquidityPool{lp})
	tt.Assert.NoError(err)
	return lp
}

func TestGetAccountsHandlerPageResultsByLiquidityPool(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)

	q := &history.Q{tt.AuroraSession()}
	handler := &GetAccountsHandler{}

	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{account1, account2})
	assert.NoError(t, err)

	ledgerCloseTime := time.Now().Unix()
	_, err = q.InsertLedger(tt.Ctx, xdr.LedgerHeaderHistoryEntry{
		Header: xdr.LedgerHeader{
			LedgerSeq: 1234,
			ScpValue: xdr.DiamcircleValue{
				CloseTime: xdr.TimePoint(ledgerCloseTime),
			},
		},
	}, 0, 0, 0, 0, 0)
	assert.NoError(t, err)
	var assetType, code, issuer string
	usd.MustExtract(&assetType, &code, &issuer)
	params := map[string]string{
		"liquidity_pool": "cafebabedeadbeef000000000000000000000000000000000000000000000000",
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			params,
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(0, len(records))

	lp := createLP(tt, q)
	err = q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
		{
			AccountID:          account1.AccountID,
			AssetType:          xdr.AssetTypeAssetTypePoolShare,
			Balance:            10,
			LedgerKey:          "pool-share-1",
			Limit:              100,
			LiquidityPoolID:    lp.PoolID,
			Flags:              uint32(xdr.TrustLineFlagsAuthorizedFlag),
			LastModifiedLedger: lp.LastModifiedLedger,
		},
	})
	assert.NoError(t, err)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t,
			params,
			map[string]string{},
			q,
		),
	)

	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(records))
	result := records[0].(protocol.Account)
	tt.Assert.Equal(accountOne, result.AccountID)
	tt.Assert.NotNil(result.LastModifiedTime)
	tt.Assert.Equal(ledgerCloseTime, result.LastModifiedTime.Unix())
	tt.Assert.Len(result.Balances, 2)
	tt.Assert.True(*result.Balances[0].IsAuthorized)
	tt.Assert.True(*result.Balances[0].IsAuthorizedToMaintainLiabilities)
	result.Balances[0].IsAuthorized = nil
	result.Balances[0].IsAuthorizedToMaintainLiabilities = nil
	tt.Assert.Equal(
		protocol.Balance{
			Balance:                           "0.0000010",
			LiquidityPoolId:                   "cafebabedeadbeef000000000000000000000000000000000000000000000000",
			Limit:                             "0.0000100",
			BuyingLiabilities:                 "",
			SellingLiabilities:                "",
			Sponsor:                           "",
			LastModifiedLedger:                123,
			IsAuthorized:                      nil,
			IsAuthorizedToMaintainLiabilities: nil,
			IsClawbackEnabled:                 nil,
			Asset: base.Asset{
				Type: "liquidity_pool_shares",
			},
		},
		result.Balances[0],
	)
}

func TestGetAccountsHandlerInvalidParams(t *testing.T) {
	testCases := []struct {
		desc                    string
		params                  map[string]string
		expectedInvalidField    string
		expectedErr             string
		isInvalidAccountsParams bool
	}{
		{
			desc:                    "empty filters",
			isInvalidAccountsParams: true,
		},
		{
			desc: "signer and asset",
			params: map[string]string{
				"signer": accountOne,
				"asset":  "USD" + ":" + accountOne,
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "signer and liquidity pool",
			params: map[string]string{
				"signer":         accountOne,
				"liquidity_pool": "48672641c88264272787837f5c306f5ce93be3c2c7df68a092fbea55f5f4aa1d",
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "signer and sponsor",
			params: map[string]string{
				"signer":  accountOne,
				"sponsor": accountTwo,
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "asset and sponsor",
			params: map[string]string{
				"asset":   "USD" + ":" + accountOne,
				"sponsor": accountTwo,
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "asset and liquidity pool",
			params: map[string]string{
				"asset":          "USD" + ":" + accountOne,
				"liquidity_pool": "48672641c88264272787837f5c306f5ce93be3c2c7df68a092fbea55f5f4aa1d",
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "sponsor and liquidity pool",
			params: map[string]string{
				"sponsor":        accountTwo,
				"liquidity_pool": "48672641c88264272787837f5c306f5ce93be3c2c7df68a092fbea55f5f4aa1d",
			},
			isInvalidAccountsParams: true,
		},
		{
			desc: "filtering by native asset",
			params: map[string]string{
				"asset": "native",
			},
			expectedInvalidField: "asset",
			expectedErr:          "you can't filter by asset: native",
		},
		{
			desc: "invalid asset",
			params: map[string]string{
				"asset_issuer": accountOne,
				"asset":        "USDCOP:someissuer",
			},
			expectedInvalidField: "asset",
			expectedErr:          customTagsErrorMessages["asset"],
		},
		{
			desc: "invalid liquidity pool",
			params: map[string]string{
				"liquidity_pool": "USDCOP:someissuer",
			},
			expectedInvalidField: "liquidity_pool",
			expectedErr:          "USDCOP:someissuer does not validate as sha256",
		},
		{
			desc: "liquidity pool too short",
			params: map[string]string{
				"liquidity_pool": "48672641c882642727",
			},
			expectedInvalidField: "liquidity_pool",
			expectedErr:          "48672641c882642727 does not validate as sha256",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tt := test.Start(t)
			defer tt.Finish()
			q := &history.Q{tt.AuroraSession()}
			handler := &GetAccountsHandler{}

			_, err := handler.GetResourcePage(
				httptest.NewRecorder(),
				makeRequest(
					t,
					tc.params,
					map[string]string{},
					q,
				),
			)
			tt.Assert.Error(err)
			if tc.isInvalidAccountsParams {
				tt.Assert.Equal(invalidAccountsParams, err)
			} else {
				if tt.Assert.IsType(&problem.P{}, err) {
					p := err.(*problem.P)
					tt.Assert.Equal("bad_request", p.Type)
					tt.Assert.Equal(tc.expectedInvalidField, p.Extras["invalid_field"])
					tt.Assert.Equal(
						tc.expectedErr,
						p.Extras["reason"],
					)
				}
			}
		})
	}
}

func TestAccountQueryURLTemplate(t *testing.T) {
	tt := assert.New(t)
	expected := "/accounts{?signer,sponsor,asset,liquidity_pool,cursor,limit,order}"
	accountsQuery := AccountsQuery{}
	tt.Equal(expected, accountsQuery.URITemplate())
}
