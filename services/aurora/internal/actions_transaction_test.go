package aurora

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"go/protocols/aurora"
	"go/services/aurora/internal/corestate"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest"
	"go/services/aurora/internal/test"
	"go/xdr"
)

func TestTransactionActions_Show(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get("/transactions/2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.Equal(
			"2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d",
			actual.Hash,
		)
	}

	// missing tx
	w = ht.Get("/transactions/ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	ht.Assert.Equal(404, w.Code)

	// uppercase tx hash not accepted
	w = ht.Get("/transactions/2374E99349B9EF7DBA9A5DB3339B78FDA8F34777B1AF33BA468AD5C0DF946D4D")
	ht.Assert.Equal(400, w.Code)

	// badly formated tx hash not accepted
	w = ht.Get("/transactions/%00%1E4%5E%EF%BF%BD%EF%BF%BD%EF%BF%BDpVP%EF%BF%BDI&R%0BK%EF%BF%BD%1D%EF%BF%BD%EF%BF%BD=%EF%BF%BD%3F%23%EF%BF%BD%EF%BF%BDl%EF%BF%BD%1El%EF%BF%BD%EF%BF%BD")
	ht.Assert.Equal(400, w.Code)
}

func TestTransactionActions_Show_Failed(t *testing.T) {
	ht := StartHTTPTest(t, "failed_transactions")
	defer ht.Finish()

	// Failed single
	w := ht.Get("/transactions/aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.Equal(
			"aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf",
			actual.Hash,
		)

		ht.Assert.False(actual.Successful)
	}

	// Successful single
	w = ht.Get("/transactions/56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.Equal(
			"56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
			actual.Hash,
		)

		ht.Assert.True(actual.Successful)
	}

	// Should show successful transactions only
	w = ht.Get("/transactions?limit=200")

	if ht.Assert.Equal(200, w.Code) {
		records := []aurora.Transaction{}
		ht.UnmarshalPage(w.Body, &records)

		successful := 0
		failed := 0

		for _, tx := range records {
			if tx.Successful {
				successful++
			} else {
				failed++
			}
		}

		ht.Assert.Equal(8, successful)
		ht.Assert.Equal(0, failed)
	}

	// Should show all transactions: both successful and failed
	w = ht.Get("/transactions?limit=200&include_failed=true")

	if ht.Assert.Equal(200, w.Code) {
		records := []aurora.Transaction{}
		ht.UnmarshalPage(w.Body, &records)

		successful := 0
		failed := 0

		for _, tx := range records {
			if tx.Successful {
				successful++
			} else {
				failed++
			}
		}

		ht.Assert.Equal(8, successful)
		ht.Assert.Equal(1, failed)
	}

	w = ht.Get("/transactions/aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.False(actual.Successful)
	}

	w = ht.Get("/transactions/56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.True(actual.Successful)
	}

	// NULL value
	_, err := ht.auroraSession().ExecRaw(ht.Ctx,
		`UPDATE history_transactions SET successful = NULL WHERE transaction_hash = ?`,
		"56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
	)
	ht.Require.NoError(err)

	w = ht.Get("/transactions/56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1")

	if ht.Assert.Equal(200, w.Code) {
		var actual aurora.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		ht.Require.NoError(err)

		ht.Assert.True(actual.Successful)
	}
}

func TestTransactionActions_Index(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	w := ht.Get("/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(4, w.Body)
	}

	// filtered by ledger
	w = ht.Get("/ledgers/1/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}

	w = ht.Get("/ledgers/2/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/ledgers/3/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// missing ledger
	w = ht.Get("/ledgers/100/transactions")
	ht.Assert.Equal(404, w.Code)

	// Makes StateMiddleware happy
	q := history.Q{ht.auroraSession()}
	err := q.UpdateLastLedgerIngest(ht.Ctx, 100)
	ht.Assert.NoError(err)
	err = q.UpdateIngestVersion(ht.Ctx, ingest.CurrentVersion)
	ht.Assert.NoError(err)

	// checks if empty param returns 404 instead of all payments
	w = ht.Get("/accounts//transactions")
	ht.Assert.NotEqual(404, w.Code)

	// filtering by account
	w = ht.Get("/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(2, w.Body)
	}

	// filtering by claimable balance
	w = ht.Get("/claimable_balances/00000000178826fbfe339e1f5c53417c6fedfe2c05e8bec14303143ec46b38981b09c3f9/transactions")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(2, w.Body)
	}

	// Check extra params
	w = ht.Get("/ledgers/100/transactions?account_id=GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H")
	ht.Assert.Equal(400, w.Code)
	w = ht.Get("/ledgers/100/transactions?claimable_balance_id=00000000178826fbfe339e1f5c53417c6fedfe2c05e8bec14303143ec46b38981b09c3f9")
	ht.Assert.Equal(400, w.Code)
	w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/transactions?ledger_id=5")
	ht.Assert.Equal(400, w.Code)
	w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/transactions?cursor=limit=order=")
	ht.Assert.Equal(400, w.Code)

	// regression: https://go/services/aurora/internal/issues/365
	w = ht.Get("/transactions?limit=200")
	ht.Require.Equal(200, w.Code)
	w = ht.Get("/transactions?limit=201")
	ht.Assert.Equal(400, w.Code)
	w = ht.Get("/transactions?limit=0")
	ht.Assert.Equal(400, w.Code)

}

func TestTransactionActions_Post(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// Pass Synced check
	ht.App.coreState.SetState(corestate.State{Synced: true})

	tx := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTxV0,
		V0: &xdr.TransactionV0Envelope{
			Tx: xdr.TransactionV0{
				SourceAccountEd25519: *xdr.MustAddress("GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H").Ed25519,
				Fee:                  100,
				SeqNum:               1,
				Operations: []xdr.Operation{
					{
						Body: xdr.OperationBody{
							Type: xdr.OperationTypeCreateAccount,
							CreateAccountOp: &xdr.CreateAccountOp{
								Destination:     xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"),
								StartingBalance: 1000000000,
							},
						},
					},
				},
			},
			Signatures: []xdr.DecoratedSignature{
				{
					Hint:      xdr.SignatureHint{86, 252, 5, 247},
					Signature: xdr.Signature{131, 206, 171, 228, 64, 20, 40, 52, 2, 98, 124, 244, 87, 14, 130, 225, 190, 220, 156, 79, 121, 69, 60, 36, 57, 214, 9, 29, 176, 81, 218, 4, 213, 176, 211, 148, 191, 86, 21, 180, 94, 9, 43, 208, 32, 79, 19, 131, 90, 21, 93, 138, 153, 203, 55, 103, 2, 230, 137, 190, 19, 70, 179, 11},
				},
			},
		},
	}

	txStr, err := xdr.MarshalBase64(tx)
	assert.NoError(t, err)
	form := url.Values{"tx": []string{txStr}}

	// existing transaction
	w := ht.Post("/transactions", form)
	ht.Assert.Equal(200, w.Code)
}

func TestTransactionActions_PostSuccessful(t *testing.T) {
	ht := StartHTTPTest(t, "failed_transactions")
	defer ht.Finish()

	// Pass Synced check
	ht.App.coreState.SetState(corestate.State{Synced: true})

	destAID := xdr.MustAddress("GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON")
	tx2 := xdr.TransactionEnvelope{
		Type: xdr.EnvelopeTypeEnvelopeTypeTxV0,
		V0: &xdr.TransactionV0Envelope{
			Tx: xdr.TransactionV0{
				SourceAccountEd25519: *xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU").Ed25519,
				Fee:                  100,
				SeqNum:               8589934593,
				Operations: []xdr.Operation{
					{
						Body: xdr.OperationBody{
							Type: xdr.OperationTypePayment,
							PaymentOp: &xdr.PaymentOp{
								Destination: destAID.ToMuxedAccount(),
								Asset: xdr.Asset{
									Type: xdr.AssetTypeAssetTypeCreditAlphanum4,
									AlphaNum4: &xdr.AlphaNum4{
										AssetCode: xdr.AssetCode4{85, 83, 68, 0},
										Issuer:    xdr.MustAddress("GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"),
									},
								},
								Amount: 1000000000,
							},
						},
					},
				},
			},
			Signatures: []xdr.DecoratedSignature{
				{
					Hint:      xdr.SignatureHint{174, 228, 190, 76},
					Signature: xdr.Signature{73, 202, 13, 176, 216, 188, 169, 9, 141, 130, 180, 106, 187, 225, 22, 89, 254, 24, 173, 62, 236, 12, 186, 131, 70, 190, 214, 24, 209, 69, 233, 68, 1, 238, 48, 154, 55, 170, 53, 196, 96, 218, 110, 2, 159, 187, 120, 2, 50, 115, 2, 192, 208, 35, 72, 151, 106, 17, 155, 160, 147, 200, 52, 12},
				},
			},
		},
	}

	txStr, err := xdr.MarshalBase64(tx2)
	assert.NoError(t, err)

	// 56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1
	form := url.Values{"tx": []string{txStr}}

	w := ht.Post("/transactions", form)
	ht.Assert.Equal(200, w.Code)
	ht.Assert.Contains(string(w.Body.Bytes()), `"result_xdr": "AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA="`)
}

