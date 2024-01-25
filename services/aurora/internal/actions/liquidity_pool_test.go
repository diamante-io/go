package actions

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"go/keypair"
	protocol "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/test"
	"go/support/errors"
	"go/support/render/problem"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestGetLiquidityPoolByID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	lp := history.MakeTestPool(xdr.MustNewNativeAsset(), 100, usdAsset, 200)
	err := q.UpsertLiquidityPools(tt.Ctx, []history.LiquidityPool{lp})
	tt.Assert.NoError(err)

	handler := GetLiquidityPoolByIDHandler{}
	response, err := handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"liquidity_pool_id": lp.PoolID},
		q,
	))
	tt.Assert.NoError(err)

	resource := response.(protocol.LiquidityPool)
	tt.Assert.Equal(lp.PoolID, resource.ID)
	tt.Assert.Equal("constant_product", resource.Type)
	tt.Assert.Equal(uint32(30), resource.FeeBP)
	tt.Assert.Equal(uint64(12345), resource.TotalTrustlines)
	tt.Assert.Equal("0.0067890", resource.TotalShares)
	tt.Assert.Equal("native", resource.Reserves[0].Asset)
	tt.Assert.Equal("0.0000100", resource.Reserves[0].Amount)

	tt.Assert.Equal(usdAsset.StringCanonical(), resource.Reserves[1].Asset)
	tt.Assert.Equal("0.0000200", resource.Reserves[1].Amount)

	// try to fetch pool which does not exist
	_, err = handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"liquidity_pool_id": "123816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		q,
	))
	tt.Assert.Error(err)
	tt.Assert.True(q.NoRows(errors.Cause(err)))

	// try to fetch a random invalid hex id
	_, err = handler.GetResource(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{"liquidity_pool_id": "0000001112122"},
		q,
	))
	tt.Assert.Error(err)
	p := err.(*problem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("liquidity_pool_id", p.Extras["invalid_field"])
	tt.Assert.Equal("0000001112122 does not validate as sha256", p.Extras["reason"])
}

func TestGetLiquidityPools(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &history.Q{tt.auroraSession()}

	lp1 := history.MakeTestPool(nativeAsset, 100, usdAsset, 200)
	lp2 := history.MakeTestPool(eurAsset, 300, usdAsset, 400)
	err := q.UpsertLiquidityPools(tt.Ctx, []history.LiquidityPool{lp1, lp2})
	tt.Assert.NoError(err)

	handler := GetLiquidityPoolsHandler{}
	response, err := handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
		t,
		map[string]string{},
		map[string]string{},
		q,
	))
	tt.Assert.NoError(err)
	tt.Assert.Len(response, 2)

	resource := response[0].(protocol.LiquidityPool)
	tt.Assert.Equal(lp1.PoolID, resource.ID)
	tt.Assert.Equal("constant_product", resource.Type)
	tt.Assert.Equal(uint32(30), resource.FeeBP)
	tt.Assert.Equal(uint64(12345), resource.TotalTrustlines)
	tt.Assert.Equal("0.0067890", resource.TotalShares)

	tt.Assert.Equal("native", resource.Reserves[0].Asset)
	tt.Assert.Equal("0.0000100", resource.Reserves[0].Amount)

	tt.Assert.Equal(usdAsset.StringCanonical(), resource.Reserves[1].Asset)
	tt.Assert.Equal("0.0000200", resource.Reserves[1].Amount)

	resource = response[1].(protocol.LiquidityPool)
	tt.Assert.Equal(lp2.PoolID, resource.ID)
	tt.Assert.Equal("constant_product", resource.Type)
	tt.Assert.Equal(uint32(30), resource.FeeBP)
	tt.Assert.Equal(uint64(12345), resource.TotalTrustlines)
	tt.Assert.Equal("0.0067890", resource.TotalShares)

	tt.Assert.Equal(eurAsset.StringCanonical(), resource.Reserves[0].Asset)
	tt.Assert.Equal("0.0000300", resource.Reserves[0].Amount)

	tt.Assert.Equal(usdAsset.StringCanonical(), resource.Reserves[1].Asset)
	tt.Assert.Equal("0.0000400", resource.Reserves[1].Amount)

	t.Run("filtering by reserves", func(t *testing.T) {
		response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
			t,
			map[string]string{"reserves": "native"},
			map[string]string{},
			q,
		))
		assert.NoError(t, err)
		assert.Len(t, response, 1)
	})

	t.Run("paging via cursor", func(t *testing.T) {
		response, err = handler.GetResourcePage(httptest.NewRecorder(), makeRequest(
			t,
			map[string]string{"cursor": lp1.PoolID},
			map[string]string{},
			q,
		))
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		resource = response[0].(protocol.LiquidityPool)
		assert.Equal(t, lp2.PoolID, resource.ID)
	})

	t.Run("filtering by participating account", func(t *testing.T) {
		// we need to add trustlines to filter by account
		accountId := keypair.MustRandom().Address()
		assert.NoError(t, q.UpsertTrustLines(tt.Ctx, []history.TrustLine{
			history.MakeTestTrustline(accountId, nativeAsset, ""),
			history.MakeTestTrustline(accountId, eurAsset, ""),
			history.MakeTestTrustline(accountId, xdr.Asset{}, lp1.PoolID),
		}))

		request := makeRequest(
			t,
			map[string]string{"account": accountId},
			map[string]string{},
			q,
		)
		assert.Contains(t, request.URL.String(), fmt.Sprintf("account=%s", accountId))

		handler := GetLiquidityPoolsHandler{}
		response, err := handler.GetResourcePage(httptest.NewRecorder(), request)
		assert.NoError(t, err)
		assert.Len(t, response, 1)

		assert.IsType(t, protocol.LiquidityPool{}, response[0])
		resource = response[0].(protocol.LiquidityPool)
		assert.Equal(t, lp1.PoolID, resource.ID)
	})
}
