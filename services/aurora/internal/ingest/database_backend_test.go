package ingest

import (
	"testing"

	"go/ingest/ledgerbackend"
	"go/network"
	"go/services/aurora/internal/test"
)

func TestGetLatestLedger(t *testing.T) {
	tt := test.Start(t)
	tt.ScenarioWithoutaurora("base")
	defer tt.Finish()

	backend, err := ledgerbackend.NewDatabaseBackendFromSession(tt.CoreSession(), network.TestNetworkPassphrase)
	tt.Assert.NoError(err)
	seq, err := backend.GetLatestLedgerSequence(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Equal(uint32(3), seq)
}

func TestGetLatestLedgerNotFound(t *testing.T) {
	tt := test.Start(t)
	tt.ScenarioWithoutaurora("base")
	defer tt.Finish()

	_, err := tt.CoreDB.Exec(`DELETE FROM ledgerheaders`)
	tt.Assert.NoError(err, "failed to remove ledgerheaders")

	backend, err := ledgerbackend.NewDatabaseBackendFromSession(tt.CoreSession(), network.TestNetworkPassphrase)
	tt.Assert.NoError(err)
	_, err = backend.GetLatestLedgerSequence(tt.Ctx)
	tt.Assert.EqualError(err, "no ledgers exist in ledgerheaders table")
}
