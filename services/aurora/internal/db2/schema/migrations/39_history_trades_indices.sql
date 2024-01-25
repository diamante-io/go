-- +migrate Up

CREATE INDEX htrd_pair_pid ON history_trades USING BTREE(base_asset_id, counter_asset_id, history_operation_id, "order");

-- +migrate Down

DROP INDEX htrd_pair_pid;
