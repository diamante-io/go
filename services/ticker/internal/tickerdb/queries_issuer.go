package tickerdb

import (
	"context"
	"strings"

	"github.com/diamcircle/go/services/ticker/internal/utils"
)

// InsertOrUpdateIssuer inserts an Issuer on the database (if new),
// or updates an existing one
func (s *TickerSession) InsertOrUpdateIssuer(ctx context.Context, issuer *Issuer, preserveFields []string) (id int32, err error) {
	dbFields := getDBFieldTags(*issuer, true)
	dbFieldsString := strings.Join(dbFields, ", ")
	dbValues := getDBFieldValues(*issuer, true)

	cleanPreservedFields := sanitizeFieldNames(preserveFields)
	toUpdateFields := utils.SliceDiff(dbFields, cleanPreservedFields)

	qs := "INSERT INTO issuers (" + dbFieldsString + ")"
	qs += " VALUES (" + generatePlaceholders(dbValues) + ")"
	qs += " " + createOnConflictFragment("public_key_unique", toUpdateFields)
	qs += " RETURNING id;"

	rows, err := s.QueryRaw(ctx, qs, dbValues...)
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&id)
	}
	return
}

// GetAllIssuers returns a slice with all issuers in the database
func (s *TickerSession) GetAllIssuers(ctx context.Context) (issuers []Issuer, err error) {
	err = s.SelectRaw(ctx, &issuers, "SELECT * FROM issuers")
	return
}
