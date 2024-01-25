-- +migrate Up

ALTER TABLE accounts ADD sponsor TEXT;
CREATE INDEX accounts_by_sponsor ON accounts USING BTREE(sponsor);

ALTER TABLE accounts_data ADD sponsor TEXT;
CREATE INDEX accounts_data_by_sponsor ON accounts_data USING BTREE(sponsor);

ALTER TABLE accounts_signers ADD sponsor TEXT;
CREATE INDEX accounts_signers_by_sponsor ON accounts_signers USING BTREE(sponsor);

ALTER TABLE trust_lines ADD sponsor TEXT;
CREATE INDEX trust_lines_by_sponsor ON trust_lines USING BTREE(sponsor);

ALTER TABLE offers ADD sponsor TEXT;
CREATE INDEX offers_by_sponsor ON offers USING BTREE(sponsor);

-- +migrate Down

ALTER TABLE accounts DROP sponsor;
ALTER TABLE accounts_data DROP sponsor;
ALTER TABLE accounts_signers DROP sponsor;
ALTER TABLE trust_lines DROP sponsor;
ALTER TABLE offers DROP sponsor;
