-- +migrate Up
ALTER TABLE ONLY history_ledgers 
  ADD COLUMN protocol_version integer DEFAULT 0 NOT NULL;

-- +migrate Down
ALTER TABLE ONLY history_ledgers DROP COLUMN protocol_version;
