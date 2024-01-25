-- +migrate Up

ALTER TABLE history_operation_participants
    DROP COLUMN id;

ALTER TABLE history_transaction_participants
    DROP COLUMN id;

-- +migrate Down

ALTER TABLE history_operation_participants
    ADD COLUMN id integer;

ALTER TABLE history_transaction_participants
    ADD COLUMN id integer;