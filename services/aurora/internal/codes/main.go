//Package codes is a helper package to help convert to transaction and operation result codes
//to strings used in aurora.
package codes

import (
	"go/xdr"

	"github.com/go-errors/errors"
)

// ErrUnknownCode is returned when an unexepcted value is provided to `String`
var ErrUnknownCode = errors.New("Unknown result code")

const (
	// OpSuccess is the string code used to specify the operation was successful
	OpSuccess = "op_success"
	// OpMalformed is the string code used to specify the operation was malformed
	// in some way.
	OpMalformed = "op_malformed"
	// OpUnderfunded is the string code used to specify the operation failed
	// due to a lack of funds.
	OpUnderfunded = "op_underfunded"
	// OpLowReserve is the string code used to specify the operation failed
	// because the account in question does not have enough balance to satisfy
	// what their new minimum balance would be.
	OpLowReserve = "op_low_reserve"
	// OpLineFull occurs when a payment would cause a destination account to
	// exceed their declared trust limit for the asset being sent.
	OpLineFull = "op_line_full"
	// OpNoIssuer occurs when a operation does not correctly specify an issuing
	// asset
	OpNoIssuer = "op_no_issuer"
	// OpNoTrust occurs when there is no trust line to a given asset
	OpNoTrust = "op_no_trust"
	// OpNotAuthorized occurs when a trust line is not authorized
	OpNotAuthorized = "op_not_authorized"
	// OpDoesNotExist occurs when claimable balance or sponsorship does not exist
	OpDoesNotExist = "op_does_not_exist"
)

