-- +migrate Up

CREATE TABLE claimable_balances (
    id TEXT NOT NULL, -- ClaimableBalanceID in base64
    claimants jsonb NOT NULL,
    asset text NOT NULL,
    amount bigint NOT NULL CHECK (amount > 0),
    sponsor TEXT,
    last_modified_ledger integer NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX claimable_balances_by_asset ON claimable_balances USING BTREE(asset);
CREATE INDEX claimable_balances_by_sponsor ON claimable_balances USING BTREE(sponsor);
CREATE INDEX claimabable_balances_by_claimants ON claimable_balances USING gin(claimants jsonb_path_ops);
CREATE INDEX claimable_balances_by_last_modified_ledger_and_id ON claimable_balances USING BTREE(last_modified_ledger, id);

-- +migrate Down

DROP TABLE claimable_balances cascade;
