package history

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	sq "github.com/Masterminds/squirrel"

	"github.com/diamcircle/go/protocols/aurora/effects"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
)

// UnmarshalDetails unmarshals the details of this effect into `dest`
func (r *Effect) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}

	err := errors.Wrap(json.Unmarshal([]byte(r.DetailsString.String), &dest), "unmarshal effect details failed")
	if err == nil {
		// In 2.9.0 a new `asset_type` was introduced to include liquidity
		// pools. Instead of reingesting entire history, let's fill the
		// `asset_type` here if it's empty.
		// (I hate to convert to `protocol` types here but there's no other way
		// without larger refactor.)
		switch dest := dest.(type) {
		case *effects.TrustlineSponsorshipCreated:
			if dest.Type == "" {
				dest.Type = getAssetTypeForCanonicalAsset(dest.Asset)
			}
		case *effects.TrustlineSponsorshipUpdated:
			if dest.Type == "" {
				dest.Type = getAssetTypeForCanonicalAsset(dest.Asset)
			}
		case *effects.TrustlineSponsorshipRemoved:
			if dest.Type == "" {
				dest.Type = getAssetTypeForCanonicalAsset(dest.Asset)
			}
		}
	}
	return err
}

func getAssetTypeForCanonicalAsset(canonicalAsset string) string {
	if len(canonicalAsset) <= 61 {
		return "credit_alphanum4"
	} else {
		return "credit_alphanum12"
	}
}

// ID returns a lexically ordered id for this effect record
func (r *Effect) ID() string {
	return fmt.Sprintf("%019d-%010d", r.HistoryOperationID, r.Order)
}

// LedgerSequence return the ledger in which the effect occurred.
func (r *Effect) LedgerSequence() int32 {
	id := toid.Parse(r.HistoryOperationID)
	return id.LedgerSequence
}

// PagingToken returns a cursor for this effect
func (r *Effect) PagingToken() string {
	return fmt.Sprintf("%d-%d", r.HistoryOperationID, r.Order)
}

// Effects provides a helper to filter rows from the `history_effects`
// table with pre-defined filters.  See `TransactionsQ` methods for the
// available filters.
func (q *Q) Effects() *EffectsQ {
	return &EffectsQ{
		parent: q,
		sql:    selectEffect,
	}
}

// ForAccount filters the operations collection to a specific account
func (q *EffectsQ) ForAccount(ctx context.Context, aid string) *EffectsQ {
	var account Account
	q.Err = q.parent.AccountByAddress(ctx, &account, aid)
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("heff.history_account_id = ?", account.ID)

	return q
}

// ForLedger filters the query to only effects in a specific ledger,
// specified by its sequence.
func (q *EffectsQ) ForLedger(ctx context.Context, seq int32) *EffectsQ {
	var ledger Ledger
	q.Err = q.parent.LedgerBySequence(ctx, &ledger, seq)
	if q.Err != nil {
		return q
	}

	start := toid.ID{LedgerSequence: seq}
	end := toid.ID{LedgerSequence: seq + 1}
	q.sql = q.sql.Where(
		"heff.history_operation_id >= ? AND heff.history_operation_id < ?",
		start.ToInt64(),
		end.ToInt64(),
	)

	return q
}

// ForOperation filters the query to only effects in a specific operation,
// specified by its id.
func (q *EffectsQ) ForOperation(id int64) *EffectsQ {
	start := toid.Parse(id)
	end := start
	end.IncOperationOrder()
	q.sql = q.sql.Where(
		"heff.history_operation_id >= ? AND heff.history_operation_id < ?",
		start.ToInt64(),
		end.ToInt64(),
	)

	return q
}

// ForLiquidityPool filters the query to only effects in a specific liquidity pool,
// specified by its id.
func (q *EffectsQ) ForLiquidityPool(ctx context.Context, page db2.PageQuery, id string) *EffectsQ {
	if q.Err != nil {
		return q
	}

	op, _, err := page.CursorInt64Pair(db2.DefaultPairSep)
	if err != nil {
		q.Err = err
		return q
	}

	query := `SELECT holp.history_operation_id
	FROM history_operation_liquidity_pools holp
	WHERE holp.history_liquidity_pool_id = (SELECT id FROM history_liquidity_pools WHERE liquidity_pool_id =  ?)
	`
	switch page.Order {
	case "asc":
		query += "AND holp.history_operation_id >= ? ORDER BY holp.history_operation_id asc LIMIT ?"
	case "desc":
		query += "AND holp.history_operation_id <= ? ORDER BY holp.history_operation_id desc LIMIT ?"
	default:
		q.Err = errors.Errorf("invalid paging order: %s", page.Order)
		return q
	}

	var liquidityPoolOperationIDs []int64
	err = q.parent.SelectRaw(ctx, &liquidityPoolOperationIDs, query, id, op, page.Limit)
	if err != nil {
		q.Err = err
		return q
	}

	q.sql = q.sql.Where(map[string]interface{}{
		"heff.history_operation_id": liquidityPoolOperationIDs,
	})
	return q
}

// ForTransaction filters the query to only effects in a specific
// transaction, specified by the transactions's hex-encoded hash.
func (q *EffectsQ) ForTransaction(ctx context.Context, hash string) *EffectsQ {
	var tx Transaction
	q.Err = q.parent.TransactionByHash(ctx, &tx, hash)
	if q.Err != nil {
		return q
	}

	start := toid.Parse(tx.ID)
	end := start
	end.TransactionOrder++
	q.sql = q.sql.Where(
		"heff.history_operation_id >= ? AND heff.history_operation_id < ?",
		start.ToInt64(),
		end.ToInt64(),
	)

	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *EffectsQ) Page(page db2.PageQuery) *EffectsQ {
	if q.Err != nil {
		return q
	}

	op, idx, err := page.CursorInt64Pair(db2.DefaultPairSep)
	if err != nil {
		q.Err = err
		return q
	}

	if idx > math.MaxInt32 {
		idx = math.MaxInt32
	}

	// NOTE: Remember to test the queries below with EXPLAIN / EXPLAIN ANALYZE
	// before changing them.
	// This condition is using multicolumn index and it's easy to write it in a way that
	// DB will perform a full table scan.
	switch page.Order {
	case "asc":
		q.sql = q.sql.
			Where(`(
					 heff.history_operation_id >= ?
				AND (
					 heff.history_operation_id > ? OR
					(heff.history_operation_id = ? AND heff.order > ?)
				))`, op, op, op, idx).
			OrderBy("heff.history_operation_id asc, heff.order asc")
	case "desc":
		q.sql = q.sql.
			Where(`(
					 heff.history_operation_id <= ?
				AND (
					 heff.history_operation_id < ? OR
					(heff.history_operation_id = ? AND heff.order < ?)
				))`, op, op, op, idx).
			OrderBy("heff.history_operation_id desc, heff.order desc")
	}

	q.sql = q.sql.Limit(page.Limit)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *EffectsQ) Select(ctx context.Context, dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(ctx, dest, q.sql)
	return q.Err
}

// QEffects defines history_effects related queries.
type QEffects interface {
	QCreateAccountsHistory
	NewEffectBatchInsertBuilder(maxBatchSize int) EffectBatchInsertBuilder
}

var selectEffect = sq.Select("heff.*, hacc.address").
	From("history_effects heff").
	LeftJoin("history_accounts hacc ON hacc.id = heff.history_account_id")
