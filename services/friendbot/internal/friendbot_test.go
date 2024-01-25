package internal

import (
	"sync"
	"testing"

	"github.com/diamcircle/go/txnbuild"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/stretchr/testify/assert"
)

func TestFriendbot_Pay(t *testing.T) {
	mockSubmitTransaction := func(minion *Minion, hclient auroraclient.ClientInterface, tx string) (*hProtocol.Transaction, error) {
		// Instead of submitting the tx, we emulate a success.
		txSuccess := hProtocol.Transaction{EnvelopeXdr: tx, Successful: true}
		return &txSuccess, nil
	}

	// Public key: GD25B4QI6KWVDWXDW25CIM7EKR6A6PBSWE2RCNSAC4NJQDQJXZJYMMKR
	botSeed := "SCWNLYELENPBXN46FHYXETT5LJCYBZD5VUQQVW4KZPHFO2YTQJUWT4D5"
	botKeypair, err := keypair.Parse(botSeed)
	if !assert.NoError(t, err) {
		return
	}
	botAccount := Account{AccountID: botKeypair.Address()}

	// Public key: GD4AGPPDFFHKK3Z2X4XZDRXX6GZQKP4FMLVQ5T55NDEYGG3GIP7BQUHM
	minionSeed := "SDTNSEERJPJFUE2LSDNYBFHYGVTPIWY7TU2IOJZQQGLWO2THTGB7NU5A"
	minionKeypair, err := keypair.Parse(minionSeed)
	if !assert.NoError(t, err) {
		return
	}

	minion := Minion{
		Account: Account{
			AccountID: minionKeypair.Address(),
			Sequence:  1,
		},
		Keypair:              minionKeypair.(*keypair.Full),
		BotAccount:           botAccount,
		BotKeypair:           botKeypair.(*keypair.Full),
		Network:              "Test SDF Network ; September 2015",
		StartingBalance:      "10000.00",
		SubmitTransaction:    mockSubmitTransaction,
		CheckSequenceRefresh: CheckSequenceRefresh,
		BaseFee:              txnbuild.MinBaseFee,
	}
	fb := &Bot{Minions: []Minion{minion}}

	recipientAddress := "GDJIN6W6PLTPKLLM57UW65ZH4BITUXUMYQHIMAZFYXF45PZVAWDBI77Z"
	txSuccess, err := fb.Pay(recipientAddress)
	if !assert.NoError(t, err) {
		return
	}
	expectedTxn := "AAAAAgAAAAD4Az3jKU6lbzq/L5HG9/GzBT+FYusOz71oyYMbZkP+GAAAAGQAAAAAAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPXQ8gjyrVHa47a6JDPkVHwPPDKxNRE2QBcamA4JvlOGAAAAAAAAAADShvreeub1LWzv6W93J+BROl6MxA6GAyXFy86/NQWGFAAAABdIdugAAAAAAAAAAAJmQ/4YAAAAQDRLEljDVYALnTk9mDceQEd5PrjQyE3LUAjstIyTWH5t/TP909F66TgEfBFKMxSKF6fka7ZuPcSs40ix4AomEgoJvlOGAAAAQPSGs88OwXubz7UT6nFhvhF47EQfaOsmiIsOkjgzUrmBoypJQTmMMbgeix0kdbfHqS75+iefJpdXLNFDreGnxgE="
	assert.Equal(t, expectedTxn, txSuccess.EnvelopeXdr)

	// Don't assert on tx values below, since the completion order is unknown.
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		_, err := fb.Pay(recipientAddress)
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		_, err := fb.Pay(recipientAddress)
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Wait()
}
