-- +migrate Up
-- exp_asset_stats is like the existing asset_stats table first defined in 8_create_asset_stats_table.sql
-- except that exp_asset_stats is populated by the experimental ingestion system while asset_stats is populated
-- by the legacy ingestion system. Once the experimental ingestion system replaces the legacy system then
-- we can remove asset_stats

CREATE TABLE exp_asset_stats (
    asset_type      INT NOT NULL,
    asset_code      VARCHAR(12) NOT NULL,
    asset_issuer    VARCHAR(56) NOT NULL,
    amount          TEXT        NOT NULL,
    num_accounts    INTEGER     NOT NULL,
    PRIMARY KEY(asset_code, asset_issuer, asset_type)
);

CREATE INDEX exp_asset_stats_by_issuer ON exp_asset_stats USING btree (asset_issuer);
CREATE INDEX exp_asset_stats_by_code   ON exp_asset_stats USING btree (asset_code);

-- +migrate Down
DROP TABLE exp_asset_stats cascade;