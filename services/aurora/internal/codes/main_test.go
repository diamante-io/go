package codes

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestForOperationResultCoversForAllOpTypes(t *testing.T) {
	for typ, s := range xdr.OperationTypeToStringMap {
		result := xdr.OperationResult{
			Code: xdr.OperationResultCodeOpInner,
			Tr: &xdr.OperationResultTr{
				Type: xdr.OperationType(typ),
			},
		}
		f := func() {
			ForOperationResult(result)
		}
		// it must panic because the operation result is not set
		assert.Panics(t, f, s)
	}

	// Check if all operations' codes are covered
	resultTypes := map[xdr.OperationType]reflect.Type{
		xdr.OperationTypeCreateAccount:                 reflect.TypeOf(xdr.CreateAccountResultCode(0)),
		xdr.OperationTypePayment:                       reflect.TypeOf(xdr.PaymentResultCode(0)),
		xdr.OperationTypePathPaymentStrictReceive:      reflect.TypeOf(xdr.PathPaymentStrictReceiveResultCode(0)),
		xdr.OperationTypeManageSellOffer:               reflect.TypeOf(xdr.ManageSellOfferResultCode(0)),
		xdr.OperationTypeCreatePassiveSellOffer:        reflect.TypeOf(xdr.ManageSellOfferResultCode(0)),
		xdr.OperationTypeSetOptions:                    reflect.TypeOf(xdr.SetOptionsResultCode(0)),
		xdr.OperationTypeChangeTrust:                   reflect.TypeOf(xdr.ChangeTrustResultCode(0)),
		xdr.OperationTypeAllowTrust:                    reflect.TypeOf(xdr.AllowTrustResultCode(0)),
		xdr.OperationTypeAccountMerge:                  reflect.TypeOf(xdr.AccountMergeResultCode(0)),
		xdr.OperationTypeInflation:                     reflect.TypeOf(xdr.InflationResultCode(0)),
		xdr.OperationTypeManageData:                    reflect.TypeOf(xdr.ManageDataResultCode(0)),
		xdr.OperationTypeBumpSequence:                  reflect.TypeOf(xdr.BumpSequenceResultCode(0)),
		xdr.OperationTypeManageBuyOffer:                reflect.TypeOf(xdr.ManageBuyOfferResultCode(0)),
		xdr.OperationTypePathPaymentStrictSend:         reflect.TypeOf(xdr.PathPaymentStrictSendResultCode(0)),
		xdr.OperationTypeCreateClaimableBalance:        reflect.TypeOf(xdr.CreateClaimableBalanceResultCode(0)),
		xdr.OperationTypeClaimClaimableBalance:         reflect.TypeOf(xdr.ClaimClaimableBalanceResultCode(0)),
		xdr.OperationTypeBeginSponsoringFutureReserves: reflect.TypeOf(xdr.BeginSponsoringFutureReservesResultCode(0)),
		xdr.OperationTypeEndSponsoringFutureReserves:   reflect.TypeOf(xdr.EndSponsoringFutureReservesResultCode(0)),
		xdr.OperationTypeRevokeSponsorship:             reflect.TypeOf(xdr.RevokeSponsorshipResultCode(0)),
		xdr.OperationTypeClawback:                      reflect.TypeOf(xdr.ClawbackResultCode(0)),
		xdr.OperationTypeClawbackClaimableBalance:      reflect.TypeOf(xdr.ClawbackClaimableBalanceResultCode(0)),
		xdr.OperationTypeSetTrustLineFlags:             reflect.TypeOf(xdr.SetTrustLineFlagsResultCode(0)),
		xdr.OperationTypeLiquidityPoolDeposit:          reflect.TypeOf(xdr.LiquidityPoolDepositResultCode(0)),
		xdr.OperationTypeLiquidityPoolWithdraw:         reflect.TypeOf(xdr.LiquidityPoolWithdrawResultCode(0)),
	}
	// If this is not equal it means one or more result struct is missing in resultTypes map.
	assert.Equal(t, len(xdr.OperationTypeToStringMap), len(resultTypes))

	type validEnum interface {
		ValidEnum(v int32) bool
	}

	for _, resultCode := range resultTypes {
		integerCode := int32(0)
		for {
			// Create a new variable of result code type and set it to the current
			// integer Code.
			val := reflect.New(resultCode).Elem()
			val.SetInt(int64(integerCode))

			// Then check if integer value of the code is valid. If it's not, break.
			// We exploit the fact that the code's integer values are a sequence:
			// [0, -1, -2, ...].
			iValue := val.Interface()
			valid := iValue.(validEnum).ValidEnum(integerCode)
			if !valid {
				break
			}

			res, err := String(iValue)
			if assert.NoError(t, err, fmt.Sprintf("type=%T code=%d not implemented", iValue, iValue)) {
				// Ensure value is not empty even when implemented
				assert.NotEmpty(t, res, fmt.Sprintf("type=%T code=%d empty", iValue, iValue))
			}
			integerCode--
		}
	}

	// make sure the check works for an unknown operation type
	result := xdr.OperationResult{
		Code: xdr.OperationResultCodeOpInner,
		Tr: &xdr.OperationResultTr{
			Type: xdr.OperationType(200000),
		},
	}
	f := func() {
		ForOperationResult(result)
	}
	// it doesn't panic because it doesn't branch out into the operation type
	assert.NotPanics(t, f)
}

func TestString(t *testing.T) {
	tests := []struct {
		Input    interface{}
		Expected string
		Err      error
	}{
		{xdr.TransactionResultCodeTxSuccess, "tx_success", nil},
		{xdr.OperationResultCodeOpBadAuth, "op_bad_auth", nil},
		{xdr.CreateAccountResultCodeCreateAccountLowReserve, "op_low_reserve", nil},
		{xdr.PaymentResultCodePaymentSrcNoTrust, "op_src_no_trust", nil},
		{xdr.SetOptionsResultCodeSetOptionsAuthRevocableRequired, "op_auth_revocable_required", nil},
		{xdr.ClawbackResultCodeClawbackNotClawbackEnabled, "op_not_clawback_enabled", nil},
		{0, "", ErrUnknownCode},
	}

	for _, test := range tests {
		actual, err := String(test.Input)

		if test.Err != nil {
			assert.NotNil(t, err)
			assert.Equal(t, test.Err.Error(), err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.Expected, actual)
		}
	}
}

//TODO: op_inner refers to inner result code
//TODO: non op_inner uses the outer result code
//TODO: one test for each operation type
