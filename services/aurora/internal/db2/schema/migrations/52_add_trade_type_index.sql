-- +migrate Up

ALTER TABLE history_trades ADD trade_type smallint DEFAULT 1 CHECK(trade_type > 0);
UPDATE history_trades SET trade_type = 2 WHERE base_liquidity_pool_id IS NOT NULL OR counter_liquidity_pool_id IS NOT NULL;
CREATE INDEX htrd_by_trade_type ON history_trades USING BTREE(trade_type, history_operation_id, "order");

-- +migrate Down

DROP INDEX htrd_by_trade_type;
ALTER TABLE history_trades DROP trade_type;