//String returns the appropriate string representation of the provided result code
func String(code interface{}) (string, error) {
	switch code := code.(type) {
	case xdr.TransactionResultCode:
		switch code {
		case xdr.TransactionResultCodeTxFeeBumpInnerSuccess:
			return "tx_fee_bump_inner_success", nil
		case xdr.TransactionResultCodeTxFeeBumpInnerFailed:
			return "tx_fee_bump_inner_failed", nil
		case xdr.TransactionResultCodeTxNotSupported:
			return "tx_not_supported", nil
		case xdr.TransactionResultCodeTxSuccess:
			return "tx_success", nil
		case xdr.TransactionResultCodeTxFailed:
			return "tx_failed", nil
		case xdr.TransactionResultCodeTxTooEarly:
			return "tx_too_early", nil
		case xdr.TransactionResultCodeTxTooLate:
			return "tx_too_late", nil
		case xdr.TransactionResultCodeTxMissingOperation:
			return "tx_missing_operation", nil
		case xdr.TransactionResultCodeTxBadSeq:
			return "tx_bad_seq", nil
		case xdr.TransactionResultCodeTxBadAuth:
			return "tx_bad_auth", nil
		case xdr.TransactionResultCodeTxInsufficientBalance:
			return "tx_insufficient_balance", nil
		case xdr.TransactionResultCodeTxNoAccount:
			return "tx_no_source_account", nil
		case xdr.TransactionResultCodeTxInsufficientFee:
			return "tx_insufficient_fee", nil
		case xdr.TransactionResultCodeTxBadAuthExtra:
			return "tx_bad_auth_extra", nil
		case xdr.TransactionResultCodeTxInternalError:
			return "tx_internal_error", nil
		case xdr.TransactionResultCodeTxBadSponsorship:
			return "tx_bad_sponsorship", nil
		}
	case xdr.OperationResultCode:
		switch code {
		case xdr.OperationResultCodeOpInner:
			return "op_inner", nil
		case xdr.OperationResultCodeOpBadAuth:
			return "op_bad_auth", nil
		case xdr.OperationResultCodeOpNoAccount:
			return "op_no_source_account", nil
		case xdr.OperationResultCodeOpNotSupported:
			return "op_not_supported", nil
		case xdr.OperationResultCodeOpTooManySubentries:
			return "op_too_many_subentries", nil
		case xdr.OperationResultCodeOpExceededWorkLimit:
			return "op_exceeded_work_limit", nil
		}
	case xdr.CreateAccountResultCode:
		switch code {
		case xdr.CreateAccountResultCodeCreateAccountSuccess:
			return OpSuccess, nil
		case xdr.CreateAccountResultCodeCreateAccountMalformed:
			return OpMalformed, nil
		case xdr.CreateAccountResultCodeCreateAccountUnderfunded:
			return OpUnderfunded, nil
		case xdr.CreateAccountResultCodeCreateAccountLowReserve:
			return OpLowReserve, nil
		case xdr.CreateAccountResultCodeCreateAccountAlreadyExist:
			return "op_already_exists", nil
		}
	case xdr.PaymentResultCode:
		switch code {
		case xdr.PaymentResultCodePaymentSuccess:
			return OpSuccess, nil
		case xdr.PaymentResultCodePaymentMalformed:
			return OpMalformed, nil
		case xdr.PaymentResultCodePaymentUnderfunded:
			return OpUnderfunded, nil
		case xdr.PaymentResultCodePaymentSrcNoTrust:
			return "op_src_no_trust", nil
		case xdr.PaymentResultCodePaymentSrcNotAuthorized:
			return "op_src_not_authorized", nil
		case xdr.PaymentResultCodePaymentNoDestination:
			return "op_no_destination", nil
		case xdr.PaymentResultCodePaymentNoTrust:
			return OpNoTrust, nil
		case xdr.PaymentResultCodePaymentNotAuthorized:
			return OpNotAuthorized, nil
		case xdr.PaymentResultCodePaymentLineFull:
			return OpLineFull, nil
		case xdr.PaymentResultCodePaymentNoIssuer:
			return OpNoIssuer, nil
		}
	case xdr.PathPaymentStrictReceiveResultCode:
		switch code {
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess:
			return OpSuccess, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveMalformed:
			return OpMalformed, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveUnderfunded:
			return OpUnderfunded, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNoTrust:
			return "op_src_no_trust", nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSrcNotAuthorized:
			return "op_src_not_authorized", nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoDestination:
			return "op_no_destination", nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoTrust:
			return OpNoTrust, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNotAuthorized:
			return OpNotAuthorized, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveLineFull:
			return OpLineFull, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer:
			return OpNoIssuer, nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveTooFewOffers:
			return "op_too_few_offers", nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOfferCrossSelf:
			return "op_cross_self", nil
		case xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveOverSendmax:
			return "op_over_source_max", nil
		}
	case xdr.ManageBuyOfferResultCode:
		switch code {
		case xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess:
			return OpSuccess, nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferMalformed:
			return OpMalformed, nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferSellNoTrust:
			return "op_sell_no_trust", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferBuyNoTrust:
			return "op_buy_no_trust", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferSellNotAuthorized:
			return "sell_not_authorized", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferBuyNotAuthorized:
			return "buy_not_authorized", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferLineFull:
			return OpLineFull, nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferUnderfunded:
			return OpUnderfunded, nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferCrossSelf:
			return "op_cross_self", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferSellNoIssuer:
			return "op_sell_no_issuer", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferBuyNoIssuer:
			return "buy_no_issuer", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferNotFound:
			return "op_offer_not_found", nil
		case xdr.ManageBuyOfferResultCodeManageBuyOfferLowReserve:
			return OpLowReserve, nil
		}
	case xdr.ManageSellOfferResultCode:
		switch code {
		case xdr.ManageSellOfferResultCodeManageSellOfferSuccess:
			return OpSuccess, nil
		case xdr.ManageSellOfferResultCodeManageSellOfferMalformed:
			return OpMalformed, nil
		case xdr.ManageSellOfferResultCodeManageSellOfferSellNoTrust:
			return "op_sell_no_trust", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferBuyNoTrust:
			return "op_buy_no_trust", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferSellNotAuthorized:
			return "sell_not_authorized", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferBuyNotAuthorized:
			return "buy_not_authorized", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferLineFull:
			return OpLineFull, nil
		case xdr.ManageSellOfferResultCodeManageSellOfferUnderfunded:
			return OpUnderfunded, nil
		case xdr.ManageSellOfferResultCodeManageSellOfferCrossSelf:
			return "op_cross_self", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferSellNoIssuer:
			return "op_sell_no_issuer", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferBuyNoIssuer:
			return "buy_no_issuer", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferNotFound:
			return "op_offer_not_found", nil
		case xdr.ManageSellOfferResultCodeManageSellOfferLowReserve:
			return OpLowReserve, nil
		}
	case xdr.SetOptionsResultCode:
		switch code {
		case xdr.SetOptionsResultCodeSetOptionsSuccess:
			return OpSuccess, nil
		case xdr.SetOptionsResultCodeSetOptionsLowReserve:
			return OpLowReserve, nil
		case xdr.SetOptionsResultCodeSetOptionsTooManySigners:
			return "op_too_many_signers", nil
		case xdr.SetOptionsResultCodeSetOptionsBadFlags:
			return "op_bad_flags", nil
		case xdr.SetOptionsResultCodeSetOptionsInvalidInflation:
			return "op_invalid_inflation", nil
		case xdr.SetOptionsResultCodeSetOptionsCantChange:
			return "op_cant_change", nil
		case xdr.SetOptionsResultCodeSetOptionsUnknownFlag:
			return "op_unknown_flag", nil
		case xdr.SetOptionsResultCodeSetOptionsThresholdOutOfRange:
			return "op_threshold_out_of_range", nil
		case xdr.SetOptionsResultCodeSetOptionsBadSigner:
			return "op_bad_signer", nil
		case xdr.SetOptionsResultCodeSetOptionsInvalidHomeDomain:
			return "op_invalid_home_domain", nil
		case xdr.SetOptionsResultCodeSetOptionsAuthRevocableRequired:
			return "op_auth_revocable_required", nil
		}
	case xdr.ChangeTrustResultCode:
		switch code {
		case xdr.ChangeTrustResultCodeChangeTrustSuccess:
			return OpSuccess, nil
		case xdr.ChangeTrustResultCodeChangeTrustMalformed:
			return OpMalformed, nil
		case xdr.ChangeTrustResultCodeChangeTrustNoIssuer:
			return OpNoIssuer, nil
		case xdr.ChangeTrustResultCodeChangeTrustInvalidLimit:
			return "op_invalid_limit", nil
		case xdr.ChangeTrustResultCodeChangeTrustLowReserve:
			return OpLowReserve, nil
		case xdr.ChangeTrustResultCodeChangeTrustSelfNotAllowed:
			return "op_self_not_allowed", nil
		case xdr.ChangeTrustResultCodeChangeTrustTrustLineMissing:
			return "op_trust_line_missing", nil
		case xdr.ChangeTrustResultCodeChangeTrustCannotDelete:
			return "op_cannot_delete", nil
		case xdr.ChangeTrustResultCodeChangeTrustNotAuthMaintainLiabilities:
			return "op_not_aut_maintain_liabilities", nil
		}
	case xdr.AllowTrustResultCode:
		switch code {
		case xdr.AllowTrustResultCodeAllowTrustSuccess:
			return OpSuccess, nil
		case xdr.AllowTrustResultCodeAllowTrustMalformed:
			return OpMalformed, nil
		case xdr.AllowTrustResultCodeAllowTrustNoTrustLine:
			return OpNoTrust, nil
		case xdr.AllowTrustResultCodeAllowTrustTrustNotRequired:
			return "op_not_required", nil
		case xdr.AllowTrustResultCodeAllowTrustCantRevoke:
			return "op_cant_revoke", nil
		case xdr.AllowTrustResultCodeAllowTrustSelfNotAllowed:
			return "op_self_not_allowed", nil
		case xdr.AllowTrustResultCodeAllowTrustLowReserve:
			return OpLowReserve, nil
		}
	case xdr.AccountMergeResultCode:
		switch code {
		case xdr.AccountMergeResultCodeAccountMergeSuccess:
			return OpSuccess, nil
		case xdr.AccountMergeResultCodeAccountMergeMalformed:
			return OpMalformed, nil
		case xdr.AccountMergeResultCodeAccountMergeNoAccount:
			return "op_no_account", nil
		case xdr.AccountMergeResultCodeAccountMergeImmutableSet:
			return "op_immutable_set", nil
		case xdr.AccountMergeResultCodeAccountMergeHasSubEntries:
			return "op_has_sub_entries", nil
		case xdr.AccountMergeResultCodeAccountMergeSeqnumTooFar:
			return "op_seq_num_too_far", nil
		case xdr.AccountMergeResultCodeAccountMergeDestFull:
			return "op_dest_full", nil
		case xdr.AccountMergeResultCodeAccountMergeIsSponsor:
			return "op_is_sponsor", nil
		}
	case xdr.InflationResultCode:
		switch code {
		case xdr.InflationResultCodeInflationSuccess:
			return OpSuccess, nil
		case xdr.InflationResultCodeInflationNotTime:
			return "op_not_time", nil
		}
	case xdr.ManageDataResultCode:
		switch code {
		case xdr.ManageDataResultCodeManageDataSuccess:
			return OpSuccess, nil
		case xdr.ManageDataResultCodeManageDataNotSupportedYet:
			return "op_not_supported_yet", nil
		case xdr.ManageDataResultCodeManageDataNameNotFound:
			return "op_data_name_not_found", nil
		case xdr.ManageDataResultCodeManageDataLowReserve:
			return OpLowReserve, nil
		case xdr.ManageDataResultCodeManageDataInvalidName:
			return "op_data_invalid_name", nil
		}
	case xdr.BumpSequenceResultCode:
		switch code {
		case xdr.BumpSequenceResultCodeBumpSequenceSuccess:
			return OpSuccess, nil
		case xdr.BumpSequenceResultCodeBumpSequenceBadSeq:
			return "op_bad_seq", nil
		}
	case xdr.PathPaymentStrictSendResultCode:
		switch code {
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess:
			return OpSuccess, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendMalformed:
			return OpMalformed, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendUnderfunded:
			return OpUnderfunded, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNoTrust:
			return "op_src_no_trust", nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSrcNotAuthorized:
			return "op_src_not_authorized", nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendNoDestination:
			return "op_no_destination", nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendNoTrust:
			return OpNoTrust, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendNotAuthorized:
			return OpNotAuthorized, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendLineFull:
			return OpLineFull, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer:
			return OpNoIssuer, nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendTooFewOffers:
			return "op_too_few_offers", nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendOfferCrossSelf:
			return "op_cross_self", nil
		case xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendUnderDestmin:
			return "op_under_dest_min", nil
		}
	case xdr.CreateClaimableBalanceResultCode:
		switch code {
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess:
			return OpSuccess, nil
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceMalformed:
			return OpMalformed, nil
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceLowReserve:
			return OpLowReserve, nil
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceNoTrust:
			return OpNoTrust, nil
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceNotAuthorized:
			return OpNotAuthorized, nil
		case xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceUnderfunded:
			return "op_underfunded", nil
		}
	case xdr.ClaimClaimableBalanceResultCode:
		switch code {
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceSuccess:
			return OpSuccess, nil
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceDoesNotExist:
			return OpDoesNotExist, nil
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceCannotClaim:
			return "op_cannot_claim", nil
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceLineFull:
			return OpLineFull, nil
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceNoTrust:
			return OpNoTrust, nil
		case xdr.ClaimClaimableBalanceResultCodeClaimClaimableBalanceNotAuthorized:
			return OpNotAuthorized, nil
		}
	case xdr.BeginSponsoringFutureReservesResultCode:
		switch code {
		case xdr.BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesSuccess:
			return OpSuccess, nil
		case xdr.BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesMalformed:
			return OpMalformed, nil
		case xdr.BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesAlreadySponsored:
			return "op_already_sponsored", nil
		case xdr.BeginSponsoringFutureReservesResultCodeBeginSponsoringFutureReservesRecursive:
			return "op_recursive", nil
		}
	case xdr.EndSponsoringFutureReservesResultCode:
		switch code {
		case xdr.EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesSuccess:
			return OpSuccess, nil
		case xdr.EndSponsoringFutureReservesResultCodeEndSponsoringFutureReservesNotSponsored:
			return "op_not_sponsored", nil
		}
	case xdr.RevokeSponsorshipResultCode:
		switch code {
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipSuccess:
			return OpSuccess, nil
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipDoesNotExist:
			return OpDoesNotExist, nil
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipNotSponsor:
			return "op_not_sponsor", nil
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipLowReserve:
			return OpLowReserve, nil
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipOnlyTransferable:
			return "op_only_transferable", nil
		case xdr.RevokeSponsorshipResultCodeRevokeSponsorshipMalformed:
			return OpMalformed, nil
		}
	case xdr.ClawbackResultCode:
		switch code {
		case xdr.ClawbackResultCodeClawbackSuccess:
			return OpSuccess, nil
		case xdr.ClawbackResultCodeClawbackMalformed:
			return OpMalformed, nil
		case xdr.ClawbackResultCodeClawbackNotClawbackEnabled:
			return "op_not_clawback_enabled", nil
		case xdr.ClawbackResultCodeClawbackNoTrust:
			return OpNoTrust, nil
		case xdr.ClawbackResultCodeClawbackUnderfunded:
			return OpUnderfunded, nil
		}
	case xdr.ClawbackClaimableBalanceResultCode:
		switch code {
		case xdr.ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceSuccess:
			return OpSuccess, nil
		case xdr.ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceDoesNotExist:
			return OpDoesNotExist, nil
		case xdr.ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotIssuer:
			return OpNoIssuer, nil
		case xdr.ClawbackClaimableBalanceResultCodeClawbackClaimableBalanceNotClawbackEnabled:
			return "op_not_clawback_enabled", nil
		}
	case xdr.SetTrustLineFlagsResultCode:
		switch code {
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsSuccess:
			return OpSuccess, nil
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsMalformed:
			return OpMalformed, nil
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsNoTrustLine:
			return OpNoTrust, nil
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsCantRevoke:
			return "op_cant_revoke", nil
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsInvalidState:
			return "op_invalid_state", nil
		case xdr.SetTrustLineFlagsResultCodeSetTrustLineFlagsLowReserve:
			return OpLowReserve, nil
		}
	case xdr.LiquidityPoolDepositResultCode:
		switch code {
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositSuccess:
			return OpSuccess, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositMalformed:
			return OpMalformed, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositNoTrust:
			return OpNoTrust, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositNotAuthorized:
			return OpNotAuthorized, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositUnderfunded:
			return OpUnderfunded, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositLineFull:
			return OpLineFull, nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositBadPrice:
			return "op_bad_price", nil
		case xdr.LiquidityPoolDepositResultCodeLiquidityPoolDepositPoolFull:
			return "op_pool_full", nil
		}
	case xdr.LiquidityPoolWithdrawResultCode:
		switch code {
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawSuccess:
			return OpSuccess, nil
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawMalformed:
			return OpMalformed, nil
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawNoTrust:
			return OpNoTrust, nil
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderfunded:
			return OpUnderfunded, nil
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawLineFull:
			return OpLineFull, nil
		case xdr.LiquidityPoolWithdrawResultCodeLiquidityPoolWithdrawUnderMinimum:
			return "op_under_minimum", nil
		}
	}

	return "", errors.New(ErrUnknownCode)
}

