//Package transactions offers common infrastructure for testing Transactions
package transactions

import (
	"encoding/hex"
	"testing"

	"github.com/diamcircle/go/ingest"
	"github.com/diamcircle/go/xdr"
	"github.com/stretchr/testify/assert"
)

// TestTransaction transaction meta
type TestTransaction struct {
	Index         uint32
	EnvelopeXDR   string
	ResultXDR     string
	FeeChangesXDR string
	MetaXDR       string
	Hash          string
}

// BuildLedgerTransaction builds a ledger transaction
func BuildLedgerTransaction(t *testing.T, tx TestTransaction) ingest.LedgerTransaction {
	transaction := ingest.LedgerTransaction{
		Index:      tx.Index,
		Envelope:   xdr.TransactionEnvelope{},
		Result:     xdr.TransactionResultPair{},
		FeeChanges: xdr.LedgerEntryChanges{},
		UnsafeMeta: xdr.TransactionMeta{},
	}

	tt := assert.New(t)

	err := xdr.SafeUnmarshalBase64(tx.EnvelopeXDR, &transaction.Envelope)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.ResultXDR, &transaction.Result.Result)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.MetaXDR, &transaction.UnsafeMeta)
	tt.NoError(err)
	err = xdr.SafeUnmarshalBase64(tx.FeeChangesXDR, &transaction.FeeChanges)
	tt.NoError(err)

	_, err = hex.Decode(transaction.Result.TransactionHash[:], []byte(tx.Hash))
	tt.NoError(err)

	return transaction
}
