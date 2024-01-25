package integration

import (
	"math"
	"testing"

	"github.com/diamcircle/go/clients/auroraclient"
	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/protocols/aurora/effects"
	"github.com/diamcircle/go/protocols/aurora/operations"
	"github.com/diamcircle/go/services/aurora/internal/test/integration"
	"github.com/diamcircle/go/txnbuild"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

func TestMuxedAccountDetails(t *testing.T) {
	tt := assert.New(t)
	itest := integration.NewTest(t, integration.Config{ProtocolVersion: 17})
	master := itest.Master()
	masterStr := master.Address()
	masterAcID := xdr.MustAddress(masterStr)

	accs, _ := itest.CreateAccounts(1, "100")
	destionationStr := accs[0].Address()
	destinationAcID := xdr.MustAddress(destionationStr)

	source := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			Id:      0xcafebabecafebabe,
			Ed25519: *masterAcID.Ed25519,
		},
	}

	destination := xdr.MuxedAccount{
		Type: xdr.CryptoKeyTypeKeyTypeMuxedEd25519,
		Med25519: &xdr.MuxedAccountMed25519{
			// Make sure we cover the full uint64 range
			Id:      math.MaxUint64,
			Ed25519: *destinationAcID.Ed25519,
		},
	}

	// Submit a simple payment tx
	op := txnbuild.Payment{
		SourceAccount: source.Address(),
		Destination:   destination.Address(),
		Amount:        "10",
		Asset:         txnbuild.NativeAsset{},
	}

	txSource := itest.MasterAccount().(*hProtocol.Account)
	txSource.AccountID = source.Address()
	txResp := itest.MustSubmitOperations(txSource, master, &op)

	// check the transaction details
	txDetails, err := itest.Client().TransactionDetail(txResp.Hash)
	tt.NoError(err)
	tt.Equal(source.Address(), txDetails.AccountMuxed)
	tt.Equal(uint64(source.Med25519.Id), txDetails.AccountMuxedID)
	tt.Equal(source.Address(), txDetails.FeeAccountMuxed)
	tt.Equal(uint64(source.Med25519.Id), txDetails.FeeAccountMuxedID)

	// check the operation details
	opsResp, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: txResp.Hash,
	})
	tt.NoError(err)
	opDetails := opsResp.Embedded.Records[0].(operations.Payment)
	tt.Equal(source.Address(), opDetails.SourceAccountMuxed)
	tt.Equal(uint64(source.Med25519.Id), opDetails.SourceAccountMuxedID)
	tt.Equal(source.Address(), opDetails.FromMuxed)
	tt.Equal(uint64(source.Med25519.Id), opDetails.FromMuxedID)
	tt.Equal(destination.Address(), opDetails.ToMuxed)
	tt.Equal(uint64(destination.Med25519.Id), opDetails.ToMuxedID)

	// check the effect details
	effectsResp, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForTransaction: txResp.Hash,
	})
	tt.NoError(err)
	records := effectsResp.Embedded.Records

	credited := records[0].(effects.AccountCredited)
	tt.Equal(destination.Address(), credited.AccountMuxed)
	tt.Equal(uint64(destination.Med25519.Id), credited.AccountMuxedID)

	debited := records[1].(effects.AccountDebited)
	tt.Equal(source.Address(), debited.AccountMuxed)
	tt.Equal(uint64(source.Med25519.Id), debited.AccountMuxedID)
}
