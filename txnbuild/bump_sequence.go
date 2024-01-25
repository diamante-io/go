package txnbuild

import (
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// BumpSequence represents the Diamcircle bump sequence operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type BumpSequence struct {
	BumpTo        int64
	SourceAccount string
}

// BuildXDR for BumpSequence returns a fully configured XDR Operation.
func (bs *BumpSequence) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	opType := xdr.OperationTypeBumpSequence
	xdrOp := xdr.BumpSequenceOp{BumpTo: xdr.SequenceNumber(bs.BumpTo)}
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, bs.SourceAccount)
	} else {
		SetOpSourceAccount(&op, bs.SourceAccount)
	}
	return op, nil
}

// FromXDR for BumpSequence initialises the txnbuild struct from the corresponding xdr Operation.
func (bs *BumpSequence) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetBumpSequenceOp()
	if !ok {
		return errors.New("error parsing bump_sequence operation from xdr")
	}

	bs.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	bs.BumpTo = int64(result.BumpTo)
	return nil
}

// Validate for BumpSequence validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (bs *BumpSequence) Validate(withMuxedAccounts bool) error {
	err := validateAmount(bs.BumpTo)
	if err != nil {
		return NewValidationError("BumpTo", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (bs *BumpSequence) GetSourceAccount() string {
	return bs.SourceAccount
}
