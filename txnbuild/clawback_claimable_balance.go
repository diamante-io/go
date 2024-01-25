//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite
package txnbuild

import (
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// ClawbackClaimableBalance represents the Diamcircle clawback claimable balance operation. See
// https://developers.diamcircle.org/docs/start/list-of-operations/
type ClawbackClaimableBalance struct {
	BalanceID     string
	SourceAccount string
}

// BuildXDR for ClawbackClaimableBalance returns a fully configured XDR Operation.
func (cb *ClawbackClaimableBalance) BuildXDR(withMuxedAccounts bool) (xdr.Operation, error) {
	var xdrBalanceID xdr.ClaimableBalanceId
	err := xdr.SafeUnmarshalHex(cb.BalanceID, &xdrBalanceID)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set XDR 'ClaimableBalanceId' field")
	}
	xdrOp := xdr.ClawbackClaimableBalanceOp{
		BalanceId: xdrBalanceID,
	}

	opType := xdr.OperationTypeClawbackClaimableBalance
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	if withMuxedAccounts {
		SetOpSourceMuxedAccount(&op, cb.SourceAccount)
	} else {
		SetOpSourceAccount(&op, cb.SourceAccount)
	}
	return op, nil
}

// FromXDR for ClawbackClaimableBalance initializes the txnbuild struct from the corresponding xdr Operation.
func (cb *ClawbackClaimableBalance) FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error {
	result, ok := xdrOp.Body.GetClawbackClaimableBalanceOp()
	if !ok {
		return errors.New("error parsing clawback_claimable_balance operation from xdr")
	}

	cb.SourceAccount = accountFromXDR(xdrOp.SourceAccount, withMuxedAccounts)
	balanceID, err := xdr.MarshalHex(result.BalanceId)
	if err != nil {
		return errors.New("error parsing BalanceID in claim_claimable_balance operation from xdr")
	}
	cb.BalanceID = balanceID

	return nil
}

// Validate for ClawbackClaimableBalance validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (cb *ClawbackClaimableBalance) Validate(withMuxedAccounts bool) error {
	var xdrBalanceID xdr.ClaimableBalanceId
	err := xdr.SafeUnmarshalHex(cb.BalanceID, &xdrBalanceID)
	if err != nil {
		return NewValidationError("BalanceID", err.Error())
	}

	return nil
}

// GetSourceAccount returns the source account of the operation, or nil if not
// set.
func (cb *ClawbackClaimableBalance) GetSourceAccount() string {
	return cb.SourceAccount
}
