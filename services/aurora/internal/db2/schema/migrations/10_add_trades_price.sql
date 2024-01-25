-- +migrate Up

-- Add rational price to trades table
ALTER TABLE history_trades ADD price_n BIGINT;
ALTER TABLE history_trades ADD price_d BIGINT;

-- aggregate function for finding minimal price when price is represented as an array of two elements [n,d]
CREATE OR REPLACE FUNCTION public.min_price_agg ( NUMERIC[], NUMERIC[])
  RETURNS NUMERIC[] LANGUAGE SQL IMMUTABLE STRICT AS $$ SELECT (
  CASE WHEN $1[1]/$1[2]<$2[1]/$2[2] THEN $1 ELSE $2 END) $$;

CREATE AGGREGATE public.MIN_PRICE (
SFUNC = PUBLIC.min_price_agg,
BASETYPE = NUMERIC[],
STYPE = NUMERIC[]
);

-- aggregate function for finding maximal price when price is represented as an array of two elements [n,d]
CREATE OR REPLACE FUNCTION public.max_price_agg ( NUMERIC[], NUMERIC[])
  RETURNS NUMERIC[] LANGUAGE SQL IMMUTABLE STRICT AS $$ SELECT (
  CASE WHEN $1[1]/$1[2]>$2[1]/$2[2] THEN $1 ELSE $2 END) $$;

CREATE AGGREGATE public.MAX_PRICE (
SFUNC = PUBLIC.max_price_agg,
BASETYPE = NUMERIC[],
STYPE = NUMERIC[]
);

-- +migrate Down
ALTER TABLE history_trades DROP price_n;
ALTER TABLE history_trades DROP price_d;

DROP FUNCTION public.min_price_agg ( NUMERIC[], NUMERIC[] ) CASCADE;
DROP FUNCTION public.max_price_agg ( NUMERIC[], NUMERIC[] ) CASCADE;