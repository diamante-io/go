package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "go/clients/auroraclient"
	hEffects "go/protocols/aurora/effects"
	"go/protocols/aurora/operations"
	"go/services/aurora/internal/test/integration"
	"go/txnbuild"
	"go/xdr"
)

func TestClaimableBalanceCreationOperationsAndEffects(t *testing.T) {
	tt := assert.New(t)
	itest := integration.NewTest(t, protocol15Config)
	master := itest.Master()

	t.Run("Successful", func(t *testing.T) {
		op := txnbuild.CreateClaimableBalance{
			Destinations: []txnbuild.Claimant{
				txnbuild.NewClaimant(master.Address(), nil),
			},
			Amount: "10",
			Asset:  txnbuild.NativeAsset{},
		}

		txResp, err := itest.SubmitOperations(itest.MasterAccount(), master, &op)
		tt.NoError(err)

		var txResult xdr.TransactionResult
		err = xdr.SafeUnmarshalBase64(txResp.ResultXdr, &txResult)
		tt.NoError(err)
		tt.Equal(xdr.TransactionResultCodeTxSuccess, txResult.Result.Code)

		opResults, _ := txResult.OperationResults()
		tt.Len(opResults, 1)
		claimCreationOp := opResults[0].MustTr().CreateClaimableBalanceResult
		tt.Equal(xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess, claimCreationOp.Code)
		expectedBalanceID, err := xdr.MarshalHex(claimCreationOp.BalanceId)
		tt.NoError(err)

		response, err := itest.Client().Operations(sdk.OperationRequest{})
		ops := response.Embedded.Records
		tt.NoError(err)
		tt.Len(ops, 1)
		cb := ops[0].(operations.CreateClaimableBalance)
		tt.Equal("native", cb.Asset)
		tt.Equal("10.0000000", cb.Amount)
		tt.Equal(master.Address(), cb.SourceAccount)
		tt.Len(cb.Claimants, 1)

		claimant := cb.Claimants[0]
		tt.Equal(master.Address(), claimant.Destination)
		tt.Equal(xdr.ClaimPredicateTypeClaimPredicateUnconditional, claimant.Predicate.Type)

		eResponse, err := itest.Client().Effects(sdk.EffectRequest{ForOperation: cb.ID})
		effects := eResponse.Embedded.Records
		tt.Len(effects, 4)

		claimableBalanceCreatedEffect := effects[0].(hEffects.ClaimableBalanceCreated)
		tt.Equal("claimable_balance_created", claimableBalanceCreatedEffect.Type)
		tt.Equal("10.0000000", claimableBalanceCreatedEffect.Amount)
		tt.Equal("native", claimableBalanceCreatedEffect.Asset)
		tt.Equal(expectedBalanceID, claimableBalanceCreatedEffect.BalanceID)
		tt.Equal(master.Address(), claimableBalanceCreatedEffect.Account)

		claimableBalanceClaimantCreatedEffect := effects[1].(hEffects.ClaimableBalanceClaimantCreated)
		tt.Equal("claimable_balance_claimant_created", claimableBalanceClaimantCreatedEffect.Type)
		tt.Equal(master.Address(), claimableBalanceClaimantCreatedEffect.Account)
		tt.Equal(expectedBalanceID, claimableBalanceClaimantCreatedEffect.BalanceID)
		tt.Equal("10.0000000", claimableBalanceClaimantCreatedEffect.Amount)
		tt.Equal("native", claimableBalanceClaimantCreatedEffect.Asset)
		tt.Equal(
			xdr.ClaimPredicateTypeClaimPredicateUnconditional,
			claimableBalanceClaimantCreatedEffect.Predicate.Type,
		)

		accountDebitedEffect := effects[2].(hEffects.AccountDebited)
		tt.Equal("10.0000000", accountDebitedEffect.Amount)
		tt.Equal("native", accountDebitedEffect.Asset.Type)
		tt.Equal(master.Address(), accountDebitedEffect.Account)

		claimableBalanceSponsorshipCreated := effects[3].(hEffects.ClaimableBalanceSponsorshipCreated)
		tt.Equal("claimable_balance_sponsorship_created", claimableBalanceSponsorshipCreated.Type)
		tt.Equal(master.Address(), claimableBalanceSponsorshipCreated.Sponsor)
		tt.Equal(master.Address(), claimableBalanceSponsorshipCreated.Account)
		tt.Equal(expectedBalanceID, claimableBalanceSponsorshipCreated.BalanceID)
	})

	t.Run("Invalid", func(t *testing.T) {
		keys, accounts := itest.CreateAccounts(2, "50")
		op := txnbuild.CreateClaimableBalance{
			Destinations: []txnbuild.Claimant{
				txnbuild.NewClaimant(master.Address(), nil),
				txnbuild.NewClaimant(keys[1].Address(), nil),
			},
			Amount: "100",
			Asset:  txnbuild.NativeAsset{},
		}

		// this operation will fail because the claimable balance is trying to
		// reserve 100 XLM but the account only has 50.
		_, err := itest.SubmitOperations(accounts[0], keys[0], &op)
		tt.Error(err)

		response, err := itest.Client().Operations(sdk.OperationRequest{
			Order:         "desc",
			Limit:         1,
			IncludeFailed: true,
		})
		ops := response.Embedded.Records
		tt.NoError(err)
		tt.Len(ops, 1)
		cb := ops[0].(operations.CreateClaimableBalance)
		tt.False(cb.TransactionSuccessful)
		tt.Equal("native", cb.Asset)
		tt.Equal("100.0000000", cb.Amount)
		tt.Equal(keys[0].Address(), cb.SourceAccount)
		tt.Len(cb.Claimants, 2)

		eResponse, err := itest.Client().Effects(sdk.EffectRequest{ForOperation: cb.ID})
		effects := eResponse.Embedded.Records
		tt.Len(effects, 0)
	})
}
