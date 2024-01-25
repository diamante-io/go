-- +migrate Up

ALTER TABLE accounts
ADD COLUMN num_sponsored integer DEFAULT 0 CHECK (num_sponsored >= 0),
ADD COLUMN num_sponsoring integer DEFAULT 0 CHECK (num_sponsoring >= 0);

-- +migrate Down

ALTER TABLE accounts
DROP COLUMN num_sponsored,
DROP COLUMN num_sponsoring;
