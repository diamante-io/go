package tickerdb

import (
	"context"
	"fmt"
	"strings"
)

// RetrieveMarketData retrieves the 24h- and 7d aggregated market data for all
// markets that were active during this period.
func (s *TickerSession) RetrieveMarketData(ctx context.Context) (markets []Market, err error) {
	err = s.SelectRaw(ctx, &markets, marketQuery)
	return
}

// RetrievePartialAggMarkets retrieves the aggregated market data for all
// markets (or for a specific one if PairName != nil) for a given period.
func (s *TickerSession) RetrievePartialAggMarkets(ctx context.Context,
	pairName *string,
	numHoursAgo int,
) (partialMkts []PartialMarket, err error) {
	var bCode, cCode string
	sqlTrue := new(string)
	*sqlTrue = "TRUE"
	optVars := []optionalVar{
		optionalVar{"bAsset.is_valid", sqlTrue},
		optionalVar{"cAsset.is_valid", sqlTrue},
	}

	// parse base and asset codes and add them as SQL parameters
	if pairName != nil {
		bCode, cCode, err = getBaseAndCounterCodes(*pairName)
		if err != nil {
			return
		}
		optVars = append(optVars, []optionalVar{
			optionalVar{"bAsset.code", &bCode},
			optionalVar{"cAsset.code", &cCode},
		}...)
	}

	where, args := generateWhereClause(optVars)
	where += fmt.Sprintf(
		" AND t.ledger_close_time > now() - interval '%d hours'",
		numHoursAgo,
	)
	q := strings.Replace(aggMarketQuery, "__WHERECLAUSE__", where, -1)
	q = strings.Replace(q, "__NUMHOURS__", fmt.Sprintf("%d", numHoursAgo), -1)

	argsInterface := make([]interface{}, len(args))
	for i, v := range args {
		argsInterface[i] = v
	}

	err = s.SelectRaw(ctx, &partialMkts, q, argsInterface...)
	return
}

// RetrievePartialMarkets retrieves data in the PartialMarket format from the database.
// It optionally filters the data according to the provided base and counter asset params
// provided, as well as the numHoursAgo time offset.
func (s *TickerSession) RetrievePartialMarkets(ctx context.Context,
	baseAssetCode *string,
	baseAssetIssuer *string,
	counterAssetCode *string,
	counterAssetIssuer *string,
	numHoursAgo int,
) (partialMkts []PartialMarket, err error) {
	sqlTrue := new(string)
	*sqlTrue = "TRUE"

	where, args := generateWhereClause([]optionalVar{
		optionalVar{"bAsset.is_valid", sqlTrue},
		optionalVar{"cAsset.is_valid", sqlTrue},
		optionalVar{"bAsset.code", baseAssetCode},
		optionalVar{"bAsset.issuer_account", baseAssetIssuer},
		optionalVar{"cAsset.code", counterAssetCode},
		optionalVar{"cAsset.issuer_account", counterAssetIssuer},
	})
	where += fmt.Sprintf(
		" AND t.ledger_close_time > now() - interval '%d hours'",
		numHoursAgo,
	)

	q := strings.Replace(partialMarketQuery, "__WHERECLAUSE__", where, -1)
	q = strings.Replace(q, "__NUMHOURS__", fmt.Sprintf("%d", numHoursAgo), -1)

	argsInterface := make([]interface{}, len(args))
	for i, v := range args {
		argsInterface[i] = v
	}
	err = s.SelectRaw(ctx, &partialMkts, q, argsInterface...)
	return
}

// Retrieve7DRelevantMarkets retrieves the base and counter asset data of the markets
// that were relevant in the last 7-day period.
func (s *TickerSession) Retrieve7DRelevantMarkets(ctx context.Context) (partialMkts []PartialMarket, err error) {
	q := `
	SELECT
		ba.id as base_asset_id, ba.type AS base_asset_type, ba.code AS base_asset_code, ba.issuer_account AS base_asset_issuer,
		ca.id as counter_asset_id, ca.type AS counter_asset_type, ca.code AS counter_asset_code, ca.issuer_account AS counter_asset_issuer
	FROM trades as t
		JOIN assets AS ba ON t.base_asset_id = ba.id
		JOIN assets AS ca ON t.counter_asset_id = ca.id
	WHERE ba.is_valid = TRUE AND ca.is_valid = TRUE AND t.ledger_close_time > now() - interval '7 days'
	GROUP BY ba.id, ba.type, ba.code, ba.issuer_account, ca.id, ca.type, ca.code, ca.issuer_account
	`
	err = s.SelectRaw(ctx, &partialMkts, q)
	return
}

var marketQuery = `
SELECT
	t2.trade_pair_name,
	COALESCE(base_volume_24h, 0.0) as base_volume_24h,
	COALESCE(counter_volume_24h, 0.0) as counter_volume_24h,
	COALESCE(trade_count_24h, 0) as trade_count_24h,
	COALESCE(highest_price_24h, last_price_7d, 0.0) as highest_price_24h,
	COALESCE(lowest_price_24h, last_price_7d, 0.0) as lowest_price_24h,
	COALESCE(price_change_24h, 0.0) as price_change_24h,
	COALESCE(open_price_24h, last_price_7d, 0.0) as open_price_24h,

	COALESCE(base_volume_7d, 0) as base_volume_7d,
	COALESCE(counter_volume_7d, 0) as counter_volume_7d,
	COALESCE(trade_count_7d, 0) as trade_count_7d,
	COALESCE(highest_price_7d, 0.0) as highest_price_7d,
	COALESCE(lowest_price_7d, 0.0) as lowest_price_7d,
	COALESCE(price_change_7d, 0.0) as price_change_7d,
	COALESCE(open_price_7d, 0.0) as open_price_7d,

	COALESCE(last_price, last_price_7d, 0.0) as last_price,
	COALESCE(last_close_time_24h, last_close_time_7d) as close_time,

	COALESCE(os.num_bids, 0) as num_bids,
	COALESCE(os.bid_volume, 0.0) as bid_volume,
	COALESCE(os.highest_bid, 0.0) as highest_bid,
	COALESCE(os.num_asks, 0) as num_asks,
	COALESCE(os.ask_volume, 0.0) as ask_volume,
	COALESCE(os.lowest_ask, 0.0) as lowest_ask
FROM (
	SELECT
			-- All valid trades for 24h period
			concat(
				COALESCE(NULLIF(bAsset.anchor_asset_code, ''), bAsset.code),
				'_',
				COALESCE(NULLIF(cAsset.anchor_asset_code, ''), cAsset.code)
			) as trade_pair_name,
			sum(t.base_amount) AS base_volume_24h,
			sum(t.counter_amount) AS counter_volume_24h,
			count(t.base_amount) AS trade_count_24h,
			max(t.price) AS highest_price_24h,
			min(t.price) AS lowest_price_24h,
			(array_agg(t.price ORDER BY t.ledger_close_time ASC))[1] AS open_price_24h,
			(array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] AS last_price,
			((array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] - (array_agg(t.price ORDER BY t.ledger_close_time ASC))[1]) AS price_change_24h,
			max(t.ledger_close_time) AS last_close_time_24h
		FROM trades AS t
			JOIN assets AS bAsset ON t.base_asset_id = bAsset.id
			JOIN assets AS cAsset on t.counter_asset_id = cAsset.id
		WHERE bAsset.is_valid = TRUE
			AND cAsset.is_valid = TRUE
			AND t.ledger_close_time > now() - interval '1 day'
		GROUP BY trade_pair_name
	) t1 RIGHT JOIN (
	SELECT
			-- All valid trades for 7d period
			concat(
				COALESCE(NULLIF(bAsset.anchor_asset_code, ''), bAsset.code),
				'_',
				COALESCE(NULLIF(cAsset.anchor_asset_code, ''), cAsset.code)
			) as trade_pair_name,
			sum(t.base_amount) AS base_volume_7d,
			sum(t.counter_amount) AS counter_volume_7d,
			count(t.base_amount) AS trade_count_7d,
			max(t.price) AS highest_price_7d,
			min(t.price) AS lowest_price_7d,
			(array_agg(t.price ORDER BY t.ledger_close_time ASC))[1] AS open_price_7d,
			(array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] AS last_price_7d,
			((array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] - (array_agg(t.price ORDER BY t.ledger_close_time ASC))[1]) AS price_change_7d,
			max(t.ledger_close_time) AS last_close_time_7d
		FROM trades AS t
			LEFT JOIN orderbook_stats AS os
				ON t.base_asset_id = os.base_asset_id AND t.counter_asset_id = os.counter_asset_id
			JOIN assets AS bAsset ON t.base_asset_id = bAsset.id
			JOIN assets AS cAsset on t.counter_asset_id = cAsset.id
		WHERE bAsset.is_valid = TRUE
			AND cAsset.is_valid = TRUE
			AND t.ledger_close_time > now() - interval '7 days'
		GROUP BY trade_pair_name
	) t2 ON t1.trade_pair_name = t2.trade_pair_name
	LEFT JOIN aggregated_orderbook AS os ON t2.trade_pair_name = os.trade_pair_name;
`

