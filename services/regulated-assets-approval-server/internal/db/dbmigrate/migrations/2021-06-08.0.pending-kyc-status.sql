-- +migrate Up

ALTER TABLE public.accounts_kyc_status
    ADD COLUMN pending_at timestamp with time zone;

-- +migrate Down

ALTER TABLE public.accounts_kyc_status
    DROP COLUMN pending_at;
