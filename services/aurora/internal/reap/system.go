package reap

import (
	"context"
	"time"

	herrors "github.com/diamcircle/go/services/aurora/internal/errors"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/log"
)

// DeleteUnretainedHistory removes all data associated with unretained ledgers.
func (r *System) DeleteUnretainedHistory(ctx context.Context) error {
	// RetentionCount of 0 indicates "keep all history"
	if r.RetentionCount == 0 {
		return nil
	}

	var (
		latest      = r.ledgerState.CurrentStatus()
		targetElder = (latest.HistoryLatest - int32(r.RetentionCount)) + 1
	)

	if targetElder < latest.HistoryElder {
		return nil
	}

	err := r.clearBefore(ctx, latest.HistoryElder, targetElder)
	if err != nil {
		return err
	}

	log.
		WithField("new_elder", targetElder).
		Info("reaper succeeded")

	return nil
}

// Run triggers the reaper system to update itself, deleted unretained history
// if it is the appropriate time.
func (r *System) Run() {
	for {
		select {
		case <-time.After(1 * time.Hour):
			r.runOnce(r.ctx)
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *System) Shutdown() {
	r.cancel()
}

func (r *System) runOnce(ctx context.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			err := herrors.FromPanic(rec)
			log.Errorf("reaper panicked: %s", err)
			herrors.ReportToSentry(err, nil)
		}
	}()

	err := r.DeleteUnretainedHistory(ctx)
	if err != nil {
		log.Errorf("reaper failed: %s", err)
	}
}

// Work backwards in 100k ledger blocks to prevent using all the CPU.
//
// This runs every hour, so we need to make sure it doesn't
// run for longer than an hour.
//
// Current ledger at 2021-08-12 is 36,827,497, so 100k means 368 batches. At 1
// batch/second, that seems like a reasonable balance between running well
// under an hour, and slowing it down enough to leave some CPU for other
// processes.
var batchSize = int32(100_000)
var sleep = 1 * time.Second

func (r *System) clearBefore(ctx context.Context, startSeq, endSeq int32) error {
	for batchEndSeq := endSeq - 1; batchEndSeq >= startSeq; batchEndSeq -= batchSize {
		batchStartSeq := batchEndSeq - batchSize
		if batchStartSeq < startSeq {
			batchStartSeq = startSeq
		}
		log.WithField("start_ledger", batchStartSeq).WithField("end_ledger", batchEndSeq).Info("reaper: clearing")

		batchStart, batchEnd, err := toid.LedgerRangeInclusive(batchStartSeq, batchEndSeq)
		if err != nil {
			return err
		}

		err = r.HistoryQ.Begin()
		if err != nil {
			return errors.Wrap(err, "Error in begin")
		}
		defer r.HistoryQ.Rollback()

		err = r.HistoryQ.DeleteRangeAll(ctx, batchStart, batchEnd)
		if err != nil {
			return errors.Wrap(err, "Error in DeleteRangeAll")
		}

		err = r.HistoryQ.Commit()
		if err != nil {
			return errors.Wrap(err, "Error in commit")
		}

		time.Sleep(sleep)
	}

	return nil
}
