package history

import (
	"context"
	"fmt"
	"math"

	sq "github.com/Masterminds/squirrel"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/xdr"
)

// PagingToken returns a cursor for this trade
func (r *Trade) PagingToken() string {
	return fmt.Sprintf("%d-%d", r.HistoryOperationID, r.Order)
}

// HasPrice returns true if the trade has non-null price data
func (r *Trade) HasPrice() bool {
	return r.PriceN.Valid && r.PriceD.Valid
}

const (
	AllTrades           = "all"
	OrderbookTrades     = "orderbook"
	LiquidityPoolTrades = "liquidity_pool"
)

type tradesQuery struct {
	baseAsset     *xdr.Asset
	counterAsset  *xdr.Asset
	tradeType     string
	account       string
	liquidityPool string
	offer         int64
}

func (q *Q) GetTrades(
	ctx context.Context, page db2.PageQuery, account string, tradeType string,
) ([]Trade, error) {
	return q.getTrades(ctx, page, tradesQuery{
		account:   account,
		tradeType: tradeType,
	})
}

func (q *Q) GetTradesForOffer(
	ctx context.Context, page db2.PageQuery, offerID int64,
) ([]Trade, error) {
	return q.getTrades(ctx, page, tradesQuery{
		offer:     offerID,
		tradeType: AllTrades,
	})
}

func (q *Q) GetTradesForLiquidityPool(
	ctx context.Context, page db2.PageQuery, poolID string,
) ([]Trade, error) {
	return q.getTrades(ctx, page, tradesQuery{
		liquidityPool: poolID,
		tradeType:     AllTrades,
	})
}

func (q *Q) GetTradesForAssets(
	ctx context.Context, page db2.PageQuery, account, tradeType string, baseAsset, counterAsset xdr.Asset,
) ([]Trade, error) {
	return q.getTrades(ctx, page, tradesQuery{
		account:      account,
		baseAsset:    &baseAsset,
		counterAsset: &counterAsset,
		tradeType:    tradeType,
	})
}

type historyTradesQuery struct {
	baseAssetID    int64
	counterAssetID int64
	accountID      int64
	offerID        int64
	poolID         int64
	orderPreserved bool
	tradeType      string
}

func (q *Q) getTrades(ctx context.Context, page db2.PageQuery, query tradesQuery) ([]Trade, error) {
	// Add explicit query type for prometheus metrics, since we use raw sql.
	ctx = context.WithValue(ctx, &db.QueryTypeContextKey, db.SelectQueryType)

	internalTradesQuery, err := q.transformTradesQuery(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "invalid trade query")
	}
	rawSQL, args, err := createTradesSQL(page, internalTradesQuery)
	if err != nil {
		return nil, errors.Wrap(err, "could not create trades sql query")
	}

	var dest []Trade
	if err = q.SelectRaw(ctx, &dest, rawSQL, args...); err != nil {
		return nil, errors.Wrap(err, "could not select trades")
	}

	return dest, nil
}

func (q *Q) transformTradesQuery(ctx context.Context, query tradesQuery) (historyTradesQuery, error) {
	internalQuery := historyTradesQuery{
		orderPreserved: true,
		tradeType:      query.tradeType,
		offerID:        query.offer,
	}

	if query.account != "" {
		var account Account
		if err := q.AccountByAddress(ctx, &account, query.account); err != nil {
			return internalQuery, errors.Wrap(err, "could not get account by address")
		}
		internalQuery.accountID = account.ID
	}

	if query.baseAsset != nil {
		var err error
		internalQuery.baseAssetID, err = q.GetAssetID(ctx, *query.baseAsset)
		if err != nil {
			return internalQuery, errors.Wrap(err, "could not get base asset id")
		}

		internalQuery.counterAssetID, err = q.GetAssetID(ctx, *query.counterAsset)
		if err != nil {
			return internalQuery, errors.Wrap(err, "could not get counter asset id")
		}
		internalQuery.orderPreserved, internalQuery.baseAssetID, internalQuery.counterAssetID = getCanonicalAssetOrder(
			internalQuery.baseAssetID, internalQuery.counterAssetID,
		)
	}

	if query.liquidityPool != "" {
		historyPool, err := q.LiquidityPoolByID(ctx, query.liquidityPool)
		if err != nil {
			return internalQuery, errors.Wrap(err, "could not get pool id")
		}
		internalQuery.poolID = historyPool.InternalID
	}

	return internalQuery, nil
}

