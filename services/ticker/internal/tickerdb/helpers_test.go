package tickerdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type exampleDBModel struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	Counter int    `db:"counter"`
}

func TestGetDBFieldTags(t *testing.T) {
	m := exampleDBModel{
		ID:      10,
		Name:    "John Doe",
		Counter: 15,
	}

	fieldTags := getDBFieldTags(m, true)
	assert.Contains(t, fieldTags, "\"name\"")
	assert.Contains(t, fieldTags, "\"counter\"")
	assert.NotContains(t, fieldTags, "\"id\"")
	assert.Equal(t, 2, len(fieldTags))

	fieldTagsWithID := getDBFieldTags(m, false)
	assert.Contains(t, fieldTagsWithID, "\"name\"")
	assert.Contains(t, fieldTagsWithID, "\"counter\"")
	assert.Contains(t, fieldTagsWithID, "\"id\"")
	assert.Equal(t, 3, len(fieldTagsWithID))
}

func TestGetDBFieldValues(t *testing.T) {
	m := exampleDBModel{
		ID:      10,
		Name:    "John Doe",
		Counter: 15,
	}

	fieldValues := getDBFieldValues(m, true)
	assert.Contains(t, fieldValues, 15)
	assert.Contains(t, fieldValues, "John Doe")
	assert.NotContains(t, fieldValues, 10)
	assert.Equal(t, 2, len(fieldValues))

	fieldTagsWithID := getDBFieldValues(m, false)
	assert.Contains(t, fieldTagsWithID, 15)
	assert.Contains(t, fieldTagsWithID, "John Doe")
	assert.Contains(t, fieldTagsWithID, 10)
	assert.Equal(t, 3, len(fieldTagsWithID))
}

func TestGeneratePlaceholders(t *testing.T) {
	var p []interface{}
	p = append(p, 1)
	p = append(p, 2)
	p = append(p, 3)
	placeholder := generatePlaceholders(p)
	assert.Equal(t, "?, ?, ?", placeholder)
}

func TestGenerateWhereClause(t *testing.T) {
	baseAssetCode := new(string)
	baseAssetIssuer := new(string)
	*baseAssetCode = "baseAssetCode"
	*baseAssetIssuer = "baseAssetIssuer"

	where1, args1 := generateWhereClause([]optionalVar{
		optionalVar{"t1.base_asset_code", nil},
		optionalVar{"t1.base_asset_issuer", nil},
		optionalVar{"t1.counter_asset_code", nil},
		optionalVar{"t1.counter_asset_issuer", nil},
	})

	assert.Equal(t, "", where1)
	assert.Equal(t, 0, len(args1))

	where2, args2 := generateWhereClause([]optionalVar{
		optionalVar{"t1.base_asset_code", baseAssetCode},
		optionalVar{"t1.base_asset_issuer", nil},
		optionalVar{"t1.counter_asset_code", nil},
		optionalVar{"t1.counter_asset_issuer", nil},
	})

	assert.Equal(t, "WHERE t1.base_asset_code = ?", where2)
	assert.Equal(t, 1, len(args2))
	assert.Equal(t, *baseAssetCode, args2[0])

	where3, args3 := generateWhereClause([]optionalVar{
		optionalVar{"t1.base_asset_code", baseAssetCode},
		optionalVar{"t1.base_asset_issuer", baseAssetIssuer},
		optionalVar{"t1.counter_asset_code", nil},
		optionalVar{"t1.counter_asset_issuer", nil},
	})

	assert.Equal(t, "WHERE t1.base_asset_code = ? AND t1.base_asset_issuer = ?", where3)
	assert.Equal(t, 2, len(args3))
	assert.Equal(t, *baseAssetCode, args3[0])
	assert.Equal(t, *baseAssetIssuer, args3[1])
}

func TestGetBaseAndCounterCodes(t *testing.T) {
	a1, a2, err := getBaseAndCounterCodes("XLM_BTC")
	require.NoError(t, err)
	assert.Equal(t, "XLM", a1)
	assert.Equal(t, "BTC", a2)

	a3, a4, err := getBaseAndCounterCodes("BTC_XLM")
	require.NoError(t, err)
	assert.Equal(t, "XLM", a3)
	assert.Equal(t, "BTC", a4)

	a5, a6, err := getBaseAndCounterCodes("BTC_ETH")
	require.NoError(t, err)
	assert.Equal(t, "BTC", a5)
	assert.Equal(t, "ETH", a6)

	a7, a8, err := getBaseAndCounterCodes("ETH_BTC")
	require.NoError(t, err)
	assert.Equal(t, "BTC", a7)
	assert.Equal(t, "ETH", a8)

	_, _, err = getBaseAndCounterCodes("BTC")
	require.Error(t, err)
}
