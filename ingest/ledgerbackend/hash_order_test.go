package ledgerbackend

import (
	"testing"

	"go/network"
	"go/xdr"

	"github.com/stretchr/testify/assert"
)

func TestHashOrder(t *testing.T) {
	source := xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU")
	account := source.ToMuxedAccount()
	original := []xdr.TransactionEnvelope{
		{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: account,
					SeqNum:        1,
				},
			},
		},
		{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: account,
					SeqNum:        2,
				},
			},
		},
		{
			Type: xdr.EnvelopeTypeEnvelopeTypeTx,
			V1: &xdr.TransactionV1Envelope{
				Tx: xdr.Transaction{
					SourceAccount: account,
					SeqNum:        3,
				},
			},
		},
	}

	sortByHash(original, network.TestNetworkPassphrase)
	hashes := map[int]xdr.Hash{}

	for i, tx := range original {
		var err error
		hashes[i], err = network.HashTransactionInEnvelope(tx, network.TestNetworkPassphrase)
		if err != nil {
			assert.NoError(t, err)
		}
	}

	for i := range original {
		if i == 0 {
			continue
		}
		prev := hashes[i-1]
		cur := hashes[i]
		for j := range prev {
			if !assert.True(t, prev[j] < cur[j]) {
				break
			} else {
				break
			}
		}
	}
}
