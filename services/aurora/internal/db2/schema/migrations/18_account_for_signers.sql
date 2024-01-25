-- +migrate Up

CREATE TABLE accounts_signers (
    account character varying(64),
    signer character varying(64),
    weight integer NOT NULL,
    -- we will query by signer so that is why signer is the first item in the composite key
    PRIMARY KEY (signer, account)
);

CREATE TABLE key_value_store (
  key varchar(255) NOT NULL,
  value varchar(255) NOT NULL,
  PRIMARY KEY (key)
);

-- +migrate Down
DROP TABLE accounts_signers cascade;
DROP TABLE key_value_store cascade;
