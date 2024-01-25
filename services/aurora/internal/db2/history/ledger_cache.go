package history

import (
	"context"

	"github.com/diamcircle/go/support/errors"
)

// Queue adds `seq` to the load queue for the cache.
func (lc *LedgerCache) Queue(seq int32) {
	lc.lock.Lock()

	if lc.queued == nil {
		lc.queued = map[int32]struct{}{}
	}

	lc.queued[seq] = struct{}{}
	lc.lock.Unlock()
}

// Load loads a batch of ledgers identified by `sequences`, using `q`,
// and populates the cache with the results
func (lc *LedgerCache) Load(ctx context.Context, q *Q) error {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	if len(lc.queued) == 0 {
		return nil
	}

	sequences := make([]int32, 0, len(lc.queued))
	for seq := range lc.queued {
		sequences = append(sequences, seq)
	}

	var ledgers []Ledger
	err := q.LedgersBySequence(ctx, &ledgers, sequences...)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger batch")
	}

	lc.Records = map[int32]Ledger{}
	for _, l := range ledgers {
		lc.Records[l.Sequence] = l
	}

	lc.queued = nil
	return nil
}
