-- +migrate Up

ALTER TABLE history_transactions DROP CONSTRAINT valid_max_fee;
-- Inner tx in fee bump can have max_fee=0
ALTER TABLE history_transactions ADD CONSTRAINT valid_max_fee CHECK (max_fee >= 0) NOT VALID;

-- +migrate Down

ALTER TABLE history_transactions DROP CONSTRAINT valid_max_fee;
ALTER TABLE history_transactions ADD CONSTRAINT valid_max_fee CHECK (max_fee > 0) NOT VALID;