package history

import (
	"context"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
)

// QHistoryLiquidityPools defines account related queries.
type QHistoryLiquidityPools interface {
	CreateHistoryLiquidityPools(ctx context.Context, poolIDs []string, batchSize int) (map[string]int64, error)
	NewOperationLiquidityPoolBatchInsertBuilder(maxBatchSize int) OperationLiquidityPoolBatchInsertBuilder
	NewTransactionLiquidityPoolBatchInsertBuilder(maxBatchSize int) TransactionLiquidityPoolBatchInsertBuilder
}

// CreateHistoryLiquidityPools creates rows in the history_liquidity_pools table for a given list of ids.
// CreateHistoryLiquidityPools returns a mapping of id to its corresponding internal id in the history_liquidity_pools table
func (q *Q) CreateHistoryLiquidityPools(ctx context.Context, poolIDs []string, batchSize int) (map[string]int64, error) {
	if len(poolIDs) == 0 {
		return nil, nil
	}

	builder := &db.BatchInsertBuilder{
		Table:        q.GetTable("history_liquidity_pools"),
		MaxBatchSize: batchSize,
		Suffix:       "ON CONFLICT (liquidity_pool_id) DO NOTHING",
	}

	// sort before inserting to prevent deadlocks on acquiring a ShareLock
	// https://github.com/diamcircle/go/issues/2370
	sort.Strings(poolIDs)
	var deduped []string
	for i, id := range poolIDs {
		if i > 0 && id == poolIDs[i-1] {
			// skip duplicates
			continue
		}
		deduped = append(deduped, id)
		err := builder.Row(ctx, map[string]interface{}{
			"liquidity_pool_id": id,
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not insert history_liquidity_pools row")
		}
	}

	err := builder.Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not exec claimable balance insert builder")
	}

	var lps []HistoryLiquidityPool
	toInternalID := map[string]int64{}
	const selectBatchSize = 10000

	for i := 0; i < len(deduped); i += selectBatchSize {
		end := i + selectBatchSize
		if end > len(deduped) {
			end = len(deduped)
		}
		subset := deduped[i:end]

		lps, err = q.LiquidityPoolsByIDs(ctx, subset)
		if err != nil {
			return nil, errors.Wrap(err, "could not select claimable balances")
		}

		for _, lp := range lps {
			toInternalID[lp.PoolID] = lp.InternalID
		}
	}

	return toInternalID, nil
}

// HistoryLiquidityPool is a row of data from the `history_liquidity_pools` table
type HistoryLiquidityPool struct {
	PoolID     string `db:"liquidity_pool_id"`
	InternalID int64  `db:"id"`
}

var selectHistoryLiquidityPool = sq.Select("hlp.*").From("history_liquidity_pools hlp")

// LiquidityPoolsByIDs loads rows from `history_liquidity_pools`, by liquidity_pool_id
func (q *Q) LiquidityPoolsByIDs(ctx context.Context, poolIDs []string) (dest []HistoryLiquidityPool, err error) {
	sql := selectHistoryLiquidityPool.Where(map[string]interface{}{
		"hlp.liquidity_pool_id": poolIDs, // hlp.liquidity_pool_id IN (...)
	})
	err = q.Select(ctx, &dest, sql)
	return dest, err
}

// LiquidityPoolByID loads a row from `history_liquidity_pools`, by liquidity_pool_id
func (q *Q) LiquidityPoolByID(ctx context.Context, poolID string) (dest HistoryLiquidityPool, err error) {
	sql := selectHistoryLiquidityPool.Limit(1).Where("hlp.liquidity_pool_id = ?", poolID)
	err = q.Get(ctx, &dest, sql)
	return dest, err
}

type OperationLiquidityPoolBatchInsertBuilder interface {
	Add(ctx context.Context, operationID, internalID int64) error
	Exec(ctx context.Context) error
}

type operationLiquidityPoolBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

func (q *Q) NewOperationLiquidityPoolBatchInsertBuilder(maxBatchSize int) OperationLiquidityPoolBatchInsertBuilder {
	return &operationLiquidityPoolBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_operation_liquidity_pools"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds a new operation claimable balance to the batch
func (i *operationLiquidityPoolBatchInsertBuilder) Add(ctx context.Context, operationID, internalID int64) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"history_operation_id":      operationID,
		"history_liquidity_pool_id": internalID,
	})
}

// Exec flushes all pending operation claimable balances to the db
func (i *operationLiquidityPoolBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}

type TransactionLiquidityPoolBatchInsertBuilder interface {
	Add(ctx context.Context, transactionID, internalID int64) error
	Exec(ctx context.Context) error
}

type transactionLiquidityPoolBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

func (q *Q) NewTransactionLiquidityPoolBatchInsertBuilder(maxBatchSize int) TransactionLiquidityPoolBatchInsertBuilder {
	return &transactionLiquidityPoolBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_transaction_liquidity_pools"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds a new transaction claimable balance to the batch
func (i *transactionLiquidityPoolBatchInsertBuilder) Add(ctx context.Context, transactionID, internalID int64) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"history_transaction_id":    transactionID,
		"history_liquidity_pool_id": internalID,
	})
}

// Exec flushes all pending transaction claimable balances to the db
func (i *transactionLiquidityPoolBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}
