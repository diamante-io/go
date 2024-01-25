package ingest

import (
	"context"
	"io"

	"go/ingest/ledgerbackend"
	"go/xdr"
)

// ChangeReader provides convenient, streaming access to a sequence of Changes.
type ChangeReader interface {
	// Read should return the next `Change` in the leader. If there are no more
	// changes left it should return an `io.EOF` error.
	Read() (Change, error)
	// Close should be called when reading is finished. This is especially
	// helpful when there are still some changes available so reader can stop
	// streaming them.
	Close() error
}

// ledgerChangeReaderState defines possible states of LedgerChangeReader.
type ledgerChangeReaderState int

const (
	// feeChangesState is active when LedgerChangeReader is reading fee changes.
	feeChangesState ledgerChangeReaderState = iota
	// feeChangesState is active when LedgerChangeReader is reading transaction meta changes.
	metaChangesState
	// feeChangesState is active when LedgerChangeReader is reading upgrade changes.
	upgradeChangesState
)

// LedgerChangeReader is a ChangeReader which returns Changes from diamcircle Core
// for a single ledger
type LedgerChangeReader struct {
	*LedgerTransactionReader
	state        ledgerChangeReaderState
	pending      []Change
	pendingIndex int
	upgradeIndex int
}

// Ensure LedgerChangeReader implements ChangeReader
var _ ChangeReader = (*LedgerChangeReader)(nil)

// NewLedgerChangeReader constructs a new LedgerChangeReader instance bound to the given ledger.
// Note that the returned LedgerChangeReader is not thread safe and should not be shared
// by multiple goroutines.
func NewLedgerChangeReader(ctx context.Context, backend ledgerbackend.LedgerBackend, networkPassphrase string, sequence uint32) (*LedgerChangeReader, error) {
	transactionReader, err := NewLedgerTransactionReader(ctx, backend, networkPassphrase, sequence)
	if err != nil {
		return nil, err
	}

	return &LedgerChangeReader{
		LedgerTransactionReader: transactionReader,
		state:                   feeChangesState,
	}, nil
}

// NewLedgerChangeReaderFromLedgerCloseMeta constructs a new LedgerChangeReader instance bound to the given ledger.
// Note that the returned LedgerChangeReader is not thread safe and should not be shared
// by multiple goroutines.
func NewLedgerChangeReaderFromLedgerCloseMeta(networkPassphrase string, ledger xdr.LedgerCloseMeta) (*LedgerChangeReader, error) {
	transactionReader, err := NewLedgerTransactionReaderFromLedgerCloseMeta(networkPassphrase, ledger)
	if err != nil {
		return nil, err
	}

	return &LedgerChangeReader{
		LedgerTransactionReader: transactionReader,
		state:                   feeChangesState,
	}, nil
}

// Read returns the next change in the stream.
// If there are no changes remaining io.EOF is returned as an error.
func (r *LedgerChangeReader) Read() (Change, error) {
	// Changes within a ledger should be read in the following order:
	// - fee changes of all transactions,
	// - transaction meta changes of all transactions,
	// - upgrade changes.
	// Because a single transaction can introduce many changes we read all the
	// changes from a single transaction  and save them in r.pending.
	// When Read() is called we stream pending changes first. We also call Read()
	// recursively after adding some changes (what will return them from r.pending)
	// to not duplicate the code.
	if r.pendingIndex < len(r.pending) {
		next := r.pending[r.pendingIndex]
		r.pendingIndex++
		if r.pendingIndex == len(r.pending) {
			r.pendingIndex = 0
			r.pending = r.pending[:0]
		}
		return next, nil
	}

	switch r.state {
	case feeChangesState, metaChangesState:
		tx, err := r.LedgerTransactionReader.Read()
		if err != nil {
			if err == io.EOF {
				// If done streaming fee changes rewind to stream meta changes
				if r.state == feeChangesState {
					r.LedgerTransactionReader.Rewind()
				}
				r.state++
				return r.Read()
			}
			return Change{}, err
		}

		switch r.state {
		case feeChangesState:
			r.pending = append(r.pending, tx.GetFeeChanges()...)
		case metaChangesState:
			metaChanges, err := tx.GetChanges()
			if err != nil {
				return Change{}, err
			}
			r.pending = append(r.pending, metaChanges...)
		}
		return r.Read()
	case upgradeChangesState:
		// Get upgrade changes
		if r.upgradeIndex < len(r.LedgerTransactionReader.ledgerCloseMeta.V0.UpgradesProcessing) {
			changes := GetChangesFromLedgerEntryChanges(
				r.LedgerTransactionReader.ledgerCloseMeta.V0.UpgradesProcessing[r.upgradeIndex].Changes,
			)
			r.pending = append(r.pending, changes...)
			r.upgradeIndex++
			return r.Read()
		}
	}

	return Change{}, io.EOF
}

// Close should be called when reading is finished.
func (r *LedgerChangeReader) Close() error {
	r.pending = nil
	return r.LedgerTransactionReader.Close()
}
