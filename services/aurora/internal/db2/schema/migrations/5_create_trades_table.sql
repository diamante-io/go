-- +migrate Up
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

-- +migrate Down
DROP TABLE history_trades cascade;
