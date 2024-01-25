package main

import (
	"net/http"
	"testing"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/friendbot/internal"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitFriendbot_createMinionAccounts_success(t *testing.T) {

	randSecretKey := "SDLNA2YUQSFIWVEB57M6D3OOCJHFVCVQZJ33LPA656KJESVRK5DQUZOH"
	botKP, err := keypair.Parse(randSecretKey)
	assert.NoError(t, err)

	botKeypair := botKP.(*keypair.Full)
	botAccountID := botKeypair.Address()
	botAccountMock := aurora.Account{
		AccountID: botAccountID,
		Sequence:  "1",
	}
	botAccount := internal.Account{AccountID: botAccountID, Sequence: 1}

	auroraClientMock := auroraclient.MockClient{}
	auroraClientMock.
		On("AccountDetail", auroraclient.AccountRequest{
			AccountID: botAccountID,
		}).
		Return(botAccountMock, nil)
	auroraClientMock.
		On("SubmitTransactionXDR", mock.Anything).
		Return(aurora.Transaction{}, nil)

	numMinion := 1000
	minionBatchSize := 50
	submitTxRetriesAllowed := 5
	createdMinions, err := createMinionAccounts(botAccount, botKeypair, "Test SDF Network ; September 2015", "10000", "101", numMinion, minionBatchSize, submitTxRetriesAllowed, 1000, &auroraClientMock)
	assert.NoError(t, err)

	assert.Equal(t, 1000, len(createdMinions))
}

func TestInitFriendbot_createMinionAccounts_timeoutError(t *testing.T) {
	randSecretKey := "SDLNA2YUQSFIWVEB57M6D3OOCJHFVCVQZJ33LPA656KJESVRK5DQUZOH"
	botKP, err := keypair.Parse(randSecretKey)
	assert.NoError(t, err)

	botKeypair := botKP.(*keypair.Full)
	botAccountID := botKeypair.Address()
	botAccountMock := aurora.Account{
		AccountID: botAccountID,
		Sequence:  "1",
	}
	botAccount := internal.Account{AccountID: botAccountID, Sequence: 1}

	auroraClientMock := auroraclient.MockClient{}
	auroraClientMock.
		On("AccountDetail", auroraclient.AccountRequest{
			AccountID: botAccountID,
		}).
		Return(botAccountMock, nil)

	// Successful on first 3 calls only, and then a timeout error occurs
	auroraClientMock.
		On("SubmitTransactionXDR", mock.Anything).
		Return(aurora.Transaction{}, nil).Times(3)
	hError := &auroraclient.Error{
		Problem: problem.P{
			Type:   "timeout",
			Title:  "Timeout",
			Status: http.StatusGatewayTimeout,
		},
	}
	auroraClientMock.
		On("SubmitTransactionXDR", mock.Anything).
		Return(aurora.Transaction{}, hError)

	numMinion := 1000
	minionBatchSize := 50
	submitTxRetriesAllowed := 5
	createdMinions, err := createMinionAccounts(botAccount, botKeypair, "Test SDF Network ; September 2015", "10000", "101", numMinion, minionBatchSize, submitTxRetriesAllowed, 1000, &auroraClientMock)
	assert.Equal(t, 150, len(createdMinions))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "after retrying 5 times: submitting create accounts tx:")
}
