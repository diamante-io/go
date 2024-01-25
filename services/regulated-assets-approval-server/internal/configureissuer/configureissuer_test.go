package configureissuer

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetup_accountAlreadyConfigured(t *testing.T) {
	// declare a logging buffer to validate output logs
	buf := new(strings.Builder)
	log.DefaultLogger.SetOutput(buf)
	log.DefaultLogger.SetLevel(log.InfoLevel)

	issuerKP := keypair.MustRandom()
	opts := Options{
		AssetCode:           "FOO",
		BaseURL:             "https://domain.test.com/",
		AuroraURL:          auroraclient.DefaultTestNetClient.AuroraURL,
		IssuerAccountSecret: issuerKP.Seed(),
		NetworkPassphrase:   network.TestNetworkPassphrase,
	}

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: issuerKP.Address()}).
		Return(aurora.Account{
			AccountID: issuerKP.Address(),
			Flags: aurora.AccountFlags{
				AuthRequired:  true,
				AuthRevocable: true,
			},
			HomeDomain: "domain.test.com",
			Sequence:   "10",
		}, nil)
	auroraMock.
		On("Assets", auroraclient.AssetRequest{
			ForAssetCode:   opts.AssetCode,
			ForAssetIssuer: issuerKP.Address(),
			Limit:          1,
		}).
		Return(aurora.AssetsPage{
			Embedded: struct{ Records []aurora.AssetStat }{
				Records: []aurora.AssetStat{
					{Amount: "0.0000001"},
				},
			},
		}, nil)

	err := setup(opts, &auroraMock)
	require.NoError(t, err)

	require.Contains(t, buf.String(), "Account already configured. Aborting without performing any action.")
}

func TestGetOrFundIssuerAccount_failsIfNotDefaultTesntet(t *testing.T) {
	issuerKP := keypair.MustRandom()

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: issuerKP.Address()}).
		Return(aurora.Account{}, problem.NotFound)

	_, err := getOrFundIssuerAccount(issuerKP.Address(), &auroraMock)
	wantErrMsg := fmt.Sprintf("getting detail for account %s: problem: not_found", issuerKP.Address())
	require.EqualError(t, err, wantErrMsg)
}

