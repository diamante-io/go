-- +migrate Up
ALTER TABLE asset_stats
ALTER COLUMN amount SET DATA TYPE character varying;

-- +migrate Down
ALTER TABLE asset_stats
ALTER COLUMN amount SET DATA TYPE bigint USING amount::bigint;
