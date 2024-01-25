-- +migrate Up

ALTER TABLE accounts_signers RENAME COLUMN account TO account_id;

ALTER TABLE offers RENAME COLUMN sellerid TO seller_id;
ALTER TABLE offers RENAME COLUMN offerid TO offer_id;
ALTER TABLE offers RENAME COLUMN sellingasset TO selling_asset;
ALTER TABLE offers RENAME COLUMN buyingasset TO buying_asset;

-- +migrate Down

ALTER TABLE accounts_signers RENAME COLUMN account_id TO account;

ALTER TABLE offers RENAME COLUMN seller_id TO sellerid;
ALTER TABLE offers RENAME COLUMN offer_id TO offerid;
ALTER TABLE offers RENAME COLUMN selling_asset TO sellingasset;
ALTER TABLE offers RENAME COLUMN buying_asset TO buyingasset;
