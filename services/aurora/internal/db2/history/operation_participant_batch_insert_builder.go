package history

import (
	"context"

	"github.com/diamcircle/go/support/db"
)

// OperationParticipantBatchInsertBuilder is used to insert a transaction's operations into the
// history_operations table
type OperationParticipantBatchInsertBuilder interface {
	Add(
		ctx context.Context,
		operationID int64,
		accountID int64,
	) error
	Exec(ctx context.Context) error
}

// operationParticipantBatchInsertBuilder is a simple wrapper around db.BatchInsertBuilder
type operationParticipantBatchInsertBuilder struct {
	builder db.BatchInsertBuilder
}

// NewOperationParticipantBatchInsertBuilder constructs a new TransactionBatchInsertBuilder instance
func (q *Q) NewOperationParticipantBatchInsertBuilder(maxBatchSize int) OperationParticipantBatchInsertBuilder {
	return &operationParticipantBatchInsertBuilder{
		builder: db.BatchInsertBuilder{
			Table:        q.GetTable("history_operation_participants"),
			MaxBatchSize: maxBatchSize,
		},
	}
}

// Add adds an operation participant to the batch
func (i *operationParticipantBatchInsertBuilder) Add(
	ctx context.Context,
	operationID int64,
	accountID int64,
) error {
	return i.builder.Row(ctx, map[string]interface{}{
		"history_operation_id": operationID,
		"history_account_id":   accountID,
	})
}

func (i *operationParticipantBatchInsertBuilder) Exec(ctx context.Context) error {
	return i.builder.Exec(ctx)
}
