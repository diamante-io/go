
-- +migrate Up
CREATE TABLE trades (
    id serial NOT NULL PRIMARY KEY,
    aurora_id text NOT NULL UNIQUE,

    ledger_close_time timestamptz NOT NULL,
    offer_id text NOT NULL,

    base_offer_id text NOT NULL,
    base_account text NOT NULL,
    base_amount double precision NOT NULL,
    base_asset_id integer REFERENCES assets (id),

    counter_offer_id text NOT NULL,
    counter_account text NOT NULL,
    counter_amount double precision NOT NULL,
    counter_asset_id integer REFERENCES assets (id),

    base_is_seller boolean NOT NULL,
    price double precision NOT NULL
);

-- +migrate Down
DROP TABLE trades;