func TestTransactionActions_PostFailed(t *testing.T) {
	ht := StartHTTPTest(t, "failed_transactions")
	defer ht.Finish()

	// Pass Synced check
	ht.App.coreState.SetState(corestate.State{Synced: true})

	// aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf
	form := url.Values{"tx": []string{"AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAB3NZQAAAAAAAAAAAFvFIhaAAAAQKcGS9OsVnVHCVIH04C9ZKzzKYBRdCmy+Jwmzld7QcALOxZUcAgkuGfoSdvXpH38mNvrqQiaMsSNmTJWYRzHvgo="}}

	w := ht.Post("/transactions", form)
	ht.Assert.Equal(400, w.Code)
	ht.Assert.Contains(string(w.Body.Bytes()), "op_underfunded")
	ht.Assert.Contains(string(w.Body.Bytes()), `"result_xdr": "AAAAAAAAAGT/////AAAAAQAAAAAAAAAB/////gAAAAA="`)
}

func TestPostFeeBumpTransaction(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	// Pass Synced check
	ht.App.coreState.SetState(corestate.State{Synced: true})

	test.ResetauroraDB(t, ht.auroraDB)
	q := &history.Q{ht.auroraSession()}
	fixture := history.FeeBumpScenario(ht.T, q, true)

	form := url.Values{"tx": []string{fixture.Transaction.TxEnvelope}}
	w := ht.Post("/transactions", form)
	ht.Assert.Equal(200, w.Code)
	var response aurora.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	ht.Assert.NoError(err)

	ht.Assert.Equal(fixture.Transaction.TxResult, response.ResultXdr)
	ht.Assert.Equal(fixture.Transaction.TxMeta, response.ResultMetaXdr)
	ht.Assert.Equal(fixture.Transaction.TransactionHash, response.Hash)
	ht.Assert.Equal(fixture.Transaction.TxEnvelope, response.EnvelopeXdr)
	ht.Assert.Equal(fixture.Transaction.LedgerSequence, response.Ledger)

	innerTxEnvelope, err := xdr.MarshalBase64(fixture.Envelope.FeeBump.Tx.InnerTx.V1)
	ht.Assert.NoError(err)
	form = url.Values{"tx": []string{innerTxEnvelope}}
	w = ht.Post("/transactions", form)
	ht.Assert.Equal(200, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &response)
	ht.Assert.NoError(err)

	ht.Assert.Equal(fixture.Transaction.TxResult, response.ResultXdr)
	ht.Assert.Equal(fixture.Transaction.TxMeta, response.ResultMetaXdr)
	ht.Assert.Equal(fixture.InnerHash, response.Hash)
	ht.Assert.Equal(fixture.Transaction.TxEnvelope, response.EnvelopeXdr)
	ht.Assert.Equal(fixture.Transaction.LedgerSequence, response.Ledger)
}

func TestPostFailedFeeBumpTransaction(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()

	// Pass Synced check
	ht.App.coreState.SetState(corestate.State{Synced: true})

	test.ResetauroraDB(t, ht.auroraDB)
	q := &history.Q{ht.auroraSession()}
	fixture := history.FeeBumpScenario(ht.T, q, false)

	form := url.Values{"tx": []string{fixture.Transaction.TxEnvelope}}
	w := ht.Post("/transactions", form)
	ht.Assert.Equal(400, w.Code)
	ht.Assert.Contains(string(w.Body.Bytes()), "tx_fee_bump_inner_failed")
	ht.Assert.Contains(string(w.Body.Bytes()), "tx_bad_auth")

	innerTxEnvelope, err := xdr.MarshalBase64(fixture.Envelope.FeeBump.Tx.InnerTx.V1)
	ht.Assert.NoError(err)
	form = url.Values{"tx": []string{innerTxEnvelope}}
	w = ht.Post("/transactions", form)
	ht.Assert.Equal(400, w.Code)
	ht.Assert.Contains(string(w.Body.Bytes()), "tx_bad_auth")
	ht.Assert.NotContains(string(w.Body.Bytes()), "tx_fee_bump_inner_failed")
}
