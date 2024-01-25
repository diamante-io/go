-- +migrate Up

ALTER TABLE history_ledgers ADD tx_set_operation_count integer DEFAULT NULL;

-- +migrate Down

ALTER TABLE history_ledgers DROP COLUMN tx_set_operation_count;
