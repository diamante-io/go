package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/pkg/errors"
)

// Exec executes the query represented by the builder, inserting each val
// provided to the builder into the database.
func (ib *InsertBuilder) Exec(ctx context.Context) (sql.Result, error) {
	if len(ib.rows) == 0 {
		return nil, &NoRowsError{}
	}

	template := ib.rows[0]
	cols := ColumnsForStruct(template)

	if ib.ignoredCols != nil {
		finalCols := make([]string, 0, len(cols))
		for _, col := range cols {
			if ib.ignoredCols[col] {
				continue
			}
			finalCols = append(finalCols, col)
		}
		cols = finalCols
	}

	sql := ib.sql.Columns(cols...)

	// add rows onto the builder
	for _, row := range ib.rows {

		// extract field values
		rrow := reflect.ValueOf(row)
		rvals := mapper.FieldsByName(rrow, cols)

		// convert fields values to interface{}
		vals := make([]interface{}, len(cols))
		for i, rval := range rvals {
			vals[i] = rval.Interface()
		}

		// append row to insert statement
		sql = sql.Values(vals...)
	}

	// TODO: support return inserted id

	r, err := ib.Table.Session.Exec(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "insert failed")
	}

	return r, nil
}

// IgnoreCols adds colums to ignore list (will not be inserted)
func (ib *InsertBuilder) IgnoreCols(cols ...string) *InsertBuilder {
	if ib.ignoredCols == nil {
		ib.ignoredCols = make(map[string]bool)
	}
	for _, col := range cols {
		ib.ignoredCols[col] = true
	}
	return ib
}

// Rows appends more rows onto the insert statement
func (ib *InsertBuilder) Rows(rows ...interface{}) *InsertBuilder {
	ib.rows = append(ib.rows, rows...)
	return ib
}
