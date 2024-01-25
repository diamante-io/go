-- +migrate Up

ALTER TABLE history_transactions ADD CONSTRAINT valid_ledger_sequence CHECK (ledger_sequence > 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_application_order CHECK (application_order >= 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_account_sequence CHECK (account_sequence >= 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_max_fee CHECK (max_fee > 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_operation_count CHECK (operation_count >= 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_fee_charged CHECK (fee_charged > 0) NOT VALID;
ALTER TABLE history_transactions ADD CONSTRAINT valid_new_max_fee CHECK (new_max_fee > 0) NOT VALID;

ALTER TABLE offers ADD CONSTRAINT valid_amount CHECK (amount >= 0) NOT VALID;
ALTER TABLE offers ADD CONSTRAINT valid_pricen CHECK (pricen >= 0) NOT VALID;
ALTER TABLE offers ADD CONSTRAINT valid_priced CHECK (priced >= 0) NOT VALID;
ALTER TABLE offers ADD CONSTRAINT valid_price CHECK (price >= 0) NOT VALID;
ALTER TABLE offers ADD CONSTRAINT valid_last_modified_ledger CHECK (last_modified_ledger > 0) NOT VALID;

ALTER TABLE trust_lines ADD CONSTRAINT valid_balance CHECK (balance >= 0) NOT VALID;
ALTER TABLE trust_lines ADD CONSTRAINT valid_trust_line_limit CHECK (trust_line_limit >= 0) NOT VALID;
ALTER TABLE trust_lines ADD CONSTRAINT valid_buying_liabilities CHECK (buying_liabilities >= 0) NOT VALID;
ALTER TABLE trust_lines ADD CONSTRAINT valid_selling_liabilities CHECK (selling_liabilities >= 0) NOT VALID;
ALTER TABLE trust_lines ADD CONSTRAINT valid_last_modified_ledger CHECK (last_modified_ledger > 0) NOT VALID;

ALTER TABLE history_trades ADD CONSTRAINT valid_price_n CHECK (price_n >= 0) NOT VALID;
ALTER TABLE history_trades ADD CONSTRAINT valid_price_d CHECK (price_d >= 0) NOT VALID;
ALTER TABLE history_trades ADD CONSTRAINT valid_order CHECK ("order" >= 0) NOT VALID;

ALTER TABLE history_operations ADD CONSTRAINT valid_application_order CHECK (application_order >= 0) NOT VALID;

ALTER TABLE history_ledgers ADD CONSTRAINT valid_sequence CHECK (sequence > 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_transaction_count CHECK (transaction_count >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_operation_count CHECK (operation_count >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_total_coins CHECK (total_coins >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_fee_pool CHECK (fee_pool >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_base_fee CHECK (base_fee >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_base_reserve CHECK (base_reserve >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_max_tx_set_size CHECK (max_tx_set_size >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_successful_transaction_count CHECK (successful_transaction_count >= 0) NOT VALID;
ALTER TABLE history_ledgers ADD CONSTRAINT valid_failed_transaction_count CHECK (failed_transaction_count >= 0) NOT VALID;

ALTER TABLE history_effects ADD CONSTRAINT valid_order CHECK ("order" >= 0) NOT VALID;

ALTER TABLE exp_asset_stats ADD CONSTRAINT valid_num_accounts CHECK (num_accounts >= 0) NOT VALID;

ALTER TABLE accounts_signers ADD CONSTRAINT valid_weight CHECK (weight >= 0) NOT VALID;

ALTER TABLE accounts_data ADD CONSTRAINT valid_last_modified_ledger CHECK (last_modified_ledger > 0) NOT VALID;

ALTER TABLE accounts ADD CONSTRAINT valid_balance CHECK (balance >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_buying_liabilities CHECK (buying_liabilities >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_selling_liabilities CHECK (selling_liabilities >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_last_modified_ledger CHECK (last_modified_ledger > 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_sequence_number CHECK (sequence_number >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_num_subentries CHECK (num_subentries >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_master_weight CHECK (master_weight >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_threshold_low CHECK (threshold_low >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_threshold_medium CHECK (threshold_medium >= 0) NOT VALID;
ALTER TABLE accounts ADD CONSTRAINT valid_threshold_high CHECK (threshold_high >= 0) NOT VALID;

-- +migrate Down

ALTER TABLE history_transactions DROP CONSTRAINT valid_ledger_sequence;
ALTER TABLE history_transactions DROP CONSTRAINT valid_application_order;
ALTER TABLE history_transactions DROP CONSTRAINT valid_account_sequence;
ALTER TABLE history_transactions DROP CONSTRAINT valid_max_fee;
ALTER TABLE history_transactions DROP CONSTRAINT valid_operation_count;
ALTER TABLE history_transactions DROP CONSTRAINT valid_fee_charged;
ALTER TABLE history_transactions DROP CONSTRAINT valid_new_max_fee;

ALTER TABLE offers DROP CONSTRAINT valid_amount;
ALTER TABLE offers DROP CONSTRAINT valid_pricen;
ALTER TABLE offers DROP CONSTRAINT valid_priced;
ALTER TABLE offers DROP CONSTRAINT valid_price;
ALTER TABLE offers DROP CONSTRAINT valid_last_modified_ledger;

ALTER TABLE trust_lines DROP CONSTRAINT valid_balance;
ALTER TABLE trust_lines DROP CONSTRAINT valid_trust_line_limit;
ALTER TABLE trust_lines DROP CONSTRAINT valid_buying_liabilities;
ALTER TABLE trust_lines DROP CONSTRAINT valid_selling_liabilities;
ALTER TABLE trust_lines DROP CONSTRAINT valid_last_modified_ledger;

ALTER TABLE history_trades DROP CONSTRAINT valid_price_n;
ALTER TABLE history_trades DROP CONSTRAINT valid_price_d;
ALTER TABLE history_trades DROP CONSTRAINT valid_order;

ALTER TABLE history_operations DROP CONSTRAINT valid_application_order;

ALTER TABLE history_ledgers DROP CONSTRAINT valid_sequence;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_transaction_count;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_operation_count;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_total_coins;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_fee_pool;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_base_fee;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_base_reserve;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_max_tx_set_size;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_successful_transaction_count;
ALTER TABLE history_ledgers DROP CONSTRAINT valid_failed_transaction_count;

ALTER TABLE history_effects DROP CONSTRAINT valid_order;

ALTER TABLE exp_asset_stats DROP CONSTRAINT valid_num_accounts;

ALTER TABLE accounts_signers DROP CONSTRAINT valid_weight;

ALTER TABLE accounts_data DROP CONSTRAINT valid_last_modified_ledger;

ALTER TABLE accounts DROP CONSTRAINT valid_balance;
ALTER TABLE accounts DROP CONSTRAINT valid_buying_liabilities;
ALTER TABLE accounts DROP CONSTRAINT valid_selling_liabilities;
ALTER TABLE accounts DROP CONSTRAINT valid_last_modified_ledger;
ALTER TABLE accounts DROP CONSTRAINT valid_sequence_number;
ALTER TABLE accounts DROP CONSTRAINT valid_num_subentries;
ALTER TABLE accounts DROP CONSTRAINT valid_master_weight;
ALTER TABLE accounts DROP CONSTRAINT valid_threshold_low;
ALTER TABLE accounts DROP CONSTRAINT valid_threshold_medium;
ALTER TABLE accounts DROP CONSTRAINT valid_threshold_high;