func TestSetup(t *testing.T) {
	issuerKP := keypair.MustRandom()
	opts := Options{
		AssetCode:           "FOO",
		BaseURL:             "https://domain.test.com/",
		AuroraURL:          auroraclient.DefaultTestNetClient.AuroraURL,
		IssuerAccountSecret: issuerKP.Seed(),
		NetworkPassphrase:   network.TestNetworkPassphrase,
	}

	auroraMock := auroraclient.MockClient{}
	auroraMock.
		On("AccountDetail", auroraclient.AccountRequest{AccountID: issuerKP.Address()}).
		Return(aurora.Account{
			AccountID: issuerKP.Address(),
			Sequence:  "10",
		}, nil)
	auroraMock.
		On("Assets", auroraclient.AssetRequest{
			ForAssetCode:   opts.AssetCode,
			ForAssetIssuer: issuerKP.Address(),
			Limit:          1,
		}).
		Return(aurora.AssetsPage{}, nil)

	var didTestSubmitTransaction bool
	auroraMock.
		On("SubmitTransaction", mock.AnythingOfType("*txnbuild.Transaction")).
		Run(func(args mock.Arguments) {
			tx, ok := args.Get(0).(*txnbuild.Transaction)
			require.True(t, ok)

			issuerSimpleAcc := txnbuild.SimpleAccount{
				AccountID: issuerKP.Address(),
				Sequence:  11,
			}
			assert.Equal(t, issuerSimpleAcc, tx.SourceAccount())

			assert.Equal(t, int64(11), tx.SequenceNumber())
			assert.Equal(t, int64(300), tx.BaseFee())
			assert.Equal(t, int64(0), tx.Timebounds().MinTime)
			assert.LessOrEqual(t, time.Now().UTC().Unix()+299, tx.Timebounds().MaxTime)
			assert.GreaterOrEqual(t, time.Now().UTC().Unix()+301, tx.Timebounds().MaxTime)

			beginSponsorOp, ok := tx.Operations()[1].(*txnbuild.BeginSponsoringFutureReserves)
			require.True(t, ok)
			trustorAccKP := beginSponsorOp.SponsoredID
			homeDomain := "domain.test.com"
			testAsset := txnbuild.CreditAsset{
				Code:   opts.AssetCode,
				Issuer: issuerKP.Address(),
			}

			wantOps := []txnbuild.Operation{
				&txnbuild.SetOptions{
					SetFlags: []txnbuild.AccountFlag{
						txnbuild.AuthRequired,
						txnbuild.AuthRevocable,
					},
					HomeDomain: &homeDomain,
				},
				&txnbuild.BeginSponsoringFutureReserves{
					SponsoredID:   trustorAccKP,
					SourceAccount: issuerKP.Address(),
				},
				&txnbuild.CreateAccount{
					Destination:   trustorAccKP,
					Amount:        "0",
					SourceAccount: issuerKP.Address(),
				},
				// a trustline is generated to the desired so aurora creates entry at `{aurora-url}/assets`. This was added as many Wallets reach that endpoint to check if a given asset exists.
				&txnbuild.ChangeTrust{
					Line:          testAsset.MustToChangeTrustAsset(),
					SourceAccount: trustorAccKP,
					Limit:         "922337203685.4775807",
				},
				&txnbuild.SetOptions{
					MasterWeight:    txnbuild.NewThreshold(0),
					LowThreshold:    txnbuild.NewThreshold(1),
					MediumThreshold: txnbuild.NewThreshold(1),
					HighThreshold:   txnbuild.NewThreshold(1),
					Signer:          &txnbuild.Signer{Address: issuerKP.Address(), Weight: txnbuild.Threshold(10)},
					SourceAccount:   trustorAccKP,
				},
				&txnbuild.EndSponsoringFutureReserves{
					SourceAccount: trustorAccKP,
				},
			}
			// SetOptions operation is validated separatedly because the value returned from tx.Operations()[0] contains the unexported field `xdrOp` that prevents a proper comparision.
			require.Equal(t, wantOps[0].(*txnbuild.SetOptions).SetFlags, tx.Operations()[0].(*txnbuild.SetOptions).SetFlags)
			require.Equal(t, wantOps[0].(*txnbuild.SetOptions).HomeDomain, tx.Operations()[0].(*txnbuild.SetOptions).HomeDomain)

			require.Equal(t, wantOps[1:4], tx.Operations()[1:4])

			// SetOptions operation is validated separatedly because the value returned from tx.Operations()[4] contains the unexported field `xdrOp` that prevents a proper comparision.
			require.Equal(t, wantOps[4].(*txnbuild.SetOptions).SetFlags, tx.Operations()[4].(*txnbuild.SetOptions).SetFlags)
			require.Equal(t, wantOps[4].(*txnbuild.SetOptions).HomeDomain, tx.Operations()[4].(*txnbuild.SetOptions).HomeDomain)

			require.Equal(t, wantOps[5:], tx.Operations()[5:])

			txHash, err := tx.Hash(opts.NetworkPassphrase)
			require.NoError(t, err)

			err = issuerKP.Verify(txHash[:], tx.Signatures()[0].Signature)
			require.NoError(t, err)

			trustorKP, err := keypair.ParseAddress(trustorAccKP)
			require.NoError(t, err)
			err = trustorKP.Verify(txHash[:], tx.Signatures()[1].Signature)
			require.NoError(t, err)

			didTestSubmitTransaction = true
		}).
		Return(aurora.Transaction{}, nil)

	err := setup(opts, &auroraMock)
	require.NoError(t, err)

	require.True(t, didTestSubmitTransaction)
}
