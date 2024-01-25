package dbmigrate

import (
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata@v3.18.0+incompatible -nometadata -ignore .+\.(go|swp)$ -pkg dbmigrate -o dbmigrate_generated.go ./migrations

var migrationSource = &migrate.AssetMigrationSource{
	Asset:    Asset,
	AssetDir: AssetDir,
	Dir:      "migrations",
}

// PlanMigration finds the migrations that would be applied if Migrate was to
// be run now.
func PlanMigration(db *sqlx.DB, dir migrate.MigrationDirection, count int) ([]string, error) {
	migrations, _, err := migrate.PlanMigration(db.DB, db.DriverName(), migrationSource, dir, count)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(migrations))
	for _, m := range migrations {
		ids = append(ids, m.Id)
	}
	return ids, nil
}

// Migrate runs all the migrations to get the database to the state described
// by the migration files in the direction specified. Count is the maximum
// number of migrations to apply or rollback.
func Migrate(db *sqlx.DB, dir migrate.MigrationDirection, count int) (int, error) {
	return migrate.ExecMax(db.DB, db.DriverName(), migrationSource, dir, count)
}
