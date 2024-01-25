package txnbuild

import (
	"go/support/errors"
	"go/xdr"
)

// Inflation represents the diamcircle inflation operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type Inflation struct {
	SourceAccount string
}

// BuildXDR for Inflation returns a fully configured XDR Operation.
func (inf *Inflation) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	opType := xdr.OperationTypeInflation
	body, err := xdr.NewOperationBody(opType, nil)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, inf.SourceAccount)
	} else {
		SetOpSourceAccount(&op, inf.SourceAccount)
	}
	return op, nil
}

// FromXDR for Inflation initialises the txnbuild struct from the corresponding xdr Operation.
func (inf *Inflation) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	if xdrOp.Body.Type != xdr.OperationTypeInflation {
		return errors.New("error parsing inflation operation from xdr")
	}
	inf.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	return nil
}

// Validate for Inflation is just a method that implements the Operation interface. No logic is actually performed
// because the inflation operation does not have any required field. Nil is always returned.
func (inf *Inflation) Validate(withMuxedAccounts bool) error {
	// no required fields, return nil.
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (inf *Inflation) GetSourceAccount() string {
	return inf.SourceAccount
}
