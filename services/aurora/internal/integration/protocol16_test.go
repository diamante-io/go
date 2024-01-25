package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go/clients/auroraclient"
	"go/keypair"
	protocol "go/protocols/aurora"
	"go/protocols/aurora/effects"
	"go/protocols/aurora/operations"
	"go/services/aurora/internal/codes"
	"go/services/aurora/internal/test/integration"
	"go/txnbuild"
	"go/xdr"
)

func NewProtocol17Test(t *testing.T) *integration.Test {
	config := integration.Config{ProtocolVersion: 17}
	return integration.NewTest(t, config)
}

func TestProtocol16Basics(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	t.Run("Sanity", func(t *testing.T) {
		root, err := itest.Client().Root()
		tt.NoError(err)
		tt.LessOrEqual(int32(17), root.CoreSupportedProtocolVersion)
		tt.Equal(int32(17), root.CurrentProtocolVersion)

		// Submit a simple tx
		op := txnbuild.Payment{
			Destination: master.Address(),
			Amount:      "10",
			Asset:       txnbuild.NativeAsset{},
		}

		txResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &op)
		tt.Equal(master.Address(), txResp.Account)
		tt.Equal("1", txResp.AccountSequence)
	})
}

func TestHappyClawbackAccount(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	asset, fromKey, _ := setupClawbackAccountTest(tt, itest, master)

	// Clawback all of the asset
	submissionResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &txnbuild.Clawback{
		From:   fromKey.Address(),
		Amount: "10",
		Asset:  asset,
	})

	assertClawbackAccountSuccess(tt, itest, master, fromKey, "0.0000000", submissionResp)

	// Check the operation details
	opDetailsResponse, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: submissionResp.Hash,
	})
	tt.NoError(err)
	if tt.Len(opDetailsResponse.Embedded.Records, 1) {
		clawbackOp := opDetailsResponse.Embedded.Records[0].(operations.Clawback)
		tt.Equal("PTS", clawbackOp.Code)
		tt.Equal(fromKey.Address(), clawbackOp.From)
		tt.Equal("10.0000000", clawbackOp.Amount)
	}

	// Check the operation effects
	effectsResponse, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForTransaction: submissionResp.Hash,
	})
	tt.NoError(err)

	if tt.Len(effectsResponse.Embedded.Records, 2) {
		accountCredited := effectsResponse.Embedded.Records[0].(effects.AccountCredited)
		tt.Equal(master.Address(), accountCredited.Account)
		tt.Equal("10.0000000", accountCredited.Amount)
		tt.Equal(master.Address(), accountCredited.Issuer)
		tt.Equal("PTS", accountCredited.Code)
		accountDebited := effectsResponse.Embedded.Records[1].(effects.AccountDebited)
		tt.Equal(fromKey.Address(), accountDebited.Account)
		tt.Equal("10.0000000", accountDebited.Amount)
		tt.Equal(master.Address(), accountDebited.Issuer)
		tt.Equal("PTS", accountDebited.Code)
	}
}

func TestHappyClawbackAccountPartial(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	asset, fromKey, _ := setupClawbackAccountTest(tt, itest, master)

	// Partial clawback of the asset
	submissionResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &txnbuild.Clawback{
		From:   fromKey.Address(),
		Amount: "5",
		Asset:  asset,
	})

	assertClawbackAccountSuccess(tt, itest, master, fromKey, "5.0000000", submissionResp)
}

func TestHappyClawbackAccountSellingLiabilities(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	asset, fromKey, fromAccount := setupClawbackAccountTest(tt, itest, master)

	// Add a selling liability
	submissionResp := itest.MustSubmitOperations(fromAccount, fromKey, &txnbuild.ManageSellOffer{
		Buying:        txnbuild.NativeAsset{},
		Selling:       asset,
		Amount:        "5",
		Price:         "1",
		SourceAccount: fromAccount.GetAccountID(),
	})
	tt.True(submissionResp.Successful)

	// Full clawback of the asset, with a deauthorize/reauthorize sandwich
	submissionResp = itest.MustSubmitOperations(itest.MasterAccount(), master,
		&txnbuild.SetTrustLineFlags{
			Trustor:    fromAccount.GetAccountID(),
			Asset:      asset,
			ClearFlags: []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
		},
		&txnbuild.Clawback{
			From:   fromKey.Address(),
			Amount: "10",
			Asset:  asset,
		},
		&txnbuild.SetTrustLineFlags{
			Trustor:  fromAccount.GetAccountID(),
			Asset:    asset,
			SetFlags: []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
		},
	)

	assertClawbackAccountSuccess(tt, itest, master, fromKey, "0.0000000", submissionResp)
}

