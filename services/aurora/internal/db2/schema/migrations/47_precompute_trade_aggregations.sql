-- +migrate Up notransaction

-- Create the new table
CREATE TABLE history_trades_60000 (
  timestamp bigint not null,
  base_asset_id bigint not null,
  counter_asset_id bigint not null,
  count integer not null,
  base_volume numeric not null,
  counter_volume numeric not null,
  avg numeric not null,
  high_n numeric not null,
  high_d numeric not null,
  low_n numeric not null,
  low_d numeric not null,
  open_ledger_toid bigint not null,
  open_n numeric not null,
  open_d numeric not null,
  close_ledger_toid bigint not null,
  close_n numeric not null,
  close_d numeric not null,

  PRIMARY KEY(base_asset_id, counter_asset_id, timestamp)
);

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION to_millis(t timestamp without time zone, trun numeric DEFAULT 1)
  RETURNS bigint AS $$
  BEGIN
    RETURN div(cast((extract(epoch from t) * 1000 ) as bigint), trun)*trun;
  END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION to_millis(t timestamp with time zone, trun numeric DEFAULT 1)
  RETURNS bigint AS $$
  BEGIN
    RETURN to_millis(t::timestamp, trun);
  END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +migrate StatementEnd

CREATE INDEX CONCURRENTLY htrd_agg_bucket_lookup ON history_trades
  USING btree (to_millis(ledger_closed_at, '60000'::numeric));

CREATE INDEX CONCURRENTLY htrd_agg_open_ledger_toid ON history_trades_60000 USING btree (open_ledger_toid);

-- +migrate Down

DROP INDEX htrd_agg_open_ledger_toid;
DROP INDEX htrd_agg_bucket_lookup;
DROP TABLE history_trades_60000;
DROP FUNCTION to_millis(timestamp with time zone, numeric);
DROP FUNCTION to_millis(timestamp without time zone, numeric);
