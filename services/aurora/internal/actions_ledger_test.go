package aurora

import (
	"encoding/json"
	"testing"

	"github.com/diamcircle/go/protocols/aurora"
)

func TestLedgerActions_Index(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// default params
	w := ht.Get("/ledgers")

	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	// with limit
	w = ht.RH.Get("/ledgers?limit=1")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}
}

func TestLedgerActions_Show(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get("/ledgers/2")
	ht.Assert.Equal(200, w.Code)

	var result aurora.Ledger
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if ht.Assert.NoError(err) {
		ht.Assert.Equal(int32(2), result.Sequence)
		ht.Assert.NotEmpty(result.HeaderXDR)
		ht.Assert.Equal(int32(3), result.SuccessfulTransactionCount)
		ht.Assert.Equal(int32(0), *result.FailedTransactionCount)
		ht.Assert.Nil(result.TxSetOperationCount)
	}

	// There's no way to test previous versions of ingestion right now
	// so let's manually update the state to look like version 14 of ingesiton
	// only the latest gap is considered for determining the elder ledger
	_, err = ht.AuroraDB.Exec(`
		UPDATE history_ledgers SET successful_transaction_count = NULL, failed_transaction_count = NULL, tx_set_operation_count = 5 WHERE sequence = 2
	`)
	ht.Require.NoError(err, "failed to update history_ledgers")

	w = ht.Get("/ledgers/2")
	ht.Assert.Equal(200, w.Code)

	result = aurora.Ledger{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if ht.Assert.NoError(err) {
		ht.Assert.Equal(int32(2), result.Sequence)
		ht.Assert.NotEmpty(result.HeaderXDR)
		ht.Assert.Equal(int32(3), result.SuccessfulTransactionCount)
		ht.Assert.Nil(result.FailedTransactionCount)
		ht.Assert.Equal(int32(5), *result.TxSetOperationCount)
	}

	// ledger higher than history
	w = ht.Get("/ledgers/100")
	ht.Assert.Equal(404, w.Code)

	// ledger that was reaped
	ht.ReapHistory(1)

	w = ht.Get("/ledgers/1")
	ht.Assert.Equal(410, w.Code)
}
