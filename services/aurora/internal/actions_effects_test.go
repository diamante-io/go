package aurora

import (
	"testing"

	"go/protocols/aurora/effects"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest"
	"go/services/aurora/internal/test"
)

func TestEffectActions_Index(t *testing.T) {

	t.Run("omnibus test", func(t *testing.T) {

		ht := StartHTTPTest(t, "base")
		defer ht.Finish()

		w := ht.Get("/effects?limit=20")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(11, w.Body)
		}

		// test streaming, regression for https://go/services/aurora/internal/issues/147
		w = ht.Get("/effects?limit=2", test.RequestHelperStreaming)
		ht.Assert.Equal(200, w.Code)

		// filtered by ledger
		w = ht.Get("/ledgers/1/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(0, w.Body)
		}

		w = ht.Get("/ledgers/2/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(9, w.Body)
		}

		w = ht.Get("/ledgers/3/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(2, w.Body)
		}

		// Makes StateMiddleware happy
		q := history.Q{ht.auroraSession()}
		err := q.UpdateLastLedgerIngest(ht.Ctx, 3)
		ht.Assert.NoError(err)
		err = q.UpdateIngestVersion(ht.Ctx, ingest.CurrentVersion)
		ht.Assert.NoError(err)

		// checks if empty param returns 404 instead of all payments
		w = ht.Get("/accounts//effects")
		ht.Assert.NotEqual(404, w.Code)

		// filtered by account
		w = ht.Get("/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(3, w.Body)
		}

		w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(2, w.Body)
		}

		w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(3, w.Body)
		}

		// filtered by transaction
		w = ht.Get("/transactions/2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(3, w.Body)
		}
		// missing tx
		w = ht.Get("/transactions/ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff/effects")
		ht.Assert.Equal(404, w.Code)
		// uppercase tx hash not accepted
		w = ht.Get("/transactions/2374E99349B9EF7DBA9A5DB3339B78FDA8F34777B1AF33BA468AD5C0DF946D4D/effects")
		ht.Assert.Equal(400, w.Code)
		// badly formated tx hash not accepted
		w = ht.Get("/transactions/%00%1E4%5E%EF%BF%BD%EF%BF%BD%EF%BF%BDpVP%EF%BF%BDI&R%0BK%EF%BF%BD%1D%EF%BF%BD%EF%BF%BD=%EF%BF%BD%3F%23%EF%BF%BD%EF%BF%BDl%EF%BF%BD%1El%EF%BF%BD%EF%BF%BD/effects")
		ht.Assert.Equal(400, w.Code)

		w = ht.Get("/transactions/2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(3, w.Body)
		}

		// filtered by operation
		w = ht.Get("/operations/8589938689/effects")
		if ht.Assert.Equal(200, w.Code) {
			ht.Assert.PageOf(3, w.Body)
		}

		// Check extra params
		w = ht.Get("/ledgers/100/effects?account_id=GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
		ht.Assert.Equal(400, w.Code)
		w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/effects?ledger_id=5")
		ht.Assert.Equal(400, w.Code)

		// before history
		ht.ReapHistory(1)
		w = ht.Get("/effects?order=desc&cursor=8589938689-1")
		ht.Assert.Equal(410, w.Code)
		ht.T.T.Log(w.Body.String())
	})

	t.Run("Effect resource props", func(t *testing.T) {
		ht := StartHTTPTest(t, "base")
		defer ht.Finish()

		// created_at
		w := ht.Get("/ledgers/2/effects")
		if ht.Assert.Equal(200, w.Code) {
			var result []effects.Base
			_ = ht.UnmarshalPage(w.Body, &result)
			ht.Require.NotEmpty(result, "unexpected empty response")

			e1 := result[0]

			var ledger2 history.Ledger
			err := ht.auroraDB.Get(&ledger2, "SELECT * FROM history_ledgers WHERE sequence = 2")
			ht.Require.NoError(err, "failed to load ledger")

			ht.Assert.Equal(ledger2.ClosedAt.UTC(), e1.LedgerCloseTime.UTC())
		}
	})
}

func TestEffectsForFeeBumpTransaction(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	test.ResetauroraDB(t, ht.auroraDB)
	q := &history.Q{ht.auroraSession()}
	fixture := history.FeeBumpScenario(ht.T, q, true)

	w := ht.Get("/transactions/" + fixture.OuterHash + "/effects")
	ht.Assert.Equal(200, w.Code)
	var byOuterHash []effects.Base
	ht.UnmarshalPage(w.Body, &byOuterHash)
	ht.Assert.Len(byOuterHash, 1)

	w = ht.Get("/transactions/" + fixture.InnerHash + "/effects")
	ht.Assert.Equal(200, w.Code)
	var byInnerHash []effects.Base
	ht.UnmarshalPage(w.Body, &byInnerHash)
	ht.Assert.Len(byInnerHash, 1)

	ht.Assert.Equal(byOuterHash, byInnerHash)
}