var partialMarketQuery = `
SELECT
	concat(bAsset.code, ':', bAsset.issuer_account, ' / ', cAsset.code, ':', cAsset.issuer_account) as trade_pair_name,
	bAsset.id AS base_asset_id,
	bAsset.code AS base_asset_code,
	bAsset.issuer_account as base_asset_issuer,
	bAsset.type as base_asset_type,
	cAsset.id AS counter_asset_id,
	cAsset.code AS counter_asset_code,
	cAsset.issuer_account AS counter_asset_issuer,
	cAsset.type as counter_asset_type,
	sum(t.base_amount) AS base_volume,
	sum(t.counter_amount) AS counter_volume,
	count(t.base_amount) AS trade_count,
	max(t.price) AS highest_price,
	min(t.price) AS lowest_price,
	(array_agg(t.price ORDER BY t.ledger_close_time ASC))[1] AS open_price,
	(array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] AS last_price,
	((array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] - (array_agg(t.price ORDER BY t.ledger_close_time ASC))[1]) AS price_change,
	(now() - interval '__NUMHOURS__ hours') AS interval_start,
	min(t.ledger_close_time) AS first_ledger_close_time,
	max(t.ledger_close_time) AS last_ledger_close_time,
	COALESCE((array_agg(os.num_bids))[1], 0) AS num_bids,
	COALESCE((array_agg(os.bid_volume))[1], 0.0) AS bid_volume,
	COALESCE((array_agg(os.highest_bid))[1], 0.0) AS highest_bid,
	COALESCE((array_agg(os.num_asks))[1], 0) AS num_asks,
	COALESCE((array_agg(os.ask_volume))[1], 0.0) AS ask_volume,
	COALESCE((array_agg(os.lowest_ask))[1], 0.0) AS lowest_ask
FROM trades AS t
	LEFT JOIN orderbook_stats AS os ON t.base_asset_id = os.base_asset_id AND t.counter_asset_id = os.counter_asset_id
	JOIN assets AS bAsset ON t.base_asset_id = bAsset.id
	JOIN assets AS cAsset on t.counter_asset_id = cAsset.id
__WHERECLAUSE__
GROUP BY bAsset.id, bAsset.code, bAsset.issuer_account, bAsset.type, cAsset.id, cAsset.code, cAsset.issuer_account, cAsset.type;
`

var aggMarketQuery = `
SELECT
	t1.trade_pair_name,
	t1.base_volume,
	t1.counter_volume,
	t1.trade_count,
	t1.highest_price,
	t1.lowest_price,
	t1.open_price,
	t1.last_price,
	t1.price_change,
	t1.interval_start,
	t1.first_ledger_close_time,
	t1.last_ledger_close_time,
	COALESCE(aob.base_asset_code, '') as base_asset_code,
	COALESCE(aob.counter_asset_code, '') as counter_asset_code,
	COALESCE(aob.num_bids, 0) AS num_bids,
	COALESCE(aob.bid_volume, 0.0) AS bid_volume,
	COALESCE(aob.highest_bid, 0.0) AS highest_bid,
	COALESCE(aob.num_asks, 0) AS num_asks,
	COALESCE(aob.ask_volume, 0.0) AS ask_volume,
	COALESCE(aob.lowest_ask, 0.0) AS lowest_ask
FROM (
	SELECT
		concat(
			COALESCE(NULLIF(bAsset.anchor_asset_code, ''), bAsset.code),
			'_',
			COALESCE(NULLIF(cAsset.anchor_asset_code, ''), cAsset.code)
		) as trade_pair_name,
		sum(t.base_amount) AS base_volume,
		sum(t.counter_amount) AS counter_volume,
		count(t.base_amount) AS trade_count,
		max(t.price) AS highest_price,
		min(t.price) AS lowest_price,
		(array_agg(t.price ORDER BY t.ledger_close_time ASC))[1] AS open_price,
		(array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] AS last_price,
		((array_agg(t.price ORDER BY t.ledger_close_time DESC))[1] - (array_agg(t.price ORDER BY t.ledger_close_time ASC))[1]) AS price_change,
		(now() - interval '__NUMHOURS__ hours') AS interval_start,
		min(t.ledger_close_time) AS first_ledger_close_time,
		max(t.ledger_close_time) AS last_ledger_close_time
	FROM trades AS t
		LEFT JOIN orderbook_stats AS os ON t.base_asset_id = os.base_asset_id AND t.counter_asset_id = os.counter_asset_id
		JOIN assets AS bAsset ON t.base_asset_id = bAsset.id
		JOIN assets AS cAsset on t.counter_asset_id = cAsset.id
	__WHERECLAUSE__
	GROUP BY trade_pair_name
) t1 LEFT JOIN aggregated_orderbook AS aob ON t1.trade_pair_name = aob.trade_pair_name;`
