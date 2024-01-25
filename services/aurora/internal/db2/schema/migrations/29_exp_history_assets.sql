-- +migrate Up

CREATE TABLE exp_history_assets (
    LIKE history_assets
    including defaults
    including constraints
    including indexes
);

-- +migrate Down

DROP TABLE exp_history_assets cascade;