func TestSadClawbackAccountInsufficientFunds(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	asset, fromKey, _ := setupClawbackAccountTest(tt, itest, master)
	// Attempt to clawback more than the account holds.
	submissionResp, err := itest.SubmitOperations(itest.MasterAccount(), master, &txnbuild.Clawback{
		From:   fromKey.Address(),
		Amount: "20",
		Asset:  asset,
	})
	tt.Error(err)
	assertClawbackAccountFailed(tt, itest, master, fromKey, submissionResp)
}

func TestSadClawbackAccountSufficientFundsSellingLiabilities(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	asset, fromKey, fromAccount := setupClawbackAccountTest(tt, itest, master)

	// Add a selling liability
	submissionResp := itest.MustSubmitOperations(fromAccount, fromKey, &txnbuild.ManageSellOffer{
		Buying:        txnbuild.NativeAsset{},
		Selling:       asset,
		Amount:        "5",
		Price:         "1",
		SourceAccount: fromAccount.GetAccountID(),
	})

	// Attempt to clawback more than is available.
	submissionResp, err := itest.SubmitOperations(itest.MasterAccount(), master, &txnbuild.Clawback{
		From:   fromKey.Address(),
		Amount: "10",
		Asset:  asset,
	})
	tt.Error(err)

	assertClawbackAccountFailed(tt, itest, master, fromKey, submissionResp)
}

func setupClawbackAccountTest(tt *assert.Assertions, itest *integration.Test, master *keypair.Full) (txnbuild.CreditAsset, *keypair.Full, txnbuild.Account) {
	// Give the master account the revocable flag (needed to set the clawback flag)
	setRevocableFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthRevocable,
		},
	}

	itest.MustSubmitOperations(itest.MasterAccount(), master, &setRevocableFlag)

	// Give the master account the clawback flag
	setClawBackFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthClawbackEnabled,
		},
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &setClawBackFlag)

	// Make sure the clawback flag was set
	accountDetails := itest.MustGetAccount(master)
	tt.True(accountDetails.Flags.AuthClawbackEnabled)

	// Create another account from which to claw an asset back
	keyPairs, accounts := itest.CreateAccounts(1, "100")
	accountKeyPair := keyPairs[0]
	account := accounts[0]

	// Add some assets to the account with asset which allows clawback

	// Time machine to Spain before Euros were a thing
	pesetasAsset := txnbuild.CreditAsset{Code: "PTS", Issuer: master.Address()}
	itest.MustEstablishTrustline(accountKeyPair, account, pesetasAsset)
	pesetasPayment := txnbuild.Payment{
		Destination: accountKeyPair.Address(),
		Amount:      "10",
		Asset:       pesetasAsset,
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &pesetasPayment)

	accountDetails = itest.MustGetAccount(accountKeyPair)
	if tt.Len(accountDetails.Balances, 2) {
		pts := accountDetails.Balances[0]
		tt.Equal("PTS", pts.Code)
		if tt.NotNil(pts.IsClawbackEnabled) {
			tt.True(*pts.IsClawbackEnabled)
		}
		tt.Equal("10.0000000", pts.Balance)
	}

	return pesetasAsset, accountKeyPair, account
}

func assertClawbackAccountSuccess(tt *assert.Assertions, itest *integration.Test, master, accountKeyPair *keypair.Full, expectedBalance string, submissionResp protocol.Transaction) {
	tt.True(submissionResp.Successful)
	assertAccountBalance(tt, itest, accountKeyPair, expectedBalance)
}

func assertClawbackAccountFailed(tt *assert.Assertions, itest *integration.Test, master, accountKeyPair *keypair.Full, submissionResp protocol.Transaction) {
	tt.False(submissionResp.Successful)
	assertAccountBalance(tt, itest, accountKeyPair, "10.0000000")
}

func assertAccountBalance(tt *assert.Assertions, itest *integration.Test, accountKeyPair *keypair.Full, expectedBalance string) {
	accountDetails := itest.MustGetAccount(accountKeyPair)
	if tt.Len(accountDetails.Balances, 2) {
		pts := accountDetails.Balances[0]
		tt.Equal("PTS", pts.Code)
		if tt.NotNil(pts.IsClawbackEnabled) {
			tt.True(*pts.IsClawbackEnabled)
		}
		tt.Equal(expectedBalance, pts.Balance)
	}
}

