package txnbuild

import (
	"fmt"

	"go/xdr"
)

// Operation represents the operation types of the diamcircle network.
type Operation interface {
	BuildXDR(withMuxedAccounts bool) (xdr.Operation, error)
	FromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) error
	Validate(withMuxedAccounts bool) error
	GetSourceAccount() string
}

// SetOpSourceAccount sets the source account ID on an Operation.
func SetOpSourceAccount(op *xdr.Operation, sourceAccount string) {
	if sourceAccount == "" {
		return
	}
	var opSourceAccountID xdr.MuxedAccount
	opSourceAccountID.SetEd25519Address(sourceAccount)
	op.SourceAccount = &opSourceAccountID
}

// SetOpSourceAccount sets the source account ID on an Operation, allowing M-strkeys (as defined in SEP23).
func SetOpSourceMuxedAccount(op *xdr.Operation, sourceAccount string) {
	if sourceAccount == "" {
		return
	}
	var opSourceAccountID xdr.MuxedAccount
	opSourceAccountID.SetAddress(sourceAccount)
	op.SourceAccount = &opSourceAccountID
}

// operationFromXDR returns a txnbuild Operation from its corresponding XDR operation
func operationFromXDR(xdrOp xdr.Operation, withMuxedAccounts bool) (Operation, error) {
	var newOp Operation
	switch xdrOp.Body.Type {
	case xdr.OperationTypeCreateAccount:
		newOp = &CreateAccount{}
	case xdr.OperationTypePayment:
		newOp = &Payment{}
	case xdr.OperationTypePathPaymentStrictReceive:
		newOp = &PathPayment{}
	case xdr.OperationTypeManageSellOffer:
		newOp = &ManageSellOffer{}
	case xdr.OperationTypeCreatePassiveSellOffer:
		newOp = &CreatePassiveSellOffer{}
	case xdr.OperationTypeSetOptions:
		newOp = &SetOptions{}
	case xdr.OperationTypeChangeTrust:
		newOp = &ChangeTrust{}
	case xdr.OperationTypeAllowTrust:
		newOp = &AllowTrust{}
	case xdr.OperationTypeAccountMerge:
		newOp = &AccountMerge{}
	case xdr.OperationTypeInflation:
		newOp = &Inflation{}
	case xdr.OperationTypeManageData:
		newOp = &ManageData{}
	case xdr.OperationTypeBumpSequence:
		newOp = &BumpSequence{}
	case xdr.OperationTypeManageBuyOffer:
		newOp = &ManageBuyOffer{}
	case xdr.OperationTypePathPaymentStrictSend:
		newOp = &PathPaymentStrictSend{}
	case xdr.OperationTypeBeginSponsoringFutureReserves:
		newOp = &BeginSponsoringFutureReserves{}
	case xdr.OperationTypeEndSponsoringFutureReserves:
		newOp = &EndSponsoringFutureReserves{}
	case xdr.OperationTypeCreateClaimableBalance:
		newOp = &CreateClaimableBalance{}
	case xdr.OperationTypeClaimClaimableBalance:
		newOp = &ClaimClaimableBalance{}
	case xdr.OperationTypeRevokeSponsorship:
		newOp = &RevokeSponsorship{}
	case xdr.OperationTypeClawback:
		newOp = &Clawback{}
	case xdr.OperationTypeClawbackClaimableBalance:
		newOp = &ClawbackClaimableBalance{}
	case xdr.OperationTypeSetTrustLineFlags:
		newOp = &SetTrustLineFlags{}
	case xdr.OperationTypeLiquidityPoolDeposit:
		newOp = &LiquidityPoolDeposit{}
	case xdr.OperationTypeLiquidityPoolWithdraw:
		newOp = &LiquidityPoolWithdraw{}
	default:
		return nil, fmt.Errorf("unknown operation type: %d", xdrOp.Body.Type)
	}

	err := newOp.FromXDR(xdrOp, withMuxedAccounts)
	return newOp, err
}

func accountFromXDR(account *xdr.MuxedAccount, withMuxedAccounts bool) string {
	if account != nil {
		if withMuxedAccounts {
			return account.Address()
		} else {
			aid := account.ToAccountId()
			return aid.Address()
		}
	}
	return ""
}
