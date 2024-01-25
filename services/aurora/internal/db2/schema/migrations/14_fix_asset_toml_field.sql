-- +migrate Up

ALTER TABLE asset_stats ALTER COLUMN toml TYPE varchar(255);

-- +migrate Down

ALTER TABLE asset_stats ALTER COLUMN toml TYPE varchar(64);
