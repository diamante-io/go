-- +migrate Up

CREATE INDEX signers_by_account ON accounts_signers USING BTREE(account);

-- +migrate Down

DROP INDEX signers_by_account;
