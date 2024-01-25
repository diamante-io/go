
-- +migrate Up
CREATE TABLE orderbook_stats (
    id serial NOT NULL PRIMARY KEY,

    base_asset_id integer REFERENCES assets (id) NOT NULL,
    counter_asset_id integer REFERENCES assets (id) NOT NULL,

    num_bids bigint NOT NULL,
    bid_volume double precision NOT NULL,
    highest_bid double precision NOT NULL,

    num_asks bigint NOT NULL,
    ask_volume double precision NOT NULL,
    lowest_ask double precision NOT NULL,

    spread double precision NOT NULL,
    spread_mid_point double precision NOT NULL,

    updated_at timestamptz NOT NULL
);
ALTER TABLE ONLY public.orderbook_stats
    ADD CONSTRAINT orderbook_stats_base_counter_asset_key UNIQUE (base_asset_id, counter_asset_id);

-- +migrate Down
DROP TABLE orderbook_stats;
