package auroraclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetRequestBuildUrl(t *testing.T) {
	er := AssetRequest{}
	endpoint, err := er.BuildURL()

	// It should return valid all assets endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "assets", endpoint)

	er = AssetRequest{ForAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	endpoint, err = er.BuildURL()

	// It should return valid assets endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "assets?asset_issuer=GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU", endpoint)

	er = AssetRequest{ForAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU", ForAssetCode: "ABC", Order: OrderDesc}
	endpoint, err = er.BuildURL()

	// It should return valid assets endpoint and no errors
	require.NoError(t, err)
	assert.Equal(t, "assets?asset_code=ABC&asset_issuer=GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU&order=desc", endpoint)
}
