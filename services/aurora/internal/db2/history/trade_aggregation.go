package history

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/diamcircle/go/services/aurora/internal/db2"
	"github.com/diamcircle/go/services/aurora/internal/toid"
	"github.com/diamcircle/go/support/errors"
	strtime "github.com/diamcircle/go/support/time"
)

// AllowedResolutions is the set of trade aggregation time windows allowed to be used as the
// `resolution` parameter.
var AllowedResolutions = map[time.Duration]struct{}{
	time.Minute:        {}, //1 minute
	time.Minute * 5:    {}, //5 minutes
	time.Minute * 15:   {}, //15 minutes
	time.Hour:          {}, //1 hour
	time.Hour * 24:     {}, //day
	time.Hour * 24 * 7: {}, //week
}

// StrictResolutionFiltering represents a simple feature flag to determine whether only
// predetermined resolutions of trade aggregations are allowed.
var StrictResolutionFiltering = true

// TradeAggregation represents an aggregation of trades from the trades table
type TradeAggregation struct {
	Timestamp     int64   `db:"timestamp"`
	TradeCount    int64   `db:"count"`
	BaseVolume    string  `db:"base_volume"`
	CounterVolume string  `db:"counter_volume"`
	Average       float64 `db:"avg"`
	HighN         int64   `db:"high_n"`
	HighD         int64   `db:"high_d"`
	LowN          int64   `db:"low_n"`
	LowD          int64   `db:"low_d"`
	OpenN         int64   `db:"open_n"`
	OpenD         int64   `db:"open_d"`
	CloseN        int64   `db:"close_n"`
	CloseD        int64   `db:"close_d"`
}

// TradeAggregationsQ is a helper struct to aid in configuring queries to
// bucket and aggregate trades
type TradeAggregationsQ struct {
	baseAssetID    int64
	counterAssetID int64
	resolution     int64
	offset         int64
	startTime      strtime.Millis
	endTime        strtime.Millis
	pagingParams   db2.PageQuery
}

// GetTradeAggregationsQ initializes a TradeAggregationsQ query builder based on the required parameters
func (q Q) GetTradeAggregationsQ(baseAssetID int64, counterAssetID int64, resolution int64,
	offset int64, pagingParams db2.PageQuery) (*TradeAggregationsQ, error) {

	//convert resolution to a duration struct
	resolutionDuration := time.Duration(resolution) * time.Millisecond
	offsetDuration := time.Duration(offset) * time.Millisecond

	//check if resolution allowed
	if StrictResolutionFiltering {
		if _, ok := AllowedResolutions[resolutionDuration]; !ok {
			return &TradeAggregationsQ{}, errors.New("resolution is not allowed")
		}
	}
	// check if offset is allowed. Offset must be 1) a multiple of an hour 2) less than the resolution and 3)
	// less than 24 hours
	if offsetDuration%time.Hour != 0 || offsetDuration >= time.Hour*24 || offsetDuration > resolutionDuration {
		return &TradeAggregationsQ{}, errors.New("offset is not allowed.")
	}

	return &TradeAggregationsQ{
		baseAssetID:    baseAssetID,
		counterAssetID: counterAssetID,
		resolution:     resolution,
		offset:         offset,
		pagingParams:   pagingParams,
	}, nil
}

// WithStartTime adds an optional lower time boundary filter to the trades being aggregated.
func (q *TradeAggregationsQ) WithStartTime(startTime strtime.Millis) (*TradeAggregationsQ, error) {
	offsetMillis := strtime.MillisFromInt64(q.offset)
	var adjustedStartTime strtime.Millis
	// Round up to offset if the provided start time is less than the offset.
	if startTime < offsetMillis {
		adjustedStartTime = offsetMillis
	} else {
		adjustedStartTime = (startTime - offsetMillis).RoundUp(q.resolution) + offsetMillis
	}
	if !q.endTime.IsNil() && adjustedStartTime > q.endTime {
		return &TradeAggregationsQ{}, errors.New("start time is not allowed")
	} else {
		q.startTime = adjustedStartTime
		return q, nil
	}
}

// WithEndTime adds an upper optional time boundary filter to the trades being aggregated.
func (q *TradeAggregationsQ) WithEndTime(endTime strtime.Millis) (*TradeAggregationsQ, error) {
	// Round upper boundary down, to not deliver partial bucket
	offsetMillis := strtime.MillisFromInt64(q.offset)
	var adjustedEndTime strtime.Millis
	// the end time isn't allowed to be less than the offset
	if endTime < offsetMillis {
		return &TradeAggregationsQ{}, errors.New("end time is not allowed")
	} else {
		adjustedEndTime = (endTime - offsetMillis).RoundDown(q.resolution) + offsetMillis
	}
	if adjustedEndTime < q.startTime {
		return &TradeAggregationsQ{}, errors.New("end time is not allowed")
	} else {
		q.endTime = adjustedEndTime
		return q, nil
	}
}

// GetSql generates a sql statement to aggregate Trades based on given parameters
func (q *TradeAggregationsQ) GetSql() sq.SelectBuilder {
	var orderPreserved bool
	orderPreserved, q.baseAssetID, q.counterAssetID = getCanonicalAssetOrder(q.baseAssetID, q.counterAssetID)

	var bucketSQL sq.SelectBuilder
	if orderPreserved {
		bucketSQL = bucketTrades(q.resolution, q.offset)
	} else {
		bucketSQL = reverseBucketTrades(q.resolution, q.offset)
	}

	bucketSQL = bucketSQL.From("history_trades_60000").
		Where(sq.Eq{"base_asset_id": q.baseAssetID, "counter_asset_id": q.counterAssetID})

	//adjust time range and apply time filters
	bucketSQL = bucketSQL.Where(sq.GtOrEq{"timestamp": q.startTime})
	if !q.endTime.IsNil() {
		bucketSQL = bucketSQL.Where(sq.Lt{"timestamp": q.endTime})
	}

	if q.resolution != 60000 {
		//ensure open/close order for cases when multiple trades occur in the same ledger
		bucketSQL = bucketSQL.OrderBy("timestamp ASC", "open_ledger_toid ASC")
		// Do on-the-fly aggregation for higher resolutions.
		bucketSQL = aggregate(bucketSQL)
	}

	return bucketSQL.
		Limit(q.pagingParams.Limit).
		OrderBy("timestamp " + q.pagingParams.Order)
}

// formatBucketTimestampSelect formats a sql select clause for a bucketed timestamp, based on given resolution
// and the offset. Given a time t, it gives it a timestamp defined by
// f(t) = ((t - offset)/resolution)*resolution + offset.
func formatBucketTimestampSelect(resolution int64, offset int64) string {
	return fmt.Sprintf("((timestamp - %d) / %d) * %d + %d as timestamp", offset, resolution, resolution, offset)
}

// bucketTrades generates a select statement to filter rows from the `history_trades` table in
// a compact form, with a timestamp rounded to resolution and reversed base/counter.
func bucketTrades(resolution int64, offset int64) sq.SelectBuilder {
	return sq.Select(
		formatBucketTimestampSelect(resolution, offset),
		"count",
		"base_volume",
		"counter_volume",
		"avg",
		"high_n",
		"high_d",
		"low_n",
		"low_d",
		"open_n",
		"open_d",
		"close_n",
		"close_d",
	)
}

// reverseBucketTrades generates a select statement to filter rows from the `history_trades` table in
// a compact form, with a timestamp rounded to resolution and reversed base/counter.
func reverseBucketTrades(resolution int64, offset int64) sq.SelectBuilder {
	return sq.Select(
		formatBucketTimestampSelect(resolution, offset),
		"count",
		"base_volume as counter_volume",
		"counter_volume as base_volume",
		"(base_volume::numeric/counter_volume::numeric) as avg",
		"low_n as high_d",
		"low_d as high_n",
		"high_n as low_d",
		"high_d as low_n",
		"open_n as open_d",
		"open_d as open_n",
		"close_n as close_d",
		"close_d as close_n",
	)
}

func aggregate(query sq.SelectBuilder) sq.SelectBuilder {
	return sq.Select(
		"timestamp",
		"sum(\"count\") as count",
		"sum(base_volume) as base_volume",
		"sum(counter_volume) as counter_volume",
		"sum(counter_volume::numeric)/sum(base_volume::numeric) as avg",
		"(max_price(ARRAY[high_n, high_d]))[1] as high_n",
		"(max_price(ARRAY[high_n, high_d]))[2] as high_d",
		"(min_price(ARRAY[low_n, low_d]))[1] as low_n",
		"(min_price(ARRAY[low_n, low_d]))[2] as low_d",
		"(first(ARRAY[open_n, open_d]))[1] as open_n",
		"(first(ARRAY[open_n, open_d]))[2] as open_d",
		"(last(ARRAY[close_n, close_d]))[1] as close_n",
		"(last(ARRAY[close_n, close_d]))[2] as close_d",
	).FromSelect(query, "htrd").GroupBy("timestamp")
}