func TestHappyClawbackClaimableBalance(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	// Give the master account the revocable flag (needed to set the clawback flag)
	setRevocableFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthRevocable,
		},
	}

	itest.MustSubmitOperations(itest.MasterAccount(), master, &setRevocableFlag)

	// Give the master account the clawback flag
	setClawBackFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthClawbackEnabled,
		},
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &setClawBackFlag)

	// Make sure the clawback flag was set
	accountDetails := itest.MustGetAccount(master)
	tt.True(accountDetails.Flags.AuthClawbackEnabled)

	// Create another account as a claimable balance claimant
	keyPairs, accounts := itest.CreateAccounts(1, "100")
	accountKeyPair := keyPairs[0]
	account := accounts[0]

	// Time machine to Spain before Euros were a thing
	pesetasAsset := txnbuild.CreditAsset{Code: "PTS", Issuer: master.Address()}
	itest.MustEstablishTrustline(accountKeyPair, account, pesetasAsset)

	// Make a claimable balance from the master account (and asset issuer) to the account with an asset which allows clawback
	pesetasCreateCB := txnbuild.CreateClaimableBalance{
		Amount: "10",
		Asset:  pesetasAsset,
		Destinations: []txnbuild.Claimant{
			txnbuild.NewClaimant(accountKeyPair.Address(), nil),
		},
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &pesetasCreateCB)

	// Check that the claimable balance was created, clawback is enabled and obtain the id to claw it back later on
	listCBResp, err := itest.Client().ClaimableBalances(auroraclient.ClaimableBalanceRequest{
		Claimant: accountKeyPair.Address(),
	})
	tt.NoError(err)
	cbID := ""
	if tt.Len(listCBResp.Embedded.Records, 1) {
		cb := listCBResp.Embedded.Records[0]
		tt.True(cb.Flags.ClawbackEnabled)
		cbID = cb.BalanceID
		tt.Equal(master.Address(), cb.Sponsor)
	}

	// check that its operations and transactions can be obtained
	transactionsResp, err := itest.Client().Transactions(auroraclient.TransactionRequest{
		ForClaimableBalance: cbID,
	})
	assert.NoError(t, err)
	assert.Len(t, transactionsResp.Embedded.Records, 1)

	operationsResp, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForClaimableBalance: cbID,
	})
	assert.NoError(t, err)
	if assert.Len(t, operationsResp.Embedded.Records, 1) {
		assert.IsType(t, operationsResp.Embedded.Records[0], operations.CreateClaimableBalance{})
	}

	// Clawback the claimable balance
	pesetasClawbackCB := txnbuild.ClawbackClaimableBalance{
		BalanceID: cbID,
	}
	clawbackCBResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &pesetasClawbackCB)

	// Make sure the claimable balance is clawed back (gone)
	_, err = itest.Client().ClaimableBalance(cbID)
	// Not found
	tt.Error(err)

	// check that its operations and transactions can still be obtained
	transactionsResp, err = itest.Client().Transactions(auroraclient.TransactionRequest{
		ForClaimableBalance: cbID,
	})
	assert.NoError(t, err)
	assert.Len(t, transactionsResp.Embedded.Records, 2)

	operationsResp, err = itest.Client().Operations(auroraclient.OperationRequest{
		ForClaimableBalance: cbID,
	})
	assert.NoError(t, err)
	if assert.Len(t, operationsResp.Embedded.Records, 2) {
		assert.IsType(t, operationsResp.Embedded.Records[0], operations.CreateClaimableBalance{})
		assert.IsType(t, operationsResp.Embedded.Records[1], operations.ClawbackClaimableBalance{})
	}

	// Check the operation details
	opDetailsResponse, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: clawbackCBResp.Hash,
	})
	tt.NoError(err)
	if tt.Len(opDetailsResponse.Embedded.Records, 1) {
		clawbackOp := opDetailsResponse.Embedded.Records[0].(operations.ClawbackClaimableBalance)
		tt.Equal(cbID, clawbackOp.BalanceID)
	}

	// Check the operation effects
	effectsResponse, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForTransaction: clawbackCBResp.Hash,
	})
	tt.NoError(err)

	if tt.Len(effectsResponse.Embedded.Records, 3) {
		claimableBalanceClawedBack := effectsResponse.Embedded.Records[0].(effects.ClaimableBalanceClawedBack)
		tt.Equal(master.Address(), claimableBalanceClawedBack.Account)
		tt.Equal(cbID, claimableBalanceClawedBack.BalanceID)
		accountCredited := effectsResponse.Embedded.Records[1].(effects.AccountCredited)
		tt.Equal(master.Address(), accountCredited.Account)
		tt.Equal("10.0000000", accountCredited.Amount)
		tt.Equal(accountCredited.Issuer, master.Address())
		tt.Equal(accountCredited.Code, "PTS")
		cbSponsorshipRemoved := effectsResponse.Embedded.Records[2].(effects.ClaimableBalanceSponsorshipRemoved)
		tt.Equal(master.Address(), cbSponsorshipRemoved.Account)
		tt.Equal(cbID, cbSponsorshipRemoved.BalanceID)
		tt.Equal(master.Address(), cbSponsorshipRemoved.FormerSponsor)
	}
}

