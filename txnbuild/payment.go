package txnbuild

import (
	"go/amount"
	"go/support/errors"
	"go/xdr"
)

// Payment represents the diamcircle payment operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type Payment struct {
	Destination   string
	Amount        string
	Asset         Asset
	SourceAccount string
}

// BuildXDR for Payment returns a fully configured XDR Operation.

func (p *Payment) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var destMuxedAccount xdr.MuxedAccount

	var err error
	if withMuxedAccounts {
		err = destMuxedAccount.SetAddress(p.Destination)
	} else {
		err = destMuxedAccount.SetEd25519Address(p.Destination)
	}
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set destination address")
	}

	xdrAmount, err := amount.Parse(p.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse amount")
	}

	if p.Asset == nil {
		return xdr.Operation{}, errors.New("you must specify an asset for payment")
	}
	xdrAsset, err := p.Asset.ToXDR()
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set asset type")
	}

	opType := xdr.OperationTypePayment
	xdrOp := xdr.PaymentOp{
		Destination: destMuxedAccount,
		Amount:      xdrAmount,
		Asset:       xdrAsset,
	}
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR Operation")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, p.SourceAccount)
	} else {
		SetOpSourceAccount(&op, p.SourceAccount)
	}
	return op, nil
}

// FromXDR for Payment initialises the txnbuild struct from the corresponding xdr Operation.
func (p *Payment) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetPaymentOp()
	if !ok {
		return errors.New("error parsing payment operation from xdr")
	}

	p.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	if withMuxedAccounts {
		p.Destination = result.Destination.Address()
	} else {
		destAID := result.Destination.ToAccountId()
		p.Destination = destAID.Address()
	}

	p.Amount = amount.String(result.Amount)

	asset, err := assetFromXDR(result.Asset)
	if err != nil {
		return errors.Wrap(err, "error parsing asset in payment operation")
	}
	p.Asset = asset

	return nil
}

// Validate for Payment validates the required struct fields. It returns an error if any
// of the fields are invalid. Otherwise, it returns nil.
func (p *Payment) Validate(withMuxedAccounts bool) error {
	var err error
	if withMuxedAccounts {
		_, err = xdr.AddressToMuxedAccount(p.Destination)
	} else {
		_, err = xdr.AddressToAccountId(p.Destination)
	}

	if err != nil {
		return NewValidationError("Destination", err.Error())
	}

	err = validatediamcircleAsset(p.Asset)
	if err != nil {
		return NewValidationError("Asset", err.Error())
	}

	err = validateAmount(p.Amount)
	if err != nil {
		return NewValidationError("Amount", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (p *Payment) GetSourceAccount() string {
	return p.SourceAccount
}
