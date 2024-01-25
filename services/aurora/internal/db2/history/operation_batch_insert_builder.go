package history

import (
	"context"

	"github.com/guregu/null"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/xdr"
)

// OperationBatchInsertBuilder is used to insert a transaction's operations into the
// history_operations table
type OperationBatchInsertBuilder interface {
	Add(
		ctx context.Context,
		id int64,
		transactionID int64,
		applicationOrder uint32,
		operationType xdr.OperationType,
		details []byte,
		sourceAccount string,
		sourceAcccountMuxed null.String,
	) error
	Exec(ctx context.Context) error
}

// operationBatchInsertBuilder is a simple wrapper around db.BatchInsertBuilder
type operationBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

// NewOperationBatchInsertBuilder constructs a new TransactionBatchInsertBuilder instance
func (q *Q) NewOperationBatchInsertBuilder(maxBatchSize int) OperationBatchInsertBuilder {
	return &operationBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_operations"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds a transaction's operations to the batch
func (i *operationBatchInsertBuilder) Add(
	ctx context.Context,
	id int64,
	transactionID int64,
	applicationOrder uint32,
	operationType xdr.OperationType,
	details []byte,
	sourceAccount string,
	sourceAccountMuxed null.String,
) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"id":                   id,
		"transaction_id":       transactionID,
		"application_order":    applicationOrder,
		"type":                 operationType,
		"details":              details,
		"source_account":       sourceAccount,
		"source_account_muxed": sourceAccountMuxed,
	})

}

func (i *operationBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}
