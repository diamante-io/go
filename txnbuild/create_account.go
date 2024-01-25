package txnbuild

import (
	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// CreateAccount represents the Diamcircle create account operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type CreateAccount struct {
	Destination   string
	Amount        string
	SourceAccount string
}

// BuildXDR for CreateAccount returns a fully configured XDR Operation.
func (ca *CreateAccount) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var xdrOp xdr.CreateAccountOp

	err := xdrOp.Destination.SetAddress(ca.Destination)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set destination address")
	}

	xdrOp.StartingBalance, err = amount.Parse(ca.Amount)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to parse amount")
	}

	opType := xdr.OperationTypeCreateAccount
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, ca.SourceAccount)
	} else {
		SetOpSourceAccount(&op, ca.SourceAccount)
	}
	return op, nil
}

// FromXDR for CreateAccount initialises the txnbuild struct from the corresponding xdr Operation.
func (ca *CreateAccount) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetCreateAccountOp()
	if !ok {
		return errors.New("error parsing create_account operation from xdr")
	}

	ca.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	ca.Destination = result.Destination.Address()
	ca.Amount = amount.String(result.StartingBalance)

	return nil
}

// Validate for CreateAccount validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (ca *CreateAccount) Validate(withMuxedAccounts bool) error {
	err := validateDiamcirclePublicKey(ca.Destination)
	if err != nil {
		return NewValidationError("Destination", err.Error())
	}

	err = validateAmount(ca.Amount)
	if err != nil {
		return NewValidationError("Amount", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (ca *CreateAccount) GetSourceAccount() string {
	return ca.SourceAccount
}
