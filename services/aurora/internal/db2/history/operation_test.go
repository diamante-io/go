package history

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/test"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/xdr"
)

func TestOperationQueries(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("base")
	defer tt.Finish()
	q := &Q{tt.AuroraSession()}

	// Test OperationByID
	op, transaction, err := q.OperationByID(tt.Ctx, false, 8589938689)
	if tt.Assert.NoError(err) {
		tt.Assert.Equal(int64(8589938689), op.ID)
	}
	tt.Assert.Nil(transaction)

	// Test Operations()
	ops, transactions, err := q.Operations().
		ForAccount(tt.Ctx, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON").
		Fetch(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 2)
	}
	tt.Assert.Len(transactions, 0)

	// ledger filter works
	ops, transactions, err = q.Operations().ForLedger(tt.Ctx, 2).Fetch(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 3)
	}
	tt.Assert.Len(transactions, 0)

	// tx filter works
	hash := "2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d"
	ops, transactions, err = q.Operations().ForTransaction(tt.Ctx, hash).Fetch(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 1)
	}
	tt.Assert.Len(transactions, 0)

	// payment filter works
	tt.Scenario("pathed_payment")
	ops, transactions, err = q.Operations().OnlyPayments().Fetch(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 10)
	}
	tt.Assert.Len(transactions, 0)

	// payment filter includes account merges
	tt.Scenario("account_merge")
	ops, transactions, err = q.Operations().OnlyPayments().Fetch(tt.Ctx)
	if tt.Assert.NoError(err) {
		tt.Assert.Len(ops, 3)
	}
	tt.Assert.Len(transactions, 0)
}

func TestOperationByLiquidityPool(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetAuroraDB(t, tt.AuroraDB)
	q := &Q{tt.AuroraSession()}

	txIndex := int32(1)
	sequence := int32(56)
	txID := toid.New(sequence, txIndex, 0).ToInt64()
	opID1 := toid.New(sequence, txIndex, 1).ToInt64()
	opID2 := toid.New(sequence, txIndex, 2).ToInt64()

	// Insert a phony transaction
	transactionBuilder := q.NewTransactionBatchInsertBuilder(2)
	firstTransaction := buildLedgerTransaction(tt.T, testTransaction{
		index:         uint32(txIndex),
		envelopeXDR:   "AAAAACiSTRmpH6bHC6Ekna5e82oiGY5vKDEEUgkq9CB//t+rAAAAyAEXUhsAADDRAAAAAAAAAAAAAAABAAAAAAAAAAsBF1IbAABX4QAAAAAAAAAA",
		resultXDR:     "AAAAAAAAASwAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAFAAAAAAAAAAA=",
		feeChangesXDR: "AAAAAA==",
		metaXDR:       "AAAAAQAAAAAAAAAA",
		hash:          "19aaa18db88605aedec04659fb45e06f240b022eb2d429e05133e4d53cd945ba",
	})
	err := transactionBuilder.Add(tt.Ctx, firstTransaction, uint32(sequence))
	tt.Assert.NoError(err)
	err = transactionBuilder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	// Insert a two phony operations
	operationBuilder := q.NewOperationBatchInsertBuilder(2)
	err = operationBuilder.Add(
		tt.Ctx,
		opID1,
		txID,
		1,
		xdr.OperationTypeEndSponsoringFutureReserves,
		[]byte("{}"),
		"GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
		null.String{},
	)
	tt.Assert.NoError(err)
	err = operationBuilder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	err = operationBuilder.Add(
		tt.Ctx,
		opID2,
		txID,
		1,
		xdr.OperationTypeEndSponsoringFutureReserves,
		[]byte("{}"),
		"GAUJETIZVEP2NRYLUESJ3LS66NVCEGMON4UDCBCSBEVPIID773P2W6AY",
		null.String{},
	)
	tt.Assert.NoError(err)
	err = operationBuilder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	// Insert Liquidity Pool history
	liquidityPoolID := "a2f38836a839de008cf1d782c81f45e1253cc5d3dad9110b872965484fec0a49"
	toInternalID, err := q.CreateHistoryLiquidityPools(tt.Ctx, []string{liquidityPoolID}, 2)
	tt.Assert.NoError(err)
	lpOperationBuilder := q.NewOperationLiquidityPoolBatchInsertBuilder(3)
	tt.Assert.NoError(err)
	internalID, ok := toInternalID[liquidityPoolID]
	tt.Assert.True(ok)
	err = lpOperationBuilder.Add(tt.Ctx, opID1, internalID)
	tt.Assert.NoError(err)
	err = lpOperationBuilder.Add(tt.Ctx, opID2, internalID)
	tt.Assert.NoError(err)
	err = lpOperationBuilder.Exec(tt.Ctx)
	tt.Assert.NoError(err)

	// Check ascending order
	pq := db2.PageQuery{
		Cursor: "",
		Order:  "asc",
		Limit:  2,
	}
	ops, _, err := q.Operations().ForLiquidityPool(tt.Ctx, liquidityPoolID).Page(pq).Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(ops, 2)
	tt.Assert.Equal(ops[0].ID, opID1)
	tt.Assert.Equal(ops[1].ID, opID2)

	// Check descending order
	pq.Order = "desc"
	ops, _, err = q.Operations().ForLiquidityPool(tt.Ctx, liquidityPoolID).Page(pq).Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(ops, 2)
	tt.Assert.Equal(ops[0].ID, opID2)
	tt.Assert.Equal(ops[1].ID, opID1)
}

