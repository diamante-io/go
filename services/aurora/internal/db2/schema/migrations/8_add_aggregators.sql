-- +migrate Up

-- https://wiki.postgresql.org/wiki/First/last_(aggregate)
-- Create a function that always returns the first non-NULL item
CREATE FUNCTION public.first_agg ( anyelement, anyelement )
  RETURNS anyelement LANGUAGE SQL IMMUTABLE STRICT AS $$ SELECT $1 $$;

-- And then wrap an aggregate around it
CREATE AGGREGATE public.FIRST (
SFUNC = PUBLIC.first_agg,
BASETYPE = ANYELEMENT,
STYPE = ANYELEMENT
);

-- Create a function that always returns the last non-NULL item
CREATE FUNCTION public.last_agg ( anyelement, anyelement )
  RETURNS anyelement LANGUAGE SQL IMMUTABLE STRICT AS $$ SELECT $2 $$;

-- And then wrap an aggregate around it
CREATE AGGREGATE public.LAST (
sfunc    = public.last_agg,
basetype = anyelement,
stype    = anyelement
);

-- +migrate Down
DROP FUNCTION public.first_agg ( anyelement, anyelement ) CASCADE;
DROP FUNCTION public.last_agg (anyelement, anyelement ) CASCADE;