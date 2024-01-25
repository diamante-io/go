-- +migrate Up


ALTER TABLE claimable_balances ADD flags int NOT NULL DEFAULT 0;

-- +migrate Down

ALTER TABLE claimable_balances DROP flags;

