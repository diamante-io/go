-- +migrate Up

CREATE TABLE offers (
    sellerid character varying(56) NOT NULL,
    offerid bigint PRIMARY KEY,
    sellingasset text NOT NULL,
    buyingasset text NOT NULL,
    amount bigint NOT NULL,
    pricen integer NOT NULL,
    priced integer NOT NULL,
    price double precision NOT NULL,
    flags integer NOT NULL,
    last_modified_ledger INT NOT NULL
);

CREATE INDEX offers_by_seller ON offers USING BTREE(sellerid);
CREATE INDEX offers_by_selling_asset ON offers USING BTREE(sellingasset);
CREATE INDEX offers_by_buying_asset ON offers USING BTREE(buyingasset);
CREATE INDEX offers_by_last_modified_ledger ON offers USING BTREE(last_modified_ledger);

-- Distributed ingestion relies on a single value locked for updating
-- in a DB. When Aurora starts clear there is no value so we create it
-- here. If there's a conflict it means the value is already there so
-- we do nothing.
INSERT INTO key_value_store (key, value)
    VALUES ('exp_ingest_last_ledger', '0')
    ON CONFLICT (key) DO NOTHING;

-- +migrate Down

DROP TABLE offers cascade;
