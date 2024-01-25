package history

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/services/aurora/internal/test"
)

func TestAddOperationParticipants(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	builder := q.NewOperationParticipantBatchInsertBuilder(1)
	err := builder.Add(tt.Ctx, 240518172673, 1)
	tt.Assert.NoError(err)

	err = builder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	type hop struct {
		OperationID int64 `db:"history_operation_id"`
		AccountID   int64 `db:"history_account_id"`
	}

	ops := []hop{}
	err = q.Select(tt.Ctx, &ops, sq.Select(
		"hopp.history_operation_id, "+
			"hopp.history_account_id").
		From("history_operation_participants hopp"),
	)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 1)

		op := ops[0]
		tt.Assert.Equal(int64(240518172673), op.OperationID)
		tt.Assert.Equal(int64(1), op.AccountID)
	}
}
