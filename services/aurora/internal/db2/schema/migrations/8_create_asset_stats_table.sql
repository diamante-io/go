-- +migrate Up
CREATE TABLE asset_stats (
    id              BIGINT      PRIMARY KEY REFERENCES history_assets ON DELETE CASCADE ON UPDATE RESTRICT,
    amount          BIGINT      NOT NULL,
    num_accounts    INTEGER     NOT NULL,
    flags           SMALLINT    NOT NULL,
    toml            VARCHAR(64) NOT NULL
);

CREATE INDEX asset_by_code ON history_assets USING btree (asset_code);

-- +migrate Down
DROP TABLE asset_stats cascade;