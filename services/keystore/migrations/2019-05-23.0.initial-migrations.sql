-- +migrate Up

CREATE TABLE public.encrypted_keys (
    user_id text NOT NULL PRIMARY KEY,
    encrypted_keys_data bytea NOT NULL,
    salt text NOT NULL,
    encrypter_name text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    modified_at timestamp with time zone
);

-- +migrate Down

DROP TABLE public.encrypted_keys;