func createTradesSQL(page db2.PageQuery, query historyTradesQuery) (string, []interface{}, error) {
	base := selectTradeFields
	if !query.orderPreserved {
		base = selectReverseTradeFields
	}
	sql := joinTradeAssets(
		joinTradeLiquidityPools(
			joinTradeAccounts(
				base.From("history_trades htrd"),
				"history_accounts",
			),
			"history_liquidity_pools",
		),
		"history_assets",
	)

	if query.baseAssetID != 0 {
		sql = sql.Where(sq.Eq{"base_asset_id": query.baseAssetID, "counter_asset_id": query.counterAssetID})
	}

	switch query.tradeType {
	case OrderbookTrades:
		sql = sql.Where(sq.Eq{"htrd.trade_type": OrderbookTradeType})
	case LiquidityPoolTrades:
		sql = sql.Where(sq.Eq{"htrd.trade_type": LiquidityPoolTradeType})
	case AllTrades:
	default:
		return "", nil, errors.Errorf("Invalid trade type: %v", query.tradeType)
	}

	op, idx, err := page.CursorInt64Pair(db2.DefaultPairSep)
	if err != nil {
		return "", nil, errors.Wrap(err, "could not parse cursor")
	}

	// constrain the second portion of the cursor pair to 32-bits
	if idx > math.MaxInt32 {
		idx = math.MaxInt32
	}

	if query.accountID != 0 || query.offerID != 0 || query.poolID != 0 {
		// Construct UNION query
		var firstSelect, secondSelect sq.SelectBuilder
		switch {
		case query.accountID != 0:
			firstSelect = sql.Where("htrd.base_account_id = ?", query.accountID)
			secondSelect = sql.Where("htrd.counter_account_id = ?", query.accountID)
		case query.offerID != 0:
			firstSelect = sql.Where("htrd.base_offer_id = ?", query.offerID)
			secondSelect = sql.Where("htrd.counter_offer_id = ?", query.offerID)
		case query.poolID != 0:
			firstSelect = sql.Where("htrd.base_liquidity_pool_id = ?", query.poolID)
			secondSelect = sql.Where("htrd.counter_liquidity_pool_id = ?", query.poolID)
		}

		firstSelect = appendOrdering(firstSelect, op, idx, page.Order)
		secondSelect = appendOrdering(secondSelect, op, idx, page.Order)
		firstSQL, firstArgs, err := firstSelect.ToSql()
		if err != nil {
			return "", nil, errors.Wrap(err, "error building a firstSelect query")
		}
		secondSQL, secondArgs, err := secondSelect.ToSql()
		if err != nil {
			return "", nil, errors.Wrap(err, "error building a secondSelect query")
		}

		rawSQL := fmt.Sprintf("(%s) UNION (%s) ", firstSQL, secondSQL)
		args := append(firstArgs, secondArgs...)
		// Order the final UNION:
		switch page.Order {
		case "asc":
			rawSQL = rawSQL + `ORDER BY history_operation_id asc, "order" asc `
		case "desc":
			rawSQL = rawSQL + `ORDER BY history_operation_id desc, "order" desc `
		default:
			panic("Invalid order")
		}
		rawSQL = rawSQL + fmt.Sprintf("LIMIT %d", page.Limit)
		return rawSQL, args, nil
	} else {
		sql = appendOrdering(sql, op, idx, page.Order)
		sql = sql.Limit(page.Limit)
		rawSQL, args, err := sql.ToSql()
		if err != nil {
			return "", nil, errors.Wrap(err, "error building sql query")
		}
		return rawSQL, args, nil
	}
}

func appendOrdering(sel sq.SelectBuilder, op, idx int64, order string) sq.SelectBuilder {
	// NOTE: Remember to test the queries below with EXPLAIN / EXPLAIN ANALYZE
	// before changing them.
	// This condition is using multicolumn index and it's easy to write it in a way that
	// DB will perform a full table scan.
	switch order {
	case "asc":
		return sel.
			Where(`(
				htrd.history_operation_id >= ?
			AND (
				htrd.history_operation_id > ? OR
				(htrd.history_operation_id = ? AND htrd.order > ?)
			))`, op, op, op, idx).
			OrderBy("htrd.history_operation_id asc, htrd.order asc")
	case "desc":
		return sel.
			Where(`(
				htrd.history_operation_id <= ?
			AND (
				htrd.history_operation_id < ? OR
				(htrd.history_operation_id = ? AND htrd.order < ?)
			))`, op, op, op, idx).
			OrderBy("htrd.history_operation_id desc, htrd.order desc")
	default:
		panic("Invalid order")
	}
}

