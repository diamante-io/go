package aurora

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/guregu/null"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go/protocols/aurora"
	"go/services/aurora/internal/actions"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/httpx"
	"go/services/aurora/internal/paths"
	auroraProblem "go/services/aurora/internal/render/problem"
	"go/services/aurora/internal/simplepath"
	"go/services/aurora/internal/test"
	"go/support/db"
	"go/support/render/problem"
	"go/xdr"
)

func mockPathFindingClient(
	tt *test.T,
	finder paths.Finder,
	maxAssetsParamLength int,
	session *db.Session,
) test.RequestHelper {
	router := chi.NewRouter()
	findPaths := httpx.ObjectActionHandler{actions.FindPathsHandler{
		PathFinder:           finder,
		MaxAssetsParamLength: maxAssetsParamLength,
		MaxPathLength:        3,
		SetLastLedgerHeader:  true,
	}}
	findFixedPaths := httpx.ObjectActionHandler{actions.FindFixedPathsHandler{
		PathFinder:           finder,
		MaxAssetsParamLength: maxAssetsParamLength,
		MaxPathLength:        3,
		SetLastLedgerHeader:  true,
	}}

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := session.Clone()
			s.BeginTx(&sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
				ReadOnly:  true,
			})
			defer s.Rollback()

			ctx := context.WithValue(
				r.Context(),
				&auroraContext.SessionContextKey,
				s,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	router.Group(func(r chi.Router) {
		router.Method("GET", "/paths", findPaths)
		router.Method("GET", "/paths/strict-receive", findPaths)
		router.Method("GET", "/paths/strict-send", findFixedPaths)
	})

	return test.NewRequestHelper(router)
}

func TestPathActionsStillIngesting(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)

	assertions := &test.Assertions{tt.Assert}
	finder := paths.MockFinder{}
	finder.On("Find", mock.Anything, mock.Anything, uint(3)).
		Return([]paths.Path{}, uint32(0), simplepath.ErrEmptyInMemoryOrderBook).Times(2)
	finder.On("FindFixedPaths", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]paths.Path{}, uint32(0), simplepath.ErrEmptyInMemoryOrderBook).Times(1)

	rh := mockPathFindingClient(
		tt,
		&finder,
		2,
		tt.auroraSession(),
	)

	var q = make(url.Values)

	q.Add(
		"source_assets",
		"native",
	)
	q.Add(
		"destination_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	q.Add("destination_asset_type", "credit_alphanum4")
	q.Add("destination_asset_code", "EUR")
	q.Add("destination_amount", "10")

	for _, uri := range []string{"/paths", "/paths/strict-receive"} {
		w := rh.Get(uri + "?" + q.Encode())
		assertions.Equal(auroraProblem.StillIngesting.Status, w.Code)
		assertions.Problem(w.Body, auroraProblem.StillIngesting)
		assertions.Equal("", w.Header().Get(actions.LastLedgerHeaderName))
	}

	q = make(url.Values)

	q.Add("destination_assets", "native")
	q.Add("source_asset_issuer", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN")
	q.Add("source_asset_type", "credit_alphanum4")
	q.Add("source_asset_code", "EUR")
	q.Add("source_amount", "10")

	w := rh.Get("/paths/strict-send" + "?" + q.Encode())
	assertions.Equal(auroraProblem.StillIngesting.Status, w.Code)
	assertions.Problem(w.Body, auroraProblem.StillIngesting)
	assertions.Equal("", w.Header().Get(actions.LastLedgerHeaderName))

	finder.AssertExpectations(t)
}

func TestPathActionsStrictReceive(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	sourceAssets := []xdr.Asset{
		xdr.MustNewCreditAsset("AAA", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN"),
		xdr.MustNewCreditAsset("USD", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN"),
		xdr.MustNewNativeAsset(),
	}
	sourceAccount := "GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP"

	q := &history.Q{tt.auroraSession()}

	account := history.AccountEntry{
		LastModifiedLedger: 1234,
		AccountID:          sourceAccount,
		Balance:            20000,
		SequenceNumber:     223456789,
		NumSubEntries:      10,
		Flags:              1,
		MasterWeight:       1,
		ThresholdLow:       2,
		ThresholdMedium:    3,
		ThresholdHigh:      4,
		BuyingLiabilities:  3,
		SellingLiabilities: 4,
	}

	err := q.UpsertAccounts(tt.Ctx, []history.AccountEntry{account})
	assert.NoError(t, err)

	assetsByKeys := map[string]xdr.Asset{}

	for _, asset := range sourceAssets {
		code := asset.String()
		assetsByKeys[code] = asset
		if code == "native" {
			continue
		}

		var assetType, assetCode, assetIssuer string
		asset.MustExtract(&assetType, &assetCode, &assetIssuer)

		var lk xdr.LedgerKey
		var lkStr string
		assert.NoError(t, lk.SetTrustline(xdr.MustAddress(sourceAccount), asset.ToTrustLineAsset()))
		lkStr, err = lk.MarshalBinaryBase64()
		assert.NoError(t, err)

		err = q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
			{
				AccountID:          sourceAccount,
				AssetType:          asset.Type,
				AssetIssuer:        assetIssuer,
				AssetCode:          assetCode,
				Balance:            10000,
				LedgerKey:          lkStr,
				Limit:              123456789,
				LiquidityPoolID:    "",
				BuyingLiabilities:  1,
				SellingLiabilities: 2,
				Flags:              0,
				LastModifiedLedger: 1234,
				Sponsor:            null.String{},
			},
		})
		assert.NoError(t, err)
	}
	tt.Assert.NoError(q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
		{
			AccountID:          sourceAccount,
			AssetType:          xdr.AssetTypeAssetTypePoolShare,
			Balance:            9876,
			LedgerKey:          "poolshareid1",
			Limit:              123456789,
			LiquidityPoolID:    "lpid123",
			Flags:              0,
			LastModifiedLedger: 1234,
			Sponsor:            null.String{},
		},
	}))

	finder := paths.MockFinder{}
	withSourceAssetsBalance := true

	finder.On("Find", mock.Anything, mock.Anything, uint(3)).Return([]paths.Path{}, uint32(1234), nil).Run(func(args mock.Arguments) {
		query := args.Get(1).(paths.Query)
		for _, asset := range query.SourceAssets {
			var assetType, code, issuer string

			asset.MustExtract(&assetType, &code, &issuer)
			if assetType == "native" {
				tt.Assert.NotNil(assetsByKeys["native"])
			} else {
				tt.Assert.NotNil(assetsByKeys[code])
			}

		}
		tt.Assert.Equal(xdr.MustNewCreditAsset("EUR", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN"), query.DestinationAsset)
		tt.Assert.Equal(xdr.Int64(100000000), query.DestinationAmount)

		if withSourceAssetsBalance {
			tt.Assert.Equal([]xdr.Int64{10000, 10000, 20000}, query.SourceAssetBalances)
			tt.Assert.True(query.ValidateSourceBalance)
		} else {
			tt.Assert.Equal([]xdr.Int64{0, 0, 0}, query.SourceAssetBalances)
			tt.Assert.False(query.ValidateSourceBalance)
		}

	}).Times(4)

	rh := mockPathFindingClient(
		tt,
		&finder,
		len(sourceAssets),
		tt.auroraSession(),
	)

	var withSourceAccount = make(url.Values)
	withSourceAccount.Add(
		"destination_account",
		"GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V",
	)
	withSourceAccount.Add(
		"source_account",
		sourceAccount,
	)
	withSourceAccount.Add(
		"destination_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	withSourceAccount.Add("destination_asset_type", "credit_alphanum4")
	withSourceAccount.Add("destination_asset_code", "EUR")
	withSourceAccount.Add("destination_amount", "10")

	withSourceAssets, err := url.ParseQuery(
		withSourceAccount.Encode(),
	)
	tt.Assert.NoError(err)
	withSourceAssets.Del("source_account")
	withSourceAssets.Add("source_assets", assetsToURLParam(sourceAssets))

	for _, uri := range []string{"/paths", "/paths/strict-receive"} {
		w := rh.Get(uri + "?" + withSourceAccount.Encode())
		tt.Assert.Equal(http.StatusOK, w.Code)
		tt.Assert.Equal("1234", w.Header().Get(actions.LastLedgerHeaderName))

		withSourceAssetsBalance = false
		w = rh.Get(uri + "?" + withSourceAssets.Encode())
		tt.Assert.Equal(http.StatusOK, w.Code)
		tt.Assert.Equal("1234", w.Header().Get(actions.LastLedgerHeaderName))
		withSourceAssetsBalance = true
	}

	finder.AssertExpectations(t)
}

func TestPathActionsEmptySourceAcount(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	assertions := &test.Assertions{tt.Assert}
	finder := paths.MockFinder{}
	rh := mockPathFindingClient(
		tt,
		&finder,
		2,
		tt.auroraSession(),
	)
	var q = make(url.Values)

	q.Add(
		"destination_account",
		"GAEDTJ4PPEFVW5XV2S7LUXBEHNQMX5Q2GM562RJGOQG7GVCE5H3HIB4V",
	)
	q.Add(
		"source_account",
		// there is no account associated with this address
		"GD5PM5X7Q5MM54ERO2P5PXW3HD6HVZI5IRZGEDWS4OPFBGHNTF6XOWQO",
	)
	q.Add(
		"destination_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	q.Add("destination_asset_type", "credit_alphanum4")
	q.Add("destination_asset_code", "EUR")
	q.Add("destination_amount", "10")

	for _, uri := range []string{"/paths", "/paths/strict-receive"} {
		w := rh.Get(uri + "?" + q.Encode())
		assertions.Equal(http.StatusOK, w.Code)
		inMemoryResponse := []aurora.Path{}
		tt.UnmarshalPage(w.Body, &inMemoryResponse)
		assertions.Empty(inMemoryResponse)
		tt.Assert.Equal("", w.Header().Get(actions.LastLedgerHeaderName))
	}
}

func TestPathActionsSourceAssetsValidation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	assertions := &test.Assertions{tt.Assert}
	finder := paths.MockFinder{}
	rh := mockPathFindingClient(
		tt,
		&finder,
		2,
		tt.auroraSession(),
	)

	missingSourceAccountAndAssets := make(url.Values)
	missingSourceAccountAndAssets.Add(
		"destination_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	missingSourceAccountAndAssets.Add("destination_asset_type", "credit_alphanum4")
	missingSourceAccountAndAssets.Add("destination_asset_code", "USD")
	missingSourceAccountAndAssets.Add("destination_amount", "10")

	sourceAccountAndAssets, err := url.ParseQuery(
		missingSourceAccountAndAssets.Encode(),
	)
	tt.Assert.NoError(err)
	sourceAccountAndAssets.Add(
		"source_assets",
		"EUR:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	sourceAccountAndAssets.Add(
		"source_account",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)

	tooManySourceAssets, err := url.ParseQuery(
		missingSourceAccountAndAssets.Encode(),
	)
	tt.Assert.NoError(err)
	tooManySourceAssets.Add(
		"source_assets",
		"EUR:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"GBP:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"USD:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"SEK:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)

	for _, testCase := range []struct {
		name            string
		q               url.Values
		expectedProblem problem.P
	}{
		{
			"both destination asset and destination account are missing",
			missingSourceAccountAndAssets,
			actions.SourceAssetsOrSourceAccountProblem,
		},
		{
			"both destination asset and destination account are present",
			sourceAccountAndAssets,
			actions.SourceAssetsOrSourceAccountProblem,
		},
		{
			"too many assets in destination_assets",
			tooManySourceAssets,
			*problem.MakeInvalidFieldProblem(
				"source_assets",
				fmt.Errorf("list of assets exceeds maximum length of 3"),
			),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			w := rh.Get("/paths/strict-receive?" + testCase.q.Encode())
			assertions.Equal(testCase.expectedProblem.Status, w.Code)
			assertions.Problem(w.Body, testCase.expectedProblem)
			assertions.Equal("", w.Header().Get(actions.LastLedgerHeaderName))
		})
	}
}

func TestPathActionsDestinationAssetsValidation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	assertions := &test.Assertions{tt.Assert}
	finder := paths.MockFinder{}
	rh := mockPathFindingClient(
		tt,
		&finder,
		2,
		tt.auroraSession(),
	)
	missingDestinationAccountAndAssets := make(url.Values)
	missingDestinationAccountAndAssets.Add(
		"source_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	missingDestinationAccountAndAssets.Add(
		"source_account",
		"GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP",
	)
	missingDestinationAccountAndAssets.Add("source_asset_type", "credit_alphanum4")
	missingDestinationAccountAndAssets.Add("source_asset_code", "USD")
	missingDestinationAccountAndAssets.Add("source_amount", "10")

	destinationAccountAndAssets, err := url.ParseQuery(
		missingDestinationAccountAndAssets.Encode(),
	)
	tt.Assert.NoError(err)
	destinationAccountAndAssets.Add(
		"destination_assets",
		"EUR:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	destinationAccountAndAssets.Add(
		"destination_account",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)

	tooManyDestinationAssets, err := url.ParseQuery(
		missingDestinationAccountAndAssets.Encode(),
	)
	tt.Assert.NoError(err)
	tooManyDestinationAssets.Add(
		"destination_assets",
		"EUR:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"GBP:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"USD:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN,"+
			"SEK:GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)

	for _, testCase := range []struct {
		name            string
		q               url.Values
		expectedProblem problem.P
	}{
		{
			"both destination asset and destination account are missing",
			missingDestinationAccountAndAssets,
			actions.DestinationAssetsOrDestinationAccountProblem,
		},
		{
			"both destination asset and destination account are present",
			destinationAccountAndAssets,
			actions.DestinationAssetsOrDestinationAccountProblem,
		},
		{
			"too many assets in destination_assets",
			tooManyDestinationAssets,
			*problem.MakeInvalidFieldProblem(
				"destination_assets",
				fmt.Errorf("list of assets exceeds maximum length of 3"),
			),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			w := rh.Get("/paths/strict-send?" + testCase.q.Encode())
			assertions.Equal(testCase.expectedProblem.Status, w.Code)
			assertions.Problem(w.Body, testCase.expectedProblem)
			assertions.Equal("", w.Header().Get(actions.LastLedgerHeaderName))
		})
	}
}

func TestPathActionsStrictSend(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	assertions := &test.Assertions{tt.Assert}
	historyQ := &history.Q{tt.auroraSession()}
	destinationAccount := "GARSFJNXJIHO6ULUBK3DBYKVSIZE7SC72S5DYBCHU7DKL22UXKVD7MXP"
	destinationAssets := []xdr.Asset{
		xdr.MustNewCreditAsset("AAA", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN"),
		xdr.MustNewCreditAsset("USD", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN"),
		xdr.MustNewNativeAsset(),
	}

	account := history.AccountEntry{
		LastModifiedLedger: 1234,
		AccountID:          destinationAccount,
		Balance:            20000,
		SequenceNumber:     223456789,
		NumSubEntries:      10,
		Flags:              1,
		MasterWeight:       1,
		ThresholdLow:       2,
		ThresholdMedium:    3,
		ThresholdHigh:      4,
		BuyingLiabilities:  3,
		SellingLiabilities: 4,
	}

	err := historyQ.UpsertAccounts(tt.Ctx, []history.AccountEntry{account})
	assert.NoError(t, err)

	assetsByKeys := map[string]xdr.Asset{}

	for _, asset := range destinationAssets {
		code := asset.String()
		assetsByKeys[code] = asset
		if code == "native" {
			continue
		}

		var assetType, assetCode, assetIssuer string
		asset.MustExtract(&assetType, &assetCode, &assetIssuer)

		var lk xdr.LedgerKey
		var lkStr string
		assert.NoError(t, lk.SetTrustline(xdr.MustAddress(destinationAccount), asset.ToTrustLineAsset()))
		lkStr, err = lk.MarshalBinaryBase64()
		assert.NoError(t, err)

		err = historyQ.UpsertTrustLines(tt.Ctx, []history.TrustLine{
			{
				AccountID:          destinationAccount,
				AssetType:          asset.Type,
				AssetIssuer:        assetIssuer,
				AssetCode:          assetCode,
				Balance:            10000,
				LedgerKey:          lkStr,
				Limit:              123456789,
				LiquidityPoolID:    "",
				BuyingLiabilities:  1,
				SellingLiabilities: 2,
				Flags:              0,
				LastModifiedLedger: 1234,
				Sponsor:            null.String{},
			},
		})
		assert.NoError(t, err)
	}
	tt.Assert.NoError(historyQ.UpsertTrustLines(tt.Ctx, []history.TrustLine{
		{
			AccountID:          destinationAccount,
			AssetType:          xdr.AssetTypeAssetTypePoolShare,
			Balance:            9876,
			LedgerKey:          "poolshareid1",
			Limit:              123456789,
			LiquidityPoolID:    "lpid123",
			Flags:              0,
			LastModifiedLedger: 1234,
			Sponsor:            null.String{},
		},
	}))

	finder := paths.MockFinder{}
	// withSourceAssetsBalance := true
	sourceAsset := xdr.MustNewCreditAsset("USD", "GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN")

	finder.On("FindFixedPaths", mock.Anything, sourceAsset, xdr.Int64(100000000), mock.Anything, uint(3)).Return([]paths.Path{}, uint32(1234), nil).Run(func(args mock.Arguments) {
		destinationAssets := args.Get(3).([]xdr.Asset)
		for _, asset := range destinationAssets {
			var assetType, code, issuer string

			asset.MustExtract(&assetType, &code, &issuer)
			if assetType == "native" {
				tt.Assert.NotNil(assetsByKeys["native"])
			} else {
				tt.Assert.NotNil(assetsByKeys[code])
			}

		}
	}).Times(2)

	rh := mockPathFindingClient(
		tt,
		&finder,
		len(destinationAssets),
		tt.auroraSession(),
	)

	var q = make(url.Values)

	q.Add(
		"source_asset_issuer",
		"GDSBCQO34HWPGUGQSP3QBFEXVTSR2PW46UIGTHVWGWJGQKH3AFNHXHXN",
	)
	q.Add("source_asset_type", "credit_alphanum4")
	q.Add("source_asset_code", "USD")
	q.Add("source_amount", "10")
	q.Add(
		"destination_account",
		destinationAccount,
	)

	w := rh.Get("/paths/strict-send?" + q.Encode())
	assertions.Equal(http.StatusOK, w.Code)
	assertions.Equal("1234", w.Header().Get(actions.LastLedgerHeaderName))

	q.Del("destination_account")
	q.Add("destination_assets", assetsToURLParam(destinationAssets))
	w = rh.Get("/paths/strict-send?" + q.Encode())
	assertions.Equal(http.StatusOK, w.Code)
	assertions.Equal("1234", w.Header().Get(actions.LastLedgerHeaderName))

	finder.AssertExpectations(t)
}

func assetsToURLParam(xdrAssets []xdr.Asset) string {
	var assets []string
	for _, xdrAsset := range xdrAssets {
		var assetType, code, issuer string
		xdrAsset.MustExtract(&assetType, &code, &issuer)
		if assetType == "native" {
			assets = append(assets, "native")
		} else {
			assets = append(assets, fmt.Sprintf("%s:%s", code, issuer))
		}
	}

	return strings.Join(assets, ",")
}

func TestFindFixedPathsQueryQueryURLTemplate(t *testing.T) {
	tt := assert.New(t)
	params := []string{
		"destination_account",
		"destination_assets",
		"source_asset_type",
		"source_asset_issuer",
		"source_asset_code",
		"source_amount",
	}
	expected := "/paths/strict-send{?" + strings.Join(params, ",") + "}"
	qp := actions.FindFixedPathsQuery{}
	tt.Equal(expected, qp.URITemplate())
}

func TestStrictReceivePathsQueryURLTemplate(t *testing.T) {
	tt := assert.New(t)
	params := []string{
		"source_assets",
		"source_account",
		"destination_account",
		"destination_asset_type",
		"destination_asset_issuer",
		"destination_asset_code",
		"destination_amount",
	}
	expected := "/paths/strict-receive{?" + strings.Join(params, ",") + "}"
	qp := actions.StrictReceivePathsQuery{}
	tt.Equal(expected, qp.URITemplate())
}
