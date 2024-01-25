package schema

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	stdLog "log"
	"text/tabwriter"
	"time"

	migrate "github.com/rubenv/sql-migrate"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata@v3.18.0+incompatible -nometadata -pkg schema -o bindata.go migrations/

// MigrateDir represents a direction in which to perform schema migrations.
type MigrateDir string

const (
	// MigrateUp causes migrations to be run in the "up" direction.
	MigrateUp MigrateDir = "up"
	// MigrateDown causes migrations to be run in the "down" direction.
	MigrateDown MigrateDir = "down"
	// MigrateRedo causes migrations to be run down, then up
	MigrateRedo MigrateDir = "redo"
)

// Migrations represents all of the schema migration for horizon
var Migrations migrate.MigrationSource = &migrate.AssetMigrationSource{
	Asset:    Asset,
	AssetDir: AssetDir,
	Dir:      "migrations",
}

// Migrate performs schema migration.  Migrations can occur in one of three
// ways:
//
// - up: migrations are performed from the currently installed version upwards.
// If count is 0, all unapplied migrations will be run.
//
// - down: migrations are performed from the current version downard. If count
// is 0, all applied migrations will be run in a downard direction.
//
// - redo: migrations are first ran downard `count` times, and then are rand
// upward back to the current version at the start of the process. If count is
// 0, a count of 1 will be assumed.
func Migrate(db *sql.DB, dir MigrateDir, count int) (int, error) {
	switch dir {
	case MigrateUp:
		return migrate.ExecMax(db, "postgres", Migrations, migrate.Up, count)
	case MigrateDown:
		return migrate.ExecMax(db, "postgres", Migrations, migrate.Down, count)
	case MigrateRedo:

		if count == 0 {
			count = 1
		}

		down, err := migrate.ExecMax(db, "postgres", Migrations, migrate.Down, count)
		if err != nil {
			return down, err
		}

		return migrate.ExecMax(db, "postgres", Migrations, migrate.Up, down)
	default:
		return 0, errors.New("Invalid migration direction")
	}
}

// GetMigrationsUp returns a list of names of any migrations needed in the
// "up" direction (more recent schema versions).
func GetMigrationsUp(dbUrl string) (migrationIds []string) {
	// Get a DB handle
	db, dbErr := sql.Open("postgres", dbUrl)
	if dbErr != nil {
		stdLog.Fatal(dbErr)
	}
	defer db.Close()

	// Get the possible migrations
	possibleMigrations, _, migrateErr := migrate.PlanMigration(db, "postgres", Migrations, migrate.Up, 0)
	if migrateErr != nil {
		stdLog.Fatal(migrateErr)
	}

	// Extract a list of the possible migration names
	for _, m := range possibleMigrations {
		migrationIds = append(migrationIds, m.Id)
	}

	return migrationIds
}

// GetNumMigrationsDown returns the number of migrations to apply in the
// "down" direction to return to the older schema version expected by this
// version of Horizon. To keep the code simple, it does not provide a list of
// migration names.
func GetNumMigrationsDown(dbUrl string) (nMigrations int) {
	// Get a DB handle
	db, dbErr := sql.Open("postgres", dbUrl)
	if dbErr != nil {
		stdLog.Fatal(dbErr)
	}
	defer db.Close()

	// Get the set of migrations recorded in the database
	migrationRecords, recordErr := migrate.GetMigrationRecords(db, "postgres")
	if recordErr != nil {
		stdLog.Fatal(recordErr)
	}

	// Get the list of migrations needed by this version of Horizon
	allNeededMigrations, _, migrateErr := migrate.PlanMigration(db, "postgres", Migrations, migrate.Down, 0)
	if migrateErr != nil {
		stdLog.Fatal(migrateErr)
	}

	// Return the size difference between the two sets of migrations
	return len(migrationRecords) - len(allNeededMigrations)
}

// Status returns information about the current status of db migrations. Which
// ones are pending, and when past ones were applied.
//
// From: https://github.com/rubenv/sql-migrate/blob/master/sql-migrate/command_status.go
func Status(db *sql.DB) (string, error) {
	buffer := &bytes.Buffer{}
	migrations, err := Migrations.FindMigrations()
	if err != nil {
		return "", err
	}

	records, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		return "", err
	}

	table := tabwriter.NewWriter(buffer, 60, 8, 0, '\t', 0)
	fmt.Fprintln(table, "Migration\tApplied")

	rows := make(map[string]*statusRow)

	for _, m := range migrations {
		rows[m.Id] = &statusRow{
			Id:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		if rows[r.Id] == nil {
			return "", fmt.Errorf("Could not find migration file: %v", r.Id)
		}

		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range migrations {
		if rows[m.Id] != nil && rows[m.Id].Migrated {
			fmt.Fprintf(table, "%s\t%s\n", m.Id, rows[m.Id].AppliedAt.String())
		} else {
			fmt.Fprintf(table, "%s\tno\n", m.Id)
		}
	}

	if err := table.Flush(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

type statusRow struct {
	Id        string
	Migrated  bool
	AppliedAt time.Time
}