func TestHappySetTrustLineFlags(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol17Test(t)
	master := itest.Master()

	// Give the master account the revocable flag (needed to set the clawback flag)
	setRevocableFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthRevocable,
		},
	}

	itest.MustSubmitOperations(itest.MasterAccount(), master, &setRevocableFlag)

	// Give the master account the clawback flag
	setClawBackFlag := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthClawbackEnabled,
		},
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &setClawBackFlag)

	// Make sure the clawback flag was set
	accountDetails := itest.MustGetAccount(master)
	tt.True(accountDetails.Flags.AuthClawbackEnabled)

	// Create another account fot the Trustline
	keyPairs, accounts := itest.CreateAccounts(1, "100")
	accountKeyPair := keyPairs[0]
	account := accounts[0]

	// Time machine to Spain before Euros were a thing
	pesetasAsset := txnbuild.CreditAsset{Code: "PTS", Issuer: master.Address()}
	itest.MustEstablishTrustline(accountKeyPair, account, pesetasAsset)
	// Confirm that the Trustline has the clawback flag
	accountDetails = itest.MustGetAccount(accountKeyPair)
	if tt.Len(accountDetails.Balances, 2) {
		pts := accountDetails.Balances[0]
		tt.Equal("PTS", pts.Code)
		if tt.NotNil(pts.IsClawbackEnabled) {
			tt.True(*pts.IsClawbackEnabled)
		}
	}

	// Clear the clawback flag
	setTrustlineFlags := txnbuild.SetTrustLineFlags{
		Trustor: accountKeyPair.Address(),
		Asset:   pesetasAsset,
		ClearFlags: []txnbuild.TrustLineFlag{
			txnbuild.TrustLineClawbackEnabled,
		},
	}
	submissionResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &setTrustlineFlags)

	// make sure it was cleared
	accountDetails = itest.MustGetAccount(accountKeyPair)
	if tt.Len(accountDetails.Balances, 2) {
		pts := accountDetails.Balances[0]
		tt.Equal("PTS", pts.Code)
		tt.Nil(pts.IsClawbackEnabled)
	}

	// Check the operation details
	opDetailsResponse, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: submissionResp.Hash,
	})
	tt.NoError(err)
	if tt.Len(opDetailsResponse.Embedded.Records, 1) {
		setTrustlineFlagsDetails := opDetailsResponse.Embedded.Records[0].(operations.SetTrustLineFlags)
		tt.Equal("PTS", setTrustlineFlagsDetails.Code)
		tt.Equal(master.Address(), setTrustlineFlagsDetails.Issuer)
		tt.Equal(accountKeyPair.Address(), setTrustlineFlagsDetails.Trustor)
		if tt.Len(setTrustlineFlagsDetails.ClearFlags, 1) {
			tt.True(xdr.TrustLineFlags(setTrustlineFlagsDetails.ClearFlags[0]).IsClawbackEnabledFlag())
		}
		if tt.Len(setTrustlineFlagsDetails.ClearFlagsS, 1) {
			tt.Equal(setTrustlineFlagsDetails.ClearFlagsS[0], "clawback_enabled")
		}
		tt.Len(setTrustlineFlagsDetails.SetFlags, 0)
		tt.Len(setTrustlineFlagsDetails.SetFlagsS, 0)
	}

	// Check the operation effects
	effectsResponse, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForTransaction: submissionResp.Hash,
	})
	tt.NoError(err)

	if tt.Len(effectsResponse.Embedded.Records, 1) {
		trustlineFlagsUpdated := effectsResponse.Embedded.Records[0].(effects.TrustlineFlagsUpdated)
		tt.Equal(master.Address(), trustlineFlagsUpdated.Account)
		tt.Equal(master.Address(), trustlineFlagsUpdated.Issuer)
		tt.Equal("PTS", trustlineFlagsUpdated.Code)
		tt.Nil(trustlineFlagsUpdated.Authorized)
		tt.Nil(trustlineFlagsUpdated.AuthorizedToMaintainLiabilities)
		if tt.NotNil(trustlineFlagsUpdated.ClawbackEnabled) {
			tt.False(*trustlineFlagsUpdated.ClawbackEnabled)
		}
	}

	// Try to set the clawback flag (we shouldn't be able to)
	setTrustlineFlags = txnbuild.SetTrustLineFlags{
		Trustor: accountKeyPair.Address(),
		Asset:   pesetasAsset,
		SetFlags: []txnbuild.TrustLineFlag{
			txnbuild.TrustLineClawbackEnabled,
		},
	}
	_, err = itest.SubmitOperations(itest.MasterAccount(), master, &setTrustlineFlags)
	if tt.Error(err) {
		clientErr, ok := err.(*auroraclient.Error)
		if tt.True(ok) {
			tt.Equal(400, clientErr.Problem.Status)
			resCodes, err := clientErr.ResultCodes()
			tt.NoError(err)
			tt.Equal(codes.OpMalformed, resCodes.OperationCodes[0])
		}

	}

}
