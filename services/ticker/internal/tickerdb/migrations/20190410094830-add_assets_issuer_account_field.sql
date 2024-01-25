
-- +migrate Up
ALTER TABLE public.assets
    ADD COLUMN issuer_account text NOT NULL;

ALTER TABLE public.assets
    DROP CONSTRAINT assets_code_issuer_key;

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_code_issuer_account UNIQUE (code, issuer_account);


-- +migrate Down
ALTER TABLE public.assets
    DROP COLUMN issuer_account;
