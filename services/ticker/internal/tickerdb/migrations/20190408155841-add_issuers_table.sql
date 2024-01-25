
-- +migrate Up
CREATE TABLE public.issuers (
    id serial NOT NULL PRIMARY KEY,
    public_key text NOT NULL,
    name text NOT NULL,
    url text NOT NULL,
    toml_url text NOT NULL,
    federation_server text NOT NULL,
    auth_server text NOT NULL,
    transfer_server text NOT NULL,
    web_auth_endpoint text NOT NULL,
    deposit_server text NOT NULL,
    org_twitter text NOT NULL
);

-- Issuer public key should be unique
ALTER TABLE ONLY public.issuers
    ADD CONSTRAINT public_key_unique UNIQUE (public_key);

-- Add FK from assets to issuers
ALTER TABLE public.assets
    ADD COLUMN issuer_id integer NOT NULL;

ALTER TABLE public.assets
    ADD CONSTRAINT fkey_assets_issuers FOREIGN KEY (issuer_id) REFERENCES issuers (id);

-- Delete Public Key from assets
ALTER TABLE public.assets
    DROP COLUMN "public_key";

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_code_issuer_key UNIQUE (code, issuer_id);

-- +migrate Down

