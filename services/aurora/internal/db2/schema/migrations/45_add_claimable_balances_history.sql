-- +migrate Up

CREATE SEQUENCE history_claimable_balances_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE history_claimable_balances (
    id bigint NOT NULL DEFAULT nextval('history_claimable_balances_id_seq'::regclass),
    claimable_balance_id text NOT NULL
);

CREATE UNIQUE INDEX "index_history_claimable_balances_on_id" ON history_claimable_balances USING btree (id);
CREATE UNIQUE INDEX "index_history_claimable_balances_on_claimable_balance_id" ON history_claimable_balances USING btree (claimable_balance_id);

CREATE TABLE history_operation_claimable_balances (
    history_operation_id bigint NOT NULL,
    history_claimable_balance_id bigint NOT NULL
);

CREATE UNIQUE INDEX "index_history_operation_claimable_balances_on_ids" ON history_operation_claimable_balances USING btree (history_operation_id , history_claimable_balance_id);
CREATE INDEX "index_history_operation_claimable_balances_on_operation_id" ON history_operation_claimable_balances USING btree (history_operation_id);

CREATE TABLE history_transaction_claimable_balances (
    history_transaction_id bigint NOT NULL,
    history_claimable_balance_id bigint NOT NULL
);

CREATE UNIQUE INDEX "index_history_transaction_claimable_balances_on_ids" ON history_transaction_claimable_balances USING btree (history_transaction_id , history_claimable_balance_id);
CREATE INDEX "index_history_transaction_claimable_balances_on_transaction_id" ON history_transaction_claimable_balances USING btree (history_transaction_id);

-- +migrate Down

DROP INDEX "index_history_claimable_balances_on_id";
DROP INDEX "index_history_claimable_balances_on_claimable_balance_id";

DROP TABLE history_claimable_balances;

DROP SEQUENCE history_claimable_balances_id_seq;

DROP INDEX "index_history_operation_claimable_balances_on_ids";
DROP INDEX "index_history_operation_claimable_balances_on_operation_id";

DROP TABLE history_operation_claimable_balances;

DROP INDEX "index_history_transaction_claimable_balances_on_ids";
DROP INDEX "index_history_transaction_claimable_balances_on_transaction_id";

DROP TABLE history_transaction_claimable_balances;
