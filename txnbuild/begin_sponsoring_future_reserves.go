//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package txnbuild

import (
	"go/support/errors"
	"go/xdr"
)

// BeginSponsoringFutureReserves represents the diamcircle begin sponsoring future reserves operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type BeginSponsoringFutureReserves struct {
	SponsoredID   string
	SourceAccount string
}

// BuildXDR for BeginSponsoringFutureReserves returns a fully configured XDR Operation.
func (bs *BeginSponsoringFutureReserves) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	xdrOp := xdr.BeginSponsoringFutureReservesOp{}
	err := xdrOp.SponsoredId.SetAddress(bs.SponsoredID)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set sponsored id address")
	}
	opType := xdr.OperationTypeBeginSponsoringFutureReserves
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	SetOpSourceAccount(&op, bs.SourceAccount)
	return op, nil
}

// FromXDR for BeginSponsoringFutureReserves initializes the txnbuild struct from the corresponding xdr Operation.
func (bs *BeginSponsoringFutureReserves) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetBeginSponsoringFutureReservesOp()
	if !ok {
		return errors.New("error parsing begin_sponsoring_future_reserves operation from xdr")
	}
	bs.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	bs.SponsoredID = result.SponsoredId.Address()

	return nil
}

// Validate for BeginSponsoringFutureReserves validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (bs *BeginSponsoringFutureReserves) Validate(withMuxedAccounts bool) error {
	err := validatediamcirclePublicKey(bs.SponsoredID)
	if err != nil {
		return NewValidationError("SponsoredID", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (bs *BeginSponsoringFutureReserves) GetSourceAccount() string {
	return bs.SourceAccount
}