func TestOperationQueryBuilder(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("base")
	defer tt.Finish()
	q := &Q{tt.AuroraSession()}

	opsQ := q.Operations().ForAccount(tt.Ctx, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON").Page(db2.PageQuery{Cursor: "8589938689", Order: "asc", Limit: 10})
	tt.Assert.NoError(opsQ.Err)
	got, _, err := opsQ.sql.ToSql()
	tt.Assert.NoError(err)

	// Operations for account queries will use hopp.history_operation_id in their predicates.
	want := "SELECT hop.id, hop.transaction_id, hop.application_order, hop.type, hop.details, hop.source_account, hop.source_account_muxed, ht.transaction_hash, ht.tx_result, COALESCE(ht.successful, true) as transaction_successful FROM history_operations hop LEFT JOIN history_transactions ht ON ht.id = hop.transaction_id JOIN history_operation_participants hopp ON hopp.history_operation_id = hop.id WHERE hopp.history_account_id = ? AND hopp.history_operation_id > ? ORDER BY hopp.history_operation_id asc LIMIT 10"
	tt.Assert.EqualValues(want, got)

	opsQ = q.Operations().ForLedger(tt.Ctx, 2).Page(db2.PageQuery{Cursor: "8589938689", Order: "asc", Limit: 10})
	tt.Assert.NoError(opsQ.Err)
	got, _, err = opsQ.sql.ToSql()
	tt.Assert.NoError(err)

	// Other operation queries will use hop.id in their predicates.
	want = "SELECT hop.id, hop.transaction_id, hop.application_order, hop.type, hop.details, hop.source_account, hop.source_account_muxed, ht.transaction_hash, ht.tx_result, COALESCE(ht.successful, true) as transaction_successful FROM history_operations hop LEFT JOIN history_transactions ht ON ht.id = hop.transaction_id WHERE hop.id >= ? AND hop.id < ? AND hop.id > ? ORDER BY hop.id asc LIMIT 10"
	tt.Assert.EqualValues(want, got)
}

// TestOperationSuccessfulOnly tests if default query returns operations in
// successful transactions only.
// If it's not enclosed in brackets, it may return incorrect result when mixed
// with `ForAccount` or `ForLedger` filters.
func TestOperationSuccessfulOnly(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		ForAccount(tt.Ctx, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2")

	operations, transactions, err := query.Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(transactions, 0)

	tt.Assert.Equal(3, len(operations))

	for _, operation := range operations {
		tt.Assert.True(operation.TransactionSuccessful)
	}

	sql, _, err := query.sql.ToSql()
	tt.Assert.NoError(err)
	// Note: brackets around `(ht.successful = true OR ht.successful IS NULL)` are critical!
	tt.Assert.Contains(sql, "WHERE hopp.history_account_id = ? AND (ht.successful = true OR ht.successful IS NULL)")
}

// TestOperationIncludeFailed tests `IncludeFailed` method.
func TestOperationIncludeFailed(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		ForAccount(tt.Ctx, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2").
		IncludeFailed()

	operations, transactions, err := query.Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(transactions, 0)

	var failed, successful int
	for _, operation := range operations {
		if operation.TransactionSuccessful {
			successful++
		} else {
			failed++
		}
	}

	tt.Assert.Equal(3, successful)
	tt.Assert.Equal(1, failed)

	sql, _, err := query.sql.ToSql()
	tt.Assert.NoError(err)
	tt.Assert.Equal("SELECT hop.id, hop.transaction_id, hop.application_order, hop.type, hop.details, hop.source_account, hop.source_account_muxed, ht.transaction_hash, ht.tx_result, COALESCE(ht.successful, true) as transaction_successful FROM history_operations hop LEFT JOIN history_transactions ht ON ht.id = hop.transaction_id JOIN history_operation_participants hopp ON hopp.history_operation_id = hop.id WHERE hopp.history_account_id = ?", sql)
}

// TestPaymentsSuccessfulOnly tests if default query returns payments in
// successful transactions only.
// If it's not enclosed in brackets, it may return incorrect result when mixed
// with `ForAccount` or `ForLedger` filters.
func TestPaymentsSuccessfulOnly(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		OnlyPayments().
		ForAccount(tt.Ctx, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON")

	operations, transactions, err := query.Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(transactions, 0)

	tt.Assert.Equal(2, len(operations))

	for _, operation := range operations {
		tt.Assert.True(operation.TransactionSuccessful)
	}

	sql, _, err := query.sql.ToSql()
	tt.Assert.NoError(err)
	// Note: brackets around `(ht.successful = true OR ht.successful IS NULL)` are critical!
	tt.Assert.Contains(sql, "WHERE hop.type IN (?,?,?,?,?) AND hopp.history_account_id = ? AND (ht.successful = true OR ht.successful IS NULL)")
}

// TestPaymentsIncludeFailed tests `IncludeFailed` method.
func TestPaymentsIncludeFailed(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		OnlyPayments().
		ForAccount(tt.Ctx, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON").
		IncludeFailed()

	operations, transactions, err := query.Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(transactions, 0)

	var failed, successful int
	for _, operation := range operations {
		if operation.TransactionSuccessful {
			successful++
		} else {
			failed++
		}
	}

	tt.Assert.Equal(2, successful)
	tt.Assert.Equal(1, failed)

	sql, _, err := query.sql.ToSql()
	tt.Assert.NoError(err)
	tt.Assert.Equal("SELECT hop.id, hop.transaction_id, hop.application_order, hop.type, hop.details, hop.source_account, hop.source_account_muxed, ht.transaction_hash, ht.tx_result, COALESCE(ht.successful, true) as transaction_successful FROM history_operations hop LEFT JOIN history_transactions ht ON ht.id = hop.transaction_id JOIN history_operation_participants hopp ON hopp.history_operation_id = hop.id WHERE hop.type IN (?,?,?,?,?) AND hopp.history_account_id = ?", sql)
}

func TestExtraChecksOperationsTransactionSuccessfulTrueResultFalse(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	// successful `true` but tx result `false`
	_, err := tt.AuroraDB.Exec(
		`UPDATE history_transactions SET successful = true WHERE transaction_hash = 'aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf'`,
	)
	tt.Require.NoError(err)

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		ForAccount(tt.Ctx, "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2").
		IncludeFailed()

	_, _, err = query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.Contains(err.Error(), "Corrupted data! `successful=true` but returned transaction is not success")
}

func TestExtraChecksOperationsTransactionSuccessfulFalseResultTrue(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	// successful `false` but tx result `true`
	_, err := tt.AuroraDB.Exec(
		`UPDATE history_transactions SET successful = false WHERE transaction_hash = 'a2dabf4e9d1642722602272e178a37c973c9177b957da86192a99b3e9f3a9aa4'`,
	)
	tt.Require.NoError(err)

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		ForAccount(tt.Ctx, "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON").
		IncludeFailed()

	_, _, err = query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.Contains(err.Error(), "Corrupted data! `successful=false` but returned transaction is success")
}

func assertOperationMatchesTransaction(tt *test.T, operation Operation, transaction Transaction) {
	tt.Assert.Equal(operation.TransactionID, transaction.ID)
	tt.Assert.Equal(operation.TransactionHash, transaction.TransactionHash)
	tt.Assert.Equal(operation.TxResult, transaction.TxResult)
	tt.Assert.Equal(operation.TransactionSuccessful, transaction.Successful)
}

// TestOperationIncludeTransactions tests that transactions are included when fetching records from the db.
func TestOperationIncludeTransactions(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	defer tt.Finish()

	accountID := "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		IncludeTransactions().
		ForAccount(tt.Ctx, accountID)

	operations, transactions, err := query.Fetch(tt.Ctx)
	tt.Assert.NoError(err)
	tt.Assert.Len(transactions, 3)
	tt.Assert.Len(transactions, len(operations))

	for i := range operations {
		operation := operations[i]
		transaction := transactions[i]
		assertOperationMatchesTransaction(tt, operation, transaction)
	}

	withoutTransactionsQuery := (&Q{tt.AuroraSession()}).Operations().
		ForAccount(tt.Ctx, accountID)

	var expectedTransactions []Transaction
	err = (&Q{tt.AuroraSession()}).Transactions().ForAccount(tt.Ctx, accountID).Select(tt.Ctx, &expectedTransactions)
	tt.Assert.NoError(err)

	expectedOperations, _, err := withoutTransactionsQuery.Fetch(tt.Ctx)
	tt.Assert.NoError(err)

	tt.Assert.Equal(operations, expectedOperations)
	tt.Assert.Equal(transactions, expectedTransactions)

	op, transaction, err := q.OperationByID(tt.Ctx, true, expectedOperations[0].ID)
	tt.Assert.NoError(err)
	tt.Assert.Equal(op, expectedOperations[0])
	tt.Assert.Equal(*transaction, expectedTransactions[0])
	assertOperationMatchesTransaction(tt, op, *transaction)

	_, err = q.Exec(tt.Ctx, sq.Delete("history_transactions"))
	tt.Assert.NoError(err)
	_, _, err = q.OperationByID(tt.Ctx, true, 17179877377)
	tt.Assert.Error(err)
}

func TestValidateTransactionForOperation(t *testing.T) {
	tt := test.Start(t)
	tt.Scenario("failed_transactions")
	selectTransactionCopy := selectTransaction
	defer func() {
		selectTransaction = selectTransactionCopy
		tt.Finish()
	}()

	selectTransaction = sq.Select(
		"ht.transaction_hash, " +
			"ht.tx_result, " +
			"COALESCE(ht.successful, true) as successful").
		From("history_transactions ht")

	accountID := "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"

	q := &Q{tt.AuroraSession()}
	query := q.Operations().
		IncludeTransactions().
		ForAccount(tt.Ctx, accountID)

	_, _, err := query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction with id 17179877376 could not be found")

	_, _, err = q.OperationByID(tt.Ctx, true, 17179877377)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction id 0 does not match transaction id in operation 17179877376")

	selectTransaction = sq.Select(
		"ht.id, " +
			"ht.transaction_hash, " +
			"COALESCE(ht.successful, true) as successful").
		From("history_transactions ht")
	query = q.Operations().
		IncludeTransactions().
		ForAccount(tt.Ctx, accountID)

	_, _, err = query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction result  does not match transaction result in operation AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=")

	_, _, err = q.OperationByID(tt.Ctx, true, 17179877377)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction result  does not match transaction result in operation AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=")

	selectTransaction = sq.Select(
		"ht.id, " +
			"ht.tx_result, " +
			"COALESCE(ht.successful, true) as successful").
		From("history_transactions ht")
	query = q.Operations().
		IncludeTransactions().
		ForAccount(tt.Ctx, accountID)

	_, _, err = query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction hash  does not match transaction hash in operation 1c454630267aa8767ec8c8e30450cea6ba660145e9c924abb75d7a6669b6c28a")

	_, _, err = q.OperationByID(tt.Ctx, true, 17179877377)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction hash  does not match transaction hash in operation 1c454630267aa8767ec8c8e30450cea6ba660145e9c924abb75d7a6669b6c28a")

	selectTransaction = sq.Select(
		"ht.id, " +
			"ht.tx_result, " +
			"ht.transaction_hash").
		From("history_transactions ht")
	query = q.Operations().
		IncludeTransactions().
		ForAccount(tt.Ctx, accountID)

	_, _, err = query.Fetch(tt.Ctx)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction successful flag false does not match transaction successful flag in operation true")

	_, _, err = q.OperationByID(tt.Ctx, true, 17179877377)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "transaction successful flag false does not match transaction successful flag in operation true")
}
