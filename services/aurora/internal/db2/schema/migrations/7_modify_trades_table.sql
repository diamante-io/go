-- +migrate Up
-- drop the old trades table. put a new one in place
DROP TABLE history_trades;
CREATE TABLE history_trades (
    history_operation_id BIGINT NOT NULL,
    "order" INTEGER NOT NULL,
    ledger_closed_at TIMESTAMP NOT NULL,
    offer_id BIGINT NOT NULL,
    base_account_id BIGINT NOT NULL REFERENCES history_accounts(id),
    base_asset_id BIGINT NOT NULL REFERENCES history_assets(id),
    base_amount BIGINT NOT NULL CHECK (base_amount > 0),
    counter_account_id BIGINT NOT NULL REFERENCES history_accounts(id),
    counter_asset_id BIGINT NOT NULL REFERENCES history_assets(id),
    counter_amount BIGINT NOT NULL CHECK (counter_amount > 0),
    base_is_seller BOOLEAN,
    CHECK(base_asset_id < counter_asset_id)
);

CREATE UNIQUE INDEX htrd_pid ON history_trades USING btree (history_operation_id, "order");
CREATE INDEX htrd_pair_time_lookup ON history_trades USING BTREE(base_asset_id, counter_asset_id, ledger_closed_at);
CREATE INDEX htrd_counter_lookup ON history_trades USING BTREE(counter_asset_id);
CREATE INDEX htrd_time_lookup ON history_trades USING BTREE(ledger_closed_at);
CREATE INDEX htrd_by_offer ON history_trades USING btree (offer_id);


-- +migrate Down
-- drop the newer table. reinstate the old one.
DROP TABLE history_trades cascade;
CREATE TABLE history_trades (
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    offer_id bigint NOT NULL,
    seller_id bigint NOT NULL,
    buyer_id bigint NOT NULL,
    sold_asset_type character varying(64) NOT NULL,
    sold_asset_issuer character varying(56) NOT NULL,
    sold_asset_code character varying(12) NOT NULL,
    sold_amount bigint NOT NULL CHECK (sold_amount > 0),
    bought_asset_type character varying(64) NOT NULL,
    bought_asset_issuer character varying(56) NOT NULL,
    bought_asset_code character varying(12) NOT NULL,
    bought_amount bigint NOT NULL CHECK (bought_amount > 0)
);

CREATE UNIQUE INDEX htrd_pid ON history_trades USING btree (history_operation_id, "order");
CREATE INDEX htrd_by_offer ON history_trades USING btree (offer_id);
CREATE INDEX htr_by_sold ON history_trades USING btree (sold_asset_type, sold_asset_code, sold_asset_issuer);
CREATE INDEX htr_by_bought ON history_trades USING btree (bought_asset_type, bought_asset_code, bought_asset_issuer);