-- +migrate Up

CREATE TABLE trust_lines (
    -- ledger_key is a LedgerKey marshaled using MarshalBinary
    -- and base64-encoded used to boost perfomance of some queries.
    ledger_key character varying(150) NOT NULL,
    account_id character varying(56) NOT NULL,
    asset_type int NOT NULL,
    asset_issuer character varying(56) NOT NULL,
    asset_code character varying(12) NOT NULL,
    balance bigint NOT NULL,
    trust_line_limit bigint NOT NULL,
    buying_liabilities bigint NOT NULL,
    selling_liabilities bigint NOT NULL,
    flags int NOT NULL,
    last_modified_ledger INT NOT NULL,
    PRIMARY KEY (ledger_key)
);

CREATE INDEX trust_lines_by_account_id ON trust_lines USING BTREE(account_id);
CREATE INDEX trust_lines_by_type_code_issuer ON trust_lines USING BTREE(asset_type, asset_code, asset_issuer);
CREATE INDEX trust_lines_by_issuer ON trust_lines USING BTREE(asset_issuer);

-- +migrate Down

DROP TABLE trust_lines cascade;