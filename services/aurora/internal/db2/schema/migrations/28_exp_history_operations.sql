-- +migrate Up

CREATE TABLE exp_history_operations (
    LIKE history_operations
    including defaults
    including constraints
    including indexes
);

CREATE TABLE exp_history_operation_participants (
    LIKE history_operation_participants
    including defaults
    including constraints
    including indexes
);

-- +migrate Down

DROP TABLE exp_history_operations cascade;

DROP TABLE exp_history_operation_participants cascade;
