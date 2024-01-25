-- +migrate Up

ALTER TABLE history_transactions ADD account_muxed varchar(69) NULL, ADD fee_account_muxed varchar(69) NULL;
ALTER TABLE history_operations ADD source_account_muxed varchar(69) NULL;
ALTER TABLE history_effects ADD address_muxed varchar(69) NULL;

-- +migrate Down

ALTER TABLE history_transactions DROP account_muxed, DROP fee_account_muxed;
ALTER TABLE history_operations DROP source_account_muxed;
ALTER TABLE history_effects DROP address_muxed;
