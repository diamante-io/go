-- +migrate Up

ALTER TABLE history_trades DROP CONSTRAINT history_trades_base_amount_check;
ALTER TABLE history_trades DROP CONSTRAINT history_trades_counter_amount_check;

ALTER TABLE history_trades ADD CONSTRAINT history_trades_base_amount_check CHECK (base_amount >= 0);
ALTER TABLE history_trades ADD CONSTRAINT history_trades_counter_amount_check CHECK (counter_amount >= 0);

-- +migrate Down

ALTER TABLE history_trades DROP CONSTRAINT history_trades_base_amount_check;
ALTER TABLE history_trades DROP CONSTRAINT history_trades_counter_amount_check;

ALTER TABLE history_trades ADD CONSTRAINT history_trades_base_amount_check CHECK (base_amount > 0);
ALTER TABLE history_trades ADD CONSTRAINT history_trades_counter_amount_check CHECK (counter_amount > 0);
