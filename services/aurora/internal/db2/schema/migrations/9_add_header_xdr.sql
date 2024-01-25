-- +migrate Up
ALTER TABLE ONLY history_ledgers ADD COLUMN ledger_header text NULL;

-- +migrate Down
ALTER TABLE ONLY history_ledgers DROP COLUMN ledger_header;