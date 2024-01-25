-- +migrate Up

CREATE TABLE exp_history_effects (
    LIKE history_effects
    including defaults
    including constraints
    including indexes
);

-- +migrate Down

DROP TABLE exp_history_effects cascade;
