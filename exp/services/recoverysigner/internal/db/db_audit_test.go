package db

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/db/dbtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuditTables confirms that the columns for audit tables are a superset of
// the columns in the tables they are auditing.
func TestAuditTables(t *testing.T) {
	db := dbtest.Open(t)
	conn, err := Open(db.DSN)
	require.NoError(t, err)

	assertAuditTableCols(t, conn, "accounts", "accounts_audit")
	assertAuditTableCols(t, conn, "identities", "identities_audit")
	assertAuditTableCols(t, conn, "auth_methods", "auth_methods_audit")
}

// assertAuditTableCols checks that the audit table for the given
// table has the same columns as the given table, as well as the header
// columns, that all the types and columns are as we expect.
func assertAuditTableCols(t *testing.T, db *sqlx.DB, tableName, auditTableName string) {
	t.Run(tableName, func(t *testing.T) {
		cols, err := tableCols(db, tableName)
		require.NoError(t, err)

		wantAuditHeaderCols := []tableCol{
			{Name: "audit_id", DataType: "bigint", UDTName: "int8", IsNullable: "NO"},
			{Name: "audit_at", DataType: "timestamp with time zone", UDTName: "timestamptz", IsNullable: "NO"},
			{Name: "audit_user", DataType: "text", UDTName: "text", IsNullable: "NO"},
			{Name: "audit_op", DataType: "USER-DEFINED", UDTName: "audit_op", IsNullable: "NO"},
		}
		wantAuditCols := append(append([]tableCol{}, wantAuditHeaderCols...), cols...)

		auditCols, err := tableCols(db, auditTableName)
		require.NoError(t, err)
		assert.Equal(t, wantAuditCols, auditCols)
	})
}

// tableCol represents a column in a table with some of its information as
// defined by Postgres' standard information_schema table.
type tableCol struct {
	Name       string
	DataType   string
	UDTName    string
	IsNullable string
}

// tableCols returns the column names for the table.
func tableCols(db *sqlx.DB, tableName string) ([]tableCol, error) {
	cols := []tableCol{}
	err := db.Select(
		&cols,
		`SELECT
			column_name as Name,
			data_type as DataType,
			udt_name as UDTName,
			is_nullable as IsNullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = $1
		ORDER BY ordinal_position ASC;`,
		tableName,
	)
	if err != nil {
		return nil, err
	}
	return cols, nil
}
