package dbmigrate

import (
	"net/http"
	"os"
	"strings"
	"testing"

	dbpkg "go/services/regulated-assets-approval-server/internal/db"
	"go/services/regulated-assets-approval-server/internal/db/dbtest"
	supportHttp "go/support/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/shurcooL/httpfs/filter"
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
	wantMigrations := []string{"2021-05-05.0.initial.sql"}
	assert.Equal(t, wantMigrations, migrations)
}

func TestPlanMigration_upApplyAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	migrations, err := PlanMigration(session, migrate.Up, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(migrations), 3)
	wantAtLeastMigrations := []string{
		"2021-05-05.0.initial.sql",
		"2021-05-18.0.accounts-kyc-status.sql",
		"2021-06-08.0.pending-kyc-status.sql",
	}
	assert.Equal(t, wantAtLeastMigrations, migrations)
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
	wantMigrations := []string{"2021-05-18.0.accounts-kyc-status.sql"}
	assert.Equal(t, wantMigrations, migrations)
}

func TestPlanMigration_downApplyTwo(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)

	n, err := Migrate(session, migrate.Up, 3)
	require.NoError(t, err)
	require.Equal(t, 3, n)

	migrations, err := PlanMigration(session, migrate.Down, 0)
	require.NoError(t, err)
	wantMigrations := []string{
		"2021-06-08.0.pending-kyc-status.sql",
		"2021-05-18.0.accounts-kyc-status.sql",
		"2021-05-05.0.initial.sql",
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
		"2021-05-05.0.initial.sql",
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
		"2021-05-05.0.initial.sql",
		"2021-05-18.0.accounts-kyc-status.sql",
		"2021-06-08.0.pending-kyc-status.sql",
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
		"2021-05-05.0.initial.sql",
		"2021-05-18.0.accounts-kyc-status.sql",
		"2021-06-08.0.pending-kyc-status.sql",
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
		"2021-05-05.0.initial.sql",
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
