package integration

import (
	"math"
	"testing"

	"go/clients/auroraclient"
	"go/services/aurora/internal/test/integration"
	"go/txnbuild"

	"github.com/stretchr/testify/assert"
)

func TestNegativeSequenceTxSubmission(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol18Test(t)
	master := itest.Master()

	// First, bump the sequence to the maximum value -1
	op := txnbuild.BumpSequence{
		BumpTo: int64(math.MaxInt64) - 1,
	}
	itest.MustSubmitOperations(itest.MasterAccount(), master, &op)

	account := itest.MasterAccount()
	seqnum, err := account.GetSequenceNumber()
	tt.NoError(err)
	tt.Equal(int64(math.MaxInt64)-1, seqnum)

	// Submit a simple payment
	op2 := txnbuild.Payment{
		Destination: master.Address(),
		Amount:      "10",
		Asset:       txnbuild.NativeAsset{},
	}

	txResp := itest.MustSubmitOperations(account, master, &op2)
	tt.Equal(master.Address(), txResp.Account)

	// The transaction should had bumped our sequence to the maximum possible value
	seqnum, err = account.GetSequenceNumber()
	tt.NoError(err)
	tt.Equal(int64(math.MaxInt64), seqnum)

	// Using txnbuild to create another transaction should fail, since it would cause a sequence number overflow
	txResp, err = itest.SubmitOperations(account, master, &op2)
	tt.Error(err)
	tt.Contains(err.Error(), "sequence cannot be increased, it already reached MaxInt64")

	// We can enforce a negative sequence without errors by setting IncrementSequenceNum=false
	account = &txnbuild.SimpleAccount{
		AccountID: account.GetAccountID(),
		Sequence:  math.MinInt64,
	}
	txParams := txnbuild.TransactionParams{
		SourceAccount:        account,
		Operations:           []txnbuild.Operation{&op2},
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		IncrementSequenceNum: false,
		EnableMuxedAccounts:  true,
	}
	tx, err := txnbuild.NewTransaction(txParams)
	tt.NoError(err)
	tx, err = tx.Sign(integration.StandaloneNetworkPassphrase, master)
	tt.NoError(err)
	txResp, err = itest.Client().SubmitTransaction(tx)
	tt.Error(err)
	clientErr, ok := err.(*auroraclient.Error)
	tt.True(ok)
	codes, err := clientErr.ResultCodes()
	tt.NoError(err)
	tt.Equal("tx_bad_seq", codes.TransactionCode)

}
