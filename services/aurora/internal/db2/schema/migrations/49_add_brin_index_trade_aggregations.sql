-- +migrate Up notransaction

CREATE INDEX CONCURRENTLY IF NOT EXISTS htrd_agg_timestamp_brin ON history_trades_60000 USING brin(timestamp);

-- +migrate Down

DROP INDEX IF EXISTS htrd_agg_timestamp_brin;
