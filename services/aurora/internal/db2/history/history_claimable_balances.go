package history

import (
	"context"
	"sort"

	sq "github.com/Masterminds/squirrel"

	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
)

// QHistoryClaimableBalances defines account related queries.
type QHistoryClaimableBalances interface {
	CreateHistoryClaimableBalances(ctx context.Context, ids []string, batchSize int) (map[string]int64, error)
	NewOperationClaimableBalanceBatchInsertBuilder(maxBatchSize int) OperationClaimableBalanceBatchInsertBuilder
	NewTransactionClaimableBalanceBatchInsertBuilder(maxBatchSize int) TransactionClaimableBalanceBatchInsertBuilder
}

// CreateHistoryClaimableBalances creates rows in the history_claimable_balances table for a given list of ids.
// CreateHistoryClaimableBalances returns a mapping of id to its corresponding internal id in the history_claimable_balances table
func (q *Q) CreateHistoryClaimableBalances(ctx context.Context, ids []string, batchSize int) (map[string]int64, error) {
	builder := &db.BatchInsertBuilder{
		Table:        q.GetTable("history_claimable_balances"),
		MaxBatchSize: batchSize,
		Suffix:       "ON CONFLICT (claimable_balance_id) DO NOTHING",
	}

	// sort before inserting to prevent deadlocks on acquiring a ShareLock
	// https://github.com/diamcircle/go/issues/2370
	sort.Strings(ids)
	for _, id := range ids {
		err := builder.Row(ctx, map[string]interface{}{
			"claimable_balance_id": id,
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not insert history_claimable_balances row")
		}
	}

	err := builder.Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not exec claimable balance insert builder")
	}

	var cbs []HistoryClaimableBalance
	toInternalID := map[string]int64{}
	const selectBatchSize = 10000

	for i := 0; i < len(ids); i += selectBatchSize {
		end := i + selectBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		subset := ids[i:end]

		cbs, err = q.ClaimableBalancesByIDs(ctx, subset)
		if err != nil {
			return nil, errors.Wrap(err, "could not select claimable balances")
		}

		for _, cb := range cbs {
			toInternalID[cb.BalanceID] = cb.InternalID
		}
	}

	return toInternalID, nil
}

// HistoryClaimableBalance is a row of data from the `history_claimable_balances` table
type HistoryClaimableBalance struct {
	BalanceID  string `db:"claimable_balance_id"`
	InternalID int64  `db:"id"`
}

var selectHistoryClaimableBalance = sq.Select("hcb.*").From("history_claimable_balances hcb")

// ClaimableBalancesByIDs loads rows from `history_claimable_balances`, by claimable_balance_id
func (q *Q) ClaimableBalancesByIDs(ctx context.Context, ids []string) (dest []HistoryClaimableBalance, err error) {
	sql := selectHistoryClaimableBalance.Where(map[string]interface{}{
		"hcb.claimable_balance_id": ids, // hcb.claimable_balance_id IN (...)
	})
	err = q.Select(ctx, &dest, sql)
	return dest, err
}

// ClaimableBalanceByID loads a row from `history_claimable_balances`, by claimable_balance_id
func (q *Q) ClaimableBalanceByID(ctx context.Context, id string) (dest HistoryClaimableBalance, err error) {
	sql := selectHistoryClaimableBalance.Limit(1).Where("hcb.claimable_balance_id = ?", id)
	err = q.Get(ctx, &dest, sql)
	return dest, err
}

type OperationClaimableBalanceBatchInsertBuilder interface {
	Add(ctx context.Context, operationID, internalID int64) error
	Exec(ctx context.Context) error
}

type operationClaimableBalanceBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

func (q *Q) NewOperationClaimableBalanceBatchInsertBuilder(maxBatchSize int) OperationClaimableBalanceBatchInsertBuilder {
	return &operationClaimableBalanceBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_operation_claimable_balances"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds a new operation claimable balance to the batch
func (i *operationClaimableBalanceBatchInsertBuilder) Add(ctx context.Context, operationID, internalID int64) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"history_operation_id":         operationID,
		"history_claimable_balance_id": internalID,
	})
}

// Exec flushes all pending operation claimable balances to the db
func (i *operationClaimableBalanceBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}

type TransactionClaimableBalanceBatchInsertBuilder interface {
	Add(ctx context.Context, transactionID, internalID int64) error
	Exec(ctx context.Context) error
}

type transactionClaimableBalanceBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

func (q *Q) NewTransactionClaimableBalanceBatchInsertBuilder(maxBatchSize int) TransactionClaimableBalanceBatchInsertBuilder {
	return &transactionClaimableBalanceBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_transaction_claimable_balances"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds a new transaction claimable balance to the batch
func (i *transactionClaimableBalanceBatchInsertBuilder) Add(ctx context.Context, transactionID, internalID int64) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"history_transaction_id":       transactionID,
		"history_claimable_balance_id": internalID,
	})
}

// Exec flushes all pending transaction claimable balances to the db
func (i *transactionClaimableBalanceBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}
