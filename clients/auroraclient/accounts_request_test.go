package auroraclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountsRequestBuildUrl(t *testing.T) {
	// error case: No parameters
	_, err := AccountsRequest{}.BuildURL()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid request: no parameters")
	}

	// error case: too many parameters
	_, err = AccountsRequest{
		Signer: "signer",
		Asset:  "asset",
	}.BuildURL()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid request: too many parameters")
	}

	// signer
	endpoint, err := AccountsRequest{Signer: "abcdef"}.BuildURL()
	require.NoError(t, err)
	assert.Equal(t, "accounts?signer=abcdef", endpoint)

	// asset
	endpoint, err = AccountsRequest{Asset: "abcdef"}.BuildURL()
	require.NoError(t, err)
	assert.Equal(t, "accounts?asset=abcdef", endpoint)

	// sponsor
	endpoint, err = AccountsRequest{Sponsor: "abcdef"}.BuildURL()
	require.NoError(t, err)
	assert.Equal(t, "accounts?sponsor=abcdef", endpoint)

	// liquidity_pool
	endpoint, err = AccountsRequest{LiquidityPool: "abcdef"}.BuildURL()
	require.NoError(t, err)
	assert.Equal(t, "accounts?liquidity_pool=abcdef", endpoint)
}
