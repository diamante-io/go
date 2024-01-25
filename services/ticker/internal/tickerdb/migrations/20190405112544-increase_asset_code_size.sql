
-- +migrate Up
ALTER TABLE public.assets
    ALTER COLUMN code type character varying(64);

ALTER TABLE public.assets
    ALTER COLUMN anchor_asset_code type character varying(64);

-- +migrate Down
ALTER TABLE public.assets
    ALTER COLUMN code type character varying(12);

ALTER TABLE public.assets
    ALTER COLUMN anchor_asset_code type character varying(12);
