-- +migrate Up

CREATE TABLE accounts (
    account_id character varying(56) NOT NULL,
    balance bigint NOT NULL,
    buying_liabilities bigint NOT NULL,
    selling_liabilities bigint NOT NULL,
    sequence_number bigint NOT NULL,
    num_subentries int NOT NULL,
    inflation_destination character varying(56) NOT NULL,
    flags int NOT NULL,
    home_domain character varying(32) NOT NULL,
    master_weight smallint NOT NULL,
    threshold_low smallint NOT NULL,
    threshold_medium smallint NOT NULL,
    threshold_high smallint NOT NULL,
    last_modified_ledger INT NOT NULL,
    PRIMARY KEY (account_id)
);

CREATE INDEX accounts_inflation_destination ON accounts USING BTREE(inflation_destination);
CREATE INDEX accounts_home_domain ON accounts USING BTREE(home_domain);

CREATE TABLE accounts_data (
    -- ledger_key is a LedgerKey marshaled using MarshalBinary
    -- and base64-encoded used to boost perfomance of some queries.
    ledger_key character varying(150) NOT NULL,
    account_id character varying(56) NOT NULL,
    name character varying(64) NOT NULL,
    value character varying(90) NOT NULL, -- base64-encoded 64 bytes
    last_modified_ledger INT NOT NULL,
    PRIMARY KEY (ledger_key)
);

CREATE UNIQUE INDEX accounts_data_account_id_name ON accounts_data USING BTREE(account_id, name);

-- +migrate Down

DROP TABLE accounts cascade;
DROP TABLE accounts_data cascade;