-- +migrate Up
CREATE TABLE history_assets (
    id SERIAL PRIMARY KEY ,
    asset_type VARCHAR(64) NOT NULL,
    asset_code VARCHAR(12) NOT NULL,
    asset_issuer VARCHAR(56) NOT NULL,
    UNIQUE(asset_code, asset_type, asset_issuer)
);

CREATE INDEX asset_by_issuer ON history_assets USING btree (asset_issuer);

-- +migrate Down
DROP TABLE history_assets cascade;