
-- +migrate Up
CREATE OR REPLACE VIEW aggregated_orderbook AS
    SELECT
        concat(bAsset.code, '_', cAsset.code) as trade_pair_name,
        bAsset.code as base_asset_code,
        cAsset.code as counter_asset_code,
        COALESCE(sum(os.num_bids), 0) AS num_bids,
        COALESCE(sum(os.bid_volume), 0.0) AS bid_volume,
        COALESCE(max(os.highest_bid), 0.0) AS highest_bid,
        COALESCE(sum(os.num_asks), 0) AS num_asks,
        COALESCE(sum(os.ask_volume), 0.0) AS ask_volume,
        COALESCE(min(os.lowest_ask), 0.0) AS lowest_ask
    FROM orderbook_stats AS os
    JOIN assets AS bAsset ON os.base_asset_id = bAsset.id
    JOIN assets AS cAsset on os.counter_asset_id = cAsset.id
    GROUP BY trade_pair_name, base_asset_code, counter_asset_code;

-- +migrate Down
DROP VIEW IF EXISTS aggregated_orderbook;
