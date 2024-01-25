-- +migrate Up

ALTER TABLE history_ledgers ADD successful_transaction_count integer DEFAULT NULL;
ALTER TABLE history_ledgers ADD failed_transaction_count integer DEFAULT NULL;

-- +migrate Down

ALTER TABLE history_ledgers DROP COLUMN successful_transaction_count;
ALTER TABLE history_ledgers DROP COLUMN failed_transaction_count;
