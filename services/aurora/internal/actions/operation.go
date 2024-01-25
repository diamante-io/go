package actions

import (
	"context"
	"fmt"
	"net/http"

	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ledger"
	"go/services/aurora/internal/render/problem"
	"go/services/aurora/internal/resourceadapter"
	"go/services/aurora/internal/toid"
	"go/support/errors"
	"go/support/render/hal"
	supportProblem "go/support/render/problem"
)

// Joinable query struct for join query parameter
type Joinable struct {
	Join string `schema:"join" valid:"in(transactions)~Accepted values: transactions,optional"`
}

// IncludeTransactions returns extra fields to include in the response
func (qp Joinable) IncludeTransactions() bool {
	return qp.Join == "transactions"
}

// OperationsQuery query struct for operations end-points
type OperationsQuery struct {
	Joinable                  `valid:"optional"`
	AccountID                 string `schema:"account_id" valid:"accountID,optional"`
	ClaimableBalanceID        string `schema:"claimable_balance_id" valid:"claimableBalanceID,optional"`
	LiquidityPoolID           string `schema:"liquidity_pool_id" valid:"sha256,optional"`
	TransactionHash           string `schema:"tx_id" valid:"transactionHash,optional"`
	IncludeFailedTransactions bool   `schema:"include_failed" valid:"-"`
	LedgerID                  uint32 `schema:"ledger_id" valid:"-"`
}

// Validate runs extra validations on query parameters
func (qp OperationsQuery) Validate() error {
	filters, err := countNonEmpty(
		qp.AccountID,
		qp.ClaimableBalanceID,
		qp.LiquidityPoolID,
		qp.LedgerID,
		qp.TransactionHash,
	)

	if err != nil {
		return supportProblem.BadRequest
	}

	if filters > 1 {
		return supportProblem.MakeInvalidFieldProblem(
			"filters",
			errors.New("Use a single filter for operations, you can only use one of tx_id, account_id or ledger_id"),
		)
	}

	return nil
}

// GetOperationsHandler is the action handler for all end-points returning a list of operations.
type GetOperationsHandler struct {
	LedgerState  *ledger.State
	OnlyPayments bool
}

// GetResourcePage returns a page of operations.
func (handler GetOperationsHandler) GetResourcePage(w HeaderWriter, r *http.Request) ([]hal.Pageable, error) {
	ctx := r.Context()

	pq, err := GetPageQuery(handler.LedgerState, r)
	if err != nil {
		return nil, err
	}

	err = validateCursorWithinHistory(handler.LedgerState, pq)
	if err != nil {
		return nil, err
	}

	qp := OperationsQuery{}
	err = getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	query := historyQ.Operations()

	switch {
	case qp.AccountID != "":
		query.ForAccount(ctx, qp.AccountID)
	case qp.ClaimableBalanceID != "":
		query.ForClaimableBalance(ctx, qp.ClaimableBalanceID)
	case qp.LiquidityPoolID != "":
		query.ForLiquidityPool(ctx, qp.LiquidityPoolID)
	case qp.LedgerID > 0:
		query.ForLedger(ctx, int32(qp.LedgerID))
	case qp.TransactionHash != "":
		query.ForTransaction(ctx, qp.TransactionHash)
	}
	// When querying operations for transaction return both successful
	// and failed operations. We assume that because the user is querying
	// this specific transactions, they knows its status.
	if qp.TransactionHash != "" || qp.IncludeFailedTransactions {
		query.IncludeFailed()
	}

	if qp.IncludeTransactions() {
		query.IncludeTransactions()
	}

	if handler.OnlyPayments {
		query.OnlyPayments()
	}

	ops, txs, err := query.Page(pq).Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return buildOperationsPage(ctx, historyQ, ops, txs, qp.IncludeTransactions())
}

// GetOperationByIDHandler is the action handler for all end-points returning a list of operations.
type GetOperationByIDHandler struct {
	LedgerState *ledger.State
}

// OperationQuery query struct for operation/id end-point
type OperationQuery struct {
	LedgerState *ledger.State `valid:"-"`
	Joinable    `valid:"optional"`
	ID          uint64 `schema:"id" valid:"-"`
}

// Validate runs extra validations on query parameters
func (qp OperationQuery) Validate() error {
	parsed := toid.Parse(int64(qp.ID))
	if parsed.LedgerSequence < qp.LedgerState.CurrentStatus().HistoryElder {
		return problem.BeforeHistory
	}
	return nil
}

// GetResource returns an operation page.
func (handler GetOperationByIDHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	qp := OperationQuery{
		LedgerState: handler.LedgerState,
	}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}
	op, tx, err := historyQ.OperationByID(ctx, qp.IncludeTransactions(), int64(qp.ID))
	if err != nil {
		return nil, err
	}

	var ledger history.Ledger
	err = historyQ.LedgerBySequence(ctx, &ledger, op.LedgerSequence())
	if err != nil {
		return nil, err
	}

	return resourceadapter.NewOperation(
		ctx,
		op,
		op.TransactionHash,
		tx,
		ledger,
	)
}

func buildOperationsPage(ctx context.Context, historyQ *history.Q, operations []history.Operation, transactions []history.Transaction, includeTransactions bool) ([]hal.Pageable, error) {
	ledgerCache := history.LedgerCache{}
	for _, record := range operations {
		ledgerCache.Queue(record.LedgerSequence())
	}

	if err := ledgerCache.Load(ctx, historyQ); err != nil {
		return nil, errors.Wrap(err, "failed to load ledger batch")
	}

	var response []hal.Pageable
	for i, operationRecord := range operations {
		ledger, found := ledgerCache.Records[operationRecord.LedgerSequence()]
		if !found {
			msg := fmt.Sprintf("could not find ledger data for sequence %d", operationRecord.LedgerSequence())
			return nil, errors.New(msg)
		}

		var transactionRecord *history.Transaction

		if includeTransactions {
			transactionRecord = &transactions[i]
		}

		var res hal.Pageable
		res, err := resourceadapter.NewOperation(
			ctx,
			operationRecord,
			operationRecord.TransactionHash,
			transactionRecord,
			ledger,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, res)
	}

	return response, nil
}
