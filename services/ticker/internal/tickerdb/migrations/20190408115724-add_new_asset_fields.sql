
-- +migrate Up
ALTER TABLE public.assets
    RENAME COLUMN issuer TO public_key;

ALTER TABLE public.assets
    ADD COLUMN display_decimals integer NOT NULL DEFAULT 7,
    ADD COLUMN "name" text NOT NULL DEFAULT '',
    ADD COLUMN "desc" text NOT NULL DEFAULT '',
    ADD COLUMN conditions text NOT NULL DEFAULT '',
    ADD COLUMN is_asset_anchored boolean NOT NULL DEFAULT FALSE,
    ADD COLUMN fixed_number bigint NOT NULL DEFAULT 0,
    ADD COLUMN max_number bigint NOT NULL DEFAULT 0,
    ADD COLUMN is_unlimited boolean NOT NULL DEFAULT TRUE,
    ADD COLUMN redemption_instructions text NOT NULL DEFAULT '',
    ADD COLUMN collateral_addresses text NOT NULL DEFAULT '',
    ADD COLUMN collateral_address_signatures text NOT NULL DEFAULT '',
    ADD COLUMN countries text NOT NULL DEFAULT '',
    ADD COLUMN "status" text NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE public.assets
    RENAME COLUMN public_key TO issuer;

ALTER TABLE public.assets
    DROP COLUMN display_decimals,
    DROP COLUMN "name",
    DROP COLUMN "desc",
    DROP COLUMN conditions,
    DROP COLUMN is_asset_anchored,
    DROP COLUMN fixed_number,
    DROP COLUMN max_number,
    DROP COLUMN is_unlimited,
    DROP COLUMN redemption_instructions,
    DROP COLUMN collateral_addresses,
    DROP COLUMN collateral_address_signatures,
    DROP COLUMN countries,
    DROP COLUMN "status";
