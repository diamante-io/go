package txnbuild

import (
	"go/support/errors"
	"go/xdr"
)

// AccountMerge represents the diamcircle merge account operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type AccountMerge struct {
	Destination   string
	SourceAccount string
}

// BuildXDR for AccountMerge returns a fully configured XDR Operation.
func (am *AccountMerge) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var xdrOp xdr.MuxedAccount
	var err error
	if withMuxedAccounts {
		err = xdrOp.SetAddress(am.Destination)
	} else {
		err = xdrOp.SetEd25519Address(am.Destination)
	}
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set destination address")
	}

	opType := xdr.OperationTypeAccountMerge
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, am.SourceAccount)
	} else {
		SetOpSourceAccount(&op, am.SourceAccount)
	}
	return op, nil
}

// FromXDR for AccountMerge initialises the txnbuild struct from the corresponding xdr Operation.
func (am *AccountMerge) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	if xdrOp.Body.Type != xdr.OperationTypeAccountMerge {
		return errors.New("error parsing account_merge operation from xdr")
	}

	am.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	if xdrOp.Body.Destination != nil {
		if withMuxedAccounts {
			am.Destination = xdrOp.Body.Destination.Address()
		} else {
			aid := xdrOp.Body.Destination.ToAccountId()
			am.Destination = aid.Address()
		}
	}

	return nil
}

// Validate for AccountMerge validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (am *AccountMerge) Validate(withMuxedAccounts bool) error {
	var err error
	if withMuxedAccounts {
		_, err = xdr.AddressToMuxedAccount(am.Destination)
	} else {
		_, err = xdr.AddressToAccountId(am.Destination)
	}
	if err != nil {
		return NewValidationError("Destination", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (am *AccountMerge) GetSourceAccount() string {
	return am.SourceAccount
}
