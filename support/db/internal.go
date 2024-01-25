package db

import (
	"reflect"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

var mapper = reflectx.NewMapper("db")

// ColumnsForStruct returns a slice of column names for the provided value
// (which should be a struct, a slice of structs).
func ColumnsForStruct(dest interface{}) []string {
	typ := reflect.TypeOf(dest)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}

	typmap := mapper.TypeMap(typ)

	var keys []string
	for k := range typmap.Names {
		// If a struct contains another struct (ex. sql.NullString) mapper.TypeMap
		// will return the fields of an internal struct (like: "payment_id.String",
		// "payment_id.Valid").
		// This will break the query so skip these fields.
		if strings.Contains(k, ".") {
			continue
		}
		keys = append(keys, k)
	}

	// Ensure keys are sorted.  keys is populated from a map, which has no
	// defined iteration order.  Different versions of go or different
	// architectures may cause non-deterministic results to occur (and in our CI
	// environment, they have).  To make testing easier, we sort the keys.
	sort.Strings(keys)

	return keys
}
