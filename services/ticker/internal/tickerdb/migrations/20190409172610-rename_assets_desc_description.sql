
-- +migrate Up
ALTER TABLE public.assets
    RENAME COLUMN "desc" TO description;

-- +migrate Down
ALTER TABLE public.assets
    RENAME COLUMN description TO "desc";
