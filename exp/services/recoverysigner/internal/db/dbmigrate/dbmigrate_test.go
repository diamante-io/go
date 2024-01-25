package dbmigrate

import (
	"net/http"
	"os"
	"strings"
	"testing"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/shurcooL/httpfs/filter"
	dbpkg "github.com/diamcircle/go/exp/services/recoverysigner/internal/db"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/db/dbtest"
	supportHttp "github.com/diamcircle/go/support/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedAssets(t *testing.T) {
	localAssets := http.FileSystem(filter.Keep(http.Dir("."), func(path string, fi os.FileInfo) bool {
		return fi.IsDir() || strings.HasSuffix(path, ".sql")
	}))
	generatedAssets := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}

	if !supportHttp.EqualFileSystems(localAssets, generatedAssets, "/") {
		t.Fatalf("generated migrations does not match local migrations")
	}
}

func TestPlanMigration_upApplyOne(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	migrations, err := PlanMigration(session, migrate.Up, 1)
	require.NoError(t, err)
	wantMigrations := []string{"20200309000000-initial-1.sql"}
	assert.Equal(t, wantMigrations, migrations)
}

func TestPlanMigration_upApplyAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	migrations, err := PlanMigration(session, migrate.Up, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(migrations), 2)
	wantAtLeastMigrations := []string{
		"20200309000000-initial-1.sql",
		"20200309000001-initial-2.sql",
	}
	assert.Equal(t, wantAtLeastMigrations, migrations[:2])
}

func TestPlanMigration_upApplyNone(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 0)
	require.NoError(t, err)
	require.Greater(t, n, 1)

	migrations, err := PlanMigration(session, migrate.Up, 0)
	require.NoError(t, err)
	require.Empty(t, migrations)
}

func TestPlanMigration_downApplyOne(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	migrations, err := PlanMigration(session, migrate.Down, 1)
	require.NoError(t, err)
	wantMigrations := []string{"20200309000001-initial-2.sql"}
	assert.Equal(t, wantMigrations, migrations)
}

func TestPlanMigration_downApplyAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	migrations, err := PlanMigration(session, migrate.Down, 0)
	require.NoError(t, err)
	wantMigrations := []string{
		"20200309000001-initial-2.sql",
		"20200309000000-initial-1.sql",
	}
	assert.Equal(t, wantMigrations, migrations)
}

func TestPlanMigration_downApplyNone(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	n, err = Migrate(session, migrate.Down, 0)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	migrations, err := PlanMigration(session, migrate.Down, 0)
	require.NoError(t, err)
	assert.Empty(t, migrations)
}

func TestMigrate_upApplyOne(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	wantIDs := []string{
		"20200309000000-initial-1.sql",
	}
	assert.Equal(t, wantIDs, ids)
}

func TestMigrate_upApplyAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 0)
	require.NoError(t, err)
	require.Greater(t, n, 1)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	wantIDs := []string{
		"20200309000000-initial-1.sql",
		"20200309000001-initial-2.sql",
		"20200311000000-create-accounts.sql",
		"20200311000001-create-identities.sql",
		"20200311000002-create-auth-methods.sql",
		"20200320000000-create-accounts-audit.sql",
		"20200320000001-create-identities-audit.sql",
		"20200320000002-create-auth-methods-audit.sql",
	}
	assert.Equal(t, wantIDs, ids)
}

func TestMigrate_upApplyNone(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 0)
	require.NoError(t, err)
	require.Greater(t, n, 1)

	n, err = Migrate(session, migrate.Up, 0)
	require.NoError(t, err)
	require.Zero(t, n)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	wantIDs := []string{
		"20200309000000-initial-1.sql",
		"20200309000001-initial-2.sql",
		"20200311000000-create-accounts.sql",
		"20200311000001-create-identities.sql",
		"20200311000002-create-auth-methods.sql",
		"20200320000000-create-accounts-audit.sql",
		"20200320000001-create-identities-audit.sql",
		"20200320000002-create-auth-methods-audit.sql",
	}
	assert.Equal(t, wantIDs, ids)
}

func TestMigrate_downApplyOne(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	n, err = Migrate(session, migrate.Down, 1)
	require.NoError(t, err)
	require.Equal(t, 1, n)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	wantIDs := []string{
		"20200309000000-initial-1.sql",
	}
	assert.Equal(t, wantIDs, ids)
}

func TestMigrate_downApplyAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	n, err = Migrate(session, migrate.Down, 0)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	assert.Empty(t, ids)
}

func TestMigrate_downApplyNone(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	n, err = Migrate(session, migrate.Down, 0)
	require.NoError(t, err)
	require.Equal(t, 2, n)

	n, err = Migrate(session, migrate.Down, 0)
	require.NoError(t, err)
	require.Equal(t, 0, n)

	ids := []string{}
	err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
	require.NoError(t, err)
	assert.Empty(t, ids)
}
