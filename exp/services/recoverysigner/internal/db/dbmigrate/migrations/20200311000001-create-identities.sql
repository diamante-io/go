-- +migrate Up

CREATE TABLE identities (
  account_id BIGINT NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
  id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,

  role TEXT NOT NULL
);

CREATE INDEX ON identities (account_id);

-- +migrate Down

DROP TABLE identities;
