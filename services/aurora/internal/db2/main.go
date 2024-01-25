// Package db2 is the replacement for db.  It provides low level db connection
// and query capabilities.
package db2

// PageQuery represents a portion of a Query struct concerned with paging
// through a large dataset.
type PageQuery struct {
	Cursor string
	Order  string
	Limit  uint64
}
