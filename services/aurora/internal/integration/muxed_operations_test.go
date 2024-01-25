package integration

import (
	"testing"

	"go/clients/auroraclient"
	"go/keypair"
	"go/protocols/aurora/effects"
	"go/services/aurora/internal/test/integration"
	"go/txnbuild"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestMuxedOperations(t *testing.T) {
	itest := integration.NewTest(t, integration.Config{ProtocolVersion: 17})

	sponsored := keypair.MustRandom()
	// Is there an easier way?
	sponsoredMuxed := xdr.MustMuxedAddress(sponsored.Address())
	sponsoredMuxed.Type = xdr.CryptoKeyTypeKeyTypeMuxedEd25519
	sponsoredMuxed.Med25519 = &xdr.MuxedAccountMed25519{
		Ed25519: *sponsoredMuxed.Ed25519,
		Id:      100,
	}

	master := itest.Master()
	masterMuxed := xdr.MustMuxedAddress(master.Address())
	masterMuxed.Type = xdr.CryptoKeyTypeKeyTypeMuxedEd25519
	masterMuxed.Med25519 = &xdr.MuxedAccountMed25519{
		Ed25519: *masterMuxed.Ed25519,
		Id:      200,
	}

	ops := []txnbuild.Operation{
		&txnbuild.BeginSponsoringFutureReserves{
			SponsoredID: sponsored.Address(),
		},
		&txnbuild.CreateAccount{
			Destination: sponsored.Address(),
			Amount:      "100",
		},
		&txnbuild.ChangeTrust{
			SourceAccount: sponsoredMuxed.Address(),
			Line:          txnbuild.CreditAsset{"ABCD", master.Address()}.MustToChangeTrustAsset(),
			Limit:         txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.ManageSellOffer{
			SourceAccount: sponsoredMuxed.Address(),
			Selling:       txnbuild.NativeAsset{},
			Buying:        txnbuild.CreditAsset{"ABCD", master.Address()},
			Amount:        "3",
			Price:         "1",
		},
		// This will generate a trade effect:
		&txnbuild.ManageSellOffer{
			SourceAccount: masterMuxed.Address(),
			Selling:       txnbuild.CreditAsset{"ABCD", master.Address()},
			Buying:        txnbuild.NativeAsset{},
			Amount:        "3",
			Price:         "1",
		},
		&txnbuild.ManageData{
			SourceAccount: sponsoredMuxed.Address(),
			Name:          "test",
			Value:         []byte("test"),
		},
		&txnbuild.Payment{
			SourceAccount: sponsoredMuxed.Address(),
			Destination:   master.Address(),
			Amount:        "1",
			Asset:         txnbuild.NativeAsset{},
		},
		&txnbuild.CreateClaimableBalance{
			SourceAccount: sponsoredMuxed.Address(),
			Amount:        "2",
			Asset:         txnbuild.NativeAsset{},
			Destinations: []txnbuild.Claimant{
				txnbuild.NewClaimant(keypair.MustRandom().Address(), nil),
			},
		},
		&txnbuild.EndSponsoringFutureReserves{
			SourceAccount: sponsored.Address(),
		},
	}
	txResp, err := itest.SubmitMultiSigOperations(itest.MasterAccount(), []*keypair.Full{master, sponsored}, ops...)
	assert.NoError(t, err)
	assert.True(t, txResp.Successful)

	ops = []txnbuild.Operation{
		// Remove subentries to be able to merge account
		&txnbuild.Payment{
			SourceAccount: sponsoredMuxed.Address(),
			Destination:   master.Address(),
			Amount:        "3",
			Asset:         txnbuild.CreditAsset{"ABCD", master.Address()},
		},
		&txnbuild.ChangeTrust{
			SourceAccount: sponsoredMuxed.Address(),
			Line:          txnbuild.CreditAsset{"ABCD", master.Address()}.MustToChangeTrustAsset(),
			Limit:         "0",
		},
		&txnbuild.ManageData{
			SourceAccount: sponsoredMuxed.Address(),
			Name:          "test",
		},
		&txnbuild.AccountMerge{
			SourceAccount: sponsoredMuxed.Address(),
			Destination:   masterMuxed.Address(),
		},
	}
	txResp, err = itest.SubmitMultiSigOperations(itest.MasterAccount(), []*keypair.Full{master, sponsored}, ops...)
	assert.NoError(t, err)
	assert.True(t, txResp.Successful)

	// Check if no 5xx after processing the tx above
	// TODO expand it to test actual muxed fields
	_, err = itest.Client().Operations(auroraclient.OperationRequest{Limit: 200})
	assert.NoError(t, err, "/operations failed")

	_, err = itest.Client().Payments(auroraclient.OperationRequest{Limit: 200})
	assert.NoError(t, err, "/payments failed")

	effectsPage, err := itest.Client().Effects(auroraclient.EffectRequest{Limit: 200})
	assert.NoError(t, err, "/effects failed")

	for _, effect := range effectsPage.Embedded.Records {
		if effect.GetType() == "trade" {
			trade := effect.(effects.Trade)
			oneSet := trade.AccountMuxedID != 0 || trade.SellerMuxedID != 0
			assert.True(t, oneSet, "at least one of account_muxed_id, seller_muxed_id must be set")
		}
	}
}
