//Package assets is a simple helper package to help convert to/from xdr.AssetType values
package assets

import (
	"github.com/go-errors/errors"
	"github.com/diamcircle/go/xdr"
)

// ErrInvalidString gets returns when the string form of the asset type is invalid
var ErrInvalidString = errors.New("invalid asset type: was not one of 'native', 'credit_alphanum4', 'credit_alphanum12', 'liquidity_pool_shares'")

//ErrInvalidValue gets returned when the xdr.AssetType int value is not one of the valid enum values
var ErrInvalidValue = errors.New("unknown asset type, cannot convert to string")

// AssetTypeMap is the read-only (i.e. don't modify it) map from string names to
// xdr.AssetType values
var AssetTypeMap = map[string]xdr.AssetType{
	"native":            xdr.AssetTypeAssetTypeNative,
	"credit_alphanum4":  xdr.AssetTypeAssetTypeCreditAlphanum4,
	"credit_alphanum12": xdr.AssetTypeAssetTypeCreditAlphanum12,
}

//Parse creates an asset from the provided strings.  See AssetTypeMap for valid strings for aType.
func Parse(aType string) (result xdr.AssetType, err error) {
	result, ok := AssetTypeMap[aType]
	if !ok {
		err = errors.New(ErrInvalidString)
	}

	return
}

//String returns the appropriate string representation of the provided xdr.AssetType.
func String(aType xdr.AssetType) (string, error) {
	for s, v := range AssetTypeMap {
		if v == aType {
			return s, nil
		}
	}

	return "", errors.New(ErrInvalidValue)
}
