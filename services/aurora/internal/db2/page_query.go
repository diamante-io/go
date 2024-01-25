package db2

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"go/support/errors"
)

const (
	// DefaultPageSize is the default page size for db queries
	DefaultPageSize = 10
	// MaxPageSize is the max page size for db queries
	MaxPageSize = 200

	// OrderAscending is used to indicate an ascending order in request params
	OrderAscending = "asc"

	// OrderDescending is used to indicate an descending order in request params
	OrderDescending = "desc"

	// DefaultPairSep is the default separator used to separate two numbers for CursorInt64Pair
	DefaultPairSep = "-"
)

var (
	// ErrInvalidOrder is an error that occurs when a user-provided order string
	// is invalid
	ErrInvalidOrder = &InvalidFieldError{"order"}
	// ErrInvalidLimit is an error that occurs when a user-provided limit num
	// is invalid
	ErrInvalidLimit = &InvalidFieldError{"limit"}
	// ErrInvalidCursor is an error that occurs when a user-provided cursor string
	// is invalid
	ErrInvalidCursor = &InvalidFieldError{"cursor"}
)

type InvalidFieldError struct {
	Name string
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("%s: invalid value", e.Name)
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p PageQuery) ApplyTo(
	sql sq.SelectBuilder,
	col string,
) (sq.SelectBuilder, error) {
	cursor, err := p.CursorInt64()
	if err != nil {
		return sql, err
	}

	return p.ApplyToUsingCursor(sql, col, cursor)
}

// ApplyRawTo applies the raw PageQuery.Cursor cursor to the builder
func (p PageQuery) ApplyRawTo(
	sql sq.SelectBuilder,
	col string,
) (sq.SelectBuilder, error) {
	return p.ApplyToUsingCursor(sql, col, p.Cursor)
}

// ApplyToUsingCursor returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method allows any type of cursor by a single column
func (p PageQuery) ApplyToUsingCursor(
	sql sq.SelectBuilder,
	col string,
	cursor interface{},
) (sq.SelectBuilder, error) {
	sql = sql.Limit(p.Limit)

	switch p.Order {
	case "asc":
		if cursor == "" {
			sql = sql.
				OrderBy(fmt.Sprintf("%s asc", col))
		} else {
			sql = sql.
				Where(fmt.Sprintf("%s > ?", col), cursor).
				OrderBy(fmt.Sprintf("%s asc", col))
		}
	case "desc":
		if cursor == "" {
			sql = sql.
				OrderBy(fmt.Sprintf("%s desc", col))
		} else {
			sql = sql.
				Where(fmt.Sprintf("%s < ?", col), cursor).
				OrderBy(fmt.Sprintf("%s desc", col))
		}
	default:
		return sql, errors.Errorf("invalid order: %s", p.Order)
	}

	return sql, nil
}

// Invert returns a new PageQuery whose order is reversed
func (p PageQuery) Invert() PageQuery {
	switch p.Order {
	case OrderAscending:
		p.Order = OrderDescending
	case OrderDescending:
		p.Order = OrderAscending
	}

	return p
}

// CursorInt64 parses this query's Cursor string as an int64
func (p PageQuery) CursorInt64() (int64, error) {
	if p.Cursor == "" {
		switch p.Order {
		case OrderAscending:
			return 0, nil
		case OrderDescending:
			return math.MaxInt64, nil
		default:
			return 0, ErrInvalidOrder
		}
	}

	i, err := strconv.ParseInt(p.Cursor, 10, 64)

	if err != nil {
		return 0, ErrInvalidCursor
	}

	if i < 0 {
		return 0, ErrInvalidCursor
	}

	return i, nil

}

// CursorInt64Pair parses this query's Cursor string as two int64s, separated by the provided separator
func (p PageQuery) CursorInt64Pair(sep string) (l int64, r int64, err error) {
	if p.Cursor == "" {
		switch p.Order {
		case OrderAscending:
			l = 0
			r = 0
		case OrderDescending:
			l = math.MaxInt64
			r = math.MaxInt64
		default:
			err = ErrInvalidOrder
		}
		return
	}

	parts := strings.SplitN(p.Cursor, sep, 2)

	// In the event that the cursor is only a single number
	// we use maxInt as the second element.  This ensures that
	// cursors containing a single element skip past all entries
	// specified by the first element.
	//
	// As an example, this behavior ensures that an effect cursor
	// specified using only a ledger sequence will properly exclude
	// all effects originated in the sequence provided.
	if len(parts) != 2 {
		// NOTE(scott): we work around a build issue here which cases an overflow
		// when building for ARM. We assign the untyped constant explicitly to a
		// int64.
		var max int64 = math.MaxInt64
		parts = append(parts, fmt.Sprintf("%d", max))
	}

	l, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		err = errors.Wrap(err, "first component unparseable")
		return
	}

	r, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = errors.Wrap(err, "second component unparseable")
		return
	}

	if l < 0 || r < 0 {
		err = ErrInvalidCursor
	}

	return
}

// NewPageQuery creates a new PageQuery struct, ensuring the order, limit, and
// cursor are set to the appropriate defaults and are valid.
func NewPageQuery(
	cursor string,
	validateCursor bool,
	order string,
	limit uint64,
) (result PageQuery, err error) {

	// Set order
	switch order {
	case "":
		result.Order = OrderAscending
	case OrderAscending, OrderDescending:
		result.Order = order
	default:
		err = ErrInvalidOrder
		return
	}

	// Set cursor
	result.Cursor = cursor
	if validateCursor {
		_, _, err = result.CursorInt64Pair(DefaultPairSep)
		if err != nil {
			err = ErrInvalidCursor
			return
		}
	}

	// Set limit
	switch {
	case limit == 0:
		err = ErrInvalidLimit
		return
	case limit > MaxPageSize:
		err = ErrInvalidLimit
		return
	default:
		result.Limit = limit
	}

	return
}

// MustPageQuery behaves as NewPageQuery, but panics upon error
func MustPageQuery(cursor string, validateCursor bool, order string, limit uint64) PageQuery {
	r, err := NewPageQuery(cursor, validateCursor, order, limit)
	if err != nil {
		panic(err)
	}

	return r
}
