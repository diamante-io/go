package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// Exec executes the query represented by the builder, deleting any rows that
// match the queries where clauses.
func (delb *DeleteBuilder) Exec(ctx context.Context) (sql.Result, error) {
	r, err := delb.Table.Session.Exec(ctx, delb.sql)
	if err != nil {
		return nil, errors.Wrap(err, "delete failed")
	}
	return r, nil
}

// Where is a passthrough call to the squirrel.  See
// https://godoc.org/github.com/Masterminds/squirrel#DeleteBuilder.Where
func (delb *DeleteBuilder) Where(
	pred interface{},
	args ...interface{},
) *DeleteBuilder {
	delb.sql = delb.sql.Where(pred, args...)
	return delb
}
