package internal

import (
	"fmt"

	"go/clients/auroraclient"
	"go/keypair"
	hProtocol "go/protocols/aurora"
	"go/support/errors"
	"go/txnbuild"
)

const createAccountAlreadyExistXDR = "AAAAAAAAAGT/////AAAAAQAAAAAAAAAA/////AAAAAA="

var ErrAccountExists error = errors.New(fmt.Sprintf("createAccountAlreadyExist (%s)", createAccountAlreadyExistXDR))

// Minion contains a diamcircle channel account and Go channels to communicate with friendbot.
type Minion struct {
	Account         Account
	Keypair         *keypair.Full
	BotAccount      txnbuild.Account
	BotKeypair      *keypair.Full
	aurora          auroraclient.ClientInterface
	Network         string
	StartingBalance string
	BaseFee         int64

	// Mockable functions
	SubmitTransaction    func(minion *Minion, hclient auroraclient.ClientInterface, tx string) (*hProtocol.Transaction, error)
	CheckSequenceRefresh func(minion *Minion, hclient auroraclient.ClientInterface) error

	// Uninitialized.
	forceRefreshSequence bool
}

// Run reads a payment destination address and an output channel. It attempts
// to pay that address and submits the result to the channel.
func (minion *Minion) Run(destAddress string, resultChan chan SubmitResult) {
	err := minion.CheckSequenceRefresh(minion, minion.aurora)
	if err != nil {
		resultChan <- SubmitResult{
			maybeTransactionSuccess: nil,
			maybeErr:                errors.Wrap(err, "checking minion seq"),
		}
		return
	}
	txStr, err := minion.makeTx(destAddress)
	if err != nil {
		resultChan <- SubmitResult{
			maybeTransactionSuccess: nil,
			maybeErr:                errors.Wrap(err, "making payment tx"),
		}
		return
	}
	succ, err := minion.SubmitTransaction(minion, minion.aurora, txStr)
	resultChan <- SubmitResult{
		maybeTransactionSuccess: succ,
		maybeErr:                errors.Wrap(err, "submitting tx to minion"),
	}
}

// SubmitTransaction should be passed to the Minion.
func SubmitTransaction(minion *Minion, hclient auroraclient.ClientInterface, tx string) (*hProtocol.Transaction, error) {
	result, err := hclient.SubmitTransactionXDR(tx)
	if err != nil {
		errStr := "submitting tx to aurora"
		switch e := err.(type) {
		case *auroraclient.Error:
			minion.checkHandleBadSequence(e)
			resStr, resErr := e.ResultString()
			if resErr != nil {
				errStr += ": error getting aurora error code: " + resErr.Error()
			} else if resStr == createAccountAlreadyExistXDR {
				return nil, errors.Wrap(ErrAccountExists, errStr)
			} else {
				errStr += ": aurora error string: " + resStr
			}
			return nil, errors.New(errStr)
		}
		return nil, errors.Wrap(err, errStr)
	}
	return &result, nil
}

// CheckSequenceRefresh establishes the minion's initial sequence number, if needed.
// This should also be passed to the minion.
func CheckSequenceRefresh(minion *Minion, hclient auroraclient.ClientInterface) error {
	if minion.Account.Sequence != 0 && !minion.forceRefreshSequence {
		return nil
	}
	err := minion.Account.RefreshSequenceNumber(hclient)
	if err != nil {
		return errors.Wrap(err, "refreshing minion seqnum")
	}
	minion.forceRefreshSequence = false
	return nil
}

func (minion *Minion) checkHandleBadSequence(err *auroraclient.Error) {
	resCode, e := err.ResultCodes()
	isTxBadSeqCode := e == nil && resCode.TransactionCode == "tx_bad_seq"
	if !isTxBadSeqCode {
		return
	}
	minion.forceRefreshSequence = true
}

func (minion *Minion) makeTx(destAddress string) (string, error) {
	createAccountOp := txnbuild.CreateAccount{
		Destination:   destAddress,
		SourceAccount: minion.BotAccount.GetAccountID(),
		Amount:        minion.StartingBalance,
	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        minion.Account,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&createAccountOp},
			BaseFee:              minion.BaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "unable to build tx")
	}

	tx, err = tx.Sign(minion.Network, minion.Keypair, minion.BotKeypair)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	txe, err := tx.Base64()
	if err != nil {
		return "", errors.Wrap(err, "unable to serialize")
	}

	// Increment the in-memory sequence number, since the tx will be submitted.
	_, err = minion.Account.IncrementSequenceNumber()
	if err != nil {
		return "", errors.Wrap(err, "incrementing minion seq")
	}
	return txe, err
}
