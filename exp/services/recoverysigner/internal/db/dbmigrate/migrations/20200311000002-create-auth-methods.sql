-- +migrate Up

CREATE TYPE auth_method_type AS ENUM (
  'diamcircle_address',
  'phone_number',
  'email'
);

CREATE TABLE auth_methods (
  account_id BIGINT NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
  identity_id BIGINT NOT NULL REFERENCES identities (id) ON DELETE CASCADE,
  id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,

  type_ auth_method_type NOT NULL,
  value text NOT NULL
);

CREATE INDEX ON auth_methods (account_id);
CREATE INDEX ON auth_methods (identity_id);
CREATE INDEX ON auth_methods (type_, value);

-- +migrate Down

DROP TABLE auth_methods;
DROP TYPE auth_method_type;
