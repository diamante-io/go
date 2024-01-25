-- +migrate Up

ALTER TABLE history_transactions
    ALTER COLUMN fee_charged TYPE bigint,
    ALTER COLUMN max_fee TYPE bigint,
    ADD inner_transaction_hash character varying(64),
    ADD fee_account character varying(64),
    ADD inner_signatures character varying(96)[],
    ADD new_max_fee bigint;

CREATE INDEX by_inner_hash ON history_transactions USING btree (inner_transaction_hash) WHERE inner_transaction_hash IS NOT NULL;
CREATE INDEX by_fee_account ON history_transactions USING btree (fee_account) WHERE fee_account IS NOT NULL;

-- +migrate Down

DROP INDEX by_inner_hash;
DROP INDEX by_fee_account;

ALTER TABLE history_transactions
    ALTER COLUMN fee_charged TYPE integer,
    ALTER COLUMN max_fee TYPE integer,
    DROP COLUMN inner_transaction_hash,
    DROP COLUMN fee_account,
    DROP COLUMN inner_signatures,
    DROP COLUMN new_max_fee;