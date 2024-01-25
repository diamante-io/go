package actions

import (
	"net/http"

	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/render/problem"
	"github.com/diamcircle/go/services/aurora/internal/resourceadapter"
	"github.com/diamcircle/go/support/render/hal"
)

type GetLedgersHandler struct {
	LedgerState *ledger.State
}

func (handler GetLedgersHandler) GetResourcePage(w HeaderWriter, r *http.Request) ([]hal.Pageable, error) {
	pq, err := GetPageQuery(handler.LedgerState, r)
	if err != nil {
		return nil, err
	}

	err = validateCursorWithinHistory(handler.LedgerState, pq)
	if err != nil {
		return nil, err
	}

	historyQ, err := context.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	var records []history.Ledger
	if err = historyQ.Ledgers().Page(pq).Select(r.Context(), &records); err != nil {
		return nil, err
	}

	var result []hal.Pageable
	for _, record := range records {
		var ledger aurora.Ledger
		resourceadapter.PopulateLedger(r.Context(), &ledger, record)
		if err != nil {
			return nil, err
		}
		result = append(result, ledger)
	}

	return result, nil
}

// LedgerByIDQuery query struct for the ledger/{id} endpoint
type LedgerByIDQuery struct {
	LedgerID uint32 `schema:"ledger_id" valid:"-"`
}

type GetLedgerByIDHandler struct {
	LedgerState *ledger.State
}

func (handler GetLedgerByIDHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	qp := LedgerByIDQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}
	if int32(qp.LedgerID) < handler.LedgerState.CurrentStatus().HistoryElder {
		return nil, problem.BeforeHistory
	}
	historyQ, err := context.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}
	var ledger history.Ledger
	err = historyQ.LedgerBySequence(r.Context(), &ledger, int32(qp.LedgerID))
	if err != nil {
		return nil, err
	}
	var result aurora.Ledger
	resourceadapter.PopulateLedger(r.Context(), &result, ledger)
	return result, nil
}
