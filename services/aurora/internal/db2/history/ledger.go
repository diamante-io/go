package history

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"go/services/aurora/internal/db2"
	"go/services/aurora/internal/toid"
	"go/support/errors"
	"go/xdr"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"
)

// LedgerBySequence loads the single ledger at `seq` into `dest`
func (q *Q) LedgerBySequence(ctx context.Context, dest interface{}, seq int32) error {
	sql := selectLedger.
		Limit(1).
		Where("sequence = ?", seq)

	return q.Get(ctx, dest, sql)
}

// Ledgers provides a helper to filter rows from the `history_ledgers` table
// with pre-defined filters.  See `LedgersQ` methods for the available filters.
func (q *Q) Ledgers() *LedgersQ {
	return &LedgersQ{
		parent: q,
		sql:    selectLedger,
	}
}

// LedgersBySequence loads the a set of ledgers identified by the sequences
// `seqs` into `dest`.
func (q *Q) LedgersBySequence(ctx context.Context, dest interface{}, seqs ...int32) error {
	if len(seqs) == 0 {
		return errors.New("no sequence arguments provided")
	}
	in := fmt.Sprintf("sequence IN (%s)", sq.Placeholders(len(seqs)))

	whereArgs := make([]interface{}, len(seqs))
	for i, s := range seqs {
		whereArgs[i] = s
	}

	sql := selectLedger.Where(in, whereArgs...)

	return q.Select(ctx, dest, sql)
}

// LedgerCapacityUsageStats returns ledger capacity stats for the last 5 ledgers.
// Currently, we hard code the query to return the last 5 ledgers.
// TODO: make the number of ledgers configurable.
func (q *Q) LedgerCapacityUsageStats(ctx context.Context, currentSeq int32, dest *LedgerCapacityUsageStats) error {
	const ledgers int32 = 5
	return q.GetRaw(ctx, dest, `
		SELECT ROUND(SUM(CAST(operation_count as decimal))/SUM(max_tx_set_size), 2) as ledger_capacity_usage FROM
			(SELECT
			  hl.sequence, COALESCE(SUM(ht.operation_count), 0) as operation_count, hl.max_tx_set_size
			  FROM history_ledgers hl
			  LEFT JOIN history_transactions ht ON ht.ledger_sequence = hl.sequence
			  WHERE hl.sequence > $1 AND hl.sequence <= $2
			  GROUP BY hl.sequence, hl.max_tx_set_size) as a
	`, currentSeq-ledgers, currentSeq)
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *LedgersQ) Page(page db2.PageQuery) *LedgersQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "hl.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *LedgersQ) Select(ctx context.Context, dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(ctx, dest, q.sql)
	return q.Err
}

// QLedgers defines ingestion ledger related queries.
type QLedgers interface {
	InsertLedger(
		ctx context.Context,
		ledger xdr.LedgerHeaderHistoryEntry,
		successTxsCount int,
		failedTxsCount int,
		opCount int,
		txSetOpCount int,
		ingestVersion int,
	) (int64, error)
}

// InsertLedger creates a row in the history_ledgers table.
// Returns number of rows affected and error.
func (q *Q) InsertLedger(ctx context.Context,
	ledger xdr.LedgerHeaderHistoryEntry,
	successTxsCount int,
	failedTxsCount int,
	opCount int,
	txSetOpCount int,
	ingestVersion int,
) (int64, error) {
	m, err := ledgerHeaderToMap(
		ledger,
		successTxsCount,
		failedTxsCount,
		opCount,
		txSetOpCount,
		ingestVersion,
	)
	if err != nil {
		return 0, err
	}

	sql := sq.Insert("history_ledgers").SetMap(m)
	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// GetLedgerGaps obtains ingestion gaps in the history_ledgers table.
// Returns the gaps and error.
func (q *Q) GetLedgerGaps(ctx context.Context) ([]LedgerRange, error) {
	var gaps []LedgerRange
	query := `
    SELECT sequence + 1 AS start,
		next_number - 1 AS end
	FROM (
		SELECT sequence,
		LEAD(sequence) OVER (ORDER BY sequence) AS next_number
	FROM history_ledgers
	) number
	WHERE sequence + 1 <> next_number;`
	if err := q.SelectRaw(ctx, &gaps, query); err != nil {
		return nil, err
	}
	sort.Slice(gaps, func(i, j int) bool {
		return gaps[i].StartSequence < gaps[j].StartSequence
	})
	return gaps, nil
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint32) uint32 {
	if a > b {
		return b
	}
	return a
}

// GetLedgerGapsInRange obtains ingestion gaps in the history_ledgers table within the given range.
// Returns the gaps and error.
func (q *Q) GetLedgerGapsInRange(ctx context.Context, start, end uint32) ([]LedgerRange, error) {
	var result []LedgerRange
	var oldestLedger, latestLedger uint32

	if err := q.ElderLedger(ctx, &oldestLedger); err != nil {
		return nil, errors.Wrap(err, "Could not query elder ledger")
	} else if oldestLedger == 0 {
		return []LedgerRange{{
			StartSequence: start,
			EndSequence:   end,
		}}, nil
	}

	if err := q.LatestLedger(ctx, &latestLedger); err != nil {
		return nil, errors.Wrap(err, "Could not query latest ledger")
	}

	if start < oldestLedger {
		result = append(result, LedgerRange{
			StartSequence: start,
			EndSequence:   min(end, oldestLedger-1),
		})
	}
	if end <= oldestLedger {
		return result, nil
	}

	gaps, err := q.GetLedgerGaps(ctx)
	if err != nil {
		return nil, err
	}

	for _, gap := range gaps {
		if gap.EndSequence < start {
			continue
		}
		if gap.StartSequence > end {
			break
		}
		result = append(result, LedgerRange{
			StartSequence: max(gap.StartSequence, start),
			EndSequence:   min(gap.EndSequence, end),
		})
	}

	if latestLedger < end {
		result = append(result, LedgerRange{
			StartSequence: max(latestLedger+1, start),
			EndSequence:   end,
		})
	}

	return result, nil
}

func ledgerHeaderToMap(
	ledger xdr.LedgerHeaderHistoryEntry,
	successTxsCount int,
	failedTxsCount int,
	opCount int,
	txSetOpCount int,
	importerVersion int,
) (map[string]interface{}, error) {
	ledgerHeaderBase64, err := xdr.MarshalBase64(ledger.Header)
	if err != nil {
		return nil, err
	}
	closeTime := time.Unix(int64(ledger.Header.ScpValue.CloseTime), 0).UTC()
	return map[string]interface{}{
		"importer_version":             importerVersion,
		"id":                           toid.New(int32(ledger.Header.LedgerSeq), 0, 0).ToInt64(),
		"sequence":                     ledger.Header.LedgerSeq,
		"ledger_hash":                  hex.EncodeToString(ledger.Hash[:]),
		"previous_ledger_hash":         null.NewString(hex.EncodeToString(ledger.Header.PreviousLedgerHash[:]), ledger.Header.LedgerSeq > 1),
		"total_coins":                  ledger.Header.TotalCoins,
		"fee_pool":                     ledger.Header.FeePool,
		"base_fee":                     ledger.Header.BaseFee,
		"base_reserve":                 ledger.Header.BaseReserve,
		"max_tx_set_size":              ledger.Header.MaxTxSetSize,
		"closed_at":                    closeTime,
		"created_at":                   time.Now().UTC(),
		"updated_at":                   time.Now().UTC(),
		"transaction_count":            successTxsCount,
		"successful_transaction_count": successTxsCount,
		"failed_transaction_count":     failedTxsCount,
		"operation_count":              opCount,
		"tx_set_operation_count":       txSetOpCount,
		"protocol_version":             ledger.Header.LedgerVersion,
		"ledger_header":                ledgerHeaderBase64,
	}, nil
}

var selectLedger = sq.Select(
	"hl.id",
	"hl.sequence",
	"hl.importer_version",
	"hl.ledger_hash",
	"hl.previous_ledger_hash",
	"hl.transaction_count",
	"hl.successful_transaction_count",
	"hl.failed_transaction_count",
	"hl.operation_count",
	"hl.tx_set_operation_count",
	"hl.closed_at",
	"hl.created_at",
	"hl.updated_at",
	"hl.total_coins",
	"hl.fee_pool",
	"hl.base_fee",
	"hl.base_reserve",
	"hl.max_tx_set_size",
	"hl.protocol_version",
	"hl.ledger_header",
).From("history_ledgers hl")
