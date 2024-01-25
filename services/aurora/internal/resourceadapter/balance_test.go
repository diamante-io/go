package resourceadapter

import (
	"testing"

	. "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestPopulateBalance(t *testing.T) {
	testAssetCode1 := "TEST_ASSET_1"
	testAssetCode2 := "TEST_ASSET_2"
	authorizedTrustline := history.TrustLine{
		AccountID:   "testID",
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
		AssetIssuer: "",
		AssetCode:   testAssetCode1,
		Limit:       100,
		Balance:     10,
		Flags:       1,
	}
	authorizedToMaintainLiabilitiesTrustline := history.TrustLine{
		AccountID:   "testID",
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
		AssetIssuer: "",
		AssetCode:   testAssetCode1,
		Limit:       100,
		Balance:     10,
		Flags:       2,
	}
	unauthorizedTrustline := history.TrustLine{
		AccountID:   "testID",
		AssetType:   xdr.AssetTypeAssetTypeCreditAlphanum12,
		AssetIssuer: "",
		AssetCode:   testAssetCode2,
		Limit:       100,
		Balance:     10,
		Flags:       0,
	}
	poolshareTrustline := history.TrustLine{
		AccountID: "testID",
		AssetType: xdr.AssetTypeAssetTypePoolShare,
		Limit:     100,
		Balance:   10,
		Flags:     0,
	}

	want := Balance{}
	err := PopulateAssetBalance(&want, authorizedTrustline)
	assert.NoError(t, err)
	assert.Equal(t, "credit_alphanum12", want.Type)
	assert.Equal(t, "0.0000010", want.Balance)
	assert.Equal(t, "0.0000100", want.Limit)
	assert.Equal(t, "", want.Issuer)
	assert.Equal(t, testAssetCode1, want.Code)
	assert.Equal(t, true, *want.IsAuthorized)
	assert.Equal(t, true, *want.IsAuthorizedToMaintainLiabilities)

	want = Balance{}
	err = PopulateAssetBalance(&want, authorizedToMaintainLiabilitiesTrustline)
	assert.NoError(t, err)
	assert.Equal(t, "credit_alphanum12", want.Type)
	assert.Equal(t, "0.0000010", want.Balance)
	assert.Equal(t, "0.0000100", want.Limit)
	assert.Equal(t, "", want.Issuer)
	assert.Equal(t, testAssetCode1, want.Code)
	assert.Equal(t, false, *want.IsAuthorized)
	assert.Equal(t, true, *want.IsAuthorizedToMaintainLiabilities)

	want = Balance{}
	err = PopulateAssetBalance(&want, unauthorizedTrustline)
	assert.NoError(t, err)
	assert.Equal(t, "credit_alphanum12", want.Type)
	assert.Equal(t, "0.0000010", want.Balance)
	assert.Equal(t, "0.0000100", want.Limit)
	assert.Equal(t, "", want.Issuer)
	assert.Equal(t, testAssetCode2, want.Code)
	assert.Equal(t, false, *want.IsAuthorized)
	assert.Equal(t, false, *want.IsAuthorizedToMaintainLiabilities)

	want = Balance{}
	err = PopulatePoolShareBalance(&want, poolshareTrustline)
	assert.NoError(t, err)
	assert.Equal(t, "liquidity_pool_shares", want.Type)
	assert.Equal(t, "0.0000010", want.Balance)
	assert.Equal(t, "0.0000100", want.Limit)
	assert.Equal(t, "", want.Issuer)
	assert.Equal(t, "", want.Code)
	assert.Equal(t, false, *want.IsAuthorized)
	assert.Equal(t, false, *want.IsAuthorizedToMaintainLiabilities)
}

func TestPopulateNativeBalance(t *testing.T) {
	want := Balance{}
	err := PopulateNativeBalance(&want, 10, 10, 10)
	assert.NoError(t, err)
	assert.Equal(t, "native", want.Type)
	assert.Equal(t, "0.0000010", want.Balance)
	assert.Equal(t, "0.0000010", want.BuyingLiabilities)
	assert.Equal(t, "0.0000010", want.SellingLiabilities)
	assert.Equal(t, "", want.Limit)
	assert.Equal(t, "", want.Issuer)
	assert.Equal(t, "", want.Code)
	assert.Nil(t, want.IsAuthorized)
	assert.Nil(t, want.IsAuthorizedToMaintainLiabilities)
}
