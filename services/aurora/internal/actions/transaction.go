package actions

import (
	"context"
	"net/http"

	"github.com/diamcircle/go/protocols/aurora"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/resourceadapter"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/hal"
	supportProblem "github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/xdr"
)

// TransactionQuery query struct for transactions/id end-point
type TransactionQuery struct {
	TransactionHash string `schema:"tx_id" valid:"transactionHash,optional"`
}

// GetTransactionByHashHandler is the action handler for the end-point returning a transaction.
type GetTransactionByHashHandler struct {
}

// GetResource returns a transaction page.
func (handler GetTransactionByHashHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	qp := TransactionQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	var (
		record   history.Transaction
		resource aurora.Transaction
	)

	err = historyQ.TransactionByHash(ctx, &record, qp.TransactionHash)
	if err != nil {
		return resource, errors.Wrap(err, "loading transaction record")
	}

	if err = resourceadapter.PopulateTransaction(ctx, qp.TransactionHash, &resource, record); err != nil {
		return resource, errors.Wrap(err, "could not populate transaction")
	}
	return resource, nil
}

// TransactionsQuery query struct for transactions end-points
type TransactionsQuery struct {
	AccountID                 string `schema:"account_id" valid:"accountID,optional"`
	ClaimableBalanceID        string `schema:"claimable_balance_id" valid:"claimableBalanceID,optional"`
	LiquidityPoolID           string `schema:"liquidity_pool_id" valid:"sha256,optional"`
	IncludeFailedTransactions bool   `schema:"include_failed" valid:"-"`
	LedgerID                  uint32 `schema:"ledger_id" valid:"-"`
}

// Validate runs extra validations on query parameters
func (qp TransactionsQuery) Validate() error {
	filters, err := countNonEmpty(
		qp.AccountID,
		qp.ClaimableBalanceID,
		qp.LiquidityPoolID,
		qp.LedgerID,
	)

	if err != nil {
		return supportProblem.BadRequest
	}

	if filters > 1 {
		return supportProblem.MakeInvalidFieldProblem(
			"filters",
			errors.New("Use a single filter for transaction, you can only use one of account_id, claimable_balance_id or ledger_id"),
		)
	}

	return nil
}

// GetTransactionsHandler is the action handler for all end-points returning a list of transactions.
type GetTransactionsHandler struct {
	LedgerState *ledger.State
}

// GetResourcePage returns a page of transactions.
func (handler GetTransactionsHandler) GetResourcePage(w HeaderWriter, r *http.Request) ([]hal.Pageable, error) {
	ctx := r.Context()

	pq, err := GetPageQuery(handler.LedgerState, r)
	if err != nil {
		return nil, err
	}

	err = validateCursorWithinHistory(handler.LedgerState, pq)
	if err != nil {
		return nil, err
	}

	qp := TransactionsQuery{}
	err = getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	records, err := loadTransactionRecords(ctx, historyQ, qp, pq)
	if err != nil {
		return nil, errors.Wrap(err, "loading transaction records")
	}

	var response []hal.Pageable

	for _, record := range records {
		var res aurora.Transaction
		err = resourceadapter.PopulateTransaction(ctx, record.TransactionHash, &res, record)
		if err != nil {
			return nil, errors.Wrap(err, "could not populate transaction")
		}
		response = append(response, res)
	}

	return response, nil
}

// loadTransactionRecords returns a slice of transaction records of an
// account/ledger identified by accountID/ledgerID based on pq and
// includeFailedTx.
func loadTransactionRecords(ctx context.Context, hq *history.Q, qp TransactionsQuery, pq db2.PageQuery) ([]history.Transaction, error) {
	var records []history.Transaction

	txs := hq.Transactions()
	switch {
	case qp.AccountID != "":
		txs.ForAccount(ctx, qp.AccountID)
	case qp.ClaimableBalanceID != "":
		txs.ForClaimableBalance(ctx, qp.ClaimableBalanceID)
	case qp.LiquidityPoolID != "":
		txs.ForLiquidityPool(ctx, qp.LiquidityPoolID)
	case qp.LedgerID > 0:
		txs.ForLedger(ctx, int32(qp.LedgerID))
	}

	if qp.IncludeFailedTransactions {
		txs.IncludeFailed()
	}

	err := txs.Page(pq).Select(ctx, &records)
	if err != nil {
		return nil, errors.Wrap(err, "executing transaction records query")
	}

	for _, t := range records {
		if !qp.IncludeFailedTransactions {
			if !t.Successful {
				return nil, errors.Errorf("Corrupted data! `include_failed=false` but returned transaction is failed: %s", t.TransactionHash)
			}

			var resultXDR xdr.TransactionResult
			err = xdr.SafeUnmarshalBase64(t.TxResult, &resultXDR)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshalling tx result")
			}

			if !resultXDR.Successful() {
				return nil, errors.Errorf("Corrupted data! `include_failed=false` but returned transaction is failed: %s %s", t.TransactionHash, t.TxResult)
			}
		}
	}

	return records, nil
}
