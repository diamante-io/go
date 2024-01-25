-- +migrate Up

DROP TABLE asset_stats cascade;

DROP INDEX exp_asset_stats_by_code;

DROP INDEX index_history_transactions_on_id;

DROP INDEX index_history_ledgers_on_id;

DROP INDEX asset_by_code;

-- +migrate Down

CREATE TABLE asset_stats (
    id              BIGINT      PRIMARY KEY REFERENCES history_assets ON DELETE CASCADE ON UPDATE RESTRICT,
    amount          BIGINT      NOT NULL,
    num_accounts    INTEGER     NOT NULL,
    flags           SMALLINT    NOT NULL,
    toml            VARCHAR(64) NOT NULL
);

CREATE INDEX asset_by_code ON history_assets USING btree (asset_code);

CREATE INDEX exp_asset_stats_by_code ON exp_asset_stats USING btree (asset_code);

CREATE UNIQUE INDEX index_history_transactions_on_id ON history_transactions USING btree (id);

CREATE UNIQUE INDEX index_history_ledgers_on_id ON history_ledgers USING btree (id);
