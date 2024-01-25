-- +migrate Up

CREATE TABLE liquidity_pools (
    id text NOT NULL, -- hex-encoded PoolID
    type smallint NOT NULL,
    fee integer NOT NULL,
    trustline_count bigint NOT NULL CHECK (trustline_count > 0),
    share_count bigint NOT NULL DEFAULT 0 CHECK(share_count >= 0),
    asset_reserves jsonb NOT NULL,
    last_modified_ledger integer NOT NULL,
    deleted boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);

CREATE INDEX liquidity_pools_by_asset_reserves ON liquidity_pools USING gin(asset_reserves jsonb_path_ops);
CREATE INDEX live_liquidity_pools ON liquidity_pools USING BTREE (deleted, last_modified_ledger);

CREATE SEQUENCE history_liquidity_pools_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE history_liquidity_pools (
    id bigint NOT NULL DEFAULT nextval('history_liquidity_pools_id_seq'::regclass),
    liquidity_pool_id text NOT NULL
);

CREATE UNIQUE INDEX index_history_liquidity_pools_on_id ON history_liquidity_pools USING btree (id);
CREATE UNIQUE INDEX index_history_liquidity_pools_on_liquidity_pool_id ON history_liquidity_pools USING btree (liquidity_pool_id);

CREATE TABLE history_operation_liquidity_pools (
    history_operation_id bigint NOT NULL,
    history_liquidity_pool_id bigint NOT NULL
);

CREATE UNIQUE INDEX index_history_operation_liquidity_pools_on_ids ON history_operation_liquidity_pools USING btree (history_operation_id , history_liquidity_pool_id);
CREATE INDEX index_history_operation_liquidity_pools_on_operation_id ON history_operation_liquidity_pools USING btree (history_operation_id);

CREATE TABLE history_transaction_liquidity_pools (
    history_transaction_id bigint NOT NULL,
    history_liquidity_pool_id bigint NOT NULL
);

CREATE UNIQUE INDEX index_history_transaction_liquidity_pools_on_ids ON history_transaction_liquidity_pools USING btree (history_transaction_id , history_liquidity_pool_id);
CREATE INDEX index_history_transaction_liquidity_pools_on_transaction_id ON history_transaction_liquidity_pools USING btree (history_transaction_id);

ALTER TABLE trust_lines ADD liquidity_pool_id text;
CREATE INDEX trust_lines_by_liquidity_pool_id ON trust_lines USING BTREE(liquidity_pool_id);

DROP INDEX htrd_by_offer;
DROP INDEX htrd_counter_lookup;

ALTER TABLE history_trades DROP offer_id,
                           ALTER base_account_id DROP NOT NULL,
                           ALTER counter_account_id DROP NOT NULL,
                           ADD base_liquidity_pool_id bigint,
                           ADD counter_liquidity_pool_id bigint,
                           ADD liquidity_pool_fee int;

CREATE INDEX htrd_by_base_liquidity_pool_id ON history_trades USING BTREE(base_liquidity_pool_id);
CREATE INDEX htrd_by_counter_liquidity_pool_id ON history_trades USING BTREE(counter_liquidity_pool_id);

-- +migrate Down

DROP INDEX htrd_by_counter_liquidity_pool_id;
DROP INDEX htrd_by_base_liquidity_pool_id;

DELETE FROM history_trades WHERE (counter_account_id IS NULL) OR (base_account_id IS NULL);

ALTER TABLE history_trades DROP liquidity_pool_fee,
                           DROP counter_liquidity_pool_id,
                           DROP base_liquidity_pool_id,
                           ALTER counter_account_id SET NOT NULL,
                           ALTER base_account_id SET NOT NULL,
                           ADD offer_id bigint;

CREATE INDEX htrd_counter_lookup on history_trades USING BTREE(counter_asset_id);
CREATE INDEX htrd_by_offer ON history_trades USING BTREE(offer_id);

DROP INDEX trust_lines_by_liquidity_pool_id;
ALTER TABLE trust_lines DROP liquidity_pool_id;

DROP TABLE history_liquidity_pools cascade;

DROP SEQUENCE history_liquidity_pools_id_seq;

DROP TABLE history_operation_liquidity_pools cascade;

DROP TABLE history_transaction_liquidity_pools cascade;

DROP TABLE liquidity_pools cascade;

