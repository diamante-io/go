package actions

import (
	"net/http/httptest"
	"testing"

	protocol "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/test"
	"go/support/errors"
	"go/support/render/problem"
	"go/xdr"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

func TestGetClaimableBalanceByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	accountID := "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML"
	asset := xdr.MustNewCreditAsset("USD", accountID)
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 2, 3},
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	cBalance := history.ClaimableBalance{
		BalanceID: id,
		Claimants: []history.Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              asset,
		Amount:             10,
		LastModifiedLedger: 123,
	}

	err = q.UpsertClaimableBalances(tt.Ctx, []history.ClaimableBalance{cBalance})
	tt.Assert.NoError(err)

	handler := GetClaimableBalanceByIDHandler{}
	response, err := handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"id": id},
		q,
	))
	tt.Assert.NoError(err)

	resource := response.(protocol.ClaimableBalance)
	tt.Assert.Equal(id, resource.BalanceID)

	// try to fetch claimable balance which does not exist
	balanceID = xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &xdr.Hash{1, 1, 1},
	}
	id, err = xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	_, err = handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"id": id},
		q,
	))
	tt.Assert.Error(err)
	tt.Assert.True(q.NoRows(errors.Cause(err)))

	// try to fetch a random invalid hex id
	_, err = handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"id": "0000001112122"},
		q,
	))
	tt.Assert.Error(err)
	p := err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("id", p.Extras["invalid_field"])
	tt.Assert.Equal("0000001112122 does not validate as claimableBalanceID", p.Extras["reason"])

	// try to fetch an empty id
	_, err = handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"id": ""},
		q,
	))
	tt.Assert.Error(err)
	p = err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("id", p.Extras["invalid_field"])
	tt.Assert.Equal("non zero value required", p.Extras["reason"])
}

func buildClaimableBalance(tt *test.T, balanceIDHash xdr.Hash, accountID string, ledger int32, asset *xdr.Asset) history.ClaimableBalance {
	balanceAsset := xdr.MustNewNativeAsset()
	var sponsor null.String
	if asset != nil {
		balanceAsset = *asset
		sponsor = null.StringFrom(accountID)
	}
	balanceID := xdr.ClaimableBalanceId{
		Type: xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0,
		V0:   &balanceIDHash,
	}
	id, err := xdr.MarshalHex(balanceID)
	tt.Assert.NoError(err)
	return history.ClaimableBalance{
		BalanceID: id,
		Claimants: []history.Claimant{
			{
				Destination: accountID,
				Predicate: xdr.ClaimPredicate{
					Type: xdr.ClaimPredicateTypeClaimPredicateUnconditional,
				},
			},
		},
		Asset:              balanceAsset,
		Amount:             10,
		LastModifiedLedger: uint32(ledger),
		Sponsor:            sponsor,
	}
}

func TestGetClaimableBalances(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	entriesMeta := []struct {
		id        xdr.Hash
		accountID string
		ledger    int32
		asset     *xdr.Asset
	}{
		{
			xdr.Hash{4, 0, 0},
			"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			1235,
			&usd,
		},
		{
			xdr.Hash{3, 0, 0},
			"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			1235,
			&euro,
		},
		{
			xdr.Hash{2, 0, 0},
			"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			1234,
			&usd,
		},
		{
			xdr.Hash{1, 0, 0},
			"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			1233,
			nil,
		},
	}

	var hCBs []history.ClaimableBalance

	for _, e := range entriesMeta {
		cb := buildClaimableBalance(tt, e.id, e.accountID, e.ledger, e.asset)
		hCBs = append(hCBs, cb)
	}

	err := q.UpsertClaimableBalances(tt.Ctx, hCBs)
	tt.Assert.NoError(err)

	handler := GetClaimableBalancesHandler{}
	response, err := handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{},
		q,
	))
	tt.Assert.NoError(err)
	tt.Assert.Len(response, 4)

	// check response is sorted in ascending order
	for entriesIndex, responseIndex := len(hCBs)-1, 0; entriesIndex >= 0; entriesIndex, responseIndex = entriesIndex-1, responseIndex+1 {
		entry := hCBs[entriesIndex]
		tt.Assert.Equal(entry.BalanceID, response[responseIndex].(protocol.ClaimableBalance).BalanceID)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"cursor": response[3].(protocol.ClaimableBalance).PagingToken(),
		},
		map[string]string{},
		q,
	))
	tt.Assert.NoError(err)
	tt.Assert.Len(response, 0)

	// test limit
	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{"limit": "2"},
		map[string]string{},
		q,
	))
	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	// response should be the last 2 elements of entries sorted by ID
	for entriesIndex, responseIndex := len(hCBs)-1, 0; entriesIndex >= 2; entriesIndex, responseIndex = entriesIndex-1, responseIndex+1 {
		entry := hCBs[entriesIndex]
		tt.Assert.Equal(entry.BalanceID, response[responseIndex].(protocol.ClaimableBalance).BalanceID)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "2",
			"cursor": response[1].(protocol.ClaimableBalance).PagingToken(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	// response should be the first 2 elements of entries sorted by ID
	for entriesIndex, responseIndex := len(hCBs)-3, 0; entriesIndex >= 0; entriesIndex, responseIndex = entriesIndex-1, responseIndex+1 {
		entry := hCBs[entriesIndex]
		tt.Assert.Equal(entry.BalanceID, response[responseIndex].(protocol.ClaimableBalance).BalanceID)
	}

	// next page should be 0, there are no new claimable balances ingested
	lastIngestedCursor := response[1].(protocol.ClaimableBalance).PagingToken()
	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "2",
			"cursor": lastIngestedCursor,
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 0)

	// new claimable balances are ingest and one of them updated, they should appear in the next pages
	cbToBeUpdated := hCBs[3]
	cbToBeUpdated.Sponsor = null.StringFrom("GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML")
	cbToBeUpdated.LastModifiedLedger = 1238
	q.UpsertClaimableBalances(tt.Ctx, []history.ClaimableBalance{cbToBeUpdated})

	entriesMeta = []struct {
		id        xdr.Hash
		accountID string
		ledger    int32
		asset     *xdr.Asset
	}{
		{
			xdr.Hash{4, 4, 4},
			"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			1236,
			nil,
		},
		{
			xdr.Hash{1, 1, 1},
			"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			1237,
			nil,
		},
	}

	hCBs = nil
	for _, e := range entriesMeta {
		entry := buildClaimableBalance(tt, e.id, e.accountID, e.ledger, e.asset)
		hCBs = append(hCBs, entry)
	}

	err = q.UpsertClaimableBalances(tt.Ctx, hCBs)
	tt.Assert.NoError(err)

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 6)

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "2",
			"cursor": lastIngestedCursor,
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	// response should be the first 2 elements of entries
	for i, entry := range hCBs {
		tt.Assert.Equal(entry.BalanceID, response[i].(protocol.ClaimableBalance).BalanceID)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "2",
			"cursor": response[1].(protocol.ClaimableBalance).PagingToken(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 1)

	tt.Assert.Equal(cbToBeUpdated.BalanceID, response[0].(protocol.ClaimableBalance).BalanceID)

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "2",
			"cursor": response[0].(protocol.ClaimableBalance).PagingToken(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 0)

	// in descending order
	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit": "2",
			"order": "desc",
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	tt.Assert.Equal(cbToBeUpdated.BalanceID, response[0].(protocol.ClaimableBalance).BalanceID)

	tt.Assert.Equal(hCBs[1].BalanceID, response[1].(protocol.ClaimableBalance).BalanceID)

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"limit":  "1",
			"order":  "desc",
			"cursor": response[1].(protocol.ClaimableBalance).PagingToken(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 1)

	tt.Assert.Equal(hCBs[0].BalanceID, response[0].(protocol.ClaimableBalance).BalanceID)

	// filter by asset
	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"asset": native.StringCanonical(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 3)

	for _, resource := range response {
		tt.Assert.Equal(
			native.StringCanonical(),
			resource.(protocol.ClaimableBalance).Asset,
		)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"asset": usd.StringCanonical(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	for _, resource := range response {
		tt.Assert.Equal(
			usd.StringCanonical(),
			resource.(protocol.ClaimableBalance).Asset,
		)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"asset": euro.StringCanonical(),
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 1)

	for _, resource := range response {
		tt.Assert.Equal(
			euro.StringCanonical(),
			resource.(protocol.ClaimableBalance).Asset,
		)
	}

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"sponsor": "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 3)
	for _, resource := range response {
		tt.Assert.Equal(
			"GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
			resource.(protocol.ClaimableBalance).Sponsor,
		)
	}

	// filter by claimant
	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"claimant": "GC3C4AKRBQLHOJ45U4XG35ESVWRDECWO5XLDGYADO6DPR3L7KIDVUMML",
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 4)

	response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"claimant": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		},
		map[string]string{},
		q,
	))

	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)
}

func TestCursorAndOrderValidation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	handler := GetClaimableBalancesHandler{}
	_, err := handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"cursor": "-1-00000043d380c38a2f2cac46ab63674064c56fdce6b977fdef1a278ad50e1a7e6a5e18",
		},
		map[string]string{},
		q,
	))
	p := err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("cursor", p.Extras["invalid_field"])
	tt.Assert.Equal("The first part should be a number higher than 0 and the second part should be a valid claimable balance ID", p.Extras["reason"])

	_, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"cursor": "1003529-00000043d380c38a2f2cac46ab63674064c56fdce6b977fdef1a278ad50e1a7e6a5e18",
		},
		map[string]string{},
		q,
	))
	p = err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("cursor", p.Extras["invalid_field"])
	tt.Assert.Equal("The first part should be a number higher than 0 and the second part should be a valid claimable balance ID", p.Extras["reason"])

	_, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{
			"order":  "arriba",
			"cursor": "1003529-00000043d380c38a2f2cac46ab63674064c56fdce6b977fdef1a278ad50e1a7e6a5e18",
		},
		map[string]string{},
		q,
	))
	p = err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("order", p.Extras["invalid_field"])
	tt.Assert.Equal("order: invalid value", p.Extras["reason"])
}

func TestClaimableBalancesQueryURLTemplate(t *testing.T) {
	tt := assert.New(t)
	expected := "/claimable_balances{?asset,sponsor,claimant,cursor,limit,order}"
	q := ClaimableBalancesQuery{}
	tt.Equal(expected, q.URITemplate())
}
