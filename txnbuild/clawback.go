package txnbuild

import (
	"github.com/pkg/errors"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/xdr"
)

// Clawback represents the Diamcircle clawback operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type Clawback struct {
	From          string
	Amount        string
	Asset         Asset
	SourceAccount string
}

// BuildXDR for Clawback returns a fully configured XDR Operation.
func (cb *Clawback) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var fromMuxedAccount xdr.MuxedAccount
	var err error
	if withMuxedAccounts {
		err = fromMuxedAccount.SetAddress(cb.From)
	} else {
		err = fromMuxedAccount.SetEd25519Address(cb.From)
	}
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set from address")
	}

	if cb.Asset == nil {
		return xdr.Operation{}, errors.New("you must specify an asset for the clawback")
	}
	// Validate this is an issued asset
	if cb.Asset.IsNative() {
		return xdr.Operation{}, errors.New("clawbacks don't support the native (XLM) asset")
	}

	xdrAmount, err := amount.Parse(cb.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse amount")
	}

	xdrAsset, err := cb.Asset.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "can't convert asset to XDR")
	}

	opType := xdr.OperationTypeClawback
	xdrOp := xdr.ClawbackOp{
		From:   fromMuxedAccount,
		Amount: xdrAmount,
		Asset:  xdrAsset,
	}
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR Operation")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, cb.SourceAccount)
	} else {
		SetOpSourceAccount(&op, cb.SourceAccount)
	}
	return op, nil
}

// FromXDR for Clawback initialises the txnbuild struct from the corresponding xdr Operation.
func (cb *Clawback) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetClawbackOp()
	if !ok {
		return errors.New("error parsing clawback operation from xdr")
	}

	cb.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	cb.From = accountFromXDR(&result.From, withMuxedAccounts)
	cb.Amount = amount.String(result.Amount)
	asset, err := assetFromXDR(result.Asset)
	if err != nil {
		return errors.Wrap(err, "error parsing clawback operation from XDR")
	}
	cb.Asset = asset

	return nil
}

// Validate for Clawback validates the required struct fields. It returns an error if any
// of the fields are invalid. Otherwise, it returns nil.
func (cb *Clawback) Validate(withMuxedAccounts bool) error {
	var err error
	if withMuxedAccounts {
		_, err = xdr.AddressToMuxedAccount(cb.From)
	} else {
		_, err = xdr.AddressToAccountId(cb.From)
	}
	if err != nil {
		return NewValidationError("From", err.Error())
	}

	err = validateAmount(cb.Amount)
	if err != nil {
		return NewValidationError("Amount", err.Error())
	}

	err = validateAssetCode(cb.Asset)
	if err != nil {
		return NewValidationError("Asset", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or nil if not
// set.
func (cb *Clawback) GetSourceAccount() string {
	return cb.SourceAccount
}
