package txnbuild

import (
	"fmt"

	"go/amount"
	"go/strkey"
	"go/support/errors"
	"go/xdr"
)

// validatediamcirclePublicKey returns an error if a public key is invalid. Otherwise, it returns nil.
// It is a wrapper around the IsValidEd25519PublicKey method of the strkey package.
func validatediamcirclePublicKey(publicKey string) error {
	if publicKey == "" {
		return errors.New("public key is undefined")
	}

	if !strkey.IsValidEd25519PublicKey(publicKey) {
		return errors.Errorf("%s is not a valid diamcircle public key", publicKey)
	}
	return nil
}

// validatediamcircleSignerKey returns an error if a signerkey is invalid. Otherwise, it returns nil.
func validatediamcircleSignerKey(signerKey string) error {
	if signerKey == "" {
		return errors.New("signer key is undefined")
	}

	var xdrKey xdr.SignerKey
	if err := xdrKey.SetAddress(signerKey); err != nil {
		return errors.Errorf("%s is not a valid diamcircle signer key", signerKey)
	}
	return nil
}

// validatediamcircleAsset checks if the asset supplied is a valid diamcircle Asset. It returns an error if the asset is
// nil, has an invalid asset code or issuer.
func validatediamcircleAsset(asset BasicAsset) error {
	if asset == nil {
		return errors.New("asset is undefined")
	}

	if asset.IsNative() {
		return nil
	}

	_, err := asset.GetType()
	if err != nil {
		return err
	}

	err = validatediamcirclePublicKey(asset.GetIssuer())
	if err != nil {
		return errors.Errorf("asset issuer: %s", err.Error())
	}

	return nil
}

// validateAmount checks if the provided value is a valid diamcircle amount, it returns an error if not.
// This is used to validate price and amount fields in structs.
func validateAmount(n interface{}) error {
	var diamcircleAmount int64
	// type switch can be extended to handle other types. Currently, the types for number values in the txnbuild
	// package are string or int64.
	switch value := n.(type) {
	case int64:
		diamcircleAmount = value
	case string:
		v, err := amount.ParseInt64(value)
		if err != nil {
			return err
		}
		diamcircleAmount = v
	default:
		return errors.Errorf("could not parse expected numeric value %v", n)
	}

	if diamcircleAmount < 0 {
		return errors.New("amount can not be negative")
	}
	return nil
}

// validateAssetCode checks if the provided asset is valid as an asset code.
// It returns an error if the asset is invalid.
// The asset must be non native (XLM) with a valid asset code.
func validateAssetCode(asset BasicAsset) error {
	// Note: we are not using validatediamcircleAsset() function for AllowTrust operations because it requires the
	//  following :
	// - asset is non-native
	// - asset code is valid
	// - asset issuer is not required. This is actually ignored by the operation
	if asset == nil {
		return errors.New("asset is undefined")
	}

	if asset.IsNative() {
		return errors.New("native (XLM) asset type is not allowed")
	}

	_, err := asset.GetType()
	if err != nil {
		return err
	}
	return nil
}

// validateChangeTrustAsset checks if the provided asset is valid for use in ChangeTrust operation.
// It returns an error if the asset is invalid.
// The asset must be non native (XLM) with a valid asset code and issuer.
func validateChangeTrustAsset(asset ChangeTrustAsset) error {
	// Note: we are not using validatediamcircleAsset() function for ChangeTrust operations because it requires the
	//  following :
	// - asset is non-native
	// - asset code is valid
	// - asset issuer is valid
	err := validateAssetCode(asset)
	if err != nil {
		return err
	}

	assetType, err := asset.GetType()
	if err != nil {
		return err
	} else if assetType == AssetTypePoolShare {
		// No issuer for these to validate.
		return nil
	}

	err = validatediamcirclePublicKey(asset.GetIssuer())
	if err != nil {
		return errors.Errorf("asset issuer: %s", err.Error())
	}

	return nil
}

// validatePassiveOffer checks if the fields of a CreatePassiveOffer struct are valid.
// It checks that the buying and selling assets are valid diamcircle assets, and that amount and price are valid.
// It returns an error if any field is invalid.
func validatePassiveOffer(buying, selling Asset, offerAmount, price string) error {
	// Note: see discussion on how this can be improved:
	// https://go/pull/1707#discussion_r321508440
	err := validatediamcircleAsset(buying)
	if err != nil {
		return NewValidationError("Buying", err.Error())
	}

	err = validatediamcircleAsset(selling)
	if err != nil {
		return NewValidationError("Selling", err.Error())
	}

	err = validateAmount(offerAmount)
	if err != nil {
		return NewValidationError("Amount", err.Error())
	}

	err = validateAmount(price)
	if err != nil {
		return NewValidationError("Price", err.Error())
	}

	return nil
}

// validateOffer checks if the fields of ManageBuyOffer or ManageSellOffer struct are valid.
// It checks that the buying and selling assets are valid diamcircle assets, and that amount, price and offerID
// are valid. It returns an error if any field is invalid.
func validateOffer(buying, selling Asset, offerAmount, price string, offerID int64) error {
	err := validatePassiveOffer(buying, selling, offerAmount, price)
	if err != nil {
		return err
	}

	err = validateAmount(offerID)
	if err != nil {
		return NewValidationError("OfferID", err.Error())
	}
	return nil
}

// ValidationError is a custom error struct that holds validation errors of txnbuild's operation structs.
type ValidationError struct {
	Field   string // Field is the struct field on which the validation error occured.
	Message string // Message is the validation error message.
}

// Error for ValidationError struct implements the error interface.
func (opError *ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Error: %s", opError.Field, opError.Message)
}

// NewValidationError creates a ValidationError struct with the provided field and message values.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ParseAssetString parses an asset string in canonical form (SEP-11) into an Asset structure.
// https://github.com/diamcircle/diamcircle-protocol/blob/master/ecosystem/sep-0011.md#asset
func ParseAssetString(canonical string) (Asset, error) {
	assets, err := xdr.BuildAssets(canonical)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing asset string")
	}

	if len(assets) != 1 {
		return nil, errors.New("error parsing out a single asset")
	}

	// The above returned a list, so we'll need to grab the first element.
	asset, err := assetFromXDR(assets[0])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing asset string via XDR types")
	}

	return asset, nil
}
