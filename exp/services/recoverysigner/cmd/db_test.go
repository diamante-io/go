package cmd

import (
	"testing"

	dbpkg "go/exp/services/recoverysigner/internal/db"
	"go/exp/services/recoverysigner/internal/db/dbtest"
	"go/support/log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBCommand_Migrate_upDownAll(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	log := log.New()

	dbCommand := DBCommand{
		Logger:      log,
		DatabaseURL: db.DSN,
	}

	// Migrate Up
	{
		logsGet := log.StartTest(logrus.InfoLevel)

		dbCommand.Migrate(&cobra.Command{}, []string{"up"})

		session, err := dbpkg.Open(db.DSN)
		require.NoError(t, err)
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

		logs := logsGet()
		messages := []string{}
		for _, l := range logs {
			messages = append(messages, l.Message)
		}
		wantMessages := []string{
			"Migrations to apply up: 20200309000000-initial-1.sql, 20200309000001-initial-2.sql, 20200311000000-create-accounts.sql, 20200311000001-create-identities.sql, 20200311000002-create-auth-methods.sql, 20200320000000-create-accounts-audit.sql, 20200320000001-create-identities-audit.sql, 20200320000002-create-auth-methods-audit.sql",
			"Successfully applied 8 migrations up.",
		}
		assert.Equal(t, wantMessages, messages)
	}

	// Migrate Down
	{
		logsGet := log.StartTest(logrus.InfoLevel)

		dbCommand.Migrate(&cobra.Command{}, []string{"down"})

		session, err := dbpkg.Open(db.DSN)
		require.NoError(t, err)
		ids := []string{}
		err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
		require.NoError(t, err)
		assert.Empty(t, ids)

		logs := logsGet()
		messages := []string{}
		for _, l := range logs {
			messages = append(messages, l.Message)
		}
		wantMessages := []string{
			"Migrations to apply down: 20200320000002-create-auth-methods-audit.sql, 20200320000001-create-identities-audit.sql, 20200320000000-create-accounts-audit.sql, 20200311000002-create-auth-methods.sql, 20200311000001-create-identities.sql, 20200311000000-create-accounts.sql, 20200309000001-initial-2.sql, 20200309000000-initial-1.sql",
			"Successfully applied 8 migrations down.",
		}
		assert.Equal(t, wantMessages, messages)
	}
}

func TestDBCommand_Migrate_upTwoDownOne(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	log := log.New()

	dbCommand := DBCommand{
		Logger:      log,
		DatabaseURL: db.DSN,
	}

	// Migrate Up 2
	{
		logsGet := log.StartTest(logrus.InfoLevel)

		dbCommand.Migrate(&cobra.Command{}, []string{"up", "2"})

		session, err := dbpkg.Open(db.DSN)
		require.NoError(t, err)
		ids := []string{}
		err = session.Unsafe().Select(&ids, `SELECT id FROM gorp_migrations`)
		require.NoError(t, err)
		wantIDs := []string{
			"20200309000000-initial-1.sql",
			"20200309000001-initial-2.sql",
		}
		assert.Equal(t, wantIDs, ids)

		logs := logsGet()
		messages := []string{}
		for _, l := range logs {
			messages = append(messages, l.Message)
		}
		wantMessages := []string{
			"Migrations to apply up: 20200309000000-initial-1.sql, 20200309000001-initial-2.sql",
			"Successfully applied 2 migrations up.",
		}
		assert.Equal(t, wantMessages, messages)
	}

	// Migrate Down 1
	{
		logsGet := log.StartTest(logrus.InfoLevel)

		dbCommand.Migrate(&cobra.Command{}, []string{"down", "1"})

		session, err := dbpkg.Open(db.DSN)
		require.NoError(t, err)
		ids := []string{}
		err = session.Select(&ids, `SELECT id FROM gorp_migrations`)
		require.NoError(t, err)
		wantIDs := []string{
			"20200309000000-initial-1.sql",
		}
		assert.Equal(t, wantIDs, ids)

		logs := logsGet()
		messages := []string{}
		for _, l := range logs {
			messages = append(messages, l.Message)
		}
		wantMessages := []string{
			"Migrations to apply down: 20200309000001-initial-2.sql",
			"Successfully applied 1 migrations down.",
		}
		assert.Equal(t, wantMessages, messages)
	}
}

func TestDBCommand_Migrate_invalidDirection(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	log := log.New()

	dbCommand := DBCommand{
		Logger:      log,
		DatabaseURL: db.DSN,
	}

	logsGet := log.StartTest(logrus.InfoLevel)

	dbCommand.Migrate(&cobra.Command{}, []string{"invalid"})

	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)
	tables := []string{}
	err = session.Select(&tables, `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`)
	require.NoError(t, err)
	assert.Empty(t, tables)

	logs := logsGet()
	messages := []string{}
	for _, l := range logs {
		messages = append(messages, l.Message)
	}
	wantMessages := []string{
		"Invalid migration direction, must be 'up' or 'down'.",
	}
	assert.Equal(t, wantMessages, messages)
}

func TestDBCommand_Migrate_invalidCount(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	log := log.New()

	dbCommand := DBCommand{
		Logger:      log,
		DatabaseURL: db.DSN,
	}

	logsGet := log.StartTest(logrus.InfoLevel)

	dbCommand.Migrate(&cobra.Command{}, []string{"down", "invalid"})
	dbCommand.Migrate(&cobra.Command{}, []string{"up", "invalid"})

	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)
	tables := []string{}
	err = session.Select(&tables, `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`)
	require.NoError(t, err)
	assert.Empty(t, tables)

	logs := logsGet()
	messages := []string{}
	for _, l := range logs {
		messages = append(messages, l.Message)
	}
	wantMessages := []string{
		"Invalid migration count, must be a number.",
		"Invalid migration count, must be a number.",
	}
	assert.Equal(t, wantMessages, messages)
}

func TestDBCommand_Migrate_zeroCount(t *testing.T) {
	db := dbtest.OpenWithoutMigrations(t)
	log := log.New()

	dbCommand := DBCommand{
		Logger:      log,
		DatabaseURL: db.DSN,
	}

	logsGet := log.StartTest(logrus.InfoLevel)

	dbCommand.Migrate(&cobra.Command{}, []string{"down", "0"})
	dbCommand.Migrate(&cobra.Command{}, []string{"up", "0"})

	session, err := dbpkg.Open(db.DSN)
	require.NoError(t, err)
	tables := []string{}
	err = session.Select(&tables, `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`)
	require.NoError(t, err)
	assert.Empty(t, tables)

	logs := logsGet()
	messages := []string{}
	for _, l := range logs {
		messages = append(messages, l.Message)
	}
	wantMessages := []string{
		"Invalid migration count, must be a number greater than zero.",
		"Invalid migration count, must be a number greater than zero.",
	}
	assert.Equal(t, wantMessages, messages)
}
