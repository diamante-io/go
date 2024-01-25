package tickerdb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go/services/ticker/internal/utils"
)

// getDBFieldTags returns all "db" tags for a given struct, optionally excluding the "id".
func getDBFieldTags(model interface{}, excludeID bool) (fields []string) {
	r := reflect.ValueOf(model)
	for i := 0; i < r.Type().NumField(); i++ {
		dbField := r.Type().Field(i).Tag.Get("db")
		if (excludeID && dbField == "id") || dbField == "-" { // ensure fields marked with a "-" tag are ignored
			continue
		}
		fields = append(fields, dbField)
	}
	fields = sanitizeFieldNames(fields)
	return
}

// sanitizeFieldNames adds double quotes to each entry on a slice of field names.
func sanitizeFieldNames(fieldNames []string) (sanitizedFields []string) {
	for _, v := range fieldNames {
		quotedField := fmt.Sprintf("\"%s\"", v)
		sanitizedFields = append(sanitizedFields, quotedField)
	}
	return
}

// getDBFieldValues returns all "db"-tagged values, optionally excluding the "id".
func getDBFieldValues(model interface{}, excludeID bool) (values []interface{}) {
	r := reflect.ValueOf(model)
	for i := 0; i < r.Type().NumField(); i++ {
		dbField := r.Type().Field(i).Tag.Get("db")
		dbVal := r.Field(i).Interface()
		if (excludeID && dbField == "id") || dbField == "-" { // ensure fields marked with a "-" tag are ignored
			continue
		}
		values = append(values, dbVal)
	}
	return
}

// createOnConflictFragment generates a ON CONFLICT sql clause for a given constraint,
// preserving the fields listed in the fields param.
func createOnConflictFragment(constraint string, fields []string) (fragment string) {
	fragment = fmt.Sprintf("ON CONFLICT ON CONSTRAINT %s DO UPDATE SET ", constraint)
	for i, field := range fields {
		fragment += fmt.Sprintf("%s = EXCLUDED.%s", field, field)

		if i != len(fields)-1 {
			fragment += ","
		}
	}

	return
}

// generatePlaceholders generates a string formatted as (?, ?, ?, ?) of length
// equal to the size of the fields param
func generatePlaceholders(fields []interface{}) (placeholder string) {
	for i := range fields {
		placeholder += "?"

		if i != len(fields)-1 {
			placeholder += ", "
		}
	}
	return
}

// optionalVar is a simple struct to represent a query variable that should
// only be used in a statement if its value is not null
type optionalVar struct {
	name string
	val  *string
}

// generateWhereClause generates a WHERE clause in the format:
// "WHERE x = ? AND y = ? AND ..." where the number of conditions is equal
// to the number of optVars with val != nil. It also returns the valid vals
// in the args param. This function was created to take advantage of go/sql's
// sanitization and to prevent possible SQL injections.
func generateWhereClause(optVars []optionalVar) (clause string, args []string) {
	for _, ov := range optVars {
		if ov.val != nil {
			if clause == "" {
				clause = fmt.Sprintf("WHERE %s = ?", ov.name)
			} else {
				clause += fmt.Sprintf(" AND %s = ?", ov.name)
			}
			args = append(args, *ov.val)
		}
	}
	return
}

// getBaseAndCounterCodes takes an asset pair name string (e.g: XLM_BTC)
// and returns the parsed asset codes (e.g.: XLM, BTC). It also reverses
// the assets, according to the following rules:
// 1. XLM is always the base asset
// 2. If XLM is not in the pair, the assets should be ordered alphabetically
func getBaseAndCounterCodes(pairName string) (string, string, error) {
	assets := strings.Split(pairName, "_")
	if len(assets) != 2 {
		return "", "", errors.New("invalid asset pair name")
	}

	if (assets[1] == "XLM") || (assets[0] != "XLM" && assets[0] > assets[1]) {
		return assets[1], assets[0], nil
	}

	return assets[0], assets[1], nil
}

// performUpsertQuery introspects a dbStruct interface{} and performs an insert query
// (if the conflictConstraint isn't broken), otherwise it updates the instance on the
// db, preserving the old values for the fields in preserveFields
func (s *TickerSession) performUpsertQuery(ctx context.Context, dbStruct interface{}, tableName string, conflictConstraint string, preserveFields []string) error {
	dbFields := getDBFieldTags(dbStruct, true)
	dbFieldsString := strings.Join(dbFields, ", ")
	dbValues := getDBFieldValues(dbStruct, true)

	cleanPreservedFields := sanitizeFieldNames(preserveFields)
	toUpdateFields := utils.SliceDiff(dbFields, cleanPreservedFields)

	qs := fmt.Sprintf("INSERT INTO %s (", tableName) + dbFieldsString + ")"
	qs += " VALUES (" + generatePlaceholders(dbValues) + ")"
	qs += " " + createOnConflictFragment(conflictConstraint, toUpdateFields) + ";"
	_, err := s.ExecRaw(ctx, qs, dbValues...)
	return err
}
