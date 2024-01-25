package actions

import (
	"encoding/hex"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/services/aurora/internal/assets"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// Validateable allow structs to define their own custom validations.
type Validateable interface {
	Validate() error
}

func init() {
	govalidator.TagMap["accountID"] = isAccountID
	govalidator.TagMap["amount"] = isAmount
	govalidator.TagMap["assetType"] = isAssetType
	govalidator.TagMap["asset"] = isAsset
	govalidator.TagMap["claimableBalanceID"] = isClaimableBalanceID
	govalidator.TagMap["transactionHash"] = isTransactionHash
	govalidator.TagMap["sha256"] = govalidator.IsSHA256
	govalidator.TagMap["tradeType"] = isTradeType
}

var customTagsErrorMessages = map[string]string{
	"accountID":            "Account ID must start with `G` and contain 56 alphanum characters",
	"amount":               "Amount must be positive",
	"asset":                "Asset must be the string \"native\" or a string of the form \"Code:IssuerAccountID\" for issued assets.",
	"assetType":            "Asset type must be native, credit_alphanum4 or credit_alphanum12",
	"bool":                 "Filter should be true or false",
	"claimable_balance_id": "Claimable Balance ID must be the hex-encoded XDR representation of a Claimable Balance ID",
	"ledger_id":            "Ledger ID must be an integer higher than 0",
	"offer_id":             "Offer ID must be an integer higher than 0",
	"op_id":                "Operation ID must be an integer higher than 0",
	"transactionHash":      "Transaction hash must be a hex-encoded, lowercase SHA-256 hash",
	"tradeType":            "Trade type must be all, orderbook, or liquidity_pool",
}

func isTradeType(tradeType string) bool {
	return tradeType == history.AllTrades ||
		tradeType == history.OrderbookTrades ||
		tradeType == history.LiquidityPoolTrades
}

// isAsset validates if string contains a valid SEP11 asset
func isAsset(assetString string) bool {
	var asset xdr.Asset

	if strings.ToLower(assetString) == "native" {
		if err := asset.SetNative(); err != nil {
			return false
		}
	} else {

		parts := strings.Split(assetString, ":")
		if len(parts) != 2 {
			return false
		}

		code := parts[0]
		if !xdr.ValidAssetCode.MatchString(code) {
			return false
		}

		issuer, err := xdr.AddressToAccountId(parts[1])
		if err != nil {
			return false
		}

		if err := asset.SetCredit(code, issuer); err != nil {
			return false
		}
	}

	return true
}

func getSchemaErrorFieldMessage(field string, err error) error {
	if customMessage, ok := customTagsErrorMessages[field]; ok {
		return errors.New(customMessage)
	}

	if ce, ok := err.(schema.ConversionError); ok {
		customMessage, ok := customTagsErrorMessages[ce.Type.String()]
		if ok {
			return errors.New(customMessage)
		}
	}

	return err
}

func getErrorFieldMessage(err error) (string, string) {
	var field, message string

	switch err := err.(type) {
	case govalidator.Error:
		field = err.Name
		validator := err.Validator
		m, ok := customTagsErrorMessages[validator]
		// Give priority to inline custom error messages.
		// CustomErrorMessageExists when the validator is defined like:
		// `validatorName~custom message`
		if !ok || err.CustomErrorMessageExists {
			m = err.Err.Error()
		}
		message = m
	case govalidator.Errors:
		for _, item := range err.Errors() {
			field, message = getErrorFieldMessage(item)
			break
		}
	}

	return field, message
}

func isAssetType(str string) bool {
	if _, err := assets.Parse(str); err != nil {
		return false
	}

	return true
}

func isAccountID(str string) bool {
	if _, err := xdr.AddressToAccountId(str); err != nil {
		return false
	}

	return true
}

func isTransactionHash(str string) bool {
	decoded, err := hex.DecodeString(str)
	if err != nil {
		return false
	}

	return len(decoded) == 32 && strings.ToLower(str) == str
}

func isAmount(str string) bool {
	parsed, err := amount.Parse(str)
	switch {
	case err != nil:
		return false
	case parsed <= 0:
		return false
	}

	return true
}

func isClaimableBalanceID(str string) bool {
	var cbID xdr.ClaimableBalanceId
	err := xdr.SafeUnmarshalHex(str, &cbID)
	return err == nil
}