// ForOperationResult returns the strong representation used by aurora for the
// error code `opr`
func ForOperationResult(opr xdr.OperationResult) (string, error) {
	if opr.Code != xdr.OperationResultCodeOpInner {
		return String(opr.Code)
	}

	ir := opr.MustTr()
	var ic interface{}

	switch ir.Type {
	case xdr.OperationTypeCreateAccount:
		ic = ir.MustCreateAccountResult().Code
	case xdr.OperationTypePayment:
		ic = ir.MustPaymentResult().Code
	case xdr.OperationTypePathPaymentStrictReceive:
		ic = ir.MustPathPaymentStrictReceiveResult().Code
	case xdr.OperationTypeManageBuyOffer:
		ic = ir.MustManageBuyOfferResult().Code
	case xdr.OperationTypeManageSellOffer:
		ic = ir.MustManageSellOfferResult().Code
	case xdr.OperationTypeCreatePassiveSellOffer:
		ic = ir.MustCreatePassiveSellOfferResult().Code
	case xdr.OperationTypeSetOptions:
		ic = ir.MustSetOptionsResult().Code
	case xdr.OperationTypeChangeTrust:
		ic = ir.MustChangeTrustResult().Code
	case xdr.OperationTypeAllowTrust:
		ic = ir.MustAllowTrustResult().Code
	case xdr.OperationTypeAccountMerge:
		ic = ir.MustAccountMergeResult().Code
	case xdr.OperationTypeInflation:
		ic = ir.MustInflationResult().Code
	case xdr.OperationTypeManageData:
		ic = ir.MustManageDataResult().Code
	case xdr.OperationTypeBumpSequence:
		ic = ir.MustBumpSeqResult().Code
	case xdr.OperationTypePathPaymentStrictSend:
		ic = ir.MustPathPaymentStrictSendResult().Code
	case xdr.OperationTypeCreateClaimableBalance:
		ic = ir.MustCreateClaimableBalanceResult().Code
	case xdr.OperationTypeClaimClaimableBalance:
		ic = ir.MustClaimClaimableBalanceResult().Code
	case xdr.OperationTypeBeginSponsoringFutureReserves:
		ic = ir.MustBeginSponsoringFutureReservesResult().Code
	case xdr.OperationTypeEndSponsoringFutureReserves:
		ic = ir.MustEndSponsoringFutureReservesResult().Code
	case xdr.OperationTypeRevokeSponsorship:
		ic = ir.MustRevokeSponsorshipResult().Code
	case xdr.OperationTypeClawback:
		ic = ir.MustClawbackResult().Code
	case xdr.OperationTypeClawbackClaimableBalance:
		ic = ir.MustClawbackClaimableBalanceResult().Code
	case xdr.OperationTypeSetTrustLineFlags:
		ic = ir.MustSetTrustLineFlagsResult().Code
	case xdr.OperationTypeLiquidityPoolDeposit:
		ic = ir.MustLiquidityPoolDepositResult().Code
	case xdr.OperationTypeLiquidityPoolWithdraw:
		ic = ir.MustLiquidityPoolWithdrawResult().Code
	}

	return String(ic)
}
