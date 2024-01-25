package resourceadapter

import (
	"context"
	"fmt"

	"go/protocols/aurora"
	"go/protocols/aurora/operations"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/support/render/hal"
	"go/xdr"
)

// NewOperation creates a new operation resource, finding the appropriate type to use
// based upon the row's type.
func NewOperation(
	ctx context.Context,
	operationRow history.Operation,
	transactionHash string,
	transactionRow *history.Transaction,
	ledger history.Ledger,
) (result hal.Pageable, err error) {

	base := operations.Base{}
	err = PopulateBaseOperation(ctx, &base, operationRow, transactionHash, transactionRow, ledger)
	if err != nil {
		return
	}

	switch operationRow.Type {
	case xdr.OperationTypeBumpSequence:
		e := operations.BumpSequence{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeCreateAccount:
		e := operations.CreateAccount{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypePayment:
		e := operations.Payment{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypePathPaymentStrictReceive:
		e := operations.PathPayment{}
		e.Payment.Base = base
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageBuyOffer:
		e := operations.ManageBuyOffer{}
		e.Offer.Base = base
		err = operationRow.UnmarshalDetails(&e.Offer)
		if err == nil {
			hmo := history.ManageOffer{}
			err = operationRow.UnmarshalDetails(&hmo)
			e.OfferID = hmo.OfferID
		}
		result = e
	case xdr.OperationTypeManageSellOffer:
		e := operations.ManageSellOffer{}
		e.Offer.Base = base
		err = operationRow.UnmarshalDetails(&e.Offer)
		if err == nil {
			hmo := history.ManageOffer{}
			err = operationRow.UnmarshalDetails(&hmo)
			e.OfferID = hmo.OfferID
		}
		result = e
	case xdr.OperationTypeCreatePassiveSellOffer:
		e := operations.CreatePassiveSellOffer{}
		e.Offer.Base = base
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeSetOptions:
		e := operations.SetOptions{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeChangeTrust:
		e := operations.ChangeTrust{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeAllowTrust:
		e := operations.AllowTrust{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		// if the trustline is authorized, we want to reflect that it implies
		// authorized_to_maintain_liabilities to true, otherwise, we use the
		// value from details
		if e.Authorize {
			e.AuthorizeToMaintainLiabilities = e.Authorize
		}
		result = e
	case xdr.OperationTypeAccountMerge:
		e := operations.AccountMerge{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeInflation:
		e := operations.Inflation{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageData:
		e := operations.ManageData{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypePathPaymentStrictSend:
		e := operations.PathPaymentStrictSend{}
		e.Payment.Base = base
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeCreateClaimableBalance:
		e := operations.CreateClaimableBalance{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeClaimClaimableBalance:
		e := operations.ClaimClaimableBalance{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeBeginSponsoringFutureReserves:
		e := operations.BeginSponsoringFutureReserves{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeEndSponsoringFutureReserves:
		e := operations.EndSponsoringFutureReserves{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeRevokeSponsorship:
		e := operations.RevokeSponsorship{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeClawback:
		e := operations.Clawback{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeClawbackClaimableBalance:
		e := operations.ClawbackClaimableBalance{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeSetTrustLineFlags:
		e := operations.SetTrustLineFlags{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeLiquidityPoolDeposit:
		e := operations.LiquidityPoolDeposit{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeLiquidityPoolWithdraw:
		e := operations.LiquidityPoolWithdraw{Base: base}
		err = operationRow.UnmarshalDetails(&e)
		result = e
	default:
		result = base
	}

	return
}

// Populate fills out this resource using `row` as the source.
func PopulateBaseOperation(ctx context.Context, dest *operations.Base, operationRow history.Operation, transactionHash string, transactionRow *history.Transaction, ledger history.Ledger) error {
	dest.ID = fmt.Sprintf("%d", operationRow.ID)
	dest.PT = operationRow.PagingToken()
	dest.TransactionSuccessful = operationRow.TransactionSuccessful
	dest.SourceAccount = operationRow.SourceAccount
	if operationRow.SourceAccountMuxed.Valid {
		dest.SourceAccountMuxed = operationRow.SourceAccountMuxed.String
		muxedAccount := xdr.MustMuxedAddress(dest.SourceAccountMuxed)
		dest.SourceAccountMuxedID = uint64(muxedAccount.Med25519.Id)
	}
	populateOperationType(dest, operationRow)
	dest.LedgerCloseTime = ledger.ClosedAt
	dest.TransactionHash = transactionHash

	lb := hal.LinkBuilder{Base: auroraContext.BaseURL(ctx)}
	self := fmt.Sprintf("/operations/%d", operationRow.ID)
	dest.Links.Self = lb.Link(self)
	dest.Links.Succeeds = lb.Linkf("/effects?order=desc&cursor=%s", dest.PT)
	dest.Links.Precedes = lb.Linkf("/effects?order=asc&cursor=%s", dest.PT)
	dest.Links.Transaction = lb.Linkf("/transactions/%s", operationRow.TransactionHash)
	dest.Links.Effects = lb.Link(self, "effects")

	if transactionRow != nil {
		dest.Transaction = new(aurora.Transaction)
		return PopulateTransaction(ctx, transactionHash, dest.Transaction, *transactionRow)
	}
	return nil
}

func populateOperationType(dest *operations.Base, row history.Operation) {
	var ok bool
	dest.TypeI = int32(row.Type)
	dest.Type, ok = operations.TypeNames[row.Type]

	if !ok {
		dest.Type = "unknown"
	}
}
