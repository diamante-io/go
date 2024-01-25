package ingest

import (
	"context"
	"encoding/hex"
	"io"

	"github.com/diamcircle/go/ingest/ledgerbackend"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// LedgerTransactionReader reads transactions for a given ledger sequence from a backend.
// Use NewTransactionReader to create a new instance.
type LedgerTransactionReader struct {
	ledgerCloseMeta xdr.LedgerCloseMeta
	transactions    []LedgerTransaction
	readIdx         int
}

// NewLedgerTransactionReader creates a new TransactionReader instance.
// Note that TransactionReader is not thread safe and should not be shared by multiple goroutines.
func NewLedgerTransactionReader(ctx context.Context, backend ledgerbackend.LedgerBackend, networkPassphrase string, sequence uint32) (*LedgerTransactionReader, error) {
	ledgerCloseMeta, err := backend.GetLedger(ctx, sequence)
	if err != nil {
		return nil, errors.Wrap(err, "error getting ledger from the backend")
	}

	return NewLedgerTransactionReaderFromLedgerCloseMeta(networkPassphrase, ledgerCloseMeta)
}

// NewLedgerTransactionReaderFromXdr creates a new TransactionReader instance from xdr.LedgerCloseMeta.
// Note that TransactionReader is not thread safe and should not be shared by multiple goroutines.
func NewLedgerTransactionReaderFromLedgerCloseMeta(networkPassphrase string, ledgerCloseMeta xdr.LedgerCloseMeta) (*LedgerTransactionReader, error) {
	reader := &LedgerTransactionReader{ledgerCloseMeta: ledgerCloseMeta}
	if err := reader.storeTransactions(ledgerCloseMeta, networkPassphrase); err != nil {
		return nil, errors.Wrap(err, "error extracting transactions from ledger close meta")
	}
	return reader, nil
}

// GetSequence returns the sequence number of the ledger data stored by this object.
func (reader *LedgerTransactionReader) GetSequence() uint32 {
	return reader.ledgerCloseMeta.LedgerSequence()
}

// GetHeader returns the XDR Header data associated with the stored ledger.
func (reader *LedgerTransactionReader) GetHeader() xdr.LedgerHeaderHistoryEntry {
	return reader.ledgerCloseMeta.V0.LedgerHeader
}

// Read returns the next transaction in the ledger, ordered by tx number, each time
// it is called. When there are no more transactions to return, an EOF error is returned.
func (reader *LedgerTransactionReader) Read() (LedgerTransaction, error) {
	if reader.readIdx < len(reader.transactions) {
		reader.readIdx++
		return reader.transactions[reader.readIdx-1], nil
	}
	return LedgerTransaction{}, io.EOF
}

// Rewind resets the reader back to the first transaction in the ledger
func (reader *LedgerTransactionReader) Rewind() {
	reader.readIdx = 0
}

// storeTransactions maps the close meta data into a slice of LedgerTransaction structs, to provide
// a per-transaction view of the data when Read() is called.
func (reader *LedgerTransactionReader) storeTransactions(lcm xdr.LedgerCloseMeta, networkPassphrase string) error {
	byHash := map[xdr.Hash]xdr.TransactionEnvelope{}
	for i, tx := range lcm.V0.TxSet.Txs {
		hash, err := network.HashTransactionInEnvelope(tx, networkPassphrase)
		if err != nil {
			return errors.Wrapf(err, "could not hash transaction %d in TxSet", i)
		}
		byHash[hash] = tx
	}

	for i := range lcm.V0.TxProcessing {
		result := lcm.V0.TxProcessing[i].Result
		envelope, ok := byHash[result.TransactionHash]
		if !ok {
			hexHash := hex.EncodeToString(result.TransactionHash[:])
			return errors.Errorf("unknown tx hash in LedgerCloseMeta: %v", hexHash)
		}

		// We check the version only if FeeProcessing are non empty because some backends
		// (like HistoryArchiveBackend) do not return meta.
		if lcm.V0.LedgerHeader.Header.LedgerVersion < 10 && lcm.V0.TxProcessing[i].TxApplyProcessing.V != 2 &&
			len(lcm.V0.TxProcessing[i].FeeProcessing) > 0 {
			return errors.New(
				"TransactionMeta.V=2 is required in protocol version older than version 10. " +
					"Please process ledgers again using the latest diamcircle-core version.",
			)
		}

		reader.transactions = append(reader.transactions, LedgerTransaction{
			Index:      uint32(i + 1), // Transactions start at '1'
			Envelope:   envelope,
			Result:     result,
			UnsafeMeta: lcm.V0.TxProcessing[i].TxApplyProcessing,
			FeeChanges: lcm.V0.TxProcessing[i].FeeProcessing,
		})
	}
	return nil
}

// Close should be called when reading is finished. This is especially
// helpful when there are still some transactions available so reader can stop
// streaming them.
func (reader *LedgerTransactionReader) Close() error {
	reader.transactions = nil
	return nil
}
