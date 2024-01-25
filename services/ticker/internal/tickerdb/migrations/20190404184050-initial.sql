
-- +migrate Up

CREATE TABLE public.assets (
    id serial NOT NULL PRIMARY KEY,
    code character varying(12) NOT NULL,
    issuer text NOT NULL,
    type character varying(64) NOT NULL,
    num_accounts integer NOT NULL,
    auth_required boolean NOT NULL,
    auth_revocable boolean NOT NULL,
    amount double precision NOT NULL,
    asset_controlled_by_domain boolean NOT NULL,
    anchor_asset_code character varying(12) NOT NULL,
    anchor_asset_type character varying(64) NOT NULL,
    is_valid boolean NOT NULL,
    validation_error text NOT NULL,
    last_valid timestamp with time zone NOT NULL,
    last_checked timestamp with time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_code_issuer_key UNIQUE (code, issuer);


-- +migrate Down

DROP TABLE public.assets;
