package aurora

import (
	"encoding/json"
	"testing"
	"time"

	"go/protocols/aurora"
	"go/protocols/aurora/operations"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest"
)

// Moved to TestGetOperationsOnlyPayments
func TestPaymentActions(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get("/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(4, w.Body)
	}

	// filtered by ledger
	w = ht.Get("/ledgers/1/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}

	w = ht.Get("/ledgers/3/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// Makes StateMiddleware happy
	initializeStateMiddleware := func() {
		q := history.Q{ht.auroraSession()}
		err := q.UpdateLastLedgerIngest(ht.Ctx, 3)
		ht.Assert.NoError(err)
		err = q.UpdateIngestVersion(ht.Ctx, ingest.CurrentVersion)
		ht.Assert.NoError(err)
	}
	initializeStateMiddleware()

	// checks if empty param returns 404 instead of all payments
	w = ht.Get("/accounts//payments")
	ht.Assert.NotEqual(404, w.Code)

	// filtered by account
	w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// switch scenarios
	ht.T.Scenario("pathed_payment")

	// filtered by transaction
	w = ht.Get("/transactions/b52f16ffb98c047e33b9c2ec30880330cde71f85b3443dae2c5cb86c7d4d8452/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}
	// ===========================================
	// The following scenarios are handled in the action test
	// missing tx
	w = ht.Get("/transactions/ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff/payments")
	ht.Assert.Equal(404, w.Code)
	// uppercase tx hash not accepted
	w = ht.Get("/transactions/2374E99349B9EF7DBA9A5DB3339B78FDA8F34777B1AF33BA468AD5C0DF946D4D/payments")
	ht.Assert.Equal(400, w.Code)
	// badly formated tx hash not accepted
	w = ht.Get("/transactions/%00%1E4%5E%EF%BF%BD%EF%BF%BD%EF%BF%BDpVP%EF%BF%BDI&R%0BK%EF%BF%BD%1D%EF%BF%BD%EF%BF%BD=%EF%BF%BD%3F%23%EF%BF%BD%EF%BF%BDl%EF%BF%BD%1El%EF%BF%BD%EF%BF%BD/payments")
	ht.Assert.Equal(400, w.Code)
	// ===========================================

	// TODO: test at the routing level
	// 400 for invalid tx hash
	w = ht.Get("/transactions/ /payments")
	ht.Assert.Equal(400, w.Code)

	// this is handled in operations test, invalid will not match as a valid tx_id.
	w = ht.Get("/transactions/invalid/payments")
	ht.Assert.Equal(400, w.Code)

	// This is already handled in operations test
	w = ht.Get("/transactions/1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc29/payments")
	ht.Assert.Equal(400, w.Code)

	// This is already handled in operations test
	w = ht.Get("/transactions/1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc29222/payments")
	ht.Assert.Equal(400, w.Code)

	w = ht.Get("/transactions/1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc292/payments")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)

		// test for existence of source_amount in path payment details
		var records []map[string]interface{}
		ht.UnmarshalPage(w.Body, &records)
		ht.Assert.Equal("10.0000000", records[0]["source_amount"])
	}

	initializeStateMiddleware()

	// This is tested in PageQueryTest
	// Regression: negative cursor
	w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/payments?cursor=-23667108046966785&order=asc&limit=100")
	ht.Assert.Equal(400, w.Code)
}

func TestPaymentActions_Includetransactions(t *testing.T) {
	ht := StartHTTPTest(t, "base")

	defer ht.Finish()

	w := ht.Get("/payments")
	ht.Assert.Equal(200, w.Code)
	withoutTransactions := []operations.Base{}
	ht.UnmarshalPage(w.Body, &withoutTransactions)

	w = ht.Get("/payments?join=transactions")
	ht.Assert.Equal(200, w.Code)
	withTransactions := []operations.Base{}
	ht.UnmarshalPage(w.Body, &withTransactions)

	for i, operation := range withTransactions {
		getTransaction := ht.Get("/transactions/" + operation.Transaction.ID)
		ht.Assert.Equal(200, getTransaction.Code)
		var getTransactionResponse aurora.Transaction
		err := json.Unmarshal(getTransaction.Body.Bytes(), &getTransactionResponse)

		ht.Require.NoError(err, "failed to parse body")
		tx := operation.Transaction
		ht.Assert.Equal(*tx, getTransactionResponse)

		withTransactions[i].Transaction = nil
	}

	ht.Assert.Equal(withoutTransactions, withTransactions)
}

func TestPaymentActions_Show_Failed(t *testing.T) {
	ht := StartHTTPTest(t, "failed_transactions")
	defer ht.Finish()

	// Should show successful transactions only
	w := ht.Get("/payments?limit=200")

	if ht.Assert.Equal(200, w.Code) {
		records := []operations.Base{}
		ht.UnmarshalPage(w.Body, &records)

		successful := 0
		failed := 0

		for _, op := range records {
			if op.TransactionSuccessful {
				successful++
			} else {
				failed++
			}
		}

		ht.Assert.Equal(5, successful)
		ht.Assert.Equal(0, failed)
	}

	// Should show all transactions: both successful and failed
	w = ht.Get("/payments?limit=200&include_failed=true")

	if ht.Assert.Equal(200, w.Code) {
		records := []operations.Base{}
		ht.UnmarshalPage(w.Body, &records)

		successful := 0
		failed := 0

		for _, op := range records {
			if op.TransactionSuccessful {
				successful++
			} else {
				failed++
			}
		}

		ht.Assert.Equal(5, successful)
		ht.Assert.Equal(1, failed)
	}

	w = ht.Get("/transactions/aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf/payments")

	if ht.Assert.Equal(200, w.Code) {
		records := []operations.Base{}
		ht.UnmarshalPage(w.Body, &records)

		ht.Assert.Equal(1, len(records))
		for _, op := range records {
			ht.Assert.False(op.TransactionSuccessful)
		}
	}

	w = ht.Get("/transactions/56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1/payments")

	if ht.Assert.Equal(200, w.Code) {
		records := []operations.Base{}
		ht.UnmarshalPage(w.Body, &records)

		ht.Assert.Equal(1, len(records))
		for _, op := range records {
			ht.Assert.True(op.TransactionSuccessful)
		}
	}

	// NULL value
	_, err := ht.auroraSession().ExecRaw(ht.Ctx,
		`UPDATE history_transactions SET successful = NULL WHERE transaction_hash = ?`,
		"56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
	)
	ht.Require.NoError(err)

	w = ht.Get("/transactions/56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1/payments")

	if ht.Assert.Equal(200, w.Code) {
		records := []operations.Base{}
		ht.UnmarshalPage(w.Body, &records)

		ht.Assert.Equal(1, len(records))
		for _, op := range records {
			ht.Assert.True(op.TransactionSuccessful)
		}
	}
}

func TestPayment_CreatedAt(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get("/ledgers/3/payments")
	records := []operations.Base{}
	ht.UnmarshalPage(w.Body, &records)

	l := history.Ledger{}
	hq := history.Q{SessionInterface: ht.auroraSession()}
	ht.Require.NoError(hq.LedgerBySequence(ht.Ctx, &l, 3))

	ht.Assert.WithinDuration(l.ClosedAt, records[0].LedgerCloseTime, 1*time.Second)
}

func TestPaymentActions_Show_Extra_TxID(t *testing.T) {
	ht := StartHTTPTest(t, "failed_transactions")
	defer ht.Finish()

	w := ht.Get("/accounts/GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON/payments?limit=200&tx_id=aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf")

	ht.Assert.Equal(400, w.Code)
	payload := ht.UnmarshalExtras(w.Body)
	ht.Assert.Equal("filters", payload["invalid_field"])
	ht.Assert.Equal(
		"Use a single filter for operations, you can only use one of tx_id, account_id or ledger_id",
		payload["reason"],
	)
}

func TestPaymentActionsPathPaymentStrictSend(t *testing.T) {
	ht := StartHTTPTest(t, "paths_strict_send")
	defer ht.Finish()

	w := ht.Get("/payments?order=desc&limit=100")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(11, w.Body)

		var records []map[string]interface{}
		ht.UnmarshalPage(w.Body, &records)

		// Record #1
		ht.Assert.Equal("path_payment_strict_send", records[0]["type"])
		ht.Assert.Equal("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", records[0]["from"])
		ht.Assert.Equal("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", records[0]["to"])
		ht.Assert.Equal("15.8400000", records[0]["amount"])
		ht.Assert.Equal("EUR", records[0]["asset_code"])
		ht.Assert.Equal("GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", records[0]["asset_issuer"])
		ht.Assert.Equal("2.0000000", records[0]["destination_min"])
		ht.Assert.Equal("12.0000000", records[0]["source_amount"])
		ht.Assert.Equal("USD", records[0]["source_asset_code"])
		ht.Assert.Equal("GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", records[0]["source_asset_issuer"])
		ht.Assert.Equal([]interface{}{map[string]interface{}{"asset_type": "native"}}, records[0]["path"])

		// Record #2
		ht.Assert.Equal("path_payment_strict_send", records[1]["type"])
		ht.Assert.Equal("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", records[1]["from"])
		ht.Assert.Equal("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", records[1]["to"])
		ht.Assert.Equal("13.0000000", records[1]["amount"])
		ht.Assert.Equal("EUR", records[1]["asset_code"])
		ht.Assert.Equal("GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", records[1]["asset_issuer"])
		ht.Assert.Equal("1.0000000", records[1]["destination_min"])
		ht.Assert.Equal("10.0000000", records[1]["source_amount"])
		ht.Assert.Equal("USD", records[1]["source_asset_code"])
		ht.Assert.Equal("GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", records[1]["source_asset_issuer"])
		ht.Assert.Equal([]interface{}{}, records[1]["path"])
	}

	// One failed:
	w = ht.Get("/payments?order=desc&include_failed=true&limit=100")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(12, w.Body)

		var records []map[string]interface{}
		ht.UnmarshalPage(w.Body, &records)

		ht.Assert.Equal("path_payment_strict_send", records[0]["type"])
		ht.Assert.Equal("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", records[0]["from"])
		ht.Assert.Equal("GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", records[0]["to"])
		ht.Assert.Equal("0.0000000", records[0]["amount"]) // failed
		ht.Assert.Equal("EUR", records[0]["asset_code"])
		ht.Assert.Equal("GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", records[0]["asset_issuer"])
		ht.Assert.Equal("100.0000000", records[0]["destination_min"])
		ht.Assert.Equal("13.0000000", records[0]["source_amount"])
		ht.Assert.Equal("USD", records[0]["source_asset_code"])
		ht.Assert.Equal("GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", records[0]["source_asset_issuer"])
	}
}
