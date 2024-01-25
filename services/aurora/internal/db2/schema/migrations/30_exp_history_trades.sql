-- +migrate Up

-- we cannot create exp_history_trades as:

-- CREATE TABLE exp_history_trades (
--     LIKE history_trades
--     including defaults
--     including constraints
--     including indexes
-- );

-- because the history_trades table has reference constraints to history_accounts and history_assets
-- and we do not want to copy those constraints. instead, we want to reference the
-- exp_history_accounts and exp_history_assets tables

CREATE TABLE exp_history_trades (
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    ledger_closed_at timestamp without time zone NOT NULL,
    offer_id bigint NOT NULL,
    base_account_id bigint NOT NULL REFERENCES exp_history_accounts(id),
    base_asset_id bigint NOT NULL REFERENCES exp_history_assets(id),
    base_amount bigint NOT NULL,
    counter_account_id bigint NOT NULL REFERENCES exp_history_accounts(id),
    counter_asset_id bigint NOT NULL REFERENCES exp_history_assets(id),
    counter_amount bigint NOT NULL,
    base_is_seller boolean,
    price_n bigint,
    price_d bigint,
    base_offer_id bigint,
    counter_offer_id bigint,
    CONSTRAINT exp_history_trades_base_amount_check CHECK ((base_amount >= 0)),
    CONSTRAINT exp_history_trades_check CHECK ((base_asset_id < counter_asset_id)),
    CONSTRAINT exp_history_trades_counter_amount_check CHECK ((counter_amount >= 0))
);


CREATE INDEX exp_htrd_by_base_account ON exp_history_trades USING btree (base_account_id);

CREATE INDEX exp_htrd_by_base_offer ON exp_history_trades USING btree (base_offer_id);

CREATE INDEX exp_htrd_by_counter_account ON exp_history_trades USING btree (counter_account_id);

CREATE INDEX exp_htrd_by_counter_offer ON exp_history_trades USING btree (counter_offer_id);

CREATE INDEX exp_htrd_by_offer ON exp_history_trades USING btree (offer_id);

CREATE INDEX exp_htrd_counter_lookup ON exp_history_trades USING btree (counter_asset_id);

CREATE INDEX exp_htrd_pair_time_lookup ON exp_history_trades USING btree (base_asset_id, counter_asset_id, ledger_closed_at);

CREATE UNIQUE INDEX exp_htrd_pid ON exp_history_trades USING btree (history_operation_id, "order");

CREATE INDEX exp_htrd_time_lookup ON exp_history_trades USING btree (ledger_closed_at);

-- +migrate Down

DROP TABLE exp_history_trades cascade;
