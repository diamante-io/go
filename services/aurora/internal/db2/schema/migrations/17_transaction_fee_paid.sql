-- +migrate Up

ALTER TABLE history_transactions ADD fee_charged integer;
ALTER TABLE history_transactions RENAME COLUMN fee_paid TO max_fee;

-- +migrate Down

ALTER TABLE history_transactions DROP COLUMN fee_charged;
ALTER TABLE history_transactions RENAME COLUMN max_fee TO fee_paid;
