//lint:file-ignore U1001 Ignore all unused code, thinks the code is unused because of the test skips
package test

import (
	"io"
	"testing"

	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/diamcircle/go/services/aurora/internal/db2/schema"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/operationfeestats"
	"github.com/diamcircle/go/services/aurora/internal/test/scenarios"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/render/hal"
)

// CoreSession returns a db.Session instance pointing at the diamcircle core test database
func (t *T) CoreSession() *db.Session {
	return &db.Session{
		DB: t.CoreDB,
	}
}

// Finish finishes the test, logging any accumulated aurora logs to the logs
// output
func (t *T) Finish() {
	logEntries := t.testLogs()
	operationfeestats.ResetState()

	for _, entry := range logEntries {
		logString, err := entry.String()
		if err != nil {
			t.T.Logf("Error from entry.String: %v", err)
		} else {
			t.T.Log(logString)
		}
	}
}

// AuroraSession returns a db.Session instance pointing at the aurora test
// database
func (t *T) AuroraSession() *db.Session {
	return &db.Session{
		DB: t.AuroraDB,
	}
}

func (t *T) loadScenario(scenarioName string, includeAurora bool) {
	diamcircleCorePath := scenarioName + "-core.sql"

	scenarios.Load(DiamcircleCoreDatabaseURL(), diamcircleCorePath)

	if includeAurora {
		auroraPath := scenarioName + "-aurora.sql"
		scenarios.Load(DatabaseURL(), auroraPath)
	}
}

// Scenario loads the named sql scenario into the database
func (t *T) Scenario(name string) ledger.Status {
	clearAuroraDB(t.T, t.AuroraDB)
	t.loadScenario(name, true)
	return t.LoadLedgerStatus()
}

// ScenarioWithoutAurora loads the named sql scenario into the database
func (t *T) ScenarioWithoutAurora(name string) ledger.Status {
	t.loadScenario(name, false)
	ResetAuroraDB(t.T, t.AuroraDB)
	return t.LoadLedgerStatus()
}

// ResetAuroraDB sets up a new aurora database with empty tables
func ResetAuroraDB(t *testing.T, db *sqlx.DB) {
	clearAuroraDB(t, db)
	_, err := schema.Migrate(db.DB, schema.MigrateUp, 0)
	if err != nil {
		t.Fatalf("could not run migrations up on test db: %v", err)
	}
}

func clearAuroraDB(t *testing.T, db *sqlx.DB) {
	_, err := schema.Migrate(db.DB, schema.MigrateDown, 0)
	if err != nil {
		t.Fatalf("could not run migrations down on test db: %v", err)
	}
}

// UnmarshalPage populates dest with the records contained in the json-encoded page in r
func (t *T) UnmarshalPage(r io.Reader, dest interface{}) hal.Links {
	var env struct {
		Embedded struct {
			Records json.RawMessage `json:"records"`
		} `json:"_embedded"`
		Links struct {
			Self hal.Link `json:"self"`
			Next hal.Link `json:"next"`
			Prev hal.Link `json:"prev"`
		} `json:"_links"`
	}

	err := json.NewDecoder(r).Decode(&env)
	t.Require.NoError(err, "failed to decode page")

	err = json.Unmarshal(env.Embedded.Records, dest)
	t.Require.NoError(err, "failed to decode records")

	return env.Links
}

// UnmarshalNext extracts and returns the next link
func (t *T) UnmarshalNext(r io.Reader) string {
	var env struct {
		Links struct {
			Next struct {
				Href string `json:"href"`
			} `json:"next"`
		} `json:"_links"`
	}

	err := json.NewDecoder(r).Decode(&env)
	t.Require.NoError(err, "failed to decode page")
	return env.Links.Next.Href
}

// UnmarshalExtras extracts and returns extras content
func (t *T) UnmarshalExtras(r io.Reader) map[string]string {
	var resp struct {
		Extras map[string]string `json:"extras"`
	}

	err := json.NewDecoder(r).Decode(&resp)
	t.Require.NoError(err, "failed to decode page")

	return resp.Extras
}

// LoadLedgerStatus loads ledger state from the core db(or panicing on failure).
func (t *T) LoadLedgerStatus() ledger.Status {
	var next ledger.Status

	err := t.CoreSession().GetRaw(t.Ctx, &next, `
		SELECT
			COALESCE(MAX(ledgerseq), 0) as core_latest
		FROM ledgerheaders
	`)

	if err != nil {
		panic(err)
	}

	err = t.AuroraSession().GetRaw(t.Ctx, &next, `
			SELECT
				COALESCE(MIN(sequence), 0) as history_elder,
				COALESCE(MAX(sequence), 0) as history_latest
			FROM history_ledgers
		`)

	if err != nil {
		panic(err)
	}

	return next
}

// retrieves entries from test logger instance
func (t *T) testLogs() []logrus.Entry {
	if t.EndLogTest == nil {
		return []logrus.Entry{}
	}

	return t.EndLogTest()
}