func joinTradeAccounts(selectBuilder sq.SelectBuilder, historyAccountsTable string) sq.SelectBuilder {
	return selectBuilder.
		LeftJoin(historyAccountsTable + " base_accounts ON base_account_id = base_accounts.id").
		LeftJoin(historyAccountsTable + " counter_accounts ON counter_account_id = counter_accounts.id")
}

func joinTradeAssets(selectBuilder sq.SelectBuilder, historyAssetsTable string) sq.SelectBuilder {
	return selectBuilder.
		Join(historyAssetsTable + " base_assets ON base_asset_id = base_assets.id").
		Join(historyAssetsTable + " counter_assets ON counter_asset_id = counter_assets.id")
}

func joinTradeLiquidityPools(selectBuilder sq.SelectBuilder, historyLiquidityPoolsTable string) sq.SelectBuilder {
	return selectBuilder.
		LeftJoin(historyLiquidityPoolsTable + " blp ON base_liquidity_pool_id = blp.id").
		LeftJoin(historyLiquidityPoolsTable + " clp ON counter_liquidity_pool_id = clp.id")
}

var selectTradeFields = sq.Select(
	"history_operation_id",
	"htrd.\"order\"",
	"htrd.ledger_closed_at",
	"htrd.base_offer_id",
	"base_accounts.address as base_account",
	"base_assets.asset_type as base_asset_type",
	"base_assets.asset_code as base_asset_code",
	"base_assets.asset_issuer as base_asset_issuer",
	"blp.liquidity_pool_id as base_liquidity_pool_id",
	"htrd.base_amount",
	"htrd.counter_offer_id",
	"counter_accounts.address as counter_account",
	"counter_assets.asset_type as counter_asset_type",
	"counter_assets.asset_code as counter_asset_code",
	"counter_assets.asset_issuer as counter_asset_issuer",
	"clp.liquidity_pool_id as counter_liquidity_pool_id",
	"htrd.counter_amount",
	"liquidity_pool_fee",
	"htrd.base_is_seller",
	"htrd.price_n",
	"htrd.price_d",
	"htrd.trade_type",
)

var selectReverseTradeFields = sq.Select(
	"history_operation_id",
	"htrd.\"order\"",
	"htrd.ledger_closed_at",
	"htrd.counter_offer_id as base_offer_id",
	"counter_accounts.address as base_account",
	"counter_assets.asset_type as base_asset_type",
	"counter_assets.asset_code as base_asset_code",
	"counter_assets.asset_issuer as base_asset_issuer",
	"clp.liquidity_pool_id as base_liquidity_pool_id",
	"htrd.counter_amount as base_amount",
	"htrd.base_offer_id as counter_offer_id",
	"base_accounts.address as counter_account",
	"base_assets.asset_type as counter_asset_type",
	"base_assets.asset_code as counter_asset_code",
	"base_assets.asset_issuer as counter_asset_issuer",
	"blp.liquidity_pool_id as counter_liquidity_pool_id",
	"htrd.base_amount as counter_amount",
	"liquidity_pool_fee",
	"NOT(htrd.base_is_seller) as base_is_seller",
	"htrd.price_d as price_n",
	"htrd.price_n as price_d",
	"htrd.trade_type",
)

func getCanonicalAssetOrder(
	assetId1 int64, assetId2 int64,
) (orderPreserved bool, baseAssetId int64, counterAssetId int64) {
	if assetId1 < assetId2 {
		return true, assetId1, assetId2
	} else {
		return false, assetId2, assetId1
	}
}

type QTrades interface {
	QCreateAccountsHistory
	NewTradeBatchInsertBuilder(maxBatchSize int) TradeBatchInsertBuilder
	RebuildTradeAggregationBuckets(ctx context.Context, fromledger, toLedger uint32) error
	CreateAssets(ctx context.Context, assets []xdr.Asset, maxBatchSize int) (map[string]Asset, error)
	CreateHistoryLiquidityPools(ctx context.Context, poolIDs []string, batchSize int) (map[string]int64, error)
}
