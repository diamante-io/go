package actions

import (
	"context"
	"net/http"
	"strings"

	"go/protocols/aurora"
	protocol "go/protocols/aurora"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ledger"
	"go/services/aurora/internal/resourceadapter"
	"go/support/errors"
	"go/support/render/hal"
	"go/support/render/problem"
	"go/xdr"
)

// GetClaimableBalanceByIDHandler is the action handler for all end-points returning a claimable balance.
type GetClaimableBalanceByIDHandler struct{}

// ClaimableBalanceQuery query struct for claimables_balances/id end-point
type ClaimableBalanceQuery struct {
	ID string `schema:"id" valid:"claimableBalanceID,required"`
}

// GetResource returns an claimable balance page.
func (handler GetClaimableBalanceByIDHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	qp := ClaimableBalanceQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}
	cb, err := historyQ.FindClaimableBalanceByID(ctx, qp.ID)
	if err != nil {
		return nil, err
	}
	ledger := &history.Ledger{}
	err = historyQ.LedgerBySequence(ctx,
		ledger,
		int32(cb.LastModifiedLedger),
	)
	if historyQ.NoRows(err) {
		ledger = nil
	} else if err != nil {
		return nil, errors.Wrap(err, "LedgerBySequence error")
	}

	var resource protocol.ClaimableBalance
	err = resourceadapter.PopulateClaimableBalance(ctx, &resource, cb, ledger)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// ClaimableBalancesQuery query struct for claimable_balances end-point
type ClaimableBalancesQuery struct {
	AssetFilter    string `schema:"asset" valid:"asset,optional"`
	SponsorFilter  string `schema:"sponsor" valid:"accountID,optional"`
	ClaimantFilter string `schema:"claimant" valid:"accountID,optional"`
}

func (q ClaimableBalancesQuery) asset() *xdr.Asset {
	if len(q.AssetFilter) > 0 {
		switch q.AssetFilter {
		case "native":
			asset := xdr.MustNewNativeAsset()
			return &asset
		default:
			parts := strings.Split(q.AssetFilter, ":")
			asset := xdr.MustNewCreditAsset(parts[0], parts[1])
			return &asset
		}
	}
	return nil
}

func (q ClaimableBalancesQuery) sponsor() *xdr.AccountId {
	if q.SponsorFilter != "" {
		return xdr.MustAddressPtr(q.SponsorFilter)
	}
	return nil
}

func (q ClaimableBalancesQuery) claimant() *xdr.AccountId {
	if q.ClaimantFilter != "" {
		return xdr.MustAddressPtr(q.ClaimantFilter)
	}
	return nil
}

// URITemplate returns a rfc6570 URI template the query struct
func (q ClaimableBalancesQuery) URITemplate() string {
	return getURITemplate(&q, "claimable_balances", true)
}

type GetClaimableBalancesHandler struct {
	LedgerState *ledger.State
}

// GetResourcePage returns a page of claimable balances.
func (handler GetClaimableBalancesHandler) GetResourcePage(
	w HeaderWriter,
	r *http.Request,
) ([]hal.Pageable, error) {
	ctx := r.Context()
	qp := ClaimableBalancesQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	pq, err := GetPageQuery(handler.LedgerState, r, DisableCursorValidation)
	if err != nil {
		return nil, err
	}

	query := history.ClaimableBalancesQuery{
		PageQuery: pq,
		Asset:     qp.asset(),
		Sponsor:   qp.sponsor(),
		Claimant:  qp.claimant(),
	}

	_, _, err = query.Cursor()
	if err != nil {
		return nil, problem.MakeInvalidFieldProblem(
			"cursor",
			errors.New("The first part should be a number higher than 0 and the second part should be a valid claimable balance ID"),
		)
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	claimableBalances, err := getClaimableBalancesPage(ctx, historyQ, query)
	if err != nil {
		return nil, err
	}

	return claimableBalances, nil
}

func getClaimableBalancesPage(ctx context.Context, historyQ *history.Q, query history.ClaimableBalancesQuery) ([]hal.Pageable, error) {
	records, err := historyQ.GetClaimableBalances(ctx, query)
	if err != nil {
		return nil, err
	}

	ledgerCache := history.LedgerCache{}
	for _, record := range records {
		ledgerCache.Queue(int32(record.LastModifiedLedger))
	}
	if err := ledgerCache.Load(ctx, historyQ); err != nil {
		return nil, errors.Wrap(err, "failed to load ledger batch")
	}

	var claimableBalances []hal.Pageable
	for _, record := range records {
		var response aurora.ClaimableBalance

		var ledger *history.Ledger
		if l, ok := ledgerCache.Records[int32(record.LastModifiedLedger)]; ok {
			ledger = &l
		}

		resourceadapter.PopulateClaimableBalance(ctx, &response, record, ledger)
		claimableBalances = append(claimableBalances, response)
	}

	return claimableBalances, nil
}