// RebuildTradeAggregationTimes rebuilds a specific set of trade aggregation
// buckets, (specified by start and end times) to ensure complete data in case
// of partial reingestion.
func (q Q) RebuildTradeAggregationTimes(ctx context.Context, from, to strtime.Millis) error {
	from = from.RoundDown(60_000)
	to = to.RoundDown(60_000)
	// Clear out the old bucket values.
	_, err := q.Exec(ctx, sq.Delete("history_trades_60000").Where(
		sq.GtOrEq{"timestamp": from},
	).Where(
		sq.LtOrEq{"timestamp": to},
	))
	if err != nil {
		return errors.Wrap(err, "could not rebuild trade aggregation bucket")
	}

	// find all related trades
	trades := sq.Select(
		"to_millis(ledger_closed_at, 60000) as timestamp",
		"history_operation_id",
		"\"order\"",
		"base_asset_id",
		"base_amount",
		"counter_asset_id",
		"counter_amount",
		"ARRAY[price_n, price_d] as price",
	).From("history_trades").Where(
		sq.GtOrEq{"to_millis(ledger_closed_at, 60000)": from},
	).Where(
		sq.LtOrEq{"to_millis(ledger_closed_at, 60000)": to},
	).OrderBy("base_asset_id", "counter_asset_id", "history_operation_id", "\"order\"")

	// figure out the new bucket values
	rebuilt := sq.Select(
		"timestamp",
		"base_asset_id",
		"counter_asset_id",
		"count(*) as count",
		"sum(base_amount) as base_volume",
		"sum(counter_amount) as counter_volume",
		"sum(counter_amount::numeric)/sum(base_amount::numeric) as avg",
		"(max_price(price))[1] as high_n",
		"(max_price(price))[2] as high_d",
		"(min_price(price))[1] as low_n",
		"(min_price(price))[2] as low_d",
		"first(history_operation_id) as open_ledger_toid",
		"(first(price))[1] as open_n",
		"(first(price))[2] as open_d",
		"last(history_operation_id) as close_ledger_toid",
		"(last(price))[1] as close_n",
		"(last(price))[2] as close_d",
	).FromSelect(trades, "trades").GroupBy("base_asset_id", "counter_asset_id", "timestamp")

	// Insert the new bucket values.
	_, err = q.Exec(ctx, sq.Insert("history_trades_60000").Select(rebuilt))
	if err != nil {
		return errors.Wrap(err, "could not rebuild trade aggregation bucket")
	}
	return nil
}

// RebuildTradeAggregationBuckets rebuilds a specific set of trade aggregation
// buckets, (specified by start and end ledger seq) to ensure complete data in
// case of partial reingestion.
func (q Q) RebuildTradeAggregationBuckets(ctx context.Context, fromSeq, toSeq uint32) error {
	fromLedgerToid := toid.New(int32(fromSeq), 0, 0).ToInt64()
	// toLedger should be inclusive here.
	toLedgerToid := toid.New(int32(toSeq+1), 0, 0).ToInt64()

	// Get the affected timestamp buckets
	timestamps := sq.Select(
		"to_millis(closed_at, 60000)",
	).From("history_ledgers").Where(
		sq.GtOrEq{"id": fromLedgerToid},
	).Where(
		sq.Lt{"id": toLedgerToid},
	)

	// Get first bucket timestamp in the ledger range
	var from strtime.Millis
	err := q.Get(ctx, &from, timestamps.OrderBy("id").Limit(1))
	if err != nil {
		return errors.Wrap(err, "could not rebuild trade aggregation bucket")
	}

	// Get last bucket timestamp in the ledger range
	var to strtime.Millis
	err = q.Get(ctx, &to, timestamps.OrderBy("id DESC").Limit(1))
	if err != nil {
		return errors.Wrap(err, "could not rebuild trade aggregation bucket")
	}

	return q.RebuildTradeAggregationTimes(ctx, from, to)
}
