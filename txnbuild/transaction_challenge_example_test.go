package txnbuild_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/txnbuild"
)

var serverAccount, _ = keypair.ParseFull("SCDXPYDGKV5HOAGVZN3FQSS5FKUPP5BAVBWH4FXKTAWAC24AE4757JSI")
var clientAccount, _ = keypair.ParseFull("SANVNCABRBVISCV7KH4SZVBKPJWWTT4424OVWUHUHPH2MVSF6RC7HPGN")
var clientSigner1, _ = keypair.ParseFull("SBPQUZ6G4FZNWFHKUWC5BEYWF6R52E3SEP7R3GWYSM2XTKGF5LNTWW4R")
var clientSigner2, _ = keypair.ParseFull("SBMSVD4KKELKGZXHBUQTIROWUAPQASDX7KEJITARP4VMZ6KLUHOGPTYW")
var auroraClient = func() auroraclient.ClientInterface {
	client := &auroraclient.MockClient{}
	client.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: clientAccount.Address()}).
		Return(
			aurora.Account{
				Thresholds: aurora.AccountThresholds{LowThreshold: 1, MedThreshold: 10, HighThreshold: 100},
				Signers: []aurora.Signer{
					{Key: clientSigner1.Address(), Weight: 40},
					{Key: clientSigner2.Address(), Weight: 60},
				},
			},
			nil,
		)
	return client
}()

func ExampleVerifyChallengeTxThreshold() {
	// Server builds challenge transaction
	var challengeTx string
	{
		tx, err := txnbuild.BuildChallengeTx(serverAccount.Seed(), clientAccount.Address(), "webauthdomain.diamcircle.org", "test", network.TestNetworkPassphrase, time.Minute)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		challengeTx, err = tx.Base64()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Client reads and signs challenge transaction
	var signedChallengeTx string
	{
		tx, txClientAccountID, _, err := txnbuild.ReadChallengeTx(challengeTx, serverAccount.Address(), network.TestNetworkPassphrase, "webauthdomain.diamcircle.org", []string{"test"})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if txClientAccountID != clientAccount.Address() {
			fmt.Println("Error: challenge tx is not for expected client account")
			return
		}
		tx, err = tx.Sign(network.TestNetworkPassphrase, clientSigner1, clientSigner2)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		signedChallengeTx, err = tx.Base64()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Server verifies signed challenge transaction
	{
		_, txClientAccountID, _, err := txnbuild.ReadChallengeTx(challengeTx, serverAccount.Address(), network.TestNetworkPassphrase, "webauthdomain.diamcircle.org", []string{"test"})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Server gets account
		clientAccountExists := false
		auroraClientAccount, err := auroraClient.AccountDetail(auroraclient.AccountRequest{AccountID: txClientAccountID})
		if auroraclient.IsNotFoundError(err) {
			clientAccountExists = false
			fmt.Println("Account does not exist, use master key to verify")
		} else if err == nil {
			clientAccountExists = true
		} else {
			fmt.Println("Error:", err)
			return
		}

		if clientAccountExists {
			// Server gets list of signers from account
			signerSummary := auroraClientAccount.SignerSummary()

			// Server chooses the threshold to require: low, med or high
			threshold := txnbuild.Threshold(auroraClientAccount.Thresholds.MedThreshold)

			// Server verifies threshold is met
			signers, err := txnbuild.VerifyChallengeTxThreshold(signedChallengeTx, serverAccount.Address(), network.TestNetworkPassphrase, "webauthdomain.diamcircle.org", []string{"test"}, threshold, signerSummary)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Client Signers Verified:")
			sort.Strings(signers)
			for _, signer := range signers {
				fmt.Println(signer, "weight:", signerSummary[signer])
			}
		} else {
			// Server verifies that master key has signed challenge transaction
			signersFound, err := txnbuild.VerifyChallengeTxSigners(signedChallengeTx, serverAccount.Address(), network.TestNetworkPassphrase, "webauthdomain.diamcircle.org", []string{"test"}, txClientAccountID)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Client Master Key Verified:")
			for _, signerFound := range signersFound {
				fmt.Println(signerFound)
			}
		}
	}

	// Output:
	// Client Signers Verified:
	// GAS4V4O2B7DW5T7IQRPEEVCRXMDZESKISR7DVIGKZQYYV3OSQ5SH5LVP weight: 60
	// GDQNY3PBOJOKYZSRMK2S7LHHGWZIUISD4QORETLMXEWXBI7KFZZMKTL3 weight: 40
}
