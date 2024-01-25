package db

import sq "github.com/Masterminds/squirrel"

// Delete returns a new query builder configured to delete rows from the table.
//
func (tbl *Table) Delete(
	pred interface{},
	args ...interface{},
) *DeleteBuilder {
	return &DeleteBuilder{
		Table: tbl,
		sql:   sq.Delete(tbl.Name).Where(pred, args...),
	}
}

// Get returns a new query builder configured to select into the provided
// `dest`.
//
// Get behaves the same was as Select, but automatically limits the query
// generated to a single value and only populates a single struct.
func (tbl *Table) Get(
	dest interface{},
	pred interface{},
	args ...interface{},
) *GetBuilder {

	cols := ColumnsForStruct(dest)
	sql := sq.Select(cols...).From(tbl.Name).Where(pred, args...).Limit(1)

	return &GetBuilder{
		Table: tbl,
		dest:  dest,
		sql:   sql,
	}
}

// Insert returns a new query builder configured to insert structs into the
// table.
//
// Insert takes one or more struct (or pointer to struct) values, each of which
// represents a single row to be created in the table.  The first value provided
// in a call to this function will operate as the template for the insert and
// will determine what columns are populated in the query.   For this reason, it
// is highly recommmended that you always use the same struct type for any
// single call this function.
//
// An InsertBuilder uses the "db" struct tag to determine the column names that
// a given struct should be mapped to, and by default the unmofdified name of
// the field will be used.  Similar to other struct tags, the value "-" will
// cause the field to be skipped.
//
// NOTE:  using the omitempty option, such as used with json struct tags, is not
// supported.
func (tbl *Table) Insert(rows ...interface{}) *InsertBuilder {
	return &InsertBuilder{
		Table: tbl,
		sql:   sq.Insert(tbl.Name),
		rows:  rows,
	}
}

// Select returns a new query builder configured to select into the provided
// `dest`.
func (tbl *Table) Select(
	dest interface{},
	pred interface{},
	args ...interface{},
) *SelectBuilder {

	cols := ColumnsForStruct(dest)
	sql := sq.Select(cols...).From(tbl.Name).Where(pred, args...)

	return &SelectBuilder{
		Table: tbl,
		dest:  dest,
		sql:   sql,
	}
}

// Update returns a new query builder configured to update the table.
// See docs for `UpdateBuilderExec` for more documentation.
func (tbl *Table) Update() *UpdateBuilder {
	return &UpdateBuilder{
		Table: tbl,
		sql:   sq.Update(tbl.Name),
	}
}
