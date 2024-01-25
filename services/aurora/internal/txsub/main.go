package txsub

import (
	"time"

	"context"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/xdr"
)

// Listener represents some client who is interested in retrieving the result
// of a specific transaction.
type Listener chan<- Result

// OpenSubmissionList represents the structure that tracks pending transactions
// and forwards Result structs on to listeners as they become available.
//
// NOTE:  An implementation of this interface will be called from multiple go-routines
// concurrently.
//
// NOTE:  A Listener must be a buffered channel.  A panic will trigger if you
// provide an unbuffered channel
type OpenSubmissionList interface {
	// Add registers the provided listener as interested in being notified when a
	// result is available for the provided transaction hash.
	Add(context.Context, string, Listener) error

	// Finish forwards the provided result on to any listeners and cleans up any
	// resources associated with the transaction that this result is for
	Finish(context.Context, string, Result) error

	// Clean removes any open submissions over the provided age.
	Clean(context.Context, time.Duration) (int, error)

	// Pending return a list of transaction hashes that have at least one
	// listener registered to them in this list.
	Pending(context.Context) []string
}

// Submitter represents the low-level "submit a transaction to diamcircle-core"
// provider.
type Submitter interface {
	// Submit sends the provided transaction envelope to diamcircle-core
	Submit(context.Context, string) SubmissionResult
}

// Result represents the response from a ResultProvider.  Given no
// Err is set, the rest of the struct should be populated appropriately.
type Result struct {
	// Any error that occurred during the retrieval of this result
	Err error

	// The full details of the transaction which was submitted
	// to Diamcircle Core
	Transaction history.Transaction
}

// SubmissionResult gets returned in response to a call to Submitter.Submit.
// It represents a single discrete submission of a transaction envelope to
// the diamcircle network.
type SubmissionResult struct {
	// Any error that occurred during the attempted submission.  A nil value
	// indicates that the submission will or already is being considered for
	// inclusion in the ledger (i.e. A successful submission).
	Err error

	// Duration records the time it took to submit a transaction
	// to diamcircle-core
	Duration time.Duration
}

func (s SubmissionResult) IsBadSeq() (bool, error) {
	if s.Err == nil {
		return false, nil
	}

	fte, ok := s.Err.(*FailedTransactionError)
	if !ok {
		return false, nil
	}

	result, err := fte.Result()
	if err != nil {
		return false, err
	}

	return result.Result.Code == xdr.TransactionResultCodeTxBadSeq, nil
}
