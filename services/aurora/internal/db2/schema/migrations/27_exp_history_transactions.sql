-- +migrate Up

CREATE TABLE exp_history_transactions (
    LIKE history_transactions
    including defaults
    including constraints
    including indexes
);

CREATE TABLE exp_history_accounts (
    LIKE history_accounts
    including defaults
    including constraints
    including indexes
);

CREATE TABLE exp_history_transaction_participants (
    LIKE history_transaction_participants
    including defaults
    including constraints
    including indexes
);


-- +migrate Down

DROP TABLE exp_history_transactions cascade;

DROP TABLE exp_history_accounts cascade;

DROP TABLE exp_history_transaction_participants cascade;