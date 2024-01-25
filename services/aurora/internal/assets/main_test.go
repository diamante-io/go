package assets

import (
	"testing"

	"go/xdr"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

// Tests for Parse
func TestParseAssets(t *testing.T) {
	var (
		result xdr.AssetType
		err    error
	)

	result, err = Parse("native")
	assert.Equal(t, xdr.AssetTypeAssetTypeNative, result)

	assert.Nil(t, err)

	result, err = Parse("credit_alphanum4")
	assert.Equal(t, xdr.AssetTypeAssetTypeCreditAlphanum4, result)
	assert.Nil(t, err)

	result, err = Parse("credit_alphanum12")
	assert.Equal(t, xdr.AssetTypeAssetTypeCreditAlphanum12, result)
	assert.Nil(t, err)

	_, err = Parse("not_real")
	assert.True(t, errors.Is(err, ErrInvalidString))

	_, err = Parse("")
	assert.True(t, errors.Is(err, ErrInvalidString))
}

// Tests for String
func TestStringAssets(t *testing.T) {
	var (
		result string
		err    error
	)

	result, err = String(xdr.AssetTypeAssetTypeNative)
	assert.Equal(t, "native", result)
	assert.Nil(t, err)

	result, err = String(xdr.AssetTypeAssetTypeCreditAlphanum4)
	assert.Equal(t, "credit_alphanum4", result)
	assert.Nil(t, err)

	result, err = String(xdr.AssetTypeAssetTypeCreditAlphanum12)
	assert.Equal(t, "credit_alphanum12", result)
	assert.Nil(t, err)

	_, err = String(xdr.AssetType(15))
	assert.True(t, errors.Is(err, ErrInvalidValue))
}
