-- +migrate Up

CREATE INDEX htrd_by_base_account ON history_trades USING BTREE(base_account_id);
CREATE INDEX htrd_by_counter_account ON history_trades USING BTREE(counter_account_id);

-- +migrate Down

DROP INDEX htrd_by_base_account;
DROP INDEX htrd_by_counter_account;