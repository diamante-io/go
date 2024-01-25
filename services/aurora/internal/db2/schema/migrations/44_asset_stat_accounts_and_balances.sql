-- +migrate Up

ALTER TABLE exp_asset_stats
ADD COLUMN accounts JSONB,
ADD COLUMN balances JSONB;
UPDATE exp_asset_stats
  SET
    accounts = jsonb_build_object('authorized', num_accounts),
    balances = jsonb_build_object('authorized', amount);

ALTER TABLE exp_asset_stats
ALTER COLUMN accounts SET NOT NULL,
ALTER COLUMN balances SET NOT NULL;

-- +migrate Down

ALTER TABLE exp_asset_stats
DROP COLUMN accounts,
DROP COLUMN balances;
