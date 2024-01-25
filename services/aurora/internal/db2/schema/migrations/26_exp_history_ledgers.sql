-- +migrate Up

CREATE TABLE exp_history_ledgers (
    LIKE history_ledgers
    including defaults
    including constraints
    including indexes
);

-- +migrate Down

DROP TABLE exp_history_ledgers cascade;
