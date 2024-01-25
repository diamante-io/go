package actions

import (
	"context"
	"net/http"
	"strings"

	protocol "github.com/diamcircle/go/protocols/aurora"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/resourceadapter"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/diamcircle/go/support/render/problem"
	"github.com/diamcircle/go/xdr"
)

// GetLiquidityPoolByIDHandler is the action handler for all end-points returning a liquidity pool.
type GetLiquidityPoolByIDHandler struct{}

// LiquidityPoolQuery query struct for liquidity_pools/id endpoint
type LiquidityPoolQuery struct {
	ID string `schema:"liquidity_pool_id" valid:"sha256"`
}

// GetResource returns an claimable balance page.
func (handler GetLiquidityPoolByIDHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	ctx := r.Context()
	qp := LiquidityPoolQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}
	cb, err := historyQ.FindLiquidityPoolByID(ctx, qp.ID)
	if err != nil {
		return nil, err
	}
	ledger := &history.Ledger{}
	err = historyQ.LedgerBySequence(ctx, ledger, int32(cb.LastModifiedLedger))
	if historyQ.NoRows(err) {
		ledger = nil
	} else if err != nil {
		return nil, errors.Wrap(err, "LedgerBySequence error")
	}

	var resource protocol.LiquidityPool
	err = resourceadapter.PopulateLiquidityPool(ctx, &resource, cb, ledger)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// LiquidityPoolsQuery query struct for liquidity_pools end-point
type LiquidityPoolsQuery struct {
	Reserves string `schema:"reserves" valid:"optional"`
	Account  string `schema:"account" valid:"optional"`

	reserves []xdr.Asset
}

// URITemplate returns a rfc6570 URI template the query struct
func (q LiquidityPoolsQuery) URITemplate() string {
	return getURITemplate(&q, "liquidity_pools", true)
}

// Validate validates and parses the query
func (q *LiquidityPoolsQuery) Validate() error {
	assets := []xdr.Asset{}
	reserves := strings.Split(q.Reserves, ",")
	reservesErr := problem.MakeInvalidFieldProblem(
		"reserves",
		errors.New("Invalid reserves, should be comma-separated list of assets in canonical form"),
	)
	for _, reserve := range reserves {
		if reserve == "" {
			continue
		}
		switch reserve {
		case "native":
			assets = append(assets, xdr.MustNewNativeAsset())
		default:
			parts := strings.Split(reserve, ":")
			if len(parts) != 2 {
				return reservesErr
			}
			asset, err := xdr.NewCreditAsset(parts[0], parts[1])
			if err != nil {
				return reservesErr
			}
			assets = append(assets, asset)
		}
	}
	q.reserves = assets
	return nil
}

type GetLiquidityPoolsHandler struct {
	LedgerState *ledger.State
}

// GetResourcePage returns a page of liquidity pools.
func (handler GetLiquidityPoolsHandler) GetResourcePage(w HeaderWriter, r *http.Request) ([]hal.Pageable, error) {
	ctx := r.Context()
	qp := LiquidityPoolsQuery{}
	err := getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	pq, err := GetPageQuery(handler.LedgerState, r, DisableCursorValidation)
	if err != nil {
		return nil, err
	}

	query := history.LiquidityPoolsQuery{
		PageQuery: pq,
		Account:   qp.Account,
		Assets:    qp.reserves,
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	liquidityPools, err := handler.getLiquidityPoolsPage(ctx, historyQ, query)
	if err != nil {
		return nil, err
	}

	return liquidityPools, nil
}

func (handler GetLiquidityPoolsHandler) getLiquidityPoolsPage(ctx context.Context, historyQ *history.Q, query history.LiquidityPoolsQuery) ([]hal.Pageable, error) {
	records, err := historyQ.GetLiquidityPools(ctx, query)
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

	var liquidityPools []hal.Pageable
	for _, record := range records {
		var response protocol.LiquidityPool

		var ledger *history.Ledger
		if l, ok := ledgerCache.Records[int32(record.LastModifiedLedger)]; ok {
			ledger = &l
		}

		resourceadapter.PopulateLiquidityPool(ctx, &response, record, ledger)
		liquidityPools = append(liquidityPools, response)
	}

	return liquidityPools, nil
}
