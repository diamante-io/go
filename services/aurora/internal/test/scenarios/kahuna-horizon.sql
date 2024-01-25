--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.1
-- Dumped by pg_dump version 9.6.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_counter_asset_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_counter_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_base_asset_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_base_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.asset_stats DROP CONSTRAINT IF EXISTS asset_stats_id_fkey;
DROP INDEX IF EXISTS public.trust_lines_by_type_code_issuer;
DROP INDEX IF EXISTS public.trust_lines_by_issuer;
DROP INDEX IF EXISTS public.trust_lines_by_account_id;
DROP INDEX IF EXISTS public.trade_effects_by_order_book;
DROP INDEX IF EXISTS public.signers_by_account;
DROP INDEX IF EXISTS public.offers_by_selling_asset;
DROP INDEX IF EXISTS public.offers_by_seller;
DROP INDEX IF EXISTS public.offers_by_last_modified_ledger;
DROP INDEX IF EXISTS public.offers_by_buying_asset;
DROP INDEX IF EXISTS public.index_history_transactions_on_id;
DROP INDEX IF EXISTS public.index_history_operations_on_type;
DROP INDEX IF EXISTS public.index_history_operations_on_transaction_id;
DROP INDEX IF EXISTS public.index_history_operations_on_id;
DROP INDEX IF EXISTS public.index_history_ledgers_on_sequence;
DROP INDEX IF EXISTS public.index_history_ledgers_on_previous_ledger_hash;
DROP INDEX IF EXISTS public.index_history_ledgers_on_ledger_hash;
DROP INDEX IF EXISTS public.index_history_ledgers_on_importer_version;
DROP INDEX IF EXISTS public.index_history_ledgers_on_id;
DROP INDEX IF EXISTS public.index_history_ledgers_on_closed_at;
DROP INDEX IF EXISTS public.index_history_effects_on_type;
DROP INDEX IF EXISTS public.index_history_accounts_on_id;
DROP INDEX IF EXISTS public.index_history_accounts_on_address;
DROP INDEX IF EXISTS public.htrd_time_lookup;
DROP INDEX IF EXISTS public.htrd_pid;
DROP INDEX IF EXISTS public.htrd_pair_time_lookup;
DROP INDEX IF EXISTS public.htrd_counter_lookup;
DROP INDEX IF EXISTS public.htrd_by_offer;
DROP INDEX IF EXISTS public.htrd_by_counter_offer;
DROP INDEX IF EXISTS public.htrd_by_counter_account;
DROP INDEX IF EXISTS public.htrd_by_base_offer;
DROP INDEX IF EXISTS public.htrd_by_base_account;
DROP INDEX IF EXISTS public.htp_by_htid;
DROP INDEX IF EXISTS public.hs_transaction_by_id;
DROP INDEX IF EXISTS public.hs_ledger_by_id;
DROP INDEX IF EXISTS public.hop_by_hoid;
DROP INDEX IF EXISTS public.hist_tx_p_id;
DROP INDEX IF EXISTS public.hist_op_p_id;
DROP INDEX IF EXISTS public.hist_e_id;
DROP INDEX IF EXISTS public.hist_e_by_order;
DROP INDEX IF EXISTS public.exp_asset_stats_by_issuer;
DROP INDEX IF EXISTS public.exp_asset_stats_by_code;
DROP INDEX IF EXISTS public.by_ledger;
DROP INDEX IF EXISTS public.by_hash;
DROP INDEX IF EXISTS public.by_account;
DROP INDEX IF EXISTS public.asset_by_issuer;
DROP INDEX IF EXISTS public.asset_by_code;
DROP INDEX IF EXISTS public.accounts_inflation_destination;
DROP INDEX IF EXISTS public.accounts_home_domain;
DROP INDEX IF EXISTS public.accounts_data_account_id_name;
ALTER TABLE IF EXISTS ONLY public.trust_lines DROP CONSTRAINT IF EXISTS trust_lines_pkey;
ALTER TABLE IF EXISTS ONLY public.offers DROP CONSTRAINT IF EXISTS offers_pkey;
ALTER TABLE IF EXISTS ONLY public.key_value_store DROP CONSTRAINT IF EXISTS key_value_store_pkey;
ALTER TABLE IF EXISTS ONLY public.history_transaction_participants DROP CONSTRAINT IF EXISTS history_transaction_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_operation_participants DROP CONSTRAINT IF EXISTS history_operation_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_asset_code_asset_type_asset_issuer_key;
ALTER TABLE IF EXISTS ONLY public.gorp_migrations DROP CONSTRAINT IF EXISTS gorp_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.exp_asset_stats DROP CONSTRAINT IF EXISTS exp_asset_stats_pkey;
ALTER TABLE IF EXISTS ONLY public.asset_stats DROP CONSTRAINT IF EXISTS asset_stats_pkey;
ALTER TABLE IF EXISTS ONLY public.accounts_signers DROP CONSTRAINT IF EXISTS accounts_signers_pkey;
ALTER TABLE IF EXISTS ONLY public.accounts DROP CONSTRAINT IF EXISTS accounts_pkey;
ALTER TABLE IF EXISTS ONLY public.accounts_data DROP CONSTRAINT IF EXISTS accounts_data_pkey;
ALTER TABLE IF EXISTS public.history_transaction_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_operation_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_assets ALTER COLUMN id DROP DEFAULT;
DROP TABLE IF EXISTS public.trust_lines;
DROP TABLE IF EXISTS public.offers;
DROP TABLE IF EXISTS public.key_value_store;
DROP TABLE IF EXISTS public.history_transactions;
DROP SEQUENCE IF EXISTS public.history_transaction_participants_id_seq;
DROP TABLE IF EXISTS public.history_transaction_participants;
DROP TABLE IF EXISTS public.history_trades;
DROP INDEX IF EXISTS public.htrd_agg_open_ledger_toid;
DROP INDEX IF EXISTS public.htrd_agg_bucket_lookup;
DROP TABLE IF EXISTS public.history_trades_60000;
DROP FUNCTION IF EXISTS public.to_millis(timestamp with time zone, numeric);
DROP FUNCTION IF EXISTS public.to_millis(timestamp without time zone, numeric);
DROP TABLE IF EXISTS public.history_operations;
DROP SEQUENCE IF EXISTS public.history_operation_participants_id_seq;
DROP TABLE IF EXISTS public.history_operation_participants;
DROP TABLE IF EXISTS public.history_ledgers;
DROP TABLE IF EXISTS public.history_effects;
DROP SEQUENCE IF EXISTS public.history_assets_id_seq;
DROP TABLE IF EXISTS public.history_assets;
DROP TABLE IF EXISTS public.history_accounts;
DROP SEQUENCE IF EXISTS public.history_accounts_id_seq;
DROP TABLE IF EXISTS public.gorp_migrations;
DROP TABLE IF EXISTS public.exp_asset_stats;
DROP TABLE IF EXISTS public.asset_stats;
DROP TABLE IF EXISTS public.accounts_signers;
DROP TABLE IF EXISTS public.accounts_data;
DROP TABLE IF EXISTS public.accounts;
DROP TABLE IF EXISTS public.history_operation_liquidity_pools;
DROP TABLE IF EXISTS public.history_transaction_liquidity_pools;
DROP AGGREGATE IF EXISTS public.min_price(numeric[]);
DROP AGGREGATE IF EXISTS public.max_price(numeric[]);
DROP AGGREGATE IF EXISTS public.last(anyelement);
DROP AGGREGATE IF EXISTS public.first(anyelement);
DROP FUNCTION IF EXISTS public.min_price_agg(numeric[], numeric[]);
DROP FUNCTION IF EXISTS public.max_price_agg(numeric[], numeric[]);
DROP FUNCTION IF EXISTS public.last_agg(anyelement, anyelement);
DROP FUNCTION IF EXISTS public.first_agg(anyelement, anyelement);
DROP EXTENSION IF EXISTS plpgsql;
DROP SCHEMA IF EXISTS public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: first_agg(anyelement, anyelement); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION first_agg(anyelement, anyelement) RETURNS anyelement
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT $1 $_$;


--
-- Name: last_agg(anyelement, anyelement); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION last_agg(anyelement, anyelement) RETURNS anyelement
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT $2 $_$;


--
-- Name: max_price_agg(numeric[], numeric[]); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION max_price_agg(numeric[], numeric[]) RETURNS numeric[]
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT (
  CASE WHEN $1[1]/$1[2]>$2[1]/$2[2] THEN $1 ELSE $2 END) $_$;


--
-- Name: min_price_agg(numeric[], numeric[]); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION min_price_agg(numeric[], numeric[]) RETURNS numeric[]
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT (
  CASE WHEN $1[1]/$1[2]<$2[1]/$2[2] THEN $1 ELSE $2 END) $_$;


--
-- Name: first(anyelement); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE first(anyelement) (
    SFUNC = first_agg,
    STYPE = anyelement
);


--
-- Name: last(anyelement); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE last(anyelement) (
    SFUNC = last_agg,
    STYPE = anyelement
);


--
-- Name: max_price(numeric[]); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE max_price(numeric[]) (
    SFUNC = max_price_agg,
    STYPE = numeric[]
);


--
-- Name: min_price(numeric[]); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE min_price(numeric[]) (
    SFUNC = min_price_agg,
    STYPE = numeric[]
);


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accounts (
    account_id character varying(56) NOT NULL,
    balance bigint NOT NULL,
    buying_liabilities bigint NOT NULL,
    selling_liabilities bigint NOT NULL,
    sequence_number bigint NOT NULL,
    num_subentries integer NOT NULL,
    inflation_destination character varying(56) NOT NULL,
    flags integer NOT NULL,
    home_domain character varying(32) NOT NULL,
    master_weight smallint NOT NULL,
    threshold_low smallint NOT NULL,
    threshold_medium smallint NOT NULL,
    threshold_high smallint NOT NULL,
    last_modified_ledger integer NOT NULL
);


--
-- Name: accounts_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accounts_data (
    ledger_key character varying(150) NOT NULL,
    account_id character varying(56) NOT NULL,
    name character varying(64) NOT NULL,
    value character varying(90) NOT NULL,
    last_modified_ledger integer NOT NULL
);


--
-- Name: accounts_signers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accounts_signers (
    account_id character varying(64) NOT NULL,
    signer character varying(64) NOT NULL,
    weight integer NOT NULL,
    sponsor character varying(56)
);

CREATE INDEX accounts_signers_by_sponsor ON accounts_signers USING BTREE(sponsor);


--
-- Name: asset_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE asset_stats (
    id bigint NOT NULL,
    amount character varying NOT NULL,
    num_accounts integer NOT NULL,
    flags smallint NOT NULL,
    toml character varying(255) NOT NULL
);


--
-- Name: exp_asset_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE exp_asset_stats (
    asset_type integer NOT NULL,
    asset_code character varying(12) NOT NULL,
    asset_issuer character varying(56) NOT NULL,
    amount text NOT NULL,
    num_accounts integer NOT NULL
);


--
-- Name: gorp_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE gorp_migrations (
    id text NOT NULL,
    applied_at timestamp with time zone
);


--
-- Name: history_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE history_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_accounts (
    id bigint DEFAULT nextval('history_accounts_id_seq'::regclass) NOT NULL,
    address character varying(64)
);


-- history_claimable_balances (manually added)
CREATE SEQUENCE history_claimable_balances_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE history_claimable_balances (
    id bigint NOT NULL DEFAULT nextval('history_claimable_balances_id_seq'::regclass),
    claimable_balance_id text NOT NULL
);

CREATE UNIQUE INDEX "index_history_claimable_balances_on_id" ON history_claimable_balances USING btree (id);
CREATE UNIQUE INDEX "index_history_claimable_balances_on_claimable_balance_id" ON history_claimable_balances USING btree (claimable_balance_id);

CREATE TABLE history_operation_claimable_balances (
    history_operation_id bigint NOT NULL,
    history_claimable_balance_id bigint NOT NULL
);

CREATE UNIQUE INDEX "index_history_operation_claimable_balances_on_ids" ON history_operation_claimable_balances USING btree (history_operation_id , history_claimable_balance_id);
CREATE INDEX "index_history_operation_claimable_balances_on_operation_id" ON history_operation_claimable_balances USING btree (history_operation_id);

CREATE TABLE history_transaction_claimable_balances (
    history_transaction_id bigint NOT NULL,
    history_claimable_balance_id bigint NOT NULL
);

CREATE UNIQUE INDEX "index_history_transaction_claimable_balances_on_ids" ON history_transaction_claimable_balances USING btree (history_transaction_id , history_claimable_balance_id);
CREATE INDEX "index_history_transaction_claimable_balances_on_transaction_id" ON history_transaction_claimable_balances USING btree (history_transaction_id);


INSERT INTO history_claimable_balances VALUES (1, '00000000178826fbfe339e1f5c53417c6fedfe2c05e8bec14303143ec46b38981b09c3f9');
SELECT pg_catalog.setval('history_claimable_balances_id_seq', 1, true);
-- The operations/transactions are going to be unrelated to claimable balances, but it doesn't matter for testing
INSERT INTO history_operation_claimable_balances VALUES (12884905985, 1);
INSERT INTO history_operation_claimable_balances VALUES (8589938689, 1);
INSERT INTO history_transaction_claimable_balances VALUES (12884905984, 1);
INSERT INTO history_transaction_claimable_balances VALUES (8589938688, 1);

--
-- Name: history_assets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_assets (
    id integer NOT NULL,
    asset_type character varying(64) NOT NULL,
    asset_code character varying(12) NOT NULL,
    asset_issuer character varying(56) NOT NULL
);


--
-- Name: history_assets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE history_assets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_assets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE history_assets_id_seq OWNED BY history_assets.id;


--
-- Name: history_effects; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_effects (
    history_account_id bigint NOT NULL,
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    type integer NOT NULL,
    details jsonb
);


--
-- Name: history_ledgers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_ledgers (
    sequence integer NOT NULL,
    ledger_hash character varying(64) NOT NULL,
    previous_ledger_hash character varying(64),
    transaction_count integer DEFAULT 0 NOT NULL,
    operation_count integer DEFAULT 0 NOT NULL,
    closed_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    id bigint,
    importer_version integer DEFAULT 1 NOT NULL,
    total_coins bigint NOT NULL,
    fee_pool bigint NOT NULL,
    base_fee integer NOT NULL,
    base_reserve integer NOT NULL,
    max_tx_set_size integer NOT NULL,
    protocol_version integer DEFAULT 0 NOT NULL,
    ledger_header text,
    successful_transaction_count integer,
    failed_transaction_count integer,
    tx_set_operation_count integer
);


--
-- Name: history_operation_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_operation_participants (
    id integer NOT NULL,
    history_operation_id bigint NOT NULL,
    history_account_id bigint NOT NULL
);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE history_operation_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE history_operation_participants_id_seq OWNED BY history_operation_participants.id;


--
-- Name: history_operations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_operations (
    id bigint NOT NULL,
    transaction_id bigint NOT NULL,
    application_order integer NOT NULL,
    type integer NOT NULL,
    details jsonb,
    source_account character varying(64) DEFAULT ''::character varying NOT NULL
);


--
-- Name: history_trades; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_trades (
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    ledger_closed_at timestamp without time zone NOT NULL,
    offer_id bigint NOT NULL,
    base_account_id bigint NOT NULL,
    base_asset_id bigint NOT NULL,
    base_amount bigint NOT NULL,
    counter_account_id bigint NOT NULL,
    counter_asset_id bigint NOT NULL,
    counter_amount bigint NOT NULL,
    base_is_seller boolean,
    price_n bigint,
    price_d bigint,
    base_offer_id bigint,
    counter_offer_id bigint,
    CONSTRAINT history_trades_base_amount_check CHECK ((base_amount >= 0)),
    CONSTRAINT history_trades_check CHECK ((base_asset_id < counter_asset_id)),
    CONSTRAINT history_trades_counter_amount_check CHECK ((counter_amount >= 0))
);

--
-- Name: history_trades; Type: TABLE; Schema: public; Owner: -
--

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

CREATE OR REPLACE FUNCTION to_millis(t timestamp without time zone, trun numeric DEFAULT 1)
  RETURNS bigint AS $$
  BEGIN
    RETURN div(cast((extract(epoch from t) * 1000 ) as bigint), trun)*trun;
  END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION to_millis(t timestamp with time zone, trun numeric DEFAULT 1)
  RETURNS bigint AS $$
  BEGIN
    RETURN to_millis(t::timestamp, trun);
  END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE INDEX htrd_agg_bucket_lookup ON history_trades
  USING btree (to_millis(ledger_closed_at, '60000'::numeric));

CREATE INDEX htrd_agg_open_ledger_toid ON history_trades_60000 USING btree (open_ledger_toid);


--
-- Name: history_transaction_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE history_transaction_participants (
    id integer NOT NULL,
    history_transaction_id bigint NOT NULL,
    history_account_id bigint NOT NULL
);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE history_transaction_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE history_transaction_participants_id_seq OWNED BY history_transaction_participants.id;


--
-- Name: history_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_transactions (
    transaction_hash character varying(64) NOT NULL,
    ledger_sequence integer NOT NULL,
    application_order integer NOT NULL,
    account character varying(64) NOT NULL,
    account_sequence bigint NOT NULL,
    max_fee integer NOT NULL,
    operation_count integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    id bigint,
    tx_envelope text NOT NULL,
    tx_result text NOT NULL,
    tx_meta text NOT NULL,
    tx_fee_meta text NOT NULL,
    signatures character varying(96)[] DEFAULT '{}'::character varying[] NOT NULL,
    memo_type character varying DEFAULT 'none'::character varying NOT NULL,
    memo character varying,
    time_bounds int8range,
    successful boolean,
    fee_charged bigint,
    inner_transaction_hash character varying(64),
    fee_account character varying(64),
    inner_signatures character varying(96)[],
    new_max_fee bigint
);

--
-- Name: key_value_store; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE key_value_store (
    key character varying(255) NOT NULL,
    value character varying(255) NOT NULL
);


--
-- Name: offers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE offers (
    seller_id character varying(56) NOT NULL,
    offer_id bigint NOT NULL,
    selling_asset text NOT NULL,
    buying_asset text NOT NULL,
    amount bigint NOT NULL,
    pricen integer NOT NULL,
    priced integer NOT NULL,
    price double precision NOT NULL,
    flags integer NOT NULL,
    last_modified_ledger integer NOT NULL
);


--
-- Name: trust_lines; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE trust_lines (
    ledger_key character varying(150) NOT NULL,
    account_id character varying(56) NOT NULL,
    asset_type integer NOT NULL,
    asset_issuer character varying(56) NOT NULL,
    asset_code character varying(12) NOT NULL,
    balance bigint NOT NULL,
    trust_line_limit bigint NOT NULL,
    buying_liabilities bigint NOT NULL,
    selling_liabilities bigint NOT NULL,
    flags integer NOT NULL,
    last_modified_ledger integer NOT NULL
);

--
-- Name: history_assets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_assets ALTER COLUMN id SET DEFAULT nextval('history_assets_id_seq'::regclass);


--
-- Name: history_operation_participants id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_operation_participants ALTER COLUMN id SET DEFAULT nextval('history_operation_participants_id_seq'::regclass);


--
-- Name: history_transaction_participants id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_transaction_participants ALTER COLUMN id SET DEFAULT nextval('history_transaction_participants_id_seq'::regclass);


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: accounts_data; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: accounts_signers; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: asset_stats; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO asset_stats VALUES (1, '200000000', 1, 0, '');
INSERT INTO asset_stats VALUES (2, '0', 1, 0, '');
INSERT INTO asset_stats VALUES (3, '0', 0, 3, '');
INSERT INTO asset_stats VALUES (4, '3000000000', 1, 0, '');
INSERT INTO asset_stats VALUES (5, '0', 1, 3, '');
INSERT INTO asset_stats VALUES (6, '100000000', 1, 0, '');


--
-- Data for Name: gorp_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO gorp_migrations VALUES ('1_initial_schema.sql', '2019-10-31 14:19:49.03833+01');
INSERT INTO gorp_migrations VALUES ('2_index_participants_by_toid.sql', '2019-10-31 14:19:49.04267+01');
INSERT INTO gorp_migrations VALUES ('3_use_sequence_in_history_accounts.sql', '2019-10-31 14:19:49.045926+01');
INSERT INTO gorp_migrations VALUES ('4_add_protocol_version.sql', '2019-10-31 14:19:49.054147+01');
INSERT INTO gorp_migrations VALUES ('5_create_trades_table.sql', '2019-10-31 14:19:49.061804+01');
INSERT INTO gorp_migrations VALUES ('6_create_assets_table.sql', '2019-10-31 14:19:49.067093+01');
INSERT INTO gorp_migrations VALUES ('7_modify_trades_table.sql', '2019-10-31 14:19:49.081047+01');
INSERT INTO gorp_migrations VALUES ('8_add_aggregators.sql', '2019-10-31 14:19:49.085128+01');
INSERT INTO gorp_migrations VALUES ('8_create_asset_stats_table.sql', '2019-10-31 14:19:49.089574+01');
INSERT INTO gorp_migrations VALUES ('9_add_header_xdr.sql', '2019-10-31 14:19:49.092366+01');
INSERT INTO gorp_migrations VALUES ('10_add_trades_price.sql', '2019-10-31 14:19:49.095671+01');
INSERT INTO gorp_migrations VALUES ('11_add_trades_account_index.sql', '2019-10-31 14:19:49.099289+01');
INSERT INTO gorp_migrations VALUES ('12_asset_stats_amount_string.sql', '2019-10-31 14:19:49.105961+01');
INSERT INTO gorp_migrations VALUES ('13_trade_offer_ids.sql', '2019-10-31 14:19:49.111757+01');
INSERT INTO gorp_migrations VALUES ('14_fix_asset_toml_field.sql', '2019-10-31 14:19:49.113736+01');
INSERT INTO gorp_migrations VALUES ('15_ledger_failed_txs.sql', '2019-10-31 14:19:49.115578+01');
INSERT INTO gorp_migrations VALUES ('16_ingest_failed_transactions.sql', '2019-10-31 14:19:49.116928+01');
INSERT INTO gorp_migrations VALUES ('17_transaction_fee_paid.sql', '2019-10-31 14:19:49.118562+01');
INSERT INTO gorp_migrations VALUES ('18_account_for_signers.sql', '2019-10-31 14:19:49.123835+01');
INSERT INTO gorp_migrations VALUES ('19_offers.sql', '2019-10-31 14:19:49.133107+01');
INSERT INTO gorp_migrations VALUES ('20_account_for_signer_index.sql', '2019-10-31 14:19:49.135499+01');
INSERT INTO gorp_migrations VALUES ('21_trades_remove_zero_amount_constraints.sql', '2019-10-31 14:19:49.138031+01');
INSERT INTO gorp_migrations VALUES ('22_trust_lines.sql', '2019-10-31 14:19:49.144708+01');
INSERT INTO gorp_migrations VALUES ('23_exp_asset_stats.sql', '2019-10-31 14:19:49.15222+01');
INSERT INTO gorp_migrations VALUES ('24_accounts.sql', '2019-10-31 14:19:49.160844+01');
INSERT INTO gorp_migrations VALUES ('25_expingest_rename_columns.sql', '2019-10-31 14:19:49.163717+01');
INSERT INTO gorp_migrations VALUES ('33_remove_unused.sql', '2019-11-30 10:19:49.163718+01');
INSERT INTO gorp_migrations VALUES ('34_fee_bump_transactions.sql', '2019-11-30 11:19:49.163718+01');
INSERT INTO gorp_migrations VALUES ('35_drop_participant_id.sql', '2019-11-30 14:19:49.163728+01');
INSERT INTO gorp_migrations VALUES ('37_add_tx_set_operation_count_to_ledgers.sql', '2019-11-30 12:19:49.163728+01');
INSERT INTO gorp_migrations VALUES ('41_add_sponsor_to_state_tables.sql', '2019-11-30 13:19:49.163718+01');
INSERT INTO gorp_migrations VALUES ('45_add_claimable_balances_history.sql', '2019-11-30 14:19:49.163718+01');
INSERT INTO gorp_migrations VALUES ('46_add_muxed_accounts.sql', '2019-12-30 14:19:49.163718+01');
INSERT INTO gorp_migrations VALUES ('47_precompute_trade_aggregations.sql', '2019-12-30 14:19:49.163719+01');
INSERT INTO gorp_migrations VALUES ('48_rebuild_trade_aggregations.sql', '2021-12-02 01:33:33.428419+00');
INSERT INTO gorp_migrations VALUES ('49_add_brin_index_trade_aggregations.sql', '2021-12-02 01:33:33.43274+00');
INSERT INTO gorp_migrations VALUES ('50_liquidity_pools.sql', '2021-12-02 01:33:33.471893+00');
INSERT INTO gorp_migrations VALUES ('51_remove_ht_unused_indexes.sql', '2021-12-02 01:33:33.47903+00');


--
-- Data for Name: history_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_accounts VALUES (1, 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN');
INSERT INTO history_accounts VALUES (2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_accounts VALUES (3, 'GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y');
INSERT INTO history_accounts VALUES (4, 'GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X');
INSERT INTO history_accounts VALUES (5, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_accounts VALUES (6, 'GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS');
INSERT INTO history_accounts VALUES (7, 'GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ');
INSERT INTO history_accounts VALUES (8, 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_accounts VALUES (9, 'GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG');
INSERT INTO history_accounts VALUES (10, 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG');
INSERT INTO history_accounts VALUES (11, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_accounts VALUES (12, 'GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q');
INSERT INTO history_accounts VALUES (13, 'GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC');
INSERT INTO history_accounts VALUES (14, 'GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD');
INSERT INTO history_accounts VALUES (15, 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_accounts VALUES (16, 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD');
INSERT INTO history_accounts VALUES (17, 'GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP');
INSERT INTO history_accounts VALUES (18, 'GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C');
INSERT INTO history_accounts VALUES (19, 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG');
INSERT INTO history_accounts VALUES (20, 'GCB7FPYGLL6RJ37HKRAYW5TAWMFBGGFGM4IM6ERBCZXI2BZ4OOOX2UAY');
INSERT INTO history_accounts VALUES (21, 'GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO');
INSERT INTO history_accounts VALUES (22, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_accounts VALUES (23, 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB');
INSERT INTO history_accounts VALUES (24, 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB');
INSERT INTO history_accounts VALUES (25, 'GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK');


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_accounts_id_seq', 25, true);


--
-- Data for Name: history_assets; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_assets VALUES (1, 'credit_alphanum4', 'USD', 'GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD');
INSERT INTO history_assets VALUES (2, 'credit_alphanum4', 'USD', 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_assets VALUES (3, 'credit_alphanum4', 'EUR', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_assets VALUES (4, 'credit_alphanum4', 'EUR', 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_assets VALUES (5, 'credit_alphanum4', 'USD', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_assets VALUES (6, 'credit_alphanum4', 'USD', 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG');
INSERT INTO history_assets VALUES (7, 'native', '', '');


--
-- Name: history_assets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_assets_id_seq', 7, true);


--
-- Data for Name: history_effects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_effects VALUES (1, 249108107265, 1, 43, '{"new_seq": 300000000000}');
INSERT INTO history_effects VALUES (1, 244813139969, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 244813139969, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 244813139969, 3, 10, '{"weight": 1, "public_key": "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN"}');
INSERT INTO history_effects VALUES (3, 240518172673, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (3, 240518172673, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (3, 236223205377, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 236223205377, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (3, 236223205377, 3, 10, '{"weight": 1, "public_key": "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y"}');
INSERT INTO history_effects VALUES (4, 231928238081, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 231928238081, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 231928238082, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (4, 231928238082, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (4, 227633270785, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 227633270785, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (4, 227633270785, 3, 10, '{"weight": 1, "public_key": "GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X"}');
INSERT INTO history_effects VALUES (5, 223338303489, 1, 42, '{"name": "name1", "value": "MDAwMA=="}');
INSERT INTO history_effects VALUES (5, 219043336193, 1, 42, '{"name": "name1", "value": "MTIzNA=="}');
INSERT INTO history_effects VALUES (5, 214748368897, 1, 41, '{"name": "name2"}');
INSERT INTO history_effects VALUES (5, 210453401601, 1, 40, '{"name": "name1", "value": "MTIzNA=="}');
INSERT INTO history_effects VALUES (5, 210453405697, 1, 40, '{"name": "name2", "value": "NTY3OA=="}');
INSERT INTO history_effects VALUES (5, 210453409793, 1, 40, '{"name": "name ", "value": "aXRzIGdvdCBzcGFjZXMh"}');
INSERT INTO history_effects VALUES (5, 206158434305, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 206158434305, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (5, 206158434305, 3, 10, '{"weight": 1, "public_key": "GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD"}');
INSERT INTO history_effects VALUES (2, 201863467009, 1, 2, '{"amount": "15257676.9536092", "asset_type": "native"}');
INSERT INTO history_effects VALUES (6, 201863467009, 2, 2, '{"amount": "3814420.0001419", "asset_type": "native"}');
INSERT INTO history_effects VALUES (6, 197568499713, 1, 7, '{"inflation_destination": "GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS"}');
INSERT INTO history_effects VALUES (2, 197568503809, 1, 7, '{"inflation_destination": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (6, 193273532417, 1, 0, '{"starting_balance": "20000000000.0000000"}');
INSERT INTO history_effects VALUES (2, 193273532417, 2, 3, '{"amount": "20000000000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (6, 193273532417, 3, 10, '{"weight": 1, "public_key": "GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS"}');
INSERT INTO history_effects VALUES (7, 188978565121, 1, 3, '{"amount": "999.9999900", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 188978565121, 2, 2, '{"amount": "999.9999900", "asset_type": "native"}');
INSERT INTO history_effects VALUES (7, 188978565121, 3, 1, '{}');
INSERT INTO history_effects VALUES (7, 184683597825, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 184683597825, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (7, 184683597825, 3, 10, '{"weight": 1, "public_key": "GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ"}');
INSERT INTO history_effects VALUES (8, 180388630529, 1, 24, '{"trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (8, 176093663233, 1, 23, '{"trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (8, 176093667329, 1, 23, '{"trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (9, 171798695937, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (9, 171798700033, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (8, 167503728641, 1, 6, '{"auth_required_flag": true, "auth_revocable_flag": true}');
INSERT INTO history_effects VALUES (9, 163208761345, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 163208761345, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (9, 163208761345, 3, 10, '{"weight": 1, "public_key": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG"}');
INSERT INTO history_effects VALUES (8, 163208765441, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 163208765441, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (8, 163208765441, 3, 10, '{"weight": 1, "public_key": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}');
INSERT INTO history_effects VALUES (10, 158913794049, 1, 21, '{"limit": "0.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (10, 154618826753, 1, 22, '{"limit": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (10, 150323859457, 1, 22, '{"limit": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (10, 146028892161, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (10, 141733924865, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 141733924865, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (10, 141733924865, 3, 10, '{"weight": 1, "public_key": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG"}');
INSERT INTO history_effects VALUES (11, 137438957569, 1, 6, '{"auth_required_flag": false, "auth_revocable_flag": false}');
INSERT INTO history_effects VALUES (11, 137438961665, 1, 12, '{"weight": 2, "public_key": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES"}');
INSERT INTO history_effects VALUES (11, 137438961665, 2, 11, '{"public_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE"}');
INSERT INTO history_effects VALUES (11, 133143990273, 1, 12, '{"weight": 2, "public_key": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES"}');
INSERT INTO history_effects VALUES (11, 133143990273, 2, 12, '{"weight": 5, "public_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE"}');
INSERT INTO history_effects VALUES (11, 120259088385, 1, 7, '{"inflation_destination": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}');
INSERT INTO history_effects VALUES (11, 120259092481, 1, 6, '{"auth_required_flag": true}');
INSERT INTO history_effects VALUES (11, 120259096577, 1, 6, '{"auth_revocable_flag": true}');
INSERT INTO history_effects VALUES (11, 120259100673, 1, 12, '{"weight": 2, "public_key": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES"}');
INSERT INTO history_effects VALUES (11, 120259104769, 1, 4, '{"low_threshold": 0, "med_threshold": 2, "high_threshold": 2}');
INSERT INTO history_effects VALUES (11, 120259108865, 1, 5, '{"home_domain": "example.com"}');
INSERT INTO history_effects VALUES (11, 120259112961, 1, 12, '{"weight": 2, "public_key": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES"}');
INSERT INTO history_effects VALUES (11, 120259112961, 2, 10, '{"weight": 1, "public_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE"}');
INSERT INTO history_effects VALUES (11, 115964121089, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 115964121089, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (11, 115964121089, 3, 10, '{"weight": 1, "public_key": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES"}');
INSERT INTO history_effects VALUES (12, 107374186497, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 107374186497, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (12, 107374186497, 3, 10, '{"weight": 1, "public_key": "GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q"}');
INSERT INTO history_effects VALUES (14, 103079219201, 1, 33, '{"seller": "GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC", "offer_id": 3, "sold_amount": "20.0000000", "bought_amount": "20.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_type": "native", "sold_asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}');
INSERT INTO history_effects VALUES (13, 103079219201, 2, 33, '{"seller": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD", "offer_id": 3, "sold_amount": "20.0000000", "bought_amount": "20.0000000", "sold_asset_type": "native", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "bought_asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}');
INSERT INTO history_effects VALUES (13, 94489284609, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}');
INSERT INTO history_effects VALUES (13, 90194317313, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 90194317313, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (13, 90194317313, 3, 10, '{"weight": 1, "public_key": "GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC"}');
INSERT INTO history_effects VALUES (14, 90194321409, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 90194321409, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (14, 90194321409, 3, 10, '{"weight": 1, "public_key": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}');
INSERT INTO history_effects VALUES (17, 85899350017, 1, 2, '{"amount": "100.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 85899350017, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (16, 85899350017, 3, 33, '{"seller": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "offer_id": 2, "sold_amount": "100.0000000", "bought_amount": "100.0000000", "sold_asset_type": "native", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "bought_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (15, 85899350017, 4, 33, '{"seller": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "offer_id": 2, "sold_amount": "100.0000000", "bought_amount": "100.0000000", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_type": "native", "sold_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (17, 81604382721, 1, 2, '{"amount": "200.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 81604382721, 2, 3, '{"amount": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 81604382721, 3, 33, '{"seller": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "offer_id": 1, "sold_amount": "100.0000000", "bought_amount": "200.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_type": "native", "sold_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (15, 81604382721, 4, 33, '{"seller": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "offer_id": 1, "sold_amount": "200.0000000", "bought_amount": "100.0000000", "sold_asset_type": "native", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "bought_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 81604382721, 5, 33, '{"seller": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "offer_id": 2, "sold_amount": "200.0000000", "bought_amount": "200.0000000", "sold_asset_type": "native", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "bought_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (15, 81604382721, 6, 33, '{"seller": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "offer_id": 2, "sold_amount": "200.0000000", "bought_amount": "200.0000000", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_type": "native", "sold_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 77309415425, 1, 2, '{"amount": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (15, 77309415425, 2, 3, '{"amount": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 73014448129, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (17, 73014452225, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (16, 68719480833, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 68719480833, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (16, 68719480833, 3, 10, '{"weight": 1, "public_key": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD"}');
INSERT INTO history_effects VALUES (17, 68719484929, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 68719484929, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (17, 68719484929, 3, 10, '{"weight": 1, "public_key": "GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP"}');
INSERT INTO history_effects VALUES (15, 68719489025, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 68719489025, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (15, 68719489025, 3, 10, '{"weight": 1, "public_key": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}');
INSERT INTO history_effects VALUES (18, 64424513537, 1, 2, '{"amount": "10.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}');
INSERT INTO history_effects VALUES (19, 64424513537, 2, 3, '{"amount": "10.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}');
INSERT INTO history_effects VALUES (18, 60129546241, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (19, 60129546241, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (18, 60129550337, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}');
INSERT INTO history_effects VALUES (19, 55834578945, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 55834578945, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (19, 55834578945, 3, 10, '{"weight": 1, "public_key": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}');
INSERT INTO history_effects VALUES (18, 55834583041, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 55834583041, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (18, 55834583041, 3, 10, '{"weight": 1, "public_key": "GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C"}');
INSERT INTO history_effects VALUES (20, 51539611649, 1, 0, '{"starting_balance": "50.0000000"}');
INSERT INTO history_effects VALUES (21, 51539611649, 2, 3, '{"amount": "50.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (20, 51539611649, 3, 10, '{"weight": 1, "public_key": "GCB7FPYGLL6RJ37HKRAYW5TAWMFBGGFGM4IM6ERBCZXI2BZ4OOOX2UAY"}');
INSERT INTO history_effects VALUES (21, 47244644353, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 47244644353, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (21, 47244644353, 3, 10, '{"weight": 1, "public_key": "GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO"}');
INSERT INTO history_effects VALUES (2, 42949677057, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (22, 42949677057, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 42949677058, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (22, 42949677058, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (22, 38654709761, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 38654709761, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (22, 38654709761, 3, 10, '{"weight": 1, "public_key": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY"}');
INSERT INTO history_effects VALUES (2, 34359742465, 1, 2, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 34359742465, 2, 3, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 34359746561, 1, 2, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 34359746561, 2, 3, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 34359750657, 1, 2, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 34359750657, 2, 3, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 34359754753, 1, 2, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 34359754753, 2, 3, '{"amount": "1.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 30064775169, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 30064775169, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (23, 30064775169, 3, 10, '{"weight": 1, "public_key": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB"}');
INSERT INTO history_effects VALUES (24, 25769807873, 1, 12, '{"weight": 2, "public_key": "GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB"}');
INSERT INTO history_effects VALUES (24, 25769807873, 2, 12, '{"weight": 1, "public_key": "GD3E7HKMRNT6HGBGHBT6I6JE4N2S4W5KZ246TGJ4KQSXJ2P4BXCUPQMP"}');
INSERT INTO history_effects VALUES (24, 21474844673, 1, 12, '{"weight": 1, "public_key": "GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB"}');
INSERT INTO history_effects VALUES (24, 21474844673, 2, 10, '{"weight": 1, "public_key": "GD3E7HKMRNT6HGBGHBT6I6JE4N2S4W5KZ246TGJ4KQSXJ2P4BXCUPQMP"}');
INSERT INTO history_effects VALUES (24, 21474848769, 1, 4, '{"low_threshold": 2, "med_threshold": 2, "high_threshold": 2}');
INSERT INTO history_effects VALUES (24, 17179873281, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 17179873281, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (24, 17179873281, 3, 10, '{"weight": 1, "public_key": "GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB"}');
INSERT INTO history_effects VALUES (25, 12884905985, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (2, 12884905985, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (25, 12884905985, 3, 10, '{"weight": 1, "public_key": "GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_ledgers VALUES (61, '74c6dfd38c6bc4b10deddbe2fdd5a143a81765123e6eb6324dfef0dcd4a55fb1', 'fc843d2143c76bd29eb3b1ffddc3d48d60c7a1dcfa13ad88191da037092f12d0', 1, 1, '2019-06-03 16:36:27', '2019-06-03 16:35:33.313992', '2019-06-03 16:35:33.313992', 261993005056, 16, 1000190721000000000, 30471289, 100, 100000000, 1000000, 11, 'AAAAC/yEPSFDx2vSnrOx/93D1I1gx6Hc+hOtiBkdoDcJLxLQwnboDZyOgcyfb4ke5wI3SosVZ6YQuxInN2U6dBd49hYAAAAAXPVMiwAAAAAAAAAAosEhvRTCvadI/J9r+8kB+lYA6OphTjHH/5o2+iklI4P4rvPHAx3asi1o47MP12rNCQSIUVDj1Vrdu81+tOLDfwAAAD0N4WQpWNjKAAAAAAAB0PR5AAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (60, 'fc843d2143c76bd29eb3b1ffddc3d48d60c7a1dcfa13ad88191da037092f12d0', '85e72664cd18352262866f05355401362ba8cea6a2f47e2716471ab80d7a8208', 1, 1, '2019-06-03 16:36:26', '2019-06-03 16:35:33.335421', '2019-06-03 16:35:33.335421', 257698037760, 16, 1000190721000000000, 30471189, 100, 100000000, 1000000, 11, 'AAAAC4XnJmTNGDUiYoZvBTVUATYrqM6movR+JxZHGrgNeoII58fVvlrThcLsRhZ/MWhMMf/UYr2OzkRkTKFYZDGQ5ewAAAAAXPVMigAAAAAAAAAAZB2X7lcCb772lP7QltIwf7Af+ESisKsoFpU1cBZvSc7chbF3PheXP+z7+tsJ5b5GYsWEay6Di6CkkkadJmFk3AAAADwN4WQpWNjKAAAAAAAB0PQVAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (59, '85e72664cd18352262866f05355401362ba8cea6a2f47e2716471ab80d7a8208', '0e91662d61bb7b28faf575b9a350576fc40f16e382177310964394bb22de658b', 1, 1, '2019-06-03 16:36:25', '2019-06-03 16:35:33.346065', '2019-06-03 16:35:33.346066', 253403070464, 16, 1000190721000000000, 30471089, 100, 100000000, 1000000, 11, 'AAAACw6RZi1hu3so+vV1uaNQV2/EDxbjghdzEJZDlLsi3mWL+ICn3BLQXRs7fl4siWFC2ldDe+71dfZ+Ot6SzuRKboAAAAAAXPVMiQAAAAAAAAAAlNMwBDUEqtiETmU+aulkMrDuVtBffL8sVeqm0xm+yFztHfj5X8CLVoSeR2gp8ACoEyIeStN4O5rt7qSSaxxhuwAAADsN4WQpWNjKAAAAAAAB0POxAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (58, '0e91662d61bb7b28faf575b9a350576fc40f16e382177310964394bb22de658b', '4f6d9aa8f0c2a0bb76221321312b0449720e047bcc47a8d47e3885621396de32', 1, 1, '2019-06-03 16:36:24', '2019-06-03 16:35:33.357252', '2019-06-03 16:35:33.357252', 249108103168, 16, 1000190721000000000, 30470989, 100, 100000000, 1000000, 11, 'AAAAC09tmqjwwqC7diITITErBElyDgR7zEeo1H44hWITlt4y5laaexZi3BeoJbBYKEHVRaXQBzZzXDjrizmBZaG28WoAAAAAXPVMiAAAAAAAAAAAGMfQBF/YqYedEUNnQuilzukzLJPeVrR9I6q3baNG83n8kdWRPqRfd6skOTdJh3gw1CLyAZfl0Kspveak/XUBrAAAADoN4WQpWNjKAAAAAAAB0PNNAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (57, '4f6d9aa8f0c2a0bb76221321312b0449720e047bcc47a8d47e3885621396de32', '44943aac60f737822cbb8ed09dad3acff60f9811db9cc56efe2fd6082ed71fb0', 1, 1, '2019-06-03 16:36:23', '2019-06-03 16:35:33.368443', '2019-06-03 16:35:33.368444', 244813135872, 16, 1000190721000000000, 30470889, 100, 100000000, 1000000, 11, 'AAAAC0SUOqxg9zeCLLuO0J2tOs/2D5gR25zFbv4v1ggu1x+wSZEueVQpVrdX4o+KMu65Z6AqPEQg3hDgMGXVzWdfg4YAAAAAXPVMhwAAAAAAAAAA14sFIRujWwg9KpW3xcPTBAGeuB6CsnSMTpP8SjSsd9oj9l8Q/bO18WpeHxO1ThKTGbAtIm5RC4yqtX4KLaRH/QAAADkN4WQpWNjKAAAAAAAB0PLpAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (56, '44943aac60f737822cbb8ed09dad3acff60f9811db9cc56efe2fd6082ed71fb0', '9fe688714ba462bb2badb522a28f435678cc3ad383dc4afd91a8c11e35a3d735', 1, 1, '2019-06-03 16:36:22', '2019-06-03 16:35:33.380494', '2019-06-03 16:35:33.380495', 240518168576, 16, 1000190721000000000, 30470789, 100, 100000000, 1000000, 11, 'AAAAC5/miHFLpGK7K621IqKPQ1Z4zDrTg9xK/ZGowR41o9c1AfruE1INcN+aRvkAlMPc22dTgggif6wloMtQ47jJqRIAAAAAXPVMhgAAAAAAAAAAudvIKbdFQO6cnqBlfAQ0qwiv8T/bow/WPTIANo/ERREqQ4CwaeHZFm+0uGoZbmYdoBwAamF5v10EmnVf4uxH6gAAADgN4WQpWNjKAAAAAAAB0PKFAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (55, '9fe688714ba462bb2badb522a28f435678cc3ad383dc4afd91a8c11e35a3d735', 'b890cfd7326083ad579a151669212f51273c080e2b67a3664eb7cbcb9f177140', 1, 1, '2019-06-03 16:36:21', '2019-06-03 16:35:33.392461', '2019-06-03 16:35:33.392462', 236223201280, 16, 1000190721000000000, 30470689, 100, 100000000, 1000000, 11, 'AAAAC7iQz9cyYIOtV5oVFmkhL1EnPAgOK2ejZk63y8ufF3FAGEeGE9B3tl7o9CfKtPVM4DeWbUf/Bo39cDHTrEIancYAAAAAXPVMhQAAAAAAAAAAHOu0BR+25TDycY7dcNSyVck1nyZBCkABFZXwjkic52S5P22ZmU+gjPDIjhAWgZBWzZZ0wRh7yORr9nijpB3wWwAAADcN4WQpWNjKAAAAAAAB0PIhAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (54, 'b890cfd7326083ad579a151669212f51273c080e2b67a3664eb7cbcb9f177140', '48645ad8430e5952553cd69569a5711bac93461b1c9c030b60b062f4853b6b14', 1, 2, '2019-06-03 16:36:20', '2019-06-03 16:35:33.405812', '2019-06-03 16:35:33.405812', 231928233984, 16, 1000190721000000000, 30470589, 100, 100000000, 1000000, 11, 'AAAAC0hkWthDDllSVTzWlWmlcRusk0YbHJwDC2CwYvSFO2sUdAzpmZ4VmiJAGIhn1C8C65yAkmtBcCkFXHDwkTPoYkwAAAAAXPVMhAAAAAAAAAAAfXVfPS9jBqhTfz4LivYYP9btnu+qwf/qpVYL7HgbAM25OdLH4LXuwnN9DpuJOk8BEc4IicijRjArwYicV3k4xQAAADYN4WQpWNjKAAAAAAAB0PG9AAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (53, '48645ad8430e5952553cd69569a5711bac93461b1c9c030b60b062f4853b6b14', 'bd7b429dbd8f8efd1b56ed702bf183e473315e2f0422db5599548c404cee6066', 1, 1, '2019-06-03 16:36:19', '2019-06-03 16:35:33.420298', '2019-06-03 16:35:33.420298', 227633266688, 16, 1000190721000000000, 30470389, 100, 100000000, 1000000, 11, 'AAAAC717Qp29j479G1btcCvxg+RzMV4vBCLbVZlUjEBM7mBmDvRi34u5RgEADibpcO4c5OBbLI5HFWx1tz1OBwx7IMYAAAAAXPVMgwAAAAAAAAAAY80asxoHVId70xk95LhsRNki3uSYhgJrvMpz6Qe35wQ6M9y+p/ptNZ6zB2qi6pqwHFy5f/dMl34F1pxfwgvBdwAAADUN4WQpWNjKAAAAAAAB0PD1AAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (52, 'bd7b429dbd8f8efd1b56ed702bf183e473315e2f0422db5599548c404cee6066', '90534350f3fea0a6a559de207695e282e898909544823ac5c7000f76abadc08e', 1, 1, '2019-06-03 16:36:18', '2019-06-03 16:35:33.430701', '2019-06-03 16:35:33.430701', 223338299392, 16, 1000190721000000000, 30470289, 100, 100000000, 1000000, 11, 'AAAAC5BTQ1Dz/qCmpVneIHaV4oLomJCVRII6xccAD3arrcCOzeFEa2j6tGG8gbM/1MR8aABY/PxuO9Zz2s9mxGnLVoQAAAAAXPVMggAAAAAAAAAAdFs0uvdku2e14ZLZNY/inkD0IgG9m+9h08GbknZs9EzZzupAXizfBXmh99nHTu7W5H1zvpkWS2CF1tqBEKwuCQAAADQN4WQpWNjKAAAAAAAB0PCRAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (51, '90534350f3fea0a6a559de207695e282e898909544823ac5c7000f76abadc08e', 'b1c7d71aeba870b162b022c31adeb26614e63070928647a1b9c6bc51aba83d0d', 1, 1, '2019-06-03 16:36:17', '2019-06-03 16:35:33.443509', '2019-06-03 16:35:33.443509', 219043332096, 16, 1000190721000000000, 30470189, 100, 100000000, 1000000, 11, 'AAAAC7HH1xrrqHCxYrAiwxresmYU5jBwkoZHobnGvFGrqD0Nt2vJTHsKyMHdVCH/gs0c3U/eQKqaWReE+4jy2gt+ukwAAAAAXPVMgQAAAAAAAAAAG9Lveo8KOMBsuTXxxarafYQzHKlqwhXWtcRjpmTOc+OB+peuNPsYlfeIN/BQ3+xOQZ7d/bV/tnJgLEpXfnSm4QAAADMN4WQpWNjKAAAAAAAB0PAtAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (50, 'b1c7d71aeba870b162b022c31adeb26614e63070928647a1b9c6bc51aba83d0d', '0dc8610ec649e718d45083b69105153b5689bec67427a72058b9357eeb0aceac', 1, 1, '2019-06-03 16:36:16', '2019-06-03 16:35:33.457157', '2019-06-03 16:35:33.457157', 214748364800, 16, 1000190721000000000, 30470089, 100, 100000000, 1000000, 11, 'AAAACw3IYQ7GSecY1FCDtpEFFTtWib7GdCenIFi5NX7rCs6si0bTsBbSdDOFoh+y5f8cl3t+mLZfzrhHG/6bxBlR1MwAAAAAXPVMgAAAAAAAAAAAe2zMUWV2ofxxSJUd0hgEpW+dk+M503UNhsqXLKL0TMgrrap+MZcPKP7WldMsB6zr/892IcwiXmAc6RON5c7/DQAAADIN4WQpWNjKAAAAAAAB0O/JAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAK62qfjGXDyj+1pXTLAes6//PdiHMIl5gHOkTjeXO/w0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (49, '0dc8610ec649e718d45083b69105153b5689bec67427a72058b9357eeb0aceac', '41b518eb3573162dba39489d5e79c1c04ee6c6d08c19311b5bf7d9b8f9aa17d6', 3, 3, '2019-06-03 16:36:15', '2019-06-03 16:35:33.472035', '2019-06-03 16:35:33.472035', 210453397504, 16, 1000190721000000000, 30469989, 100, 100000000, 1000000, 11, 'AAAAC0G1GOs1cxYtujlInV55wcBO5sbQjBkxG1v32bj5qhfWeVRF+FiW1Qlc+3u1mBcZ84ULRdo3+cTo45YRJDO+PS4AAAAAXPVMfwAAAAAAAAAAsSOLB+Fh7PI1BdmA+nlx5DdeXSAZY7OnyHcuyeMEmJA3bk41wtryM/3lghnAAVnP/iB46990qItDfm+tlOWNaQAAADEN4WQpWNjKAAAAAAAB0O9lAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0, NULL);
INSERT INTO history_ledgers VALUES (48, '41b518eb3573162dba39489d5e79c1c04ee6c6d08c19311b5bf7d9b8f9aa17d6', 'b49ea0eb1fe629fc2c9df2c976bc9971fed6fd8d257ece738e45c3f8967d1569', 1, 1, '2019-06-03 16:36:14', '2019-06-03 16:35:33.482429', '2019-06-03 16:35:33.48243', 206158430208, 16, 1000190721000000000, 30469689, 100, 100000000, 1000000, 11, 'AAAAC7SeoOsf5in8LJ3yyXa8mXH+1v2NJX7Oc45Fw/iWfRVp6aDW3LPbWQnkQqAyNZZJlDWaVAL9GmVKlMVbHoaATywAAAAAXPVMfgAAAAAAAAAAQNg4cr4Do7quX33o4rbaXa7+BDdHc8PbTEXBcuQETP8rTgcb7n3UXEG0av8iuvhxWlCZJf64MCH+txYgQXKWPAAAADAN4WQpWNjKAAAAAAAB0O45AAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (47, 'b49ea0eb1fe629fc2c9df2c976bc9971fed6fd8d257ece738e45c3f8967d1569', '706d948a65a01559c33264168e913f4b490034942f79aea96f616a03e98297c3', 1, 1, '2019-06-03 16:36:13', '2019-06-03 16:35:33.495168', '2019-06-03 16:35:33.495169', 201863462912, 16, 1000190721000000000, 30469589, 100, 100000000, 1000000, 11, 'AAAAC3BtlIploBVZwzJkFo6RP0tJADSUL3muqW9hagPpgpfD1EV61P5oEG1jpv4Ia/HUgGCpKwTQ4ZplHht7o7pA71gAAAAAXPVMfQAAAAAAAAAA8OM7CNjloxWqwlaanmBarRzmaEAjS1RTT7DMBUh/OHivMXhi18hGxAY77i6ACTEqawD8vXkvXLjBzeeYDgCPTQAAAC8N4WQpWNjKAAAAAAAB0O3VAAAAAQAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (46, '706d948a65a01559c33264168e913f4b490034942f79aea96f616a03e98297c3', '69840e81fe7e584aad0f22dfa89ee767cd3524497563e471af8af67b53f8d471', 2, 2, '2019-06-03 16:36:12', '2019-06-03 16:35:33.506546', '2019-06-03 16:35:33.506546', 197568495616, 16, 1000000000000000000, 7000, 100, 100000000, 1000000, 11, 'AAAAC2mEDoH+flhKrQ8i36ie52fNNSRJdWPkca+K9ntT+NRxAgTIPLGOynuQgOXzmXQZrLyGTAGIDWOHJRQJqUfbAiUAAAAAXPVMfAAAAAAAAAAAhkYgTYVFfiNxywZFyKUj6rUeR42SRzNzOnasgV/WuDvtwNkabdQ+tgZe43HtXEsCbQ9gWOuRX6L5cvO8ZwBRpQAAAC4N4Lazp2QAAAAAAAAAABtYAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (45, '69840e81fe7e584aad0f22dfa89ee767cd3524497563e471af8af67b53f8d471', '5eef1067a0fbdb41321e472a1c14fedc553d13539516dea4fdd319d0ca7d2af8', 1, 1, '2019-06-03 16:36:11', '2019-06-03 16:35:33.518547', '2019-06-03 16:35:33.518547', 193273528320, 16, 1000000000000000000, 6800, 100, 100000000, 1000000, 11, 'AAAAC17vEGeg+9tBMh5HKhwU/txVPRNTlRbepP3TGdDKfSr4kGCHHnA6KOdoPcSvp6Rxhk2hURlGF1VP2oSVNBkDI9cAAAAAXPVMewAAAAAAAAAA2s+Yd36GUBmFrCn8nxzR3G7G6BlWwDjWIvSuQxGMICY92RjX5DbOxYhJMM/EQc1oqGw0gsFbtqsOXBjLLV9EXAAAAC0N4Lazp2QAAAAAAAAAABqQAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (44, '5eef1067a0fbdb41321e472a1c14fedc553d13539516dea4fdd319d0ca7d2af8', '4409a871d550650b6a0a630136e80cd259a333ca55203db1548d315f19ef4a5c', 1, 1, '2019-06-03 16:36:10', '2019-06-03 16:35:33.528911', '2019-06-03 16:35:33.528911', 188978561024, 16, 1000000000000000000, 6700, 100, 100000000, 1000000, 11, 'AAAAC0QJqHHVUGULagpjATboDNJZozPKVSA9sVSNMV8Z70pcrowUdKRzA/ZRkdeo8Gu2qV4ThU16d9p8N74Hanu5UssAAAAAXPVMegAAAAAAAAAAWWc1pYhVWRoEJa6GcHIsaj1ysxI4CQ+BcJqiGW3wjmJ+vfFRLi5J2iDMXNTWqNmnft8zHFO464YgJZT7ub0qrQAAACwN4Lazp2QAAAAAAAAAABosAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (43, '4409a871d550650b6a0a630136e80cd259a333ca55203db1548d315f19ef4a5c', 'a96cbf9f3c2d3c56bd14b4050246a5e39e5920a201e9d66c40bc5aa68782fe53', 1, 1, '2019-06-03 16:36:09', '2019-06-03 16:35:33.541246', '2019-06-03 16:35:33.541246', 184683593728, 16, 1000000000000000000, 6600, 100, 100000000, 1000000, 11, 'AAAAC6lsv588LTxWvRS0BQJGpeOeWSCiAenWbEC8WqaHgv5TZNSyhXDfKEyU6jk0m8nxAE6PCJaGxAoOwqvcASMsUS4AAAAAXPVMeQAAAAAAAAAA9L2wqBam8xx6EMlQ1e/R7JPhksJqFoZDmKFcFBwf3xCwc2JLloIzptYhoNC7H0OKDeD9owZAk5Iannrejr0OAwAAACsN4Lazp2QAAAAAAAAAABnIAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (42, 'a96cbf9f3c2d3c56bd14b4050246a5e39e5920a201e9d66c40bc5aa68782fe53', '22b06b409446c1f0c278af0ce15ad99657a0ff2196a5bfeded56d69d36779c9d', 1, 1, '2019-06-03 16:36:08', '2019-06-03 16:35:33.552813', '2019-06-03 16:35:33.552813', 180388626432, 16, 1000000000000000000, 6500, 100, 100000000, 1000000, 11, 'AAAACyKwa0CURsHwwnivDOFa2ZZXoP8hlqW/7e1W1p02d5ydI1zHat+THSiG8FiJYhX2IjqcewMyqp7eILRtLxlqAV0AAAAAXPVMeAAAAAAAAAAA4mQxuIyGkMnUNprxo1u6b3NOKgS/CqRR0k9oMQPLWs8+V68qnCaFFEx4K1goWa3RDH8IZsY7i2BXWSCMWqzpxgAAACoN4Lazp2QAAAAAAAAAABlkAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (41, '22b06b409446c1f0c278af0ce15ad99657a0ff2196a5bfeded56d69d36779c9d', 'a6f5bc74da946970949cc8a58b7fc0ca52a9627fa1a45fccf567a6ad2f311a90', 2, 2, '2019-06-03 16:36:07', '2019-06-03 16:35:33.564965', '2019-06-03 16:35:33.564965', 176093659136, 16, 1000000000000000000, 6400, 100, 100000000, 1000000, 11, 'AAAAC6b1vHTalGlwlJzIpYt/wMpSqWJ/oaRfzPVnpq0vMRqQKNyYbOwL7RbiajIIn9wo1XIqxrxKXhScSJZushgn1U8AAAAAXPVMdwAAAAAAAAAA3texHWkeXlsh8yl52JzE1M7p/8uzk36DdlbvwdqWw1Y07THsoEsB7/w8RQjJ40lpM5K1EboCeT8y/j///j3AagAAACkN4Lazp2QAAAAAAAAAABkAAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (40, 'a6f5bc74da946970949cc8a58b7fc0ca52a9627fa1a45fccf567a6ad2f311a90', 'bc3a5a953848fe2cb46ed75ecebe6ba09c856d47d6fea025223858d8dc7aad95', 2, 2, '2019-06-03 16:36:06', '2019-06-03 16:35:33.580316', '2019-06-03 16:35:33.580316', 171798691840, 16, 1000000000000000000, 6200, 100, 100000000, 1000000, 11, 'AAAAC7w6WpU4SP4stG7XXs6+a6CchW1H1v6gJSI4WNjceq2V/inIeGJFrV1RVU+xprIZRH8KVaI74l/TAt+/rJeHJ50AAAAAXPVMdgAAAAAAAAAAWJlo97rfQjjINuAEkeCOb3Rwl6NYUk95Xa63Nc0UCv10U3spNQC8m5QVal2MW27Qo8OwBAnHuh8wEWKMhYtIgAAAACgN4Lazp2QAAAAAAAAAABg4AAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (39, 'bc3a5a953848fe2cb46ed75ecebe6ba09c856d47d6fea025223858d8dc7aad95', 'b8296a81cb3c468c47f100662acb337e168861732ce211ddda30218098e806ce', 1, 1, '2019-06-03 16:36:05', '2019-06-03 16:35:33.592749', '2019-06-03 16:35:33.592749', 167503724544, 16, 1000000000000000000, 6000, 100, 100000000, 1000000, 11, 'AAAAC7gpaoHLPEaMR/EAZirLM34WiGFzLOIR3dowIYCY6AbOHivLxyP8I+6k4JpB4A3PxOKsREZUMnzYBoAng8ZuENAAAAAAXPVMdQAAAAAAAAAASeHBruDgNPJb9kfcJVMyCg6meb0kov5dW9Ok8EEep2bE1Tj7RqD1/0Q6+H0uMkMVOwWWrTFNf3Y8cVz9ojCJ1QAAACcN4Lazp2QAAAAAAAAAABdwAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (38, 'b8296a81cb3c468c47f100662acb337e168861732ce211ddda30218098e806ce', 'f39961e7de53003d9a878277343091c60e1eadefba33202be08334cae6787bef', 2, 2, '2019-06-03 16:36:04', '2019-06-03 16:35:33.60356', '2019-06-03 16:35:33.60356', 163208757248, 16, 1000000000000000000, 5900, 100, 100000000, 1000000, 11, 'AAAAC/OZYefeUwA9moeCdzQwkcYOHq3vujMgK+CDNMrmeHvvvdOMXcZQv2fCMQnVMUfKBJ7XSFWbgB/wE3Bs7er87uAAAAAAXPVMdAAAAAAAAAAAN/+dCsXSVcipJLX5fSVl/dtD8pFe+iZnfnNZ2g7vXj/Ebip0yLvWR8Z7upcPNrqDXeLvVxZizB40hi2kIihnngAAACYN4Lazp2QAAAAAAAAAABcMAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (37, 'f39961e7de53003d9a878277343091c60e1eadefba33202be08334cae6787bef', 'f965f5e76e2df794975e7245831b95b86ecc9711f7b7fb144df9bcac3c8a7a2a', 1, 1, '2019-06-03 16:36:03', '2019-06-03 16:35:33.615591', '2019-06-03 16:35:33.615591', 158913789952, 16, 1000000000000000000, 5700, 100, 100000000, 1000000, 11, 'AAAAC/ll9eduLfeUl15yRYMblbhuzJcR97f7FE35vKw8inoqjeXfs5RaC/pSbMQ0UkNgOHq6KFpFmq5/xPhha/qIrikAAAAAXPVMcwAAAAAAAAAAPORbuvHNNFfsSSmddqmUP8Fk8iW9Jb1yQL/uEyYGc2uIYodbNQGTGAsLT3wVFszuQByXNpPK2sQ80gKgAo6R4QAAACUN4Lazp2QAAAAAAAAAABZEAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (36, 'f965f5e76e2df794975e7245831b95b86ecc9711f7b7fb144df9bcac3c8a7a2a', 'f090655f0d5cee583227d0088cb7af2b780f6fd91e836a4117abb11d479b5ce1', 1, 1, '2019-06-03 16:36:02', '2019-06-03 16:35:33.635589', '2019-06-03 16:35:33.63559', 154618822656, 16, 1000000000000000000, 5600, 100, 100000000, 1000000, 11, 'AAAAC/CQZV8NXO5YMifQCIy3ryt4D2/ZHoNqQRersR1Hm1zht40cpv6aCNmZXclj0Gp1g8/QTb+rEsz7O6QRevEp7pkAAAAAXPVMcgAAAAAAAAAADx7u3Tt890DDCAY2PUo3lKiU3M3nylcrqmvMlFLS+fEnxfEWjlbDeIddDXeCKT79jo4fwyV0ymPAz6/1A94F9AAAACQN4Lazp2QAAAAAAAAAABXgAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (35, 'f090655f0d5cee583227d0088cb7af2b780f6fd91e836a4117abb11d479b5ce1', 'e0905bd61697815e7dd80a6f7ad8e83077a5304ab0681825e5a820bfbcc12dbd', 1, 1, '2019-06-03 16:36:01', '2019-06-03 16:35:33.645127', '2019-06-03 16:35:33.645127', 150323855360, 16, 1000000000000000000, 5500, 100, 100000000, 1000000, 11, 'AAAAC+CQW9YWl4FefdgKb3rY6DB3pTBKsGgYJeWoIL+8wS29TrWUd6O7M4S0fe8UEcdhlcmYcyd0+obgadQjYZA/UK4AAAAAXPVMcQAAAAAAAAAAU3lWP9fTSJmRnN7PpNBaZaYziGh9Xk2I2lyAgYUdNswxQTtINH5UegpdO02kZozOfs6/jYRFpxCpYgETz6Se7gAAACMN4Lazp2QAAAAAAAAAABV8AAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (34, 'e0905bd61697815e7dd80a6f7ad8e83077a5304ab0681825e5a820bfbcc12dbd', '5273d226f3ad944db78e740d2ef1736845197a098f8c22c2a1a412838d2f2770', 1, 1, '2019-06-03 16:36:00', '2019-06-03 16:35:33.656811', '2019-06-03 16:35:33.656811', 146028888064, 16, 1000000000000000000, 5400, 100, 100000000, 1000000, 11, 'AAAAC1Jz0ibzrZRNt450DS7xc2hFGXoJj4wiwqGkEoONLydwytNJJuK4uRpBqNE6DA/twCtsyMwX/BNjTK2j1A2cLpIAAAAAXPVMcAAAAAAAAAAA/bE+PeQRbcMvSIfQyaXSWXEgzB7abwreNVQeYYIhwiZ0ZqxccyydwKiguVJdbIGZEGSMZ99UQLW9RIyPocciTAAAACIN4Lazp2QAAAAAAAAAABUYAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (33, '5273d226f3ad944db78e740d2ef1736845197a098f8c22c2a1a412838d2f2770', '05e77cb96060cb1b543745f88fe36e3246b45b26a9ee597f194fdfa332915ac6', 1, 1, '2019-06-03 16:35:59', '2019-06-03 16:35:33.665558', '2019-06-03 16:35:33.665558', 141733920768, 16, 1000000000000000000, 5300, 100, 100000000, 1000000, 11, 'AAAACwXnfLlgYMsbVDdF+I/jbjJGtFsmqe5ZfxlP36MykVrGueMGwKUy3zdCJEbcnLnxoZEU3V59aYr7+qbfMi0DwngAAAAAXPVMbwAAAAAAAAAAHdP8tBBtpeyL899tLPnNdaqUoHRWmHKjK8mobr412wEIdjVcHb3VJhQAtunNPjhphkSfh+hk4kcBwk4tfP9i6wAAACEN4Lazp2QAAAAAAAAAABS0AAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (32, '05e77cb96060cb1b543745f88fe36e3246b45b26a9ee597f194fdfa332915ac6', 'fe13366c4c53c04a65ec45e5dd7540d332a8c38456f71a46b8c197a8f61614b3', 2, 2, '2019-06-03 16:35:58', '2019-06-03 16:35:33.677273', '2019-06-03 16:35:33.677273', 137438953472, 16, 1000000000000000000, 5200, 100, 100000000, 1000000, 11, 'AAAAC/4TNmxMU8BKZexF5d11QNMyqMOEVvcaRrjBl6j2FhSzbdxxbfFDV+si6AGRvL5s+oVcsY8ZwkF5kL9EBypi7k8AAAAAXPVMbgAAAAAAAAAA2zD57d2oUw2L6murtm/oe2Dj6FuzBXnVm7l7HoDsMUzNqYbZejWmvwR44FVQPbMIdu2bjH6kX3IJ7Chi4YmgugAAACAN4Lazp2QAAAAAAAAAABRQAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (31, 'fe13366c4c53c04a65ec45e5dd7540d332a8c38456f71a46b8c197a8f61614b3', '750fedf881913e508427d7c6f1a31f032cffd8b513eab2fca3c482743c92e16c', 1, 1, '2019-06-03 16:35:57', '2019-06-03 16:35:33.694173', '2019-06-03 16:35:33.694173', 133143986176, 16, 1000000000000000000, 5000, 100, 100000000, 1000000, 11, 'AAAAC3UP7fiBkT5QhCfXxvGjHwMs/9i1E+qy/KPEgnQ8kuFscIqVmjyvqIlbFdmmCscTHwuF1xf/128KmopBLMQWNqwAAAAAXPVMbQAAAAAAAAAA8GmZ2mhGWPCcWAd1pU9SBJXfejtD7FOEoLD6oFsTOv4MKWINvUzkDiENnzMAs1lvH8EiEWJaRwrJRvN0YuRvzAAAAB8N4Lazp2QAAAAAAAAAABOIAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (30, '750fedf881913e508427d7c6f1a31f032cffd8b513eab2fca3c482743c92e16c', '936c1e1f0bf1af96380eb09e59e521d5e3ddd370d21ff4fd1a09ab123a4cfbaa', 1, 1, '2019-06-03 16:35:56', '2019-06-03 16:35:33.708481', '2019-06-03 16:35:33.708481', 128849018880, 16, 1000000000000000000, 4900, 100, 100000000, 1000000, 11, 'AAAAC5NsHh8L8a+WOA6wnlnlIdXj3dNw0h/0/RoJqxI6TPuqkzOXz5mf2LgYnaoAtXxNZLk+UfT8YR92Qb9Z7YPbL3MAAAAAXPVMbAAAAAAAAAAA1jXLaoK/IOqKNVkib2cNLrjlETHXpZlcJ4BFYO8bM3mRzqpWfgvx7nf587dxjJYVNmgRlffJEzWDBgwhfJmgUwAAAB4N4Lazp2QAAAAAAAAAABMkAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (29, '936c1e1f0bf1af96380eb09e59e521d5e3ddd370d21ff4fd1a09ab123a4cfbaa', 'c14459522fe42e2d10b29a370bbeba513e4aee429282dbe44842835d1b443663', 1, 1, '2019-06-03 16:35:55', '2019-06-03 16:35:33.718442', '2019-06-03 16:35:33.718442', 124554051584, 16, 1000000000000000000, 4800, 100, 100000000, 1000000, 11, 'AAAAC8FEWVIv5C4tELKaNwu+ulE+Su5CkoLb5EhCg10bRDZjMHDsmJoTpyfMr0n7ZRj1PtVYciwS5F25/ebV5KmPHyIAAAAAXPVMawAAAAAAAAAAbaCD4RbfEqSrI+qpIyzzUU+bYYMJsRw3pHcmHwQrvaAiHUzESb7mxWgb6rD9i+bvYadjhxMRG/exoDWp8/oUMwAAAB0N4Lazp2QAAAAAAAAAABLAAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (28, 'c14459522fe42e2d10b29a370bbeba513e4aee429282dbe44842835d1b443663', '902b0f5e1ee506b9f9d688e464673bda2faa056876cdfe0ac250ef5c9df0a8a6', 7, 7, '2019-06-03 16:35:54', '2019-06-03 16:35:33.729461', '2019-06-03 16:35:33.729461', 120259084288, 16, 1000000000000000000, 4700, 100, 100000000, 1000000, 11, 'AAAAC5ArD14e5Qa5+daI5GRnO9ovqgVods3+CsJQ71yd8KimnNDwLYHe7gIHYMyxYiWPL6AdpOmOGCRINmPcEvmBn5oAAAAAXPVMagAAAAAAAAAAdQ/T4yJNKn9hiJoTuxDsNbd9ibRWf2hNHJhsiqBKWD1btlsiWANxqEZ3UBK2IsuYFM4a68rORQwudeZsBK1uVQAAABwN4Lazp2QAAAAAAAAAABJcAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 7, 0, NULL);
INSERT INTO history_ledgers VALUES (27, '902b0f5e1ee506b9f9d688e464673bda2faa056876cdfe0ac250ef5c9df0a8a6', '925e48fe8c319e23b85014ddf6e3698dbaeb15bd84e0e4435614f3125d847eaf', 1, 1, '2019-06-03 16:35:53', '2019-06-03 16:35:33.756347', '2019-06-03 16:35:33.756348', 115964116992, 16, 1000000000000000000, 4000, 100, 100000000, 1000000, 11, 'AAAAC5JeSP6MMZ4juFAU3fbjaY266xW9hODkQ1YU8xJdhH6vQydJIXEAs1uOD6IkUMs8GlnyPDTxRqoRKW4/wgZZaNwAAAAAXPVMaQAAAAAAAAAAkpeXdNXfmh1EnUjCYIMebdzV8PlgCe6J6eJq8Mb4SXYTMU9R5pbFQB/phEpxiGicTklBEXUSQW8XrtV0H9BN+AAAABsN4Lazp2QAAAAAAAAAAA+gAAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (26, '925e48fe8c319e23b85014ddf6e3698dbaeb15bd84e0e4435614f3125d847eaf', 'f54c0097d441ff76c487d4ebb18f5e5f52c5a51ca96800bf9b818295790a3778', 2, 2, '2019-06-03 16:35:52', '2019-06-03 16:35:33.767845', '2019-06-03 16:35:33.767845', 111669149696, 16, 1000000000000000000, 3900, 100, 100000000, 1000000, 11, 'AAAAC/VMAJfUQf92xIfU67GPXl9SxaUcqWgAv5uBgpV5Cjd4ObBqces8zXZAjfLGaR+Md0gHwPV4L+NudPgSQJdsE58AAAAAXPVMaAAAAAAAAAAAGTKcjd5tsclF3P+ptDiMdUWp83tuYF+zlADp+X+gt+eBl1mZgQ2kBgOlDC3yM4k1XFSkg8nteVJMp/KhqOyW2wAAABoN4Lazp2QAAAAAAAAAAA88AAAAAAAAAAAAAAAGAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (25, 'f54c0097d441ff76c487d4ebb18f5e5f52c5a51ca96800bf9b818295790a3778', 'd71b122bfb10e1193597acf4769f99b716e79e7d1eaccf225adf5fd65cbeee37', 1, 1, '2019-06-03 16:35:51', '2019-06-03 16:35:33.779704', '2019-06-03 16:35:33.779704', 107374182400, 16, 1000000000000000000, 3700, 100, 100000000, 1000000, 11, 'AAAAC9cbEiv7EOEZNZes9HafmbcW5559HqzPIlrfX9Zcvu43C5djYweCkfubeYiyrAvxlPmuUHgxbk+pkX450lbiUBIAAAAAXPVMZwAAAAAAAAAAM9gLVDvjwMuVZ+njRjvvMxlhodRg/VJpuapiCb0y9SDsWIE0ZZXCJ72Msa/6UrUdTDk+mqNMq2KeZVKBjdtMbwAAABkN4Lazp2QAAAAAAAAAAA50AAAAAAAAAAAAAAAEAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (24, 'd71b122bfb10e1193597acf4769f99b716e79e7d1eaccf225adf5fd65cbeee37', 'd09c2755ea31109767d6ad6ad7155cff1b630edb6368a4b280bb68de24c2ebb4', 1, 1, '2019-06-03 16:35:50', '2019-06-03 16:35:33.791139', '2019-06-03 16:35:33.791139', 103079215104, 16, 1000000000000000000, 3600, 100, 100000000, 1000000, 11, 'AAAAC9CcJ1XqMRCXZ9atatcVXP8bYw7bY2iksoC7aN4kwuu03hM6c6SywARSnPb3ctF6foVdbgyxhp/R4/s+WE+te5AAAAAAXPVMZgAAAAAAAAAA1ziHajubDrj1Iu5EQ8YZB24Czm29AJPdaJnsX3+A+fV3KfdeiqJqyyhfk+7mwWUB0rbbvt6y640hL/bnoPMZ9gAAABgN4Lazp2QAAAAAAAAAAA4QAAAAAAAAAAAAAAAEAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (23, 'd09c2755ea31109767d6ad6ad7155cff1b630edb6368a4b280bb68de24c2ebb4', 'b8d81a94cc9ac1dcd4975e5bae0680a746f41ab923de9b68650422039f8df47a', 1, 1, '2019-06-03 16:35:49', '2019-06-03 16:35:33.815281', '2019-06-03 16:35:33.815282', 98784247808, 16, 1000000000000000000, 3500, 100, 100000000, 1000000, 11, 'AAAAC7jYGpTMmsHc1JdeW64GgKdG9Bq5I96baGUEIgOfjfR6OPg5AoiFUgt3GG0WOJr4C3JaQHm7MPBQisdC4hkzxS8AAAAAXPVMZQAAAAAAAAAAeu6ZLz++eOyW3/HPY616v+oF2X756WCT8G93ApdHifosUOpCb+sObcpszPtP9+Hy0MHmzxhpzAP/sTATcQpnKgAAABcN4Lazp2QAAAAAAAAAAA2sAAAAAAAAAAAAAAADAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (22, 'b8d81a94cc9ac1dcd4975e5bae0680a746f41ab923de9b68650422039f8df47a', '1e73266b9c1d6ed85ed20270b1c9c4081445d5f201189bebb5721c9b2565f749', 1, 1, '2019-06-03 16:35:48', '2019-06-03 16:35:33.827411', '2019-06-03 16:35:33.827411', 94489280512, 16, 1000000000000000000, 3400, 100, 100000000, 1000000, 11, 'AAAACx5zJmucHW7YXtICcLHJxAgURdXyARib67VyHJslZfdJ4DzzvVfPje5vC6cghCFfrZwMHMEDMF1phwp/kIju/78AAAAAXPVMZAAAAAAAAAAA0AMiyT4CyDLLvayS0AmP97hen9j1rkIsoYeN0iQrjup0aPFduVIsJ+nfz9INCJSr58NOPxPNQTKZqTTYzS5+OwAAABYN4Lazp2QAAAAAAAAAAA1IAAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (21, '1e73266b9c1d6ed85ed20270b1c9c4081445d5f201189bebb5721c9b2565f749', 'e0c98a06e3dc47183a4803f947cd00afe37d8b94dc34de429add40fa7c19c618', 2, 2, '2019-06-03 16:35:47', '2019-06-03 16:35:33.838148', '2019-06-03 16:35:33.838148', 90194313216, 16, 1000000000000000000, 3300, 100, 100000000, 1000000, 11, 'AAAAC+DJigbj3EcYOkgD+UfNAK/jfYuU3DTeQprdQPp8GcYYFo7+Glas2jipoGla9vrLLK0ooAiMI5GQeOkLV/uj9cIAAAAAXPVMYwAAAAAAAAAA3LAnNjSwmLo474ojoHgysfkXmeeBIXdgWbyOJwbnsx+QUdaKPfAXYhyvkzKQKE+TYbh3dMj8cX0gtf0+OoTtuwAAABUN4Lazp2QAAAAAAAAAAAzkAAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (20, 'e0c98a06e3dc47183a4803f947cd00afe37d8b94dc34de429add40fa7c19c618', 'ce1cd7aab0e22a1e81c89bd171d0bda5d8dec28824143025f6f0d15a1d7acc18', 1, 1, '2019-06-03 16:35:46', '2019-06-03 16:35:33.847808', '2019-06-03 16:35:33.847808', 85899345920, 16, 1000000000000000000, 3100, 100, 100000000, 1000000, 11, 'AAAAC84c16qw4ioegcib0XHQvaXY3sKIJBQwJfbw0VodeswY+fowsLcKz7J1dE35gKnYRMa2zA3cqfuoNUM1feiZPu8AAAAAXPVMYgAAAAAAAAAAx/hxaEW2XNJB3iDf23eYfGcwcQ862pSt7OJ6FF108Ia/Chyidqfse/a37FifLAlDK+EuQ3xGyaqCS5MLkWQUXgAAABQN4Lazp2QAAAAAAAAAAAwcAAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (19, 'ce1cd7aab0e22a1e81c89bd171d0bda5d8dec28824143025f6f0d15a1d7acc18', 'e6f54ea5dcfa926f5e46dd056690fb031dc392717dc0fbac0b2f411e84965210', 1, 1, '2019-06-03 16:35:45', '2019-06-03 16:35:33.86558', '2019-06-03 16:35:33.86558', 81604378624, 16, 1000000000000000000, 3000, 100, 100000000, 1000000, 11, 'AAAAC+b1TqXc+pJvXkbdBWaQ+wMdw5JxfcD7rAsvQR6EllIQaY9SIsx+2mzlbk+f9SEovenpLo/LJyUPh6A+g+VF8p4AAAAAXPVMYQAAAAAAAAAAYzC64fmWjZOtqE4JjsKqwerdnzN6ywnJSvdsDWFPcu3pJFtMmHmps2cAxoQk/qzf8QoNOBiBIGAE68eH9xRyaQAAABMN4Lazp2QAAAAAAAAAAAu4AAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (18, 'e6f54ea5dcfa926f5e46dd056690fb031dc392717dc0fbac0b2f411e84965210', '8aff00265fb98412902df9d5c30e194a01adbd792490ef58836478db9d8e3156', 3, 3, '2019-06-03 16:35:44', '2019-06-03 16:35:33.883905', '2019-06-03 16:35:33.883906', 77309411328, 16, 1000000000000000000, 2900, 100, 100000000, 1000000, 11, 'AAAAC4r/ACZfuYQSkC351cMOGUoBrb15JJDvWINkeNudjjFWp7Hi+kDBELcUYlO/VeJpvL6v5oIQgjfBH9t2VsYU0vYAAAAAXPVMYAAAAAAAAAAANKxyAMoyM/gn59hBYzLxbU28Wbix3QHy7s086vpSY4dsfsh10yperbmIfKUawooPYNqN/QB2jK4qu2tAX+tIawAAABIN4Lazp2QAAAAAAAAAAAtUAAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0, NULL);
INSERT INTO history_ledgers VALUES (17, '8aff00265fb98412902df9d5c30e194a01adbd792490ef58836478db9d8e3156', '6f4b7c957a42a32dd830d43212d5bbafeb9221538e4010a423455330ddb09b2b', 2, 2, '2019-06-03 16:35:43', '2019-06-03 16:35:33.897912', '2019-06-03 16:35:33.897912', 73014444032, 16, 1000000000000000000, 2600, 100, 100000000, 1000000, 11, 'AAAAC29LfJV6QqMt2DDUMhLVu6/rkiFTjkAQpCNFUzDdsJsrycO3J/Hp2cEvr433GPeq3XtJHPV8+zXdxLRRUx16QhEAAAAAXPVMXwAAAAAAAAAAqHWOFRdWDweK4ZIEBYTYOkzwi2q8iIoB9xFzAOqPS/ffjPLaERGDLrYC7tZYVR7nTTTU7Tn04Nc3TOs13gsC1wAAABEN4Lazp2QAAAAAAAAAAAooAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (16, '6f4b7c957a42a32dd830d43212d5bbafeb9221538e4010a423455330ddb09b2b', '8a29232992d69bb61f457c768a745b6fb2d81dadb8d6780f4960b9ee382c9ad5', 3, 3, '2019-06-03 16:35:42', '2019-06-03 16:35:33.910541', '2019-06-03 16:35:33.910541', 68719476736, 16, 1000000000000000000, 2400, 100, 100000000, 1000000, 11, 'AAAAC4opIymS1pu2H0V8dop0W2+y2B2tuNZ4D0lgue44LJrVY1gAUWBGwQQPMqY6PP8FFyppkHosKR4TD1wsYd//I+cAAAAAXPVMXgAAAAAAAAAAYKou7BkPtQtmD4O54/TToUgBAVUCmCl/O+9Tgq2hNHzlKY9Qg5CAIxdRXbs/XcVnEb3JYJf/Dk9K0NMoW5O8TwAAABAN4Lazp2QAAAAAAAAAAAlgAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0, NULL);
INSERT INTO history_ledgers VALUES (15, '8a29232992d69bb61f457c768a745b6fb2d81dadb8d6780f4960b9ee382c9ad5', 'a64af1c08fb5d28abcd1fc9dae8ccb1bfc9bbec1b7b477e14fca2cbadca4907c', 1, 1, '2019-06-03 16:35:41', '2019-06-03 16:35:33.920642', '2019-06-03 16:35:33.920642', 64424509440, 16, 1000000000000000000, 2100, 100, 100000000, 1000000, 11, 'AAAAC6ZK8cCPtdKKvNH8na6Myxv8m77Bt7R34U/KLLrcpJB8d2gda9aXjRx0ma+uG5InDPEEWb37fjZyNQyPwPiQEvEAAAAAXPVMXQAAAAAAAAAAUsjUlWJzXzXDGEu7ZX8NM4G3hXMJgTqm1o4qvEMfb6GHCRkVinBcpegw5bhvYMpwJ2bdfeb93l1bSBPVUDjwrgAAAA8N4Lazp2QAAAAAAAAAAAg0AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (14, 'a64af1c08fb5d28abcd1fc9dae8ccb1bfc9bbec1b7b477e14fca2cbadca4907c', '0c14a465f60001303a566fec15f2cbcda27abc34e9c8b51b88954a2f9b30f6b1', 2, 2, '2019-06-03 16:35:40', '2019-06-03 16:35:33.933196', '2019-06-03 16:35:33.933196', 60129542144, 16, 1000000000000000000, 2000, 100, 100000000, 1000000, 11, 'AAAACwwUpGX2AAEwOlZv7BXyy82ierw06ci1G4iVSi+bMPaxv5ctYdcMSEm83sSvokp32C1S/AIgbbctZWaWvUxN2AMAAAAAXPVMXAAAAAAAAAAAl8eMbvLJtOBoTCvorlEKxO6gpnsCesjn+1C1/HG0v1rd27oO7+VWSEG8r8fcOFYK8k3ZdplUbzIdfIxkNaBaKQAAAA4N4Lazp2QAAAAAAAAAAAfQAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (13, '0c14a465f60001303a566fec15f2cbcda27abc34e9c8b51b88954a2f9b30f6b1', '039d9303c20c9a58082f15f4875e84c1f0ce4dc20ea237866c59d597cc349a1d', 2, 2, '2019-06-03 16:35:39', '2019-06-03 16:35:33.94895', '2019-06-03 16:35:33.94895', 55834574848, 16, 1000000000000000000, 1800, 100, 100000000, 1000000, 11, 'AAAACwOdkwPCDJpYCC8V9IdehMHwzk3CDqI3hmxZ1ZfMNJodFfvdztxrYHP9xzdqb/9pcKohIXiTntWz+SxP2oy5B3sAAAAAXPVMWwAAAAAAAAAAFkWUNyF+w67R5n+HgDNebl5snV2YlBsz3xLTqBMQsQEt+GWM9UtEdrkYmmI5VHYouT3zPlYpLCOylqjq2ciR9gAAAA0N4Lazp2QAAAAAAAAAAAcIAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0, NULL);
INSERT INTO history_ledgers VALUES (12, '039d9303c20c9a58082f15f4875e84c1f0ce4dc20ea237866c59d597cc349a1d', '5e5e893b0f9078a77603f60689405d24d24247d0e19b00f815d212befb819656', 1, 1, '2019-06-03 16:35:38', '2019-06-03 16:35:33.960179', '2019-06-03 16:35:33.960179', 51539607552, 16, 1000000000000000000, 1600, 100, 100000000, 1000000, 11, 'AAAAC15eiTsPkHindgP2BolAXSTSQkfQ4ZsA+BXSEr77gZZWFSHyciAPeIm0Aaehq+XrRbfBf2YWt2omzfQnMSnMapcAAAAAXPVMWgAAAAAAAAAAc4P3o+BJiJvXXIfj/eVQPF1RRIrnE2nfX0a4J7WFySrDTzVqaL+dwGrt1jCSqZ7dkPeUD98aSKLZN1dQHYoQhQAAAAwN4Lazp2QAAAAAAAAAAAZAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (11, '5e5e893b0f9078a77603f60689405d24d24247d0e19b00f815d212befb819656', '62ad4b3aabc2a93e094118968e6114d1e15af1dfb753f5b7048c01be5fb3e9b9', 1, 1, '2019-06-03 16:35:37', '2019-06-03 16:35:33.971637', '2019-06-03 16:35:33.971637', 47244640256, 16, 1000000000000000000, 1500, 100, 100000000, 1000000, 11, 'AAAAC2KtSzqrwqk+CUEYlo5hFNHhWvHft1P1twSMAb5fs+m577MF5b1Iei3zHF9EAmsLT84l7YbAUSH/7F/wymRq98QAAAAAXPVMWQAAAAAAAAAApIxpaTPANJS3XftGXFSYrEJ0MhfN6GN1Oe5+YEP609+t1m4CDabbFVCwxaT1VIvC4WC/6n6mrrOSvo/Ob5mk1QAAAAsN4Lazp2QAAAAAAAAAAAXcAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (10, '62ad4b3aabc2a93e094118968e6114d1e15af1dfb753f5b7048c01be5fb3e9b9', '69f27e72b97f7cb7bfd7f7cd0012b7ba080d3dec181688848626b1e1bad25ea5', 1, 2, '2019-06-03 16:35:36', '2019-06-03 16:35:33.980663', '2019-06-03 16:35:33.980663', 42949672960, 16, 1000000000000000000, 1400, 100, 100000000, 1000000, 11, 'AAAAC2nyfnK5f3y3v9f3zQASt7oIDT3sGBaIhIYmseG60l6lp9U56RSocR4ZpXXZs/glcnGVHESVFFU6LiK9fweLLYkAAAAAXPVMWAAAAAAAAAAA9+jjsWK6v6g0OYMFxTo1+Yogi2yDSjXhJ86N1AxJOvEc1lxN2lpI8SegGY6tvAaspu4dHxmLVdrEkjSlIYvfcgAAAAoN4Lazp2QAAAAAAAAAAAV4AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (9, '69f27e72b97f7cb7bfd7f7cd0012b7ba080d3dec181688848626b1e1bad25ea5', '6afd133af73ffbd7ed56c265d6d3b7e374b0d17051cd1cf8e9edb76befcdb47b', 1, 1, '2019-06-03 16:35:35', '2019-06-03 16:35:33.989863', '2019-06-03 16:35:33.989863', 38654705664, 16, 1000000000000000000, 1200, 100, 100000000, 1000000, 11, 'AAAAC2r9Ezr3P/vX7VbCZdbTt+N0sNFwUc0c+Ontt2vvzbR7I1aOU+Htzv1YyoOW+NUXosjEw6RZFhwtO1PMnARAtukAAAAAXPVMVwAAAAAAAAAAbXQLWug2IZ2G4RnwUrCTvIwO6TMpu21+S2UXhqLKWTHGAwX5y0FVUEjDd5YDib7wh/bNSxZ606YM0D5P9vigSgAAAAkN4Lazp2QAAAAAAAAAAASwAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (8, '6afd133af73ffbd7ed56c265d6d3b7e374b0d17051cd1cf8e9edb76befcdb47b', '0cba5d6486b080c20944d8e2078da05854eb4d68b0c44a7a70f06dfbb28c0809', 4, 4, '2019-06-03 16:35:34', '2019-06-03 16:35:33.999565', '2019-06-03 16:35:33.999565', 34359738368, 16, 1000000000000000000, 1100, 100, 100000000, 1000000, 11, 'AAAACwy6XWSGsIDCCUTY4geNoFhU601osMRKenDwbfuyjAgJZ/gVHh4EGoRXjxBlJXgdb6o55e5LPB8/rItX+iBMmg0AAAAAXPVMVgAAAAAAAAAA5fh0DRZ+OfeA4iH3GTMbG6cqJGEx0qC57216HbTwQPU0hzOQ0tEwod0WKGD5kM6LwVHLKbIOzBmb/9SIwEo4BwAAAAgN4Lazp2QAAAAAAAAAAARMAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0, NULL);
INSERT INTO history_ledgers VALUES (7, '0cba5d6486b080c20944d8e2078da05854eb4d68b0c44a7a70f06dfbb28c0809', 'e56b9392f09492f4084789a6f7eea3274d6991c918d95096009e985fbe6c0d09', 1, 1, '2019-06-03 16:35:33', '2019-06-03 16:35:34.011488', '2019-06-03 16:35:34.011488', 30064771072, 16, 1000000000000000000, 700, 100, 100000000, 1000000, 11, 'AAAAC+Vrk5LwlJL0CEeJpvfuoydNaZHJGNlQlgCemF++bA0JrU2Bn4aTjYMLpDfffI/e6HD3GSRCEJo4ArYti94OWqYAAAAAXPVMVQAAAAAAAAAAk5HOQp3TCIyOUj8VmNGrhs2C+IXj942pn3qPVIAXXx6Y378rPisQ/PWRLazxMClLw8L2UkId9o64eIplaVk81gAAAAcN4Lazp2QAAAAAAAAAAAK8AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (6, 'e56b9392f09492f4084789a6f7eea3274d6991c918d95096009e985fbe6c0d09', 'f5f9c1fc700571da5df0d6e0ecedab78fe81e6f66054e4352f97a3f0cf3d5f37', 1, 1, '2019-06-03 16:35:32', '2019-06-03 16:35:34.0218', '2019-06-03 16:35:34.0218', 25769803776, 16, 1000000000000000000, 600, 100, 100000000, 1000000, 11, 'AAAAC/X5wfxwBXHaXfDW4Oztq3j+geb2YFTkNS+Xo/DPPV83Ny97k/0F2AWH2PP6jBxodIXMgjmRLJwhDtO8tQZgf8sAAAAAXPVMVAAAAAAAAAAAN8Q2AtCCsel4HIsK4udqXyLQqcraGjcAHVSUN815grFwXcnuHU1uUpGqTUtbI/yZP1yYvQEniAuMf2hiwiatTAAAAAYN4Lazp2QAAAAAAAAAAAJYAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (5, 'f5f9c1fc700571da5df0d6e0ecedab78fe81e6f66054e4352f97a3f0cf3d5f37', 'c26019112e7d1a307b7b6a843ad56fd73b51fbdc2e635fe7b0fcc8bae46e9417', 3, 3, '2019-06-03 16:35:31', '2019-06-03 16:35:34.03588', '2019-06-03 16:35:34.035881', 21474836480, 16, 1000000000000000000, 500, 100, 100000000, 1000000, 11, 'AAAAC8JgGREufRowe3tqhDrVb9c7UfvcLmNf57D8yLrkbpQXyz9iY2rHhwuUvOFuLE6mBvTO46hI+XxK2J5c1tuLsiUAAAAAXPVMUwAAAAAAAAAAk6vO/KRxwZ0ynJMMrKUpOSf5WhbdkNk0LIO/zZxB+qeX7aYG0z3XmFwxmr34PrxhVy8SRzlUMGjhTLWe3i+lJAAAAAUN4Lazp2QAAAAAAAAAAAH0AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0, NULL);
INSERT INTO history_ledgers VALUES (4, 'c26019112e7d1a307b7b6a843ad56fd73b51fbdc2e635fe7b0fcc8bae46e9417', '341ff9c498b9076f8f21aea78a33bc0c639ab5cf986adbc2d8d04d92b5a9d6df', 1, 1, '2019-06-03 16:35:30', '2019-06-03 16:35:34.052795', '2019-06-03 16:35:34.052795', 17179869184, 16, 1000000000000000000, 200, 100, 100000000, 1000000, 11, 'AAAACzQf+cSYuQdvjyGup4ozvAxjmrXPmGrbwtjQTZK1qdbfS+Xef4FXXIT2kc1XcfznO1AYP9JNuDCI1JFXHpnU9IsAAAAAXPVMUgAAAAAAAAAAI8FeHEzJ0/4FO8/NG/43wTEny1CW/vcH2jcrQkzbbcj+Ff6Dfgn4FQTGJXl3s+xBqyBn+uvEPEYwqrPTrPE/VAAAAAQN4Lazp2QAAAAAAAAAAADIAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (3, '341ff9c498b9076f8f21aea78a33bc0c639ab5cf986adbc2d8d04d92b5a9d6df', 'fd9268d16d7501b6eb7956b6756b28fe43ce5a6010b76e776d91f60e5eed7b2a', 1, 1, '2019-06-03 16:35:29', '2019-06-03 16:35:34.065174', '2019-06-03 16:35:34.065175', 12884901888, 16, 1000000000000000000, 100, 100, 100000000, 1000000, 11, 'AAAAC/2SaNFtdQG263lWtnVrKP5DzlpgELdud22R9g5e7XsqDZGQBE1YsotIrRXkw7x6/cKPd2zOCBuz7H2RlRoSGYEAAAAAXPVMUQAAAAAAAAAA5PnvEvMOE8kyDs3Gbu1hou8cmoww1I7xqYGObt0Zo0oZAKD5U1RxNy8ovj+4OjM4rIB0P/C3cQ+hCLArM49g+AAAAAMN4Lazp2QAAAAAAAAAAABkAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0, NULL);
INSERT INTO history_ledgers VALUES (2, 'fd9268d16d7501b6eb7956b6756b28fe43ce5a6010b76e776d91f60e5eed7b2a', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 0, 0, '2019-06-03 16:35:28', '2019-06-03 16:35:34.076291', '2019-06-03 16:35:34.076291', 8589934592, 16, 1000000000000000000, 0, 100, 100000000, 1000000, 11, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76ZuZRHr9UdXKbTKiclfOjy72YZFJUkJPVcKT5htvorm1QAAAAAXPVMUAAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERnv+Ltt/ycz/x8/+lFB80rnVx7j2Mrm29EpusUR+gv9ZAAAAAIN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0, NULL);
INSERT INTO history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-03 16:35:34.082668', '2019-06-03 16:35:34.082668', 4294967296, 16, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0, NULL);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operation_participants VALUES (1, 261993009153, 1);
INSERT INTO history_operation_participants VALUES (2, 257698041857, 1);
INSERT INTO history_operation_participants VALUES (3, 253403074561, 1);
INSERT INTO history_operation_participants VALUES (4, 249108107265, 1);
INSERT INTO history_operation_participants VALUES (5, 244813139969, 1);
INSERT INTO history_operation_participants VALUES (6, 244813139969, 2);
INSERT INTO history_operation_participants VALUES (7, 240518172673, 3);
INSERT INTO history_operation_participants VALUES (8, 236223205377, 2);
INSERT INTO history_operation_participants VALUES (9, 236223205377, 3);
INSERT INTO history_operation_participants VALUES (10, 231928238081, 2);
INSERT INTO history_operation_participants VALUES (11, 231928238081, 4);
INSERT INTO history_operation_participants VALUES (12, 231928238082, 4);
INSERT INTO history_operation_participants VALUES (13, 231928238082, 2);
INSERT INTO history_operation_participants VALUES (14, 227633270785, 2);
INSERT INTO history_operation_participants VALUES (15, 227633270785, 4);
INSERT INTO history_operation_participants VALUES (16, 223338303489, 5);
INSERT INTO history_operation_participants VALUES (17, 219043336193, 5);
INSERT INTO history_operation_participants VALUES (18, 214748368897, 5);
INSERT INTO history_operation_participants VALUES (19, 210453401601, 5);
INSERT INTO history_operation_participants VALUES (20, 210453405697, 5);
INSERT INTO history_operation_participants VALUES (21, 210453409793, 5);
INSERT INTO history_operation_participants VALUES (22, 206158434305, 2);
INSERT INTO history_operation_participants VALUES (23, 206158434305, 5);
INSERT INTO history_operation_participants VALUES (24, 201863467009, 2);
INSERT INTO history_operation_participants VALUES (25, 197568499713, 6);
INSERT INTO history_operation_participants VALUES (26, 197568503809, 2);
INSERT INTO history_operation_participants VALUES (27, 193273532417, 6);
INSERT INTO history_operation_participants VALUES (28, 193273532417, 2);
INSERT INTO history_operation_participants VALUES (29, 188978565121, 2);
INSERT INTO history_operation_participants VALUES (30, 188978565121, 7);
INSERT INTO history_operation_participants VALUES (31, 184683597825, 7);
INSERT INTO history_operation_participants VALUES (32, 184683597825, 2);
INSERT INTO history_operation_participants VALUES (33, 180388630529, 8);
INSERT INTO history_operation_participants VALUES (34, 180388630529, 9);
INSERT INTO history_operation_participants VALUES (35, 176093663233, 8);
INSERT INTO history_operation_participants VALUES (36, 176093663233, 9);
INSERT INTO history_operation_participants VALUES (37, 176093667329, 8);
INSERT INTO history_operation_participants VALUES (38, 176093667329, 9);
INSERT INTO history_operation_participants VALUES (39, 171798695937, 9);
INSERT INTO history_operation_participants VALUES (40, 171798700033, 9);
INSERT INTO history_operation_participants VALUES (41, 167503728641, 8);
INSERT INTO history_operation_participants VALUES (42, 163208761345, 2);
INSERT INTO history_operation_participants VALUES (43, 163208761345, 9);
INSERT INTO history_operation_participants VALUES (44, 163208765441, 2);
INSERT INTO history_operation_participants VALUES (45, 163208765441, 8);
INSERT INTO history_operation_participants VALUES (46, 158913794049, 10);
INSERT INTO history_operation_participants VALUES (47, 154618826753, 10);
INSERT INTO history_operation_participants VALUES (48, 150323859457, 10);
INSERT INTO history_operation_participants VALUES (49, 146028892161, 10);
INSERT INTO history_operation_participants VALUES (50, 141733924865, 2);
INSERT INTO history_operation_participants VALUES (51, 141733924865, 10);
INSERT INTO history_operation_participants VALUES (52, 137438957569, 11);
INSERT INTO history_operation_participants VALUES (53, 137438961665, 11);
INSERT INTO history_operation_participants VALUES (54, 133143990273, 11);
INSERT INTO history_operation_participants VALUES (55, 128849022977, 11);
INSERT INTO history_operation_participants VALUES (56, 124554055681, 11);
INSERT INTO history_operation_participants VALUES (57, 120259088385, 11);
INSERT INTO history_operation_participants VALUES (58, 120259092481, 11);
INSERT INTO history_operation_participants VALUES (59, 120259096577, 11);
INSERT INTO history_operation_participants VALUES (60, 120259100673, 11);
INSERT INTO history_operation_participants VALUES (61, 120259104769, 11);
INSERT INTO history_operation_participants VALUES (62, 120259108865, 11);
INSERT INTO history_operation_participants VALUES (63, 120259112961, 11);
INSERT INTO history_operation_participants VALUES (64, 115964121089, 2);
INSERT INTO history_operation_participants VALUES (65, 115964121089, 11);
INSERT INTO history_operation_participants VALUES (66, 111669153793, 12);
INSERT INTO history_operation_participants VALUES (67, 111669157889, 12);
INSERT INTO history_operation_participants VALUES (68, 107374186497, 2);
INSERT INTO history_operation_participants VALUES (69, 107374186497, 12);
INSERT INTO history_operation_participants VALUES (70, 103079219201, 14);
INSERT INTO history_operation_participants VALUES (71, 98784251905, 13);
INSERT INTO history_operation_participants VALUES (72, 94489284609, 13);
INSERT INTO history_operation_participants VALUES (73, 90194317313, 2);
INSERT INTO history_operation_participants VALUES (74, 90194317313, 13);
INSERT INTO history_operation_participants VALUES (75, 90194321409, 2);
INSERT INTO history_operation_participants VALUES (76, 90194321409, 14);
INSERT INTO history_operation_participants VALUES (77, 85899350017, 16);
INSERT INTO history_operation_participants VALUES (78, 85899350017, 17);
INSERT INTO history_operation_participants VALUES (79, 81604382721, 16);
INSERT INTO history_operation_participants VALUES (80, 81604382721, 17);
INSERT INTO history_operation_participants VALUES (81, 77309415425, 15);
INSERT INTO history_operation_participants VALUES (82, 77309415425, 16);
INSERT INTO history_operation_participants VALUES (83, 77309419521, 15);
INSERT INTO history_operation_participants VALUES (84, 77309423617, 15);
INSERT INTO history_operation_participants VALUES (85, 73014448129, 16);
INSERT INTO history_operation_participants VALUES (86, 73014452225, 17);
INSERT INTO history_operation_participants VALUES (87, 68719480833, 2);
INSERT INTO history_operation_participants VALUES (88, 68719480833, 16);
INSERT INTO history_operation_participants VALUES (89, 68719484929, 2);
INSERT INTO history_operation_participants VALUES (90, 68719484929, 17);
INSERT INTO history_operation_participants VALUES (91, 68719489025, 2);
INSERT INTO history_operation_participants VALUES (92, 68719489025, 15);
INSERT INTO history_operation_participants VALUES (93, 64424513537, 19);
INSERT INTO history_operation_participants VALUES (94, 64424513537, 18);
INSERT INTO history_operation_participants VALUES (95, 60129546241, 19);
INSERT INTO history_operation_participants VALUES (96, 60129546241, 18);
INSERT INTO history_operation_participants VALUES (97, 60129550337, 18);
INSERT INTO history_operation_participants VALUES (98, 55834578945, 2);
INSERT INTO history_operation_participants VALUES (99, 55834578945, 19);
INSERT INTO history_operation_participants VALUES (100, 55834583041, 2);
INSERT INTO history_operation_participants VALUES (101, 55834583041, 18);
INSERT INTO history_operation_participants VALUES (102, 51539611649, 21);
INSERT INTO history_operation_participants VALUES (103, 51539611649, 20);
INSERT INTO history_operation_participants VALUES (104, 47244644353, 2);
INSERT INTO history_operation_participants VALUES (105, 47244644353, 21);
INSERT INTO history_operation_participants VALUES (106, 42949677057, 22);
INSERT INTO history_operation_participants VALUES (107, 42949677057, 2);
INSERT INTO history_operation_participants VALUES (108, 42949677058, 22);
INSERT INTO history_operation_participants VALUES (109, 42949677058, 2);
INSERT INTO history_operation_participants VALUES (110, 38654709761, 2);
INSERT INTO history_operation_participants VALUES (111, 38654709761, 22);
INSERT INTO history_operation_participants VALUES (112, 34359742465, 23);
INSERT INTO history_operation_participants VALUES (113, 34359742465, 2);
INSERT INTO history_operation_participants VALUES (114, 34359746561, 23);
INSERT INTO history_operation_participants VALUES (115, 34359746561, 2);
INSERT INTO history_operation_participants VALUES (116, 34359750657, 23);
INSERT INTO history_operation_participants VALUES (117, 34359750657, 2);
INSERT INTO history_operation_participants VALUES (118, 34359754753, 23);
INSERT INTO history_operation_participants VALUES (119, 34359754753, 2);
INSERT INTO history_operation_participants VALUES (120, 30064775169, 2);
INSERT INTO history_operation_participants VALUES (121, 30064775169, 23);
INSERT INTO history_operation_participants VALUES (122, 25769807873, 24);
INSERT INTO history_operation_participants VALUES (123, 21474840577, 24);
INSERT INTO history_operation_participants VALUES (124, 21474844673, 24);
INSERT INTO history_operation_participants VALUES (125, 21474848769, 24);
INSERT INTO history_operation_participants VALUES (126, 17179873281, 24);
INSERT INTO history_operation_participants VALUES (127, 17179873281, 2);
INSERT INTO history_operation_participants VALUES (128, 12884905985, 2);
INSERT INTO history_operation_participants VALUES (129, 12884905985, 25);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_operation_participants_id_seq', 129, true);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operations VALUES (261993009153, 261993009152, 1, 11, '{"bump_to": "300000000003"}', 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN');
INSERT INTO history_operations VALUES (257698041857, 257698041856, 1, 11, '{"bump_to": "300000000001"}', 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN');
INSERT INTO history_operations VALUES (253403074561, 253403074560, 1, 11, '{"bump_to": "100"}', 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN');
INSERT INTO history_operations VALUES (249108107265, 249108107264, 1, 11, '{"bump_to": "300000000000"}', 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN');
INSERT INTO history_operations VALUES (244813139969, 244813139968, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (240518172673, 240518172672, 1, 1, '{"to": "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y", "from": "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y", "amount": "10.0000000", "asset_type": "native"}', 'GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y');
INSERT INTO history_operations VALUES (236223205377, 236223205376, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (231928238081, 231928238080, 1, 1, '{"to": "GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X", "from": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "amount": "10.0000000", "asset_type": "native"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (231928238082, 231928238080, 2, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X", "amount": "10.0000000", "asset_type": "native"}', 'GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X');
INSERT INTO history_operations VALUES (227633270785, 227633270784, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GACJPE4YUR22VP4CM2BDFDAHY3DLEF3H7NENKUQ53DT5TEI2GAHT5N4X", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (223338303489, 223338303488, 1, 10, '{"name": "name1", "value": "MDAwMA=="}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (219043336193, 219043336192, 1, 10, '{"name": "name1", "value": "MTIzNA=="}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (214748368897, 214748368896, 1, 10, '{"name": "name2", "value": null}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (210453401601, 210453401600, 1, 10, '{"name": "name1", "value": "MTIzNA=="}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (210453405697, 210453405696, 1, 10, '{"name": "name2", "value": "NTY3OA=="}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (210453409793, 210453409792, 1, 10, '{"name": "name ", "value": "aXRzIGdvdCBzcGFjZXMh"}', 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD');
INSERT INTO history_operations VALUES (206158434305, 206158434304, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (201863467009, 201863467008, 1, 9, '{}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (197568499713, 197568499712, 1, 5, '{"inflation_dest": "GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS"}', 'GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS');
INSERT INTO history_operations VALUES (197568503809, 197568503808, 1, 5, '{"inflation_dest": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (193273532417, 193273532416, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS", "starting_balance": "20000000000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (188978565121, 188978565120, 1, 8, '{"into": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ"}', 'GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ');
INSERT INTO history_operations VALUES (184683597825, 184683597824, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (180388630529, 180388630528, 1, 7, '{"trustee": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "authorize": false, "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_operations VALUES (176093663233, 176093663232, 1, 7, '{"trustee": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "authorize": true, "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_operations VALUES (176093667329, 176093667328, 1, 7, '{"trustee": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "authorize": true, "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_operations VALUES (171798695937, 171798695936, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}', 'GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG');
INSERT INTO history_operations VALUES (171798700033, 171798700032, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "trustor": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF"}', 'GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG');
INSERT INTO history_operations VALUES (167503728641, 167503728640, 1, 5, '{"set_flags": [1, 2], "set_flags_s": ["auth_required", "auth_revocable"]}', 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF');
INSERT INTO history_operations VALUES (163208761345, 163208761344, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (163208765441, 163208765440, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (158913794049, 158913794048, 1, 6, '{"limit": "0.0000000", "trustee": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "trustor": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG');
INSERT INTO history_operations VALUES (154618826753, 154618826752, 1, 6, '{"limit": "100.0000000", "trustee": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "trustor": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG');
INSERT INTO history_operations VALUES (150323859457, 150323859456, 1, 6, '{"limit": "100.0000000", "trustee": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "trustor": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG');
INSERT INTO history_operations VALUES (146028892161, 146028892160, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "trustor": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG');
INSERT INTO history_operations VALUES (141733924865, 141733924864, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (137438957569, 137438957568, 1, 5, '{"clear_flags": [1, 2], "clear_flags_s": ["auth_required", "auth_revocable"]}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (137438961665, 137438961664, 1, 5, '{"signer_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE", "signer_weight": 0}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (133143990273, 133143990272, 1, 5, '{"signer_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE", "signer_weight": 5}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (128849022977, 128849022976, 1, 5, '{"signer_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE", "signer_weight": 1}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (124554055681, 124554055680, 1, 5, '{"master_key_weight": 2}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259088385, 120259088384, 1, 5, '{"inflation_dest": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H"}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259092481, 120259092480, 1, 5, '{"set_flags": [1], "set_flags_s": ["auth_required"]}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259096577, 120259096576, 1, 5, '{"set_flags": [2], "set_flags_s": ["auth_revocable"]}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259100673, 120259100672, 1, 5, '{"master_key_weight": 2}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259104769, 120259104768, 1, 5, '{"low_threshold": 0, "med_threshold": 2, "high_threshold": 2}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259108865, 120259108864, 1, 5, '{"home_domain": "example.com"}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (120259112961, 120259112960, 1, 5, '{"signer_key": "GB6J3WOLKYQE6KVDZEA4JDMFTTONUYP3PUHNDNZRWIKA6JQWIMJZATFE", "signer_weight": 1}', 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES');
INSERT INTO history_operations VALUES (115964121089, 115964121088, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (111669153793, 111669153792, 1, 4, '{"price": "1.0000000", "amount": "200.0000000", "price_r": {"d": 1, "n": 1}, "buying_asset_type": "native", "selling_asset_code": "USD", "selling_asset_type": "credit_alphanum4", "selling_asset_issuer": "GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q"}', 'GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q');
INSERT INTO history_operations VALUES (111669157889, 111669157888, 1, 4, '{"price": "1.0000000", "amount": "200.0000000", "price_r": {"d": 1, "n": 1}, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_type": "native", "buying_asset_issuer": "GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q"}', 'GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q');
INSERT INTO history_operations VALUES (107374186497, 107374186496, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (103079219201, 103079219200, 1, 3, '{"price": "1.0000000", "amount": "30.0000000", "price_r": {"d": 1, "n": 1}, "offer_id": 0, "buying_asset_type": "native", "selling_asset_code": "USD", "selling_asset_type": "credit_alphanum4", "selling_asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}', 'GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD');
INSERT INTO history_operations VALUES (98784251905, 98784251904, 1, 3, '{"price": "1.0000000", "amount": "20.0000000", "price_r": {"d": 1, "n": 1}, "offer_id": 0, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_type": "native", "buying_asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}', 'GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC');
INSERT INTO history_operations VALUES (94489284609, 94489284608, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD", "trustor": "GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD"}', 'GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC');
INSERT INTO history_operations VALUES (90194317313, 90194317312, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (90194321409, 90194321408, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (85899350017, 85899350016, 1, 2, '{"to": "GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP", "from": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "path": [], "amount": "100.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "source_max": "100.0000000", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "source_amount": "100.0000000", "source_asset_type": "native"}', 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD');
INSERT INTO history_operations VALUES (81604382721, 81604382720, 1, 2, '{"to": "GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP", "from": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "path": [{"asset_type": "native"}], "amount": "200.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "source_max": "100.0000000", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "source_amount": "100.0000000", "source_asset_code": "USD", "source_asset_type": "credit_alphanum4", "source_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD');
INSERT INTO history_operations VALUES (77309415425, 77309415424, 1, 1, '{"to": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "from": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "amount": "100.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_operations VALUES (77309419521, 77309419520, 1, 3, '{"price": "0.5000000", "amount": "400.0000000", "price_r": {"d": 2, "n": 1}, "offer_id": 0, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_type": "native", "buying_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_operations VALUES (77309423617, 77309423616, 1, 3, '{"price": "1.0000000", "amount": "300.0000000", "price_r": {"d": 1, "n": 1}, "offer_id": 0, "buying_asset_type": "native", "selling_asset_code": "EUR", "selling_asset_type": "credit_alphanum4", "selling_asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU');
INSERT INTO history_operations VALUES (73014448129, 73014448128, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "trustor": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD');
INSERT INTO history_operations VALUES (73014452225, 73014452224, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "trustor": "GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU"}', 'GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP');
INSERT INTO history_operations VALUES (68719480833, 68719480832, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (68719484929, 68719484928, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (68719489025, 68719489024, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (64424513537, 64424513536, 1, 1, '{"to": "GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C", "from": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG", "amount": "10.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}', 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG');
INSERT INTO history_operations VALUES (60129546241, 60129546240, 1, 1, '{"to": "GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C", "from": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG", "amount": "10.0000000", "asset_type": "native"}', 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG');
INSERT INTO history_operations VALUES (60129550337, 60129550336, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG", "trustor": "GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG"}', 'GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C');
INSERT INTO history_operations VALUES (55834578945, 55834578944, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (55834583041, 55834583040, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (51539611649, 51539611648, 1, 0, '{"funder": "GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO", "account": "GCB7FPYGLL6RJ37HKRAYW5TAWMFBGGFGM4IM6ERBCZXI2BZ4OOOX2UAY", "starting_balance": "50.0000000"}', 'GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO');
INSERT INTO history_operations VALUES (47244644353, 47244644352, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (42949677057, 42949677056, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (42949677058, 42949677056, 2, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (38654709761, 38654709760, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (34359742465, 34359742464, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB", "amount": "1.0000000", "asset_type": "native"}', 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB');
INSERT INTO history_operations VALUES (34359746561, 34359746560, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB", "amount": "1.0000000", "asset_type": "native"}', 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB');
INSERT INTO history_operations VALUES (34359750657, 34359750656, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB", "amount": "1.0000000", "asset_type": "native"}', 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB');
INSERT INTO history_operations VALUES (34359754753, 34359754752, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB", "amount": "1.0000000", "asset_type": "native"}', 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB');
INSERT INTO history_operations VALUES (30064775169, 30064775168, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (25769807873, 25769807872, 1, 5, '{"master_key_weight": 2}', 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB');
INSERT INTO history_operations VALUES (21474840577, 21474840576, 1, 5, '{"master_key_weight": 1}', 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB');
INSERT INTO history_operations VALUES (21474844673, 21474844672, 1, 5, '{"signer_key": "GD3E7HKMRNT6HGBGHBT6I6JE4N2S4W5KZ246TGJ4KQSXJ2P4BXCUPQMP", "signer_weight": 1}', 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB');
INSERT INTO history_operations VALUES (21474848769, 21474848768, 1, 5, '{"low_threshold": 2, "med_threshold": 2, "high_threshold": 2}', 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB');
INSERT INTO history_operations VALUES (17179873281, 17179873280, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (12884905985, 12884905984, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GAXI33UCLQTCKM2NMRBS7XYBR535LLEVAHL5YBN4FTCB4HZHT7ZA5CVK", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_trades VALUES (103079219201, 0, '2019-06-03 16:35:50', 3, 14, 1, 200000000, 13, 7, 200000000, false, 1, 1, 4, 3);
INSERT INTO history_trades VALUES (85899350017, 0, '2019-06-03 16:35:46', 2, 15, 4, 1000000000, 16, 7, 1000000000, true, 1, 1, 2, 4611686104326737921);
INSERT INTO history_trades VALUES (81604382721, 0, '2019-06-03 16:35:45', 1, 16, 2, 1000000000, 15, 7, 2000000000, false, 2, 1, 4611686100031770625, 1);
INSERT INTO history_trades VALUES (81604382721, 1, '2019-06-03 16:35:45', 2, 15, 4, 2000000000, 16, 7, 2000000000, true, 1, 1, 2, 4611686100031770625);


--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transaction_participants VALUES (1, 261993009152, 1);
INSERT INTO history_transaction_participants VALUES (2, 257698041856, 1);
INSERT INTO history_transaction_participants VALUES (3, 253403074560, 1);
INSERT INTO history_transaction_participants VALUES (4, 249108107264, 1);
INSERT INTO history_transaction_participants VALUES (5, 244813139968, 2);
INSERT INTO history_transaction_participants VALUES (6, 244813139968, 1);
INSERT INTO history_transaction_participants VALUES (7, 240518172672, 3);
INSERT INTO history_transaction_participants VALUES (8, 236223205376, 2);
INSERT INTO history_transaction_participants VALUES (9, 236223205376, 3);
INSERT INTO history_transaction_participants VALUES (10, 231928238080, 2);
INSERT INTO history_transaction_participants VALUES (11, 231928238080, 4);
INSERT INTO history_transaction_participants VALUES (12, 227633270784, 2);
INSERT INTO history_transaction_participants VALUES (13, 227633270784, 4);
INSERT INTO history_transaction_participants VALUES (14, 223338303488, 5);
INSERT INTO history_transaction_participants VALUES (15, 219043336192, 5);
INSERT INTO history_transaction_participants VALUES (16, 214748368896, 5);
INSERT INTO history_transaction_participants VALUES (17, 210453401600, 5);
INSERT INTO history_transaction_participants VALUES (18, 210453405696, 5);
INSERT INTO history_transaction_participants VALUES (19, 210453409792, 5);
INSERT INTO history_transaction_participants VALUES (20, 206158434304, 2);
INSERT INTO history_transaction_participants VALUES (21, 206158434304, 5);
INSERT INTO history_transaction_participants VALUES (22, 201863467008, 2);
INSERT INTO history_transaction_participants VALUES (23, 197568499712, 6);
INSERT INTO history_transaction_participants VALUES (24, 197568503808, 2);
INSERT INTO history_transaction_participants VALUES (25, 193273532416, 2);
INSERT INTO history_transaction_participants VALUES (26, 193273532416, 6);
INSERT INTO history_transaction_participants VALUES (27, 188978565120, 7);
INSERT INTO history_transaction_participants VALUES (28, 188978565120, 2);
INSERT INTO history_transaction_participants VALUES (29, 184683597824, 2);
INSERT INTO history_transaction_participants VALUES (30, 184683597824, 7);
INSERT INTO history_transaction_participants VALUES (31, 180388630528, 8);
INSERT INTO history_transaction_participants VALUES (32, 180388630528, 9);
INSERT INTO history_transaction_participants VALUES (33, 176093663232, 8);
INSERT INTO history_transaction_participants VALUES (34, 176093663232, 9);
INSERT INTO history_transaction_participants VALUES (35, 176093667328, 9);
INSERT INTO history_transaction_participants VALUES (36, 176093667328, 8);
INSERT INTO history_transaction_participants VALUES (37, 171798695936, 9);
INSERT INTO history_transaction_participants VALUES (38, 171798700032, 9);
INSERT INTO history_transaction_participants VALUES (39, 167503728640, 8);
INSERT INTO history_transaction_participants VALUES (40, 163208761344, 2);
INSERT INTO history_transaction_participants VALUES (41, 163208761344, 9);
INSERT INTO history_transaction_participants VALUES (42, 163208765440, 8);
INSERT INTO history_transaction_participants VALUES (43, 163208765440, 2);
INSERT INTO history_transaction_participants VALUES (44, 158913794048, 10);
INSERT INTO history_transaction_participants VALUES (45, 154618826752, 10);
INSERT INTO history_transaction_participants VALUES (46, 150323859456, 10);
INSERT INTO history_transaction_participants VALUES (47, 146028892160, 10);
INSERT INTO history_transaction_participants VALUES (48, 141733924864, 2);
INSERT INTO history_transaction_participants VALUES (49, 141733924864, 10);
INSERT INTO history_transaction_participants VALUES (50, 137438957568, 11);
INSERT INTO history_transaction_participants VALUES (51, 137438961664, 11);
INSERT INTO history_transaction_participants VALUES (52, 133143990272, 11);
INSERT INTO history_transaction_participants VALUES (53, 128849022976, 11);
INSERT INTO history_transaction_participants VALUES (54, 124554055680, 11);
INSERT INTO history_transaction_participants VALUES (55, 120259088384, 11);
INSERT INTO history_transaction_participants VALUES (56, 120259092480, 11);
INSERT INTO history_transaction_participants VALUES (57, 120259096576, 11);
INSERT INTO history_transaction_participants VALUES (58, 120259100672, 11);
INSERT INTO history_transaction_participants VALUES (59, 120259104768, 11);
INSERT INTO history_transaction_participants VALUES (60, 120259108864, 11);
INSERT INTO history_transaction_participants VALUES (61, 120259112960, 11);
INSERT INTO history_transaction_participants VALUES (62, 115964121088, 11);
INSERT INTO history_transaction_participants VALUES (63, 115964121088, 2);
INSERT INTO history_transaction_participants VALUES (64, 111669153792, 12);
INSERT INTO history_transaction_participants VALUES (65, 111669157888, 12);
INSERT INTO history_transaction_participants VALUES (66, 107374186496, 2);
INSERT INTO history_transaction_participants VALUES (67, 107374186496, 12);
INSERT INTO history_transaction_participants VALUES (68, 103079219200, 14);
INSERT INTO history_transaction_participants VALUES (69, 98784251904, 13);
INSERT INTO history_transaction_participants VALUES (70, 94489284608, 13);
INSERT INTO history_transaction_participants VALUES (71, 90194317312, 2);
INSERT INTO history_transaction_participants VALUES (72, 90194317312, 13);
INSERT INTO history_transaction_participants VALUES (73, 90194321408, 2);
INSERT INTO history_transaction_participants VALUES (74, 90194321408, 14);
INSERT INTO history_transaction_participants VALUES (75, 85899350016, 16);
INSERT INTO history_transaction_participants VALUES (76, 85899350016, 17);
INSERT INTO history_transaction_participants VALUES (77, 81604382720, 16);
INSERT INTO history_transaction_participants VALUES (78, 81604382720, 17);
INSERT INTO history_transaction_participants VALUES (79, 77309415424, 15);
INSERT INTO history_transaction_participants VALUES (80, 77309415424, 16);
INSERT INTO history_transaction_participants VALUES (81, 77309419520, 15);
INSERT INTO history_transaction_participants VALUES (82, 77309423616, 15);
INSERT INTO history_transaction_participants VALUES (83, 73014448128, 16);
INSERT INTO history_transaction_participants VALUES (84, 73014452224, 17);
INSERT INTO history_transaction_participants VALUES (85, 68719480832, 2);
INSERT INTO history_transaction_participants VALUES (86, 68719480832, 16);
INSERT INTO history_transaction_participants VALUES (87, 68719484928, 17);
INSERT INTO history_transaction_participants VALUES (88, 68719484928, 2);
INSERT INTO history_transaction_participants VALUES (89, 68719489024, 2);
INSERT INTO history_transaction_participants VALUES (90, 68719489024, 15);
INSERT INTO history_transaction_participants VALUES (91, 64424513536, 19);
INSERT INTO history_transaction_participants VALUES (92, 64424513536, 18);
INSERT INTO history_transaction_participants VALUES (93, 60129546240, 19);
INSERT INTO history_transaction_participants VALUES (94, 60129546240, 18);
INSERT INTO history_transaction_participants VALUES (95, 60129550336, 18);
INSERT INTO history_transaction_participants VALUES (96, 55834578944, 2);
INSERT INTO history_transaction_participants VALUES (97, 55834578944, 19);
INSERT INTO history_transaction_participants VALUES (98, 55834583040, 2);
INSERT INTO history_transaction_participants VALUES (99, 55834583040, 18);
INSERT INTO history_transaction_participants VALUES (100, 51539611648, 21);
INSERT INTO history_transaction_participants VALUES (101, 51539611648, 20);
INSERT INTO history_transaction_participants VALUES (102, 47244644352, 2);
INSERT INTO history_transaction_participants VALUES (103, 47244644352, 21);
INSERT INTO history_transaction_participants VALUES (104, 42949677056, 22);
INSERT INTO history_transaction_participants VALUES (105, 42949677056, 2);
INSERT INTO history_transaction_participants VALUES (106, 38654709760, 2);
INSERT INTO history_transaction_participants VALUES (107, 38654709760, 22);
INSERT INTO history_transaction_participants VALUES (108, 34359742464, 23);
INSERT INTO history_transaction_participants VALUES (109, 34359742464, 2);
INSERT INTO history_transaction_participants VALUES (110, 34359746560, 23);
INSERT INTO history_transaction_participants VALUES (111, 34359746560, 2);
INSERT INTO history_transaction_participants VALUES (112, 34359750656, 23);
INSERT INTO history_transaction_participants VALUES (113, 34359750656, 2);
INSERT INTO history_transaction_participants VALUES (114, 34359754752, 23);
INSERT INTO history_transaction_participants VALUES (115, 34359754752, 2);
INSERT INTO history_transaction_participants VALUES (116, 30064775168, 2);
INSERT INTO history_transaction_participants VALUES (117, 30064775168, 23);
INSERT INTO history_transaction_participants VALUES (118, 25769807872, 24);
INSERT INTO history_transaction_participants VALUES (119, 21474840576, 24);
INSERT INTO history_transaction_participants VALUES (120, 21474844672, 24);
INSERT INTO history_transaction_participants VALUES (121, 21474848768, 24);
INSERT INTO history_transaction_participants VALUES (122, 17179873280, 2);
INSERT INTO history_transaction_participants VALUES (123, 17179873280, 24);
INSERT INTO history_transaction_participants VALUES (124, 12884905984, 2);
INSERT INTO history_transaction_participants VALUES (125, 12884905984, 25);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_transaction_participants_id_seq', 125, true);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transactions VALUES ('bc11b5c41de791369fd85fa1ccf01c35c20df5f98ff2f75d02ead61bfd520e21', 61, 1, 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN', 300000000003, 100, 1, '2019-06-03 16:35:33.314201', '2019-06-03 16:35:33.314201', 261993009152, 'AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAAEXZZLgDAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AwAAAAAAAAABHK0SlAAAAECcI6ex0Dq6YAh6aK14jHxuAvhvKG2+NuzboAKrfYCaC1ZSQ77BYH/5MghPX97JO9WXV17ehNK7d0umxBgaJj8A', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAPQAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvicAAAAEXZZLgCAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAPQAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvicAAAAEXZZLgDAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==', 'AAAAAgAAAAMAAAA8AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+LUAAAARdlkuAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA9AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+JwAAAARdlkuAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{nCOnsdA6umAIemiteIx8bgL4byhtvjbs26ACq32AmgtWUkO+wWB/+TIIT1/eyTvVl1de3oTSu3dLpsQYGiY/AA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c8132b95c0063cafd20b26d27f06c12e688609d2d9d3724b840821e861870b8e', 60, 1, 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN', 300000000002, 100, 1, '2019-06-03 16:35:33.335583', '2019-06-03 16:35:33.335584', 257698041856, 'AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAAEXZZLgCAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AQAAAAAAAAABHK0SlAAAAEC4H7TDntOUXDMg4MfoCPlbLRQZH7VwNpUHMvtnRWqWIiY/qnYYu0bvgYUVtoFOOeqElRKLYqtOW3Fz9iKl0WQJ', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAPAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvi1AAAAEXZZLgBAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAPAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvi1AAAAEXZZLgCAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==', 'AAAAAgAAAAMAAAA7AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+M4AAAARdlkuAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA8AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+LUAAAARdlkuAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{uB+0w57TlFwzIODH6Aj5Wy0UGR+1cDaVBzL7Z0VqliImP6p2GLtG74GFFbaBTjnqhJUSi2KrTltxc/YipdFkCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('74b62d52311ea3f47359f74790595343f976afa4fd306caaefee5efdbbb104ff', 59, 1, 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN', 300000000001, 100, 1, '2019-06-03 16:35:33.346257', '2019-06-03 16:35:33.346257', 253403074560, 'AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAAEXZZLgBAAAAAAAAAAAAAAABAAAAAAAAAAsAAAAAAAAAZAAAAAAAAAABHK0SlAAAAEAOrvZSFnT3JvmT1P5lJ/lggpZe4nxH5WvJ9K/SLOD49wfqq84suncoZIn3IAf0PExMw3etu5FiDVw3c3jYYhAL', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAOwAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjOAAAAEXZZLgAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOwAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjOAAAAEXZZLgBAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==', 'AAAAAgAAAAMAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAARdlkuAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA7AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+M4AAAARdlkuAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Dq72UhZ09yb5k9T+ZSf5YIKWXuJ8R+VryfSv0izg+PcH6qvOLLp3KGSJ9yAH9DxMTMN3rbuRYg1cN3N42GIQCw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('829d53f2dceebe10af8007564b0aefde819b95734ad431df84270651e7ed8a90', 58, 1, 'GCQZP3IU7XU6EJ63JZXKCQOYT2RNXN3HB5CNHENNUEUHSMA4VUJJJSEN', 244813135873, 100, 1, '2019-06-03 16:35:33.357399', '2019-06-03 16:35:33.357399', 249108107264, 'AAAAAKGX7RT96eIn205uoUHYnqLbt2cPRNORraEoeTAcrRKUAAAAZAAAADkAAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAABF2WS4AAAAAAAAAAABHK0SlAAAAEDq0JVhKNIq9ag0sR+R/cv3d9tEuaYEm2BazIzILRdGj9alaVMZBhxoJ3ZIpP3rraCJzyoKZO+p5HBVe10a2+UG', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOgAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvjnAAAADkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAARdlkuAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAA5AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+QAAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA6AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+OcAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{6tCVYSjSKvWoNLEfkf3L93fbRLmmBJtgWsyMyC0XRo/WpWlTGQYcaCd2SKT9662gic8qCmTvqeRwVXtdGtvlBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('0e5bd332291e3098e49886df2cdb9b5369a5f9e0a9973f0d9e1a9489c6581ba2', 57, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 26, 100, 1, '2019-06-03 16:35:33.368603', '2019-06-03 16:35:33.368603', 244813139968, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAaAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAoZftFP3p4ifbTm6hQdieotu3Zw9E05GtoSh5MBytEpQAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDHU95E9wxgETD8TqxUrkgC0/7XHyNDts6Q5huRHfDRyRcoHdv7aMp/sPvC3RPkXjOMjgbKJUX7SgExUeYB5f8F', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAOQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZY9dZxbAAAAAAAAAAZAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZY9dZxbAAAAAAAAAAaAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nFsAAAAAAAAABoAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlahyo1sAAAAAAAAABoAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAA5AAAAAAAAAAChl+0U/eniJ9tObqFB2J6i27dnD0TTka2hKHkwHK0SlAAAAAJUC+QAAAAAOQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nHQAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA5AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nFsAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{x1PeRPcMYBEw/E6sVK5IAtP+1x8jQ7bOkOYbkR3w0ckXKB3b+2jKf7D7wt0T5F4zjI4GyiVF+0oBMVHmAeX/BQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2a805712c6d10f9e74bb0ccf54ae92a2b4b1e586451fe8133a2433816f6b567c', 56, 1, 'GANFZDRBCNTUXIODCJEYMACPMCSZEVE4WZGZ3CZDZ3P2SXK4KH75IK6Y', 236223201281, 100, 1, '2019-06-03 16:35:33.380688', '2019-06-03 16:35:33.380688', 240518172672, 'AAAAABpcjiETZ0uhwxJJhgBPYKWSVJy2TZ2LI87fqV1cUf/UAAAAZAAAADcAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAAAAAAAAAX14QAAAAAAAAAAAVxR/9QAAABAK6pcXYMzAEmH08CZ1LWmvtNDKauhx+OImtP/Lk4hVTMJRVBOebVs5WEPj9iSrgGT0EswuDCZ2i5AEzwgGof9Ag==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAOAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvjnAAAADcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAA==', 'AAAAAgAAAAMAAAA3AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+QAAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA4AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+OcAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{K6pcXYMzAEmH08CZ1LWmvtNDKauhx+OImtP/Lk4hVTMJRVBOebVs5WEPj9iSrgGT0EswuDCZ2i5AEzwgGof9Ag==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('5bbbedfb52efd1d5d973e22540044a27b8115772314293e3ba8b1fb12e63ca2e', 55, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 25, 100, 1, '2019-06-03 16:35:33.392703', '2019-06-03 16:35:33.392703', 236223205376, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAZAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAGlyOIRNnS6HDEkmGAE9gpZJUnLZNnYsjzt+pXVxR/9QAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBCMMjX9xO3XKpQ6uS/U1BqdzRhSBYQ35ivmZxPBgfqQsTDma1BzOsq/bmHJ4P+fkYJRJUdZZazXJM2i4mF7nUH', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAANwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbSeJV0AAAAAAAAAAYAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbSeJV0AAAAAAAAAAZAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lXQAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatlj11nHQAAAAAAAAABkAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3AAAAAAAAAAAaXI4hE2dLocMSSYYAT2ClklSctk2diyPO36ldXFH/1AAAAAJUC+QAAAAANwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAA2AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lY0AAAAAAAAABgAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA3AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lXQAAAAAAAAABgAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{QjDI1/cTt1yqUOrkv1NQanc0YUgWEN+Yr5mcTwYH6kLEw5mtQczrKv25hyeD/n5GCUSVHWWWs1yTNouJhe51Bw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('85bbd2b558563518a38e9b749bd4b8ced60b9fbbb7a6b283e15ae98548302ac4', 54, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 24, 200, 2, '2019-06-03 16:35:33.406028', '2019-06-03 16:35:33.406028', 231928238080, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAyAAAAAAAAAAYAAAAAAAAAAAAAAACAAAAAAAAAAEAAAAABJeTmKR1qr+CZoIyjAfGxrIXZ/tI1VId2OfZkRowDz4AAAAAAAAAAAX14QAAAAABAAAAAASXk5ikdaq/gmaCMowHxsayF2f7SNVSHdjn2ZEaMA8+AAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAABfXhAAAAAAAAAAACVvwF9wAAAEDRRWwMrdLrhnl+FIP+71tTHB5rlzCsPVyGnR3scvID9NmIL3LZEo992uTvDI9QLys5bC2yRc3WYR0vFiZRs40IGjAPPgAAAEDXbXWVdzmN6NWBjYU5OvB33WTUaa2wDZX3RmFTZQQ/+7JvPdblMtNCxo8IOYePQg90RajV9rB+k8P+SEpPHCUH', 'AAAAAAAAAMgAAAAAAAAAAgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAANgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbSeJWNAAAAAAAAAAXAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbSeJWNAAAAAAAAAAYAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAACAAAABAAAAAMAAAA1AAAAAAAAAAAEl5OYpHWqv4JmgjKMB8bGshdn+0jVUh3Y59mRGjAPPgAAAAJUC+QAAAAANQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA2AAAAAAAAAAAEl5OYpHWqv4JmgjKMB8bGshdn+0jVUh3Y59mRGjAPPgAAAAJaAcUAAAAANQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAA2AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lY0AAAAAAAAABgAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA2AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltD7HU0AAAAAAAAABgAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAQAAAADAAAANgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbQ+x1NAAAAAAAAAAYAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZbSeJWNAAAAAAAAAAYAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAADAAAANgAAAAAAAAAABJeTmKR1qr+CZoIyjAfGxrIXZ/tI1VId2OfZkRowDz4AAAACWgHFAAAAADUAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANgAAAAAAAAAABJeTmKR1qr+CZoIyjAfGxrIXZ/tI1VId2OfZkRowDz4AAAACVAvkAAAAADUAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAA1AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lb8AAAAAAAAABcAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA2AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lY0AAAAAAAAABcAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{0UVsDK3S64Z5fhSD/u9bUxwea5cwrD1chp0d7HLyA/TZiC9y2RKPfdrk7wyPUC8rOWwtskXN1mEdLxYmUbONCA==,1211lXc5jejVgY2FOTrwd91k1GmtsA2V90ZhU2UEP/uybz3W5TLTQsaPCDmHj0IPdEWo1fawfpPD/khKTxwlBw==}', 'none', NULL, NULL, true, 200);
INSERT INTO history_transactions VALUES ('df5f0e8b3b533dd9cda0ff7540bef3e9e19369060f8a4b0414b0e3c1b4315b1c', 53, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 23, 100, 1, '2019-06-03 16:35:33.420452', '2019-06-03 16:35:33.420452', 227633270784, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAXAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAABJeTmKR1qr+CZoIyjAfGxrIXZ/tI1VId2OfZkRowDz4AAAACVAvkAAAAAAAAAAABVvwF9wAAAEDyHwhW9GXQVXG1qibbeqSjxYzhv5IC08K2vSkxzYTwJykvQ8l0+e4M4h2guoK89s8HUfIqIOzDmoGsNTaLcYUG', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAANQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZdne46/AAAAAAAAAAWAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZdne46/AAAAAAAAAAXAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAA1AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl2d7jr8AAAAAAAAABcAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA1AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatltJ4lb8AAAAAAAAABcAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1AAAAAAAAAAAEl5OYpHWqv4JmgjKMB8bGshdn+0jVUh3Y59mRGjAPPgAAAAJUC+QAAAAANQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAwAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl2d7jtgAAAAAAAAABYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA1AAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl2d7jr8AAAAAAAAABYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{8h8IVvRl0FVxtaom23qko8WM4b+SAtPCtr0pMc2E8CcpL0PJdPnuDOIdoLqCvPbPB1HyKiDsw5qBrDU2i3GFBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c8a28fb25d4784f37a7a078e1feef0eb30ca64e994734625ac4ea067cc621464', 52, 1, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430214, 100, 1, '2019-06-03 16:35:33.430853', '2019-06-03 16:35:33.430854', 223338303488, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTEAAAAAAAABAAAABDAwMDAAAAAAAAAAAS6Z+xkAAABA3ExJNH79wGSRYZerPP1zMYlepMsuhoJF5vHn2gCsHmDpWfgO8VKC3BRImO+ne9spUXlVHMjEuhOHoPhl1hrMCg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAANAAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvhqAAAADAAAAAFAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAANAAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvhqAAAADAAAAAGAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAzAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMQAAAAAAAAQxMjM0AAAAAAAAAAAAAAABAAAANAAAAAMAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAAFbmFtZTEAAAAAAAAEMDAwMAAAAAAAAAAA', 'AAAAAgAAAAMAAAAzAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+IMAAAAMAAAAAUAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAA0AAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+GoAAAAMAAAAAUAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{3ExJNH79wGSRYZerPP1zMYlepMsuhoJF5vHn2gCsHmDpWfgO8VKC3BRImO+ne9spUXlVHMjEuhOHoPhl1hrMCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('1d7833c4faab08e62609acf3714d1babe27621a2b328edf37465e99aaf389cab', 51, 1, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430213, 100, 1, '2019-06-03 16:35:33.443723', '2019-06-03 16:35:33.443724', 219043336192, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTEAAAAAAAABAAAABDEyMzQAAAAAAAAAAS6Z+xkAAABAIW4yrFdk66fgDDir7YFATEd2llOubzx/iaJcM2wkF3ouqJQN+Aziy2rVtK5AoyphokiwsYXvHS6UF9MhdnUADQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMwAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAviDAAAADAAAAAEAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMwAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAviDAAAADAAAAAFAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMQAAAAAAAAQxMjM0AAAAAAAAAAAAAAABAAAAMwAAAAMAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAAFbmFtZTEAAAAAAAAEMTIzNAAAAAAAAAAA', 'AAAAAgAAAAMAAAAyAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+JwAAAAMAAAAAQAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAzAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+IMAAAAMAAAAAQAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{IW4yrFdk66fgDDir7YFATEd2llOubzx/iaJcM2wkF3ouqJQN+Aziy2rVtK5AoyphokiwsYXvHS6UF9MhdnUADQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('616c609047ef8f9ca908a47a47aa4bb018449c569549ad2ca60590aab74267e8', 50, 1, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430212, 100, 1, '2019-06-03 16:35:33.457305', '2019-06-03 16:35:33.457305', 214748368896, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTIAAAAAAAAAAAAAAAAAAAEumfsZAAAAQAYRZNPhJCTwjJgAJ9beE3ZO/H3kYJhYmV1pCmy7c8Zr2sKdKOmaLn4fmA5qaL+lQMKwOShtjwkZ8JHxPUd8GAk=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMgAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvicAAAADAAAAADAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMgAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvicAAAADAAAAAEAAAAAwAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMgAAAAAAAAQ1Njc4AAAAAAAAAAAAAAACAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMgAAAAAAAAMAAAAyAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+JwAAAAMAAAAAQAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAyAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+JwAAAAMAAAAAQAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAMAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAyAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+JwAAAAMAAAAAMAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{BhFk0+EkJPCMmAAn1t4Tdk78feRgmFiZXWkKbLtzxmvawp0o6Zoufh+YDmpov6VAwrA5KG2PCRnwkfE9R3wYCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('9fff61916716fb2550043fac968ac6c13802af5176a10fc29108fcfc445ef513', 49, 1, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430209, 100, 1, '2019-06-03 16:35:33.472183', '2019-06-03 16:35:33.472183', 210453401600, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTEAAAAAAAABAAAABDEyMzQAAAAAAAAAAS6Z+xkAAABAxKiHYYNLJiW3r5+kCJm8ucaoV7BcrEnQXFb3s1RyRyUbAkDlaCvE+RKwMZoNUfbkQUGrouyVKy1ZpUeccByqDg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMQAAAAAAAAQxMjM0AAAAAAAAAAA=', 'AAAAAgAAAAMAAAAwAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+QAAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+OcAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{xKiHYYNLJiW3r5+kCJm8ucaoV7BcrEnQXFb3s1RyRyUbAkDlaCvE+RKwMZoNUfbkQUGrouyVKy1ZpUeccByqDg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('e4609180751e7702466a8845857df43e4d154ec84b6bad62ce507fe12f1daf99', 49, 2, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430210, 100, 1, '2019-06-03 16:35:33.472323', '2019-06-03 16:35:33.472324', 210453405696, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZTIAAAAAAAABAAAABDU2NzgAAAAAAAAAAS6Z+xkAAABAjxgnTRBCa0n1efZocxpEjXeITQ5sEYTVd9fowuto2kPw5eFwgVnz6OrKJwCRt5L8ylmWiATXVI3Zyfi3yTKqBA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lMgAAAAAAAAQ1Njc4AAAAAAAAAAA=', 'AAAAAgAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+OcAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+M4AAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{jxgnTRBCa0n1efZocxpEjXeITQ5sEYTVd9fowuto2kPw5eFwgVnz6OrKJwCRt5L8ylmWiATXVI3Zyfi3yTKqBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('48415cd0fda9bc9aeb1f0b419bfb2997f7a2aa1b1ef2e51a0602c61104fc23cc', 49, 3, 'GAYSCMKQY6EYLXOPTT6JPPOXDMVNBWITPTSZIVWW4LWARVBOTH5RTLAD', 206158430211, 100, 1, '2019-06-03 16:35:33.472436', '2019-06-03 16:35:33.472436', 210453409792, 'AAAAADEhMVDHiYXdz5z8l73XGyrQ2RN85ZRW1uLsCNQumfsZAAAAZAAAADAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAoAAAAFbmFtZSAAAAAAAAABAAAAD2l0cyBnb3Qgc3BhY2VzIQAAAAAAAAAAAS6Z+xkAAABANmYginYhX+6VAsl1JumfxkB57y2LHraWDUkR+KDxWW8l5pfTViLxx7J85KrOV0qNCY4RfasgqxF0FC3ErYceCQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAACAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMQAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvi1AAAADAAAAADAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAMAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAMAAAADAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAxAAAAAwAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAVuYW1lIAAAAAAAAA9pdHMgZ290IHNwYWNlcyEAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+M4AAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAxAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+LUAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{NmYginYhX+6VAsl1JumfxkB57y2LHraWDUkR+KDxWW8l5pfTViLxx7J85KrOV0qNCY4RfasgqxF0FC3ErYceCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('eb8586c9176c4cf2e864b2521948a972db5274de24673669463e0c7824cee056', 48, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 22, 100, 1, '2019-06-03 16:35:33.482587', '2019-06-03 16:35:33.482587', 206158434304, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAWAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAMSExUMeJhd3PnPyXvdcbKtDZE3zllFbW4uwI1C6Z+xkAAAACVAvkAAAAAAAAAAABVvwF9wAAAECAMOn6G4jusgpfSoHwntHQkYIDxI/VnyH/qIi+bdMWzi1T6WlwnO+yITgm2+mOaWc6zVuxiLjHllzBeQ/xKvQN', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAMAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZf8fofYAAAAAAAAAAVAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAMAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGrZf8fofYAAAAAAAAAAWAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAwAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl/x+h9gAAAAAAAAABYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAwAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl2d7jtgAAAAAAAAABYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAAAAAAAAxITFQx4mF3c+c/Je91xsq0NkTfOWUVtbi7AjULpn7GQAAAAJUC+QAAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl/x+h/EAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAwAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl/x+h9gAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{gDDp+huI7rIKX0qB8J7R0JGCA8SP1Z8h/6iIvm3TFs4tU+lpcJzvsiE4JtvpjmlnOs1bsYi4x5ZcwXkP8Sr0DQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('ea93efd8c2f4e45c0318c69ec958623a0e4374f40d569eec124d43c8a54d6256', 47, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 21, 100, 1, '2019-06-03 16:35:33.495412', '2019-06-03 16:35:33.495413', 201863467008, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAVAAAAAAAAAAAAAAABAAAAAAAAAAkAAAAAAAAAAVb8BfcAAABABUHuXY+MTgW/wDv5+NDVh9fw4meszxeXO98HEQfgXVeCZ7eObCI2orSGUNA/SK6HV9/uTVSxIQQWIso1QoxHBQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAJAAAAAAAAAAIAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAIrEjCYwXAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAIrEjfceLAAAAAA==', 'AAAAAQAAAAIAAAADAAAALwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvaAAAAAAAAAAUAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvaAAAAAAAAAAVAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+9oAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsatl/x+h/EAAAAAAAAABUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAuAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7E/+cAAAALQAAAAEAAAAAAAAAAQAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGraHekccnAAAALQAAAAEAAAAAAAAAAQAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAuAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+/MAAAAAAAAABQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAvAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+9oAAAAAAAAABQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{BUHuXY+MTgW/wDv5+NDVh9fw4meszxeXO98HEQfgXVeCZ7eObCI2orSGUNA/SK6HV9/uTVSxIQQWIso1QoxHBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('7207de5b75243e0b062c3833f587036b7e9f64453be49ff50f3f3fdc7516ec6b', 46, 1, 'GDR53WAEIKOU3ZKN34CSHAWH7HV6K63CBJRUTWUDBFSMY7RRQK3SPKOS', 193273528321, 100, 1, '2019-06-03 16:35:33.506693', '2019-06-03 16:35:33.506693', 197568499712, 'AAAAAOPd2ARCnU3lTd8FI4LH+evle2IKY0nagwlkzH4xgrcnAAAAZAAAAC0AAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAABAAAAAOPd2ARCnU3lTd8FI4LH+evle2IKY0nagwlkzH4xgrcnAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABMYK3JwAAAEAOkGOPTOBDSQ7nW2Zn+bls2PDUebk2/k3/gqHKQ8eYOFsD6nBeEvyMD858vo5BabjQwB9injABIM8esDh7bEkC', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAALgAAAAAAAAAA493YBEKdTeVN3wUjgsf56+V7YgpjSdqDCWTMfjGCtycCxorwuxP/nAAAAC0AAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALgAAAAAAAAAA493YBEKdTeVN3wUjgsf56+V7YgpjSdqDCWTMfjGCtycCxorwuxP/nAAAAC0AAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAuAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7E/+cAAAALQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAuAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7E/+cAAAALQAAAAEAAAAAAAAAAQAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAtAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7FAAAAAAALQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAuAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7E/+cAAAALQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{DpBjj0zgQ0kO51tmZ/m5bNjw1Hm5Nv5N/4KhykPHmDhbA+pwXhL8jA/OfL6OQWm40MAfYp4wASDPHrA4e2xJAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d24f486bd722fd1875b843839e880bdeea324e25db706a26af5e4daa8c5071eb', 46, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 20, 100, 1, '2019-06-03 16:35:33.506856', '2019-06-03 16:35:33.506856', 197568503808, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAUAAAAAAAAAAAAAAABAAAAAAAAAAUAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABVvwF9wAAAEBYI0TMQVWPvnC2KPbDph9Myz5UMuBRIYt2YQdtlPYC4UHamYnHsMghpIMfaS7MWdHuGY81+FBozOsS+/HGohQD', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAALgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvzAAAAAAAAAATAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcLGiubZdPvzAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAuAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+/MAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAuAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+/MAAAAAAAAABQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAtAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0/AwAAAAAAAAABMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAuAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0+/MAAAAAAAAABMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{WCNEzEFVj75wtij2w6YfTMs+VDLgUSGLdmEHbZT2AuFB2pmJx7DIIaSDH2kuzFnR7hmPNfhQaMzrEvvxxqIUAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('5b42c77042f04bf716659a05e7ca3f4703af038a7da75b10b8538707c9ff172f', 45, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 19, 100, 1, '2019-06-03 16:35:33.518779', '2019-06-03 16:35:33.51878', 193273532416, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAATAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA493YBEKdTeVN3wUjgsf56+V7YgpjSdqDCWTMfjGCtycCxorwuxQAAAAAAAAAAAABVvwF9wAAAECGClRePcAExQ/WKroo3/3dfchP/yI8TRDrrjt/chZ83ULiTc54l5wcz1AkbLa6CAapdSGpUWXk5ksTqDXLn4AA', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAALQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaMIOfwMAAAAAAAAAASAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaMIOfwMAAAAAAAAAATAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAtAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/AwAAAAAAAAABMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAtAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wsaK5tl0/AwAAAAAAAAABMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAtAAAAAAAAAADj3dgEQp1N5U3fBSOCx/nr5XtiCmNJ2oMJZMx+MYK3JwLGivC7FAAAAAAALQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAsAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/CUAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAtAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/AwAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{hgpUXj3ABMUP1iq6KN/93X3IT/8iPE0Q6647f3IWfN1C4k3OeJecHM9QJGy2uggGqXUhqVFl5OZLE6g1y5+AAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('e0773d07aba23d11e6a06b021682294be1f9f202a2926827022539662ce2c7fc', 44, 1, 'GCHPXGVDKPF5KT4CNAT7X77OXYZ7YVE4JHKFDUHCGCVWCL4K4PQ67KKZ', 184683593729, 100, 1, '2019-06-03 16:35:33.529055', '2019-06-03 16:35:33.529055', 188978565120, 'AAAAAI77mqNTy9VPgmgn+//uvjP8VJxJ1FHQ4jCrYS+K4+HvAAAAZAAAACsAAAABAAAAAAAAAAAAAAABAAAAAAAAAAgAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAYrj4e8AAABA3jJ7wBrRpsrcnqBQWjyzwvVz2v5UJ56G60IhgsaWQFSf+7om462KToc+HJ27aLVOQ83dGh1ivp+VIuREJq/SBw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAIAAAAAAAAAAJUC+OcAAAAAA==', 'AAAAAQAAAAIAAAADAAAALAAAAAAAAAAAjvuao1PL1U+CaCf7/+6+M/xUnEnUUdDiMKthL4rj4e8AAAACVAvjnAAAACsAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAALAAAAAAAAAAAjvuao1PL1U+CaCf7/+6+M/xUnEnUUdDiMKthL4rj4e8AAAACVAvjnAAAACsAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAArAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtonM3Az4AAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAsAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/CUAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAsAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+OcAAAAKwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAI77mqNTy9VPgmgn+//uvjP8VJxJ1FHQ4jCrYS+K4+Hv', 'AAAAAgAAAAMAAAArAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+QAAAAAKwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAsAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+OcAAAAKwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{3jJ7wBrRpsrcnqBQWjyzwvVz2v5UJ56G60IhgsaWQFSf+7om462KToc+HJ27aLVOQ83dGh1ivp+VIuREJq/SBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('945b6171de747ab323b3cda52290933df39edd7061f6e260762663efc51bccb0', 43, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 18, 100, 1, '2019-06-03 16:35:33.541377', '2019-06-03 16:35:33.541377', 184683597824, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAASAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAjvuao1PL1U+CaCf7/+6+M/xUnEnUUdDiMKthL4rj4e8AAAACVAvkAAAAAAAAAAABVvwF9wAAAEBFbS2c5rrYNGslNVslTHH8j8x0ggew1eHHOUTNajMPy8GYn52RSwRncwwvv1ejEfA+g/mTXMpXrBO847C46KoA', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaMIOfw+AAAAAAAAAARAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaMIOfw+AAAAAAAAAASAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAArAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/D4AAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAArAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtonM3Az4AAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAArAAAAAAAAAACO+5qjU8vVT4JoJ/v/7r4z/FScSdRR0OIwq2EviuPh7wAAAAJUC+QAAAAAKwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/FcAAAAAAAAABEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAArAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/D4AAAAAAAAABEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{RW0tnOa62DRrJTVbJUxx/I/MdIIHsNXhxzlEzWozD8vBmJ+dkUsEZ3MML79XoxHwPoP5k1zKV6wTvOOwuOiqAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d67cfb271a889e7854ffd61b08eacde76d56e758466fc37a8eec2d3a40ef8b14', 42, 1, 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF', 163208757252, 100, 1, '2019-06-03 16:35:33.552965', '2019-06-03 16:35:33.552965', 180388630528, 'AAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAZAAAACYAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAABRVVSAAAAAAAAAAAAAAAAAUpI8/gAAABAEPKcQmATGpevrtlAcZnNI/GjfLLQEp9aODGGRFV+2C4UO8dU+UAMTkCSXQLD+xPaRQxzw93ScEok6GzYCtt7Bg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKgAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvicAAAACYAAAADAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKgAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvicAAAACYAAAAEAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAApAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFFVVIAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAqAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFFVVIAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAApAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+LUAAAAJgAAAAMAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAqAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+JwAAAAJgAAAAMAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{EPKcQmATGpevrtlAcZnNI/GjfLLQEp9aODGGRFV+2C4UO8dU+UAMTkCSXQLD+xPaRQxzw93ScEok6GzYCtt7Bg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6d2e30fd57492bf2e2b132e1bc91a548a369189bebf77eb2b3d829121a9d2c50', 41, 1, 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF', 163208757250, 100, 1, '2019-06-03 16:35:33.565164', '2019-06-03 16:35:33.565164', 176093663232, 'AAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAZAAAACYAAAACAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAABVVNEAAAAAAEAAAAAAAAAAUpI8/gAAABA6O2fe1gQBwoO0fMNNEUKH0QdVXVjEWbN5VL51DmRUedYMMXtbX5JKVSzla2kIGvWgls1dXuXHZY/IOlaK01rBQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAABAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAACAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAnAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+OcAAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+M4AAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{6O2fe1gQBwoO0fMNNEUKH0QdVXVjEWbN5VL51DmRUedYMMXtbX5JKVSzla2kIGvWgls1dXuXHZY/IOlaK01rBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('a832ff67085cb9eb6f1c4b740f6e033ba9b508af725fbf203469729a64a199ff', 41, 2, 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF', 163208757251, 100, 1, '2019-06-03 16:35:33.565428', '2019-06-03 16:35:33.565429', 176093667328, 'AAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAZAAAACYAAAADAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAABRVVSAAAAAAEAAAAAAAAAAUpI8/gAAABA1Qe8ngwANz4fLqYChwRjR5xng6cIqU5WBtjkZgF4ugVhi8J6kTpACvnvXso3IVym6Rfd6JdQW8QcLkFTX1MGCg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAACAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKQAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvi1AAAACYAAAADAAAAAAAAAAAAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFFVVIAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFFVVIAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAApAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+M4AAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAApAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+LUAAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{1Qe8ngwANz4fLqYChwRjR5xng6cIqU5WBtjkZgF4ugVhi8J6kTpACvnvXso3IVym6Rfd6JdQW8QcLkFTX1MGCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6fa467b53f5386d77ad35c2502ed2cd3dd8b460a5be22b6b2818b81bcd3ed2da', 40, 1, 'GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG', 163208757249, 100, 1, '2019-06-03 16:35:33.580485', '2019-06-03 16:35:33.580485', 171798695936, 'AAAAAKturFHJX/eRt5gM6qIXAMbaXvlImqLysA6Qr9tLemxfAAAAZAAAACYAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+H//////////AAAAAAAAAAFLemxfAAAAQKN8LftAafeoAGmvpsEokqm47jAuqw4g1UWjmL0j6QPm1jxoalzDwDS3W+N2HOHdjSJlEQaTxGBfQKHhr6nNsAA=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFVU0QAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAmAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+QAAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+OcAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{o3wt+0Bp96gAaa+mwSiSqbjuMC6rDiDVRaOYvSPpA+bWPGhqXMPANLdb43Yc4d2NImURBpPEYF9AoeGvqc2wAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('0bcb67aa365446fd244fecff3a0c397f81f3a9b13428688965e776d447c0b1ea', 40, 2, 'GCVW5LCRZFP7PENXTAGOVIQXADDNUXXZJCNKF4VQB2IK7W2LPJWF73UG', 163208757250, 100, 1, '2019-06-03 16:35:33.58064', '2019-06-03 16:35:33.58064', 171798700032, 'AAAAAKturFHJX/eRt5gM6qIXAMbaXvlImqLysA6Qr9tLemxfAAAAZAAAACYAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABRVVSAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+H//////////AAAAAAAAAAFLemxfAAAAQMPVgYf+w09depDSxMcJnjVZHA2FlkBmhPmi0N66FuhAzTekWcCOMdCI0cUc+xJhywLXSMiKA6wP6K94NRlFlQE=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAKAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvjOAAAACYAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAQAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAFFVVIAAAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+OcAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAoAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+M4AAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{w9WBh/7DT116kNLExwmeNVkcDYWWQGaE+aLQ3roW6EDNN6RZwI4x0IjRxRz7EmHLAtdIyIoDrA/or3g1GUWVAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('11705f94cd65d7b673a124a85ce368c80f8458ffaedff719304d8f849535b4e0', 39, 1, 'GD4SMOE3VPSF7ZR3CTEQ3P5UNTBMEJDA2GLXTHR7MMARANKKJDZ7RPGF', 163208757249, 100, 1, '2019-06-03 16:35:33.5929', '2019-06-03 16:35:33.5929', 167503728640, 'AAAAAPkmOJur5F/mOxTJDb+0bMLCJGDRl3meP2MBEDVKSPP4AAAAZAAAACYAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABSkjz+AAAAECyjDa1e+jtXukTrHluO7x0Mx7Wj4mRoM4S5UAFmRV+2rVoxjMwqFJhtYnEAUV19+C5ycp5jOLLpWxrCeRKJQUG', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAJwAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvjnAAAACYAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJwAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvjnAAAACYAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAnAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+OcAAAAJgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAnAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+OcAAAAJgAAAAEAAAAAAAAAAAAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAmAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+QAAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAnAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+OcAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{sow2tXvo7V7pE6x5bju8dDMe1o+JkaDOEuVABZkVftq1aMYzMKhSYbWJxAFFdffgucnKeYziy6VsawnkSiUFBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('afeb8080522eba71ca328225bbcf731029edcfa254c827c45be580bae95c7231', 38, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 16, 100, 1, '2019-06-03 16:35:33.603781', '2019-06-03 16:35:33.603781', 163208761344, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAQAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAq26sUclf95G3mAzqohcAxtpe+UiaovKwDpCv20t6bF8AAAACVAvkAAAAAAAAAAABVvwF9wAAAEDnzvNgEYB1u3BGTHFDlIWnk0GOq7BMpfcyewJRsJK9lT4HTMEwMQ2jSJyrWmB7xdBxHKaNMXQaAIx6CShLXpQH', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAJgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaQyP+5XAAAAAAAAAAPAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaQyP+5XAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7lcAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gto5089VcAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAmAAAAAAAAAACrbqxRyV/3kbeYDOqiFwDG2l75SJqi8rAOkK/bS3psXwAAAAJUC+QAAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAhAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7okAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7nAAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{587zYBGAdbtwRkxxQ5SFp5NBjquwTKX3MnsCUbCSvZU+B0zBMDENo0icq1pge8XQcRymjTF0GgCMegkoS16UBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2354df802111418a999e31c2964d16b8efe8e492b7d74de54939825190e1041f', 38, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 17, 100, 1, '2019-06-03 16:35:33.603995', '2019-06-03 16:35:33.603995', 163208765440, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAARAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA+SY4m6vkX+Y7FMkNv7RswsIkYNGXeZ4/YwEQNUpI8/gAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDD6WvAYL1wilsd7zYDJt0iFO/lppQ6GJJn/A8UJl9jTjMNOjuQPBtA7fSxR5KT0BZLbtQy8qFlys0I6fTe/cwO', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAJgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaOdPPVXAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaOdPPVXAAAAAAAAAARAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gto5089VcAAAAAAAAABEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtowg5/FcAAAAAAAAABEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAmAAAAAAAAAAD5Jjibq+Rf5jsUyQ2/tGzCwiRg0Zd5nj9jARA1Skjz+AAAAAJUC+QAAAAAJgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7nAAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAmAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7lcAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{w+lrwGC9cIpbHe82AybdIhTv5aaUOhiSZ/wPFCZfY04zDTo7kDwbQO30sUeSk9AWS27UMvKhZcrNCOn03v3MDg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('52388a98e4e36c17749a94374270cc65bdb7271cb51277f095aaa8f1ca9d322c', 37, 1, 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG', 141733920772, 100, 1, '2019-06-03 16:35:33.615811', '2019-06-03 16:35:33.615811', 158913794048, 'AAAAAJKk3eoJUHc7fO9texuiGHN4NwWMTIFhj6Q30T4ctWW1AAAAZAAAACEAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAAAAAAAAAAEctWW1AAAAQM5SCoW10EJoKBBwwMu0Vw+f+bQ0GjQ9FO6w3l9Q/FIctm87248t9jXTbl0Rd4NgGcom0yoGxgcJiERwZGBMXQc=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAJQAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvicAAAACEAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJQAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvicAAAACEAAAAEAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAkAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAAAAAAAO5rKAAAAAAEAAAAAAAAAAAAAAAIAAAABAAAAAJKk3eoJUHc7fO9texuiGHN4NwWMTIFhj6Q30T4ctWW1AAAAAVVTRAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAJQAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvicAAAACEAAAAEAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJQAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvicAAAACEAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAkAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+LUAAAAIQAAAAMAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAlAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+JwAAAAIQAAAAMAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{zlIKhbXQQmgoEHDAy7RXD5/5tDQaND0U7rDeX1D8Uhy2bzvbjy32NdNuXRF3g2AZyibTKgbGBwmIRHBkYExdBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('44cb6c8ed4dbec542af1aad23001dd9d678cf19c8c461a653e762a7253eded82', 36, 1, 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG', 141733920771, 100, 1, '2019-06-03 16:35:33.635753', '2019-06-03 16:35:33.635753', 154618826752, 'AAAAAJKk3eoJUHc7fO9texuiGHN4NwWMTIFhj6Q30T4ctWW1AAAAZAAAACEAAAADAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAA7msoAAAAAAAAAAAEctWW1AAAAQO+eTIPXUZk+GAq7O6H8d1/WT5buo0apjLhGgtBeSyl37UV7LCpZfCn6DYVc7lQOVNWhBc7KDA7Ne83AR41kYAk=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAJAAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvi1AAAACEAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAJAAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvi1AAAACEAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAjAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAAAAAAAO5rKAAAAAAEAAAAAAAAAAAAAAAEAAAAkAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAAAAAAAO5rKAAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAjAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+M4AAAAIQAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAkAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+LUAAAAIQAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{755Mg9dRmT4YCrs7ofx3X9ZPlu6jRqmMuEaC0F5LKXftRXssKll8KfoNhVzuVA5U1aEFzsoMDs17zcBHjWRgCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4e2442fe2e8dd8c686570c9f537acb2f50153a9883f8d199b6f4701eb289b3a0', 35, 1, 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG', 141733920770, 100, 1, '2019-06-03 16:35:33.645263', '2019-06-03 16:35:33.645263', 150323859456, 'AAAAAJKk3eoJUHc7fO9texuiGHN4NwWMTIFhj6Q30T4ctWW1AAAAZAAAACEAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAA7msoAAAAAAAAAAAEctWW1AAAAQNugq+B30pdbzvVVGz9RO3+DMeRdWqc/Xsd2NYdg6NBu7esvOdTWQ3nvoBEJyeGz8EE9zRQiSiqorwHlm+AGfwI=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAIwAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvjOAAAACEAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAIwAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvjOAAAACEAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAiAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAjAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAAAAAAAO5rKAAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAiAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+OcAAAAIQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAjAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+M4AAAAIQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{26Cr4HfSl1vO9VUbP1E7f4Mx5F1apz9ex3Y1h2Do0G7t6y851NZDee+gEQnJ4bPwQT3NFCJKKqivAeWb4AZ/Ag==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('a05daae230b1f743474e83ab6d4817df1f3f77661a7d815f7620cee2a9809480', 34, 1, 'GCJKJXPKBFIHOO3455WXWG5CDBZXQNYFRRGICYMPUQ35CPQ4WVS3KZLG', 141733920769, 100, 1, '2019-06-03 16:35:33.656985', '2019-06-03 16:35:33.656985', 146028892160, 'AAAAAJKk3eoJUHc7fO9texuiGHN4NwWMTIFhj6Q30T4ctWW1AAAAZAAAACEAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF93//////////AAAAAAAAAAEctWW1AAAAQBYUnV3I1O35EAyay0msjg3MzZfanCtvalKGG+94pe6RxgE/kCk2kTT9HXgXjbraq//Q/0vJ0AoCAXSeT18Ujgk=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAIgAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvjnAAAACEAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAIgAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvjnAAAACEAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAiAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+OcAAAAIQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAiAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+OcAAAAIQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAiAAAAAQAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAFVU0QAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAhAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+QAAAAAIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAiAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+OcAAAAIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{FhSdXcjU7fkQDJrLSayODczNl9qcK29qUoYb73il7pHGAT+QKTaRNP0deBeNutqr/9D/S8nQCgIBdJ5PXxSOCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6d78f17fafa2317d6af679e1e5420f351207ff61cdff21c600ea8f85155b3ea1', 33, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 15, 100, 1, '2019-06-03 16:35:33.665712', '2019-06-03 16:35:33.665712', 141733924864, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAPAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAkqTd6glQdzt87217G6IYc3g3BYxMgWGPpDfRPhy1ZbUAAAACVAvkAAAAAAAAAAABVvwF9wAAAEC+mgKIzZqflQIKIqWn9LrciuyEx7XPfXGUhvyQ3sIQBnGdOWhkOt57UU/75LtUy4recT+jrY2cHKZj33puue8F', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAIQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaTHQueJAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAIQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaTHQueJAAAAAAAAAAPAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAhAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpMdC54kAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAhAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpDI/7okAAAAAAAAAA8AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAhAAAAAAAAAACSpN3qCVB3O3zvbXsbohhzeDcFjEyBYY+kN9E+HLVltQAAAAJUC+QAAAAAIQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAbAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpMdC56IAAAAAAAAAA4AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAhAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpMdC54kAAAAAAAAAA4AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{vpoCiM2an5UCCiKlp/S63IrshMe1z31xlIb8kN7CEAZxnTloZDree1FP++S7VMuK3nE/o62NnBymY996brnvBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('bb9d6654111fae501594400dc901c70d47489a67163d2a34f9b3e32a921a50dc', 32, 1, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964117003, 100, 1, '2019-06-03 16:35:33.677436', '2019-06-03 16:35:33.677436', 137438957568, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAALAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAMAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABMJ4wUAAAAEAFytUxjxN4bnJMrEJkSprnES9iGpOxAsNOFYrTP/xtGVk/PZ2oThUW+/hLRIk+hYYEgF21Gf58N/abJKFpqlsI', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAIAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvfUAAAABsAAAAKAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAABQAAAAAAAAAAAAAAAQAAACAAAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL31AAAAAbAAAACwAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAUAAAAAAAAAAAAAAAEAAAACAAAAAwAAACAAAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL31AAAAAbAAAACwAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAUAAAAAAAAAAAAAAAEAAAAgAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC99QAAAAGwAAAAsAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAAFAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAfAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+AYAAAAGwAAAAoAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAAFAAAAAAAAAAAAAAABAAAAIAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvftAAAABsAAAAKAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAABQAAAAAAAAAA', '{BcrVMY8TeG5yTKxCZEqa5xEvYhqTsQLDThWK0z/8bRlZPz2dqE4VFvv4S0SJPoWGBIBdtRn+fDf2myShaapbCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6b38cdd5c17df2013d5a5e211c4b32218b6be91025316b1aab28bc12316615d5', 32, 2, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964117004, 100, 1, '2019-06-03 16:35:33.677655', '2019-06-03 16:35:33.677656', 137438961664, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAMAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAAAAAAAAAAAATCeMFAAAABAOb0qGWnk1WrSUXS6iQFocaIOY/BDmgG1zTmlPyg0boSid3jTBK3z9U8+IPGAOELNLgkQHtgGYFgFGMio1xY+BQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAIAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvfUAAAABsAAAALAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAABQAAAAAAAAAAAAAAAQAAACAAAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL31AAAAAbAAAADAAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAUAAAAAAAAAAAAAAAEAAAACAAAAAwAAACAAAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL31AAAAAbAAAADAAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAAAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAUAAAAAAAAAAAAAAAEAAAAgAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC99QAAAAGwAAAAwAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAALZXhhbXBsZS5jb20AAgACAgAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAgAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC9+0AAAAGwAAAAoAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAAFAAAAAAAAAAAAAAABAAAAIAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvfUAAAABsAAAAKAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAABQAAAAAAAAAA', '{Ob0qGWnk1WrSUXS6iQFocaIOY/BDmgG1zTmlPyg0boSid3jTBK3z9U8+IPGAOELNLgkQHtgGYFgFGMio1xY+BQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('299dc6631d585a55ae3602f660ec5b5a0088d24a14b344c72eccc2a62d9a8938', 31, 1, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964117002, 100, 1, '2019-06-03 16:35:33.694359', '2019-06-03 16:35:33.694359', 133143990272, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAKAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAUAAAAAAAAAATCeMFAAAABA0wriernSr+5P2QCeon1uj5mrOLNTOrPYPPi5ricLug/nreEUhsgS/k3lA9JGpVbd+tacMEKmXKmFxHCEMjWPBg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHwAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvgGAAAABsAAAAJAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAAAAAAAQAAAB8AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4BgAAAAbAAAACgAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAACAAAAAwAAAB8AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4BgAAAAbAAAACgAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAAfAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+AYAAAAGwAAAAoAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAAFAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAeAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+B8AAAAGwAAAAkAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAAAAAABAAAAHwAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvgGAAAABsAAAAJAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAA', '{0wriernSr+5P2QCeon1uj5mrOLNTOrPYPPi5ricLug/nreEUhsgS/k3lA9JGpVbd+tacMEKmXKmFxHCEMjWPBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c3cd47a311e025446f72c50426b5b5444e5261431fc5760e8e57467c87cd49fc', 30, 1, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964117001, 100, 1, '2019-06-03 16:35:33.708736', '2019-06-03 16:35:33.708737', 128849022976, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAJAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAATCeMFAAAABA7ZMKq80ucQSt+55q+6VQrG3Hrv6zHtOLwkfAxxsZdYPIuk7xZsgbyhOCVXjheOQ9ygAW1vtybdXG41AxSFRtAg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHgAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvgfAAAABsAAAAIAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAAAAAAAQAAAB4AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4HwAAAAbAAAACQAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAACAAAAAwAAAB4AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4HwAAAAbAAAACQAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAAeAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+B8AAAAGwAAAAkAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAdAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+DgAAAAGwAAAAgAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAAAAAABAAAAHgAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvgfAAAABsAAAAIAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAA', '{7ZMKq80ucQSt+55q+6VQrG3Hrv6zHtOLwkfAxxsZdYPIuk7xZsgbyhOCVXjheOQ9ygAW1vtybdXG41AxSFRtAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('69f64ae0f809b08996c1f394ee795001a40eee69adb675ab63bfd1932d3aafb2', 29, 1, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964117000, 100, 1, '2019-06-03 16:35:33.718611', '2019-06-03 16:35:33.718612', 124554055680, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAEAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATCeMFAAAABAi69qDHclVS9A8GAaqyk6oIxiMC2KXXEneFijfxH5VyLGIQZNAxOOcCPpIalU6P1pYRX3K4OlKHZ4hIdxJzD6BQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHQAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvg4AAAABsAAAAHAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAAAAAAAQAAAB0AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4OAAAAAbAAAACAAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAACAAAAAwAAAB0AAAAAAAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAAlQL4OAAAAAbAAAACAAAAAEAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAwAAAAtleGFtcGxlLmNvbQACAAICAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAAAAAAAEAAAAdAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+DgAAAAGwAAAAgAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAcAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAAAAAABAAAAHQAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvg4AAAABsAAAAHAAAAAQAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAABAAAAAHyd2ctWIE8qo8kBxI2FnNzaYft9DtG3MbIUDyYWQxOQAAAAAQAAAAAAAAAA', '{i69qDHclVS9A8GAaqyk6oIxiMC2KXXEneFijfxH5VyLGIQZNAxOOcCPpIalU6P1pYRX3K4OlKHZ4hIdxJzD6BQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('fe3707fbd5c844395c598f31dc719c61218d4cea4e8dddadb6733f4866089100', 28, 1, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116993, 100, 1, '2019-06-03 16:35:33.729669', '2019-06-03 16:35:33.729669', 120259088384, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAABAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABMJ4wUAAAAEA/GIgE9sYPGwbCiIdLdhoEu25CyB0ZAcmjQonQItu6SE0gaSBVT/le355A/dw1NPaoXY9P/u0ou9D7h5Vb1fcK', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAEAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAbAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+QAAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+OcAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{PxiIBPbGDxsGwoiHS3YaBLtuQsgdGQHJo0KJ0CLbukhNIGkgVU/5Xt+eQP3cNTT2qF2PT/7tKLvQ+4eVW9X3Cg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('345ef7f85c6ea297e3f994feef279b63812628681bd173a1f615185a4368e482', 28, 2, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116994, 100, 1, '2019-06-03 16:35:33.729884', '2019-06-03 16:35:33.729884', 120259092480, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAACAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABMJ4wUAAAAEDYxq3zpaFIC2JcuJUbrQ3MFXzqvu+5G7XUi4NnHlfbLutn76ylQcjuwLgbUG2lqcQfl75doPUZyurKtFP1rkMO', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAACAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAIAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAIAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAEAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+OcAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+M4AAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{2Mat86WhSAtiXLiVG60NzBV86r7vuRu11IuDZx5X2y7rZ++spUHI7sC4G1BtpanEH5e+XaD1GcrqyrRT9a5DDg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2a14735d7b05109359444acdd87e7fe92c98e9295d2ba61b05e25d1f7ee10fd3', 28, 3, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116995, 100, 1, '2019-06-03 16:35:33.730069', '2019-06-03 16:35:33.73007', 120259096576, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAADAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABMJ4wUAAAAEAKuQ1exMu8hdf8dOPeULX2DG7DZx5WWIUFHXJMWGG9KmVrQoZDt2S6a/1uYEVJnvvY/EoJM5RpVjh2ZCs30VYA', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAACAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAABAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAADAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAABAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAMAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAEAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAMAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+M4AAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+LUAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{CrkNXsTLvIXX/HTj3lC19gxuw2ceVliFBR1yTFhhvSpla0KGQ7dkumv9bmBFSZ772PxKCTOUaVY4dmQrN9FWAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4f9598206ab17cf27b5c3eb9e906d63ebee2626654112eabdd2bce7bf12cccf2', 28, 4, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116996, 100, 1, '2019-06-03 16:35:33.730256', '2019-06-03 16:35:33.730256', 120259100672, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAEAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATCeMFAAAABAAd6MzHDjUdRtHozzDnD3jJA+uRDCar3PQtuH/43pnROzk1HkovJPQ1YyzcpOb/NeuU/LKNzseL0PJNasVX1lAQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAADAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAEAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAQAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAgAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+LUAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+JwAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Ad6MzHDjUdRtHozzDnD3jJA+uRDCar3PQtuH/43pnROzk1HkovJPQ1YyzcpOb/NeuU/LKNzseL0PJNasVX1lAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('852ba25e0e4aa149a22dc193bcb645ae9eba23e7f7432707f3b910474e9b6a5b', 28, 5, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116997, 100, 1, '2019-06-03 16:35:33.730474', '2019-06-03 16:35:33.730474', 120259104768, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAACAAAAAQAAAAIAAAAAAAAAAAAAAAAAAAABMJ4wUAAAAEAnFzc6kqweyIL4TzIDbr+8GUOGGs1W5jcX5iSNw4DeonzQARlejYJ9NOn/XkrcoC9Hvd8hc5lNx+1h991GxJUJ', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAEAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAFAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAUAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAgACAgAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+JwAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+IMAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Jxc3OpKsHsiC+E8yA26/vBlDhhrNVuY3F+YkjcOA3qJ80AEZXo2CfTTp/15K3KAvR73fIXOZTcftYffdRsSVCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('8ccc0c28c3e99a63cc59bad7dec3f5c56eb3942c548ecd40bc39c509d6b081d4', 28, 6, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116998, 100, 1, '2019-06-03 16:35:33.730663', '2019-06-03 16:35:33.730664', 120259108864, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAC2V4YW1wbGUuY29tAAAAAAAAAAAAAAAAATCeMFAAAABAkID6CkBHP9eovLQXkMQJ7QkE6NWlmdKGmLxaiI1YaVKZaKJxz5P85x+6wzpYxxbs6Bd2l4qxVjS7Q36DwRiqBA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAFAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAIAAgIAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAGAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAAAIAAgIAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAAAAgACAgAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAYAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+IMAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+GoAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{kID6CkBHP9eovLQXkMQJ7QkE6NWlmdKGmLxaiI1YaVKZaKJxz5P85x+6wzpYxxbs6Bd2l4qxVjS7Q36DwRiqBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('83201014880073f8eff6f21ae76e51c2c4faf533e550ecd3c2205b48a092960a', 28, 7, 'GCIFFRQKHMH6JD7CK5OI4XVCYCMNRNF6PYA7JTCR3FPHPJZQTYYFB5ES', 115964116999, 100, 1, '2019-06-03 16:35:33.730949', '2019-06-03 16:35:33.73095', 120259112960, 'AAAAAJBSxgo7D+SP4ldcjl6iwJjYtL5+AfTMUdled6cwnjBQAAAAZAAAABsAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAB8ndnLViBPKqPJAcSNhZzc2mH7fQ7RtzGyFA8mFkMTkAAAAAEAAAAAAAAAATCeMFAAAABAtYtlsqMReQo1UoU2GYjb3h52wEKvnouCSO6LQO1xm/ArhtQO/sX5q35St8BjaYWEiFnp+SQj2FZC89OswCldAw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAGAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAAAAAAAAAAAAAAAAAABAAAAHAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvhRAAAABsAAAAHAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAADAAAAC2V4YW1wbGUuY29tAAIAAgIAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAcAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAcAAAABAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAMAAAALZXhhbXBsZS5jb20AAgACAgAAAAEAAAAAfJ3Zy1YgTyqjyQHEjYWc3Nph+30O0bcxshQPJhZDE5AAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+GoAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAcAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+FEAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{tYtlsqMReQo1UoU2GYjb3h52wEKvnouCSO6LQO1xm/ArhtQO/sX5q35St8BjaYWEiFnp+SQj2FZC89OswCldAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('700fa44bb40e6ad2c5888656cd2e7b8d86de3d3557b653ae6874466175d64927', 27, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 14, 100, 1, '2019-06-03 16:35:33.756542', '2019-06-03 16:35:33.756543', 115964121088, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAOAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAkFLGCjsP5I/iV1yOXqLAmNi0vn4B9MxR2V53pzCeMFAAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBq3GPDVeRPfwqtW45GZNiUdQ9j6E9Nsz/lMYWcWDWGCZADSsEiEoXar1HWFK6drptsGEl9P6I9f7C2GBKb4YQM', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAGwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaVcReCiAAAAAAAAAANAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAGwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaVcReCiAAAAAAAAAAOAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAbAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpVxF4KIAAAAAAAAAA4AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAbAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpMdC56IAAAAAAAAAA4AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAbAAAAAAAAAACQUsYKOw/kj+JXXI5eosCY2LS+fgH0zFHZXnenMJ4wUAAAAAJUC+QAAAAAGwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAZAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpVxF4LsAAAAAAAAAA0AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAbAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpVxF4KIAAAAAAAAAA0AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{atxjw1XkT38KrVuORmTYlHUPY+hPTbM/5TGFnFg1hgmQA0rBIhKF2q9R1hSuna6bbBhJfT+iPX+wthgSm+GEDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('a76e0260f6b83c6ea93f545d17de721c079dc31e81ee5edc41f159ec5fb48443', 26, 1, 'GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q', 107374182401, 100, 1, '2019-06-03 16:35:33.768436', '2019-06-03 16:35:33.768437', 111669153792, 'AAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAZAAAABkAAAABAAAAAAAAAAAAAAABAAAAAAAAAAQAAAABVVNEAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAAAAAAAdzWUAAAAAAEAAAABAAAAAAAAAAEyboYYAAAAQBqzCYDuLYn/jXhfEVxEGigMCJGoOBCK92lUb3Um15PgwSJ63tNl+FpH8+y5c+mCs/rzcvdyo9uXdodd4LXWiQg=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAAAAAAAUAAAABVVNEAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAAAAAAAdzWUAAAAAAEAAAABAAAAAQAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAGgAAAAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAACVAvjOAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAGgAAAAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAACVAvjOAAAABkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+M4AAAAGQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+M4AAAAGQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAHc1lAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABoAAAACAAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAAAAAAAUAAAABVVNEAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAAAAAAAdzWUAAAAAAEAAAABAAAAAQAAAAAAAAAA', 'AAAAAgAAAAMAAAAZAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+QAAAAAGQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+OcAAAAGQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{GrMJgO4tif+NeF8RXEQaKAwIkag4EIr3aVRvdSbXk+DBInre02X4Wkfz7Llz6YKz+vNy93Kj25d2h13gtdaJCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('92a654c76966ac61acc9df0b75f91cbde3b551c9e9766730827af42d1e247cc3', 26, 2, 'GB6GN3LJUW6JYR7EDOJ47VBH7D45M4JWHXGK6LHJRAEI5JBSN2DBQY7Q', 107374182402, 100, 1, '2019-06-03 16:35:33.76906', '2019-06-03 16:35:33.76906', 111669157888, 'AAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAZAAAABkAAAACAAAAAAAAAAAAAAABAAAAAAAAAAQAAAAAAAAAAVVTRAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAAAdzWUAAAAAAEAAAABAAAAAAAAAAEyboYYAAAAQBbE9T7oBKoN0/S3AV7GoSRe+xT79SlWNCYEtL1RPExL8FLhw5EDsXLoAvIBbBvHIr9NKcPtWDyhcHlIuaZKIg8=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAAAAAAAYAAAAAAAAAAVVTRAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAAAdzWUAAAAAAEAAAABAAAAAQAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAGgAAAAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAACVAvjOAAAABkAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAB3NZQAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+M4AAAAGQAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAHc1lAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAADAAAAGgAAAAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAACVAvjOAAAABkAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAB3NZQAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+M4AAAAGQAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAHc1lAAAAAAAdzWUAAAAAAAAAAAAAAAAAAAAABoAAAACAAAAAHxm7WmlvJxH5BuTz9Qn+PnWcTY9zK8s6YgIjqQyboYYAAAAAAAAAAYAAAAAAAAAAVVTRAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAAAdzWUAAAAAAEAAAABAAAAAQAAAAAAAAAA', 'AAAAAgAAAAMAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+OcAAAAGQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAaAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+M4AAAAGQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{FsT1PugEqg3T9LcBXsahJF77FPv1KVY0JgS0vVE8TEvwUuHDkQOxcugC8gFsG8civ00pw+1YPKFweUi5pkoiDw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('5065cd7c97cfb6fbf7da8493beed47ed2c7efb3b00b77a4c92692ed487fb86a4', 25, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 13, 100, 1, '2019-06-03 16:35:33.779849', '2019-06-03 16:35:33.779849', 107374186496, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAANAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAfGbtaaW8nEfkG5PP1Cf4+dZxNj3MryzpiAiOpDJuhhgAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBthwT3JCg5IZkKRNK3pHBa/eG8zq8Af9gFPWlYvEdRo6jzA5D9fYOcDpKD3dEAuPLNNAHj9tNbZUJA3rwxN94B', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAGQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaXxSNm7AAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAGQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaXxSNm7AAAAAAAAAANAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAZAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpfFI2bsAAAAAAAAAA0AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAZAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpVxF4LsAAAAAAAAAA0AAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAAAAAAAAAAB8Zu1ppbycR+Qbk8/UJ/j51nE2PcyvLOmICI6kMm6GGAAAAAJUC+QAAAAAGQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpfFI2dQAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAZAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpfFI2bsAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{bYcE9yQoOSGZCkTSt6RwWv3hvM6vAH/YBT1pWLxHUaOo8wOQ/X2DnA6Sg93RALjyzTQB4/bTW2VCQN68MTfeAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('01346de1ca30ce03149d9f54945956a22f9cbed3d81f81c62bb59cf8cdd8b893', 24, 1, 'GB2QIYT2IAUFMRXKLSLLPRECC6OCOGJMADSPTRK7TGNT2SFR2YGWDARD', 90194313217, 100, 1, '2019-06-03 16:35:33.79149', '2019-06-03 16:35:33.79149', 103079219200, 'AAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAZAAAABUAAAABAAAAAAAAAAAAAAABAAAAAAAAAAMAAAABVVNEAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAAAAAAAEeGjAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAAbHWDWEAAABA0L+69D1hxpytAkX6cvPiBuO80ql8SQKZ15POVxx9wYl6mZrL+6UWGab/+6ng2M+a29E7ON+Xs46Y9MNqTh91AQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAEAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAAAAAAAAwAAAAAAAAAAC+vCAAAAAAFVU0QAAAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAAAvrwgAAAAAAAAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAAAAAAAQAAAABVVNEAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAAAAAAABfXhAAAAAAEAAAABAAAAAAAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAGAAAAAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAACVAvjnAAAABUAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAGAAAAAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAACVAvjnAAAABUAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACQAAAAMAAAAXAAAAAgAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAAAAAADAAAAAAAAAAFVU0QAAAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAAAvrwgAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAIAAAACAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAAAAAAMAAAADAAAAFwAAAAEAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAABVVNEAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAAAAAAAf/////////8AAAABAAAAAQAAAAAL68IAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAYAAAAAQAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAFVU0QAAAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAAAvrwgB//////////wAAAAEAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAABcAAAAAAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAlQL4zgAAAAVAAAAAgAAAAIAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAL68IAAAAAAAAAAAAAAAABAAAAGAAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACSCAhOAAAABUAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAAgAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAAAAAAEAAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAAAAAAAX14QAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAMAAAAYAAAAAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAJUC+OcAAAAFQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAYAAAAAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAJf96WcAAAAFQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAX14QAAAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAAVAAAAAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAJUC+QAAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAYAAAAAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAJUC+OcAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{0L+69D1hxpytAkX6cvPiBuO80ql8SQKZ15POVxx9wYl6mZrL+6UWGab/+6ng2M+a29E7ON+Xs46Y9MNqTh91AQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d8b2508123656b1df1ee17c2767829bc22ab41959ad25e6ccc520e849516fba1', 23, 1, 'GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC', 90194313218, 100, 1, '2019-06-03 16:35:33.815586', '2019-06-03 16:35:33.815586', 98784251904, 'AAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAZAAAABUAAAACAAAAAAAAAAAAAAABAAAAAAAAAAMAAAAAAAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAC+vCAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAARXFjncAAABATR48xYiKbu8AOoXFwvcvILZ0/pQkfGuwwAoIZNefr7ydIwlcuL44XPM7pJ/6jDSbqBudTNWdE2JRjuq7HI7IAA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAAAAAAMAAAAAAAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAC+vCAAAAAAEAAAABAAAAAAAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAFwAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACVAvjOAAAABUAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFwAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACVAvjOAAAABUAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAMAAAAXAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+M4AAAAFQAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAXAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+M4AAAAFQAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAAAAAAAAAAAC+vCAAAAAAAAAAAAAAAAAAAAABcAAAACAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAAAAAAMAAAAAAAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAC+vCAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAAwAAABYAAAABAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAAAAAAH//////////AAAAAQAAAAAAAAAAAAAAAQAAABcAAAABAAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAAVVTRAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAAAAAAAAH//////////AAAAAQAAAAEAAAAAC+vCAAAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAWAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+OcAAAAFQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAXAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+M4AAAAFQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{TR48xYiKbu8AOoXFwvcvILZ0/pQkfGuwwAoIZNefr7ydIwlcuL44XPM7pJ/6jDSbqBudTNWdE2JRjuq7HI7IAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('be05e4bd966d58689e1b6fae013e5aa77bde56e6acd2db9b96870e5e746a4ab7', 22, 1, 'GBOK7BOUSOWPHBANBYM6MIRYZJIDIPUYJPXHTHADF75UEVIVYWHHONQC', 90194313217, 100, 1, '2019-06-03 16:35:33.827572', '2019-06-03 16:35:33.827572', 94489284608, 'AAAAAFyvhdSTrPOEDQ4Z5iI4ylA0PphL7nmcAy/7QlUVxY53AAAAZAAAABUAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYX//////////AAAAAAAAAAEVxY53AAAAQDMCWfC0eGNJuYIX3s5AUNLernpcHTn8O6ygq/Nw3S5vny/W42O5G4G6UsihVU1xd5bR4im2+VzQlQYQhe0jhwg=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAFgAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACVAvjnAAAABUAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFgAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACVAvjnAAAABUAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAWAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+OcAAAAFQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAWAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+OcAAAAFQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWAAAAAQAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAFVU0QAAAAAAHUEYnpAKFZG6lyWt8SCF5wnGSwA5PnFX5mbPUix1g1hAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAVAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+QAAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAWAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+OcAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{MwJZ8LR4Y0m5ghfezkBQ0t6uelwdOfw7rKCr83DdLm+fL9bjY7kbgbpSyKFVTXF3ltHiKbb5XNCVBhCF7SOHCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2a6987a6930eab7e3becacf9b76ed7a06802668c1f1eb0f82f5671014b4b636a', 21, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 11, 100, 1, '2019-06-03 16:35:33.838303', '2019-06-03 16:35:33.838303', 90194317312, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAALAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAXK+F1JOs84QNDhnmIjjKUDQ+mEvueZwDL/tCVRXFjncAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDfpUesb4kQ/RfBx1UxqNOtZ2+4R4S0XxzggPR1C3YyhZAr/K8KyZCg4ejDTFnhu9qAh4GLZLkbBraGncT9DcYF', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAFQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LacbTsvUAAAAAAAAAAKAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LacbTsvUAAAAAAAAAALAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOy9QAAAAAAAAAAsAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpoZL0tQAAAAAAAAAAsAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVAAAAAAAAAABcr4XUk6zzhA0OGeYiOMpQND6YS+55nAMv+0JVFcWOdwAAAAJUC+QAAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOzAYAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOy+0AAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{36VHrG+JEP0XwcdVMajTrWdvuEeEtF8c4ID0dQt2MoWQK/yvCsmQoOHow0xZ4bvagIeBi2S5Gwa2hp3E/Q3GBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('96415ac1d2f79621b26b1568f963fd8dd6c50c20a22c7428cefbfe9dee867588', 21, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 12, 100, 1, '2019-06-03 16:35:33.838461', '2019-06-03 16:35:33.838461', 90194321408, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAMAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAdQRiekAoVkbqXJa3xIIXnCcZLADk+cVfmZs9SLHWDWEAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDdJGdvdZ2S4QoXdO+Odt8ZRdeVu7mBvq7FtP9okqr98pGD/jSAraklQvaRmCyMALIMD2kG8R2KjhKvy7oIL6IB', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAFQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaaGS9LUAAAAAAAAAALAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaaGS9LUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpoZL0tQAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpfFI2dQAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVAAAAAAAAAAB1BGJ6QChWRupclrfEghecJxksAOT5xV+Zmz1IsdYNYQAAAAJUC+QAAAAAFQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOy+0AAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAVAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOy9QAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{3SRnb3WdkuEKF3TvjnbfGUXXlbu5gb6uxbT/aJKq/fKRg/40gK2pJUL2kZgsjACyDA9pBvEdio4Sr8u6CC+iAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('f08dc1fec150f276562866ce4f5272f658cf0bd9fd8c1d96a22c196be2e1b25a', 20, 1, 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD', 68719476739, 100, 1, '2019-06-03 16:35:33.848132', '2019-06-03 16:35:33.848132', 85899350016, 'AAAAAONt/6wGI884Zi6sYDYC1GOV/drnh4OcRrTrqJPoOTUKAAAAZAAAABAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAIAAAAAAAAAADuaygAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAA7msoAAAAAAAAAAAAAAAAB6Dk1CgAAAEB+7jxesBKKrF343onyycjp2tiQLZiGH2ETl+9fuOqotveY2rIgvt9ng+QJ2aDP3+PnDsYEa9ZUaA+Zne2nIGgE', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAACAAAAAAAAAAEAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAAgAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAADuaygAAAAAAAAAAADuaygAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAA7msoAAAAAAA==', 'AAAAAQAAAAIAAAADAAAAFAAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvi1AAAABAAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFAAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvi1AAAABAAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAMAAAATAAAAAQAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAHc1lAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAUAAAAAQAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAALLQXgB//////////wAAAAEAAAAAAAAAAAAAAAMAAAATAAAAAgAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAACAAAAAUVVUgAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAADuaygAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAIAAAACAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAIAAAADAAAAFAAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvi1AAAABAAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAFAAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACGHEY1AAAABAAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAADAAAAEwAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAADAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAA7msoAAAAAAHc1lAAAAAAAAAAAAAAAAAEAAAAUAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAKPpqzUAAAAEAAAAAMAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAAAAAAAAAAAdzWUAAAAAAAAAAAA', 'AAAAAgAAAAMAAAATAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+M4AAAAEAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAUAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+LUAAAAEAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{fu48XrASiqxd+N6J8snI6drYkC2Yhh9hE5fvX7jqqLb3mNqyIL7fZ4PkCdmgz9/j5w7GBGvWVGgPmZ3tpyBoBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('198844c8b472daacc5b717695a4ca16ac799a13fb2cf4152d19e2117ae1c56c3', 19, 1, 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD', 68719476738, 100, 1, '2019-06-03 16:35:33.865967', '2019-06-03 16:35:33.865967', 81604382720, 'AAAAAONt/6wGI884Zi6sYDYC1GOV/drnh4OcRrTrqJPoOTUKAAAAZAAAABAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAIAAAABVVNEAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAA7msoAAAAAAAQI6ATBFImTS1I7Fly9YiufQ/dC4uMOetO+m/BysWXyAAAAAUVVUgAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAdzWUAAAAAAEAAAAAAAAAAAAAAAHoOTUKAAAAQMs9vNZ518oYUMp38TakovW//DDTbs/9oPj1RAix5ElC/d7gbWaaNNJxKQR7eMNO6rB+ntGqee4WurTJgA4k2ws=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAACAAAAAAAAAAIAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAAQAAAAAAAAAAdzWUAAAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAADuaygAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAAgAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAHc1lAAAAAAAAAAAAHc1lAAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAB3NZQAAAAAAA==', 'AAAAAQAAAAIAAAADAAAAEwAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvjOAAAABAAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEwAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvjOAAAABAAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAMAAAARAAAAAQAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAATAAAAAQAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAHc1lAB//////////wAAAAEAAAAAAAAAAAAAAAMAAAASAAAAAQAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAADuaygB//////////wAAAAEAAAAAAAAAAAAAAAEAAAATAAAAAQAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAMAAAASAAAAAgAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAACAAAAAUVVUgAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAALLQXgAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAEAAAATAAAAAgAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAACAAAAAUVVUgAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAADuaygAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAMAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAALLQXgAAAAAA7msoAAAAAAAAAAAAAAAAAQAAABMAAAAAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAlQL4tQAAAAQAAAAAwAAAAIAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAEAAAAAO5rKAAAAAAB3NZQAAAAAAAAAAAAAAAADAAAAEgAAAAIAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAAQAAAAAAAAABVVNEAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAADuaygAAAAAAQAAAAIAAAAAAAAAAAAAAAAAAAABAAAAEwAAAAIAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAAAAAAAQAAAAAAAAABVVNEAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAB3NZQAAAAAAQAAAAIAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAARAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+OcAAAAEAAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAATAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+M4AAAAEAAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{yz281nnXyhhQynfxNqSi9b/8MNNuz/2g+PVECLHkSUL93uBtZpo00nEpBHt4w07qsH6e0ap57ha6tMmADiTbCw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('902b90c2322b9e6b335e7543389a7446b86e3039ebf59ec66dffb50eaec0dc85', 18, 1, 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU', 68719476737, 100, 1, '2019-06-03 16:35:33.884097', '2019-06-03 16:35:33.884097', 77309415424, 'AAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAZAAAABAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAABVVNEAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAA7msoAAAAAAAAAAAERVFexAAAAQC9X2I3Zz1x3AQMqL4XCzePTlwnokv2BQnWGmT007oH59gai3eNu7/WVoHtW8hsgHjs1mZK709FzzRF2cbD2tQE=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAARAAAAAQAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAASAAAAAQAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAADuaygB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+OcAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{L1fYjdnPXHcBAyovhcLN49OXCeiS/YFCdYaZPTTugfn2BqLd427v9ZWge1byGyAeOzWZkrvT0XPNEXZxsPa1AQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('ca756d1519ceda79f8722042b12cea7ba004c3bd961adb62b59f88a867f86eb3', 18, 2, 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU', 68719476738, 100, 1, '2019-06-03 16:35:33.884375', '2019-06-03 16:35:33.884376', 77309419520, 'AAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAZAAAABAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAMAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAAAAAAARFUV7EAAABALuai5QxceFbtAiC5nkntNVnvSPeWR+C+FgplPAdRgRS+PPESpUiSCyuiwuhmvuDw7kwxn+A6E0M4ca1s2qzMAg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAEAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAAAAAAAAAAA7msoAAAAAAAAAAAAAAAAAAAAABIAAAACAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAEAAAAAAAAAAVVTRAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAAA7msoAAAAAAEAAAACAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+OcAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+M4AAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Luai5QxceFbtAiC5nkntNVnvSPeWR+C+FgplPAdRgRS+PPESpUiSCyuiwuhmvuDw7kwxn+A6E0M4ca1s2qzMAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('37bb79f6959c0e8e9b3d31f6c9308d8d084d9c6742cfa56ca094cfa6eae99423', 18, 3, 'GAXMF43TGZHW3QN3REOUA2U5PW5BTARXGGYJ3JIFHW3YT6QRKRL3CPPU', 68719476739, 100, 1, '2019-06-03 16:35:33.884772', '2019-06-03 16:35:33.884772', 77309423616, 'AAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAZAAAABAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAMAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAAAstBeAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAARFUV7EAAABArzp9Fxxql+yoysglDjXm9+rsJeNX2GsSa7TOy3AzHOu4Y5Z8ICx52Q885gQGQWMtEP0w6yh83d6+o6kjC/WuAg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAIAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAAAstBeAAAAAAEAAAABAAAAAAAAAAAAAAAA', 'AAAAAQAAAAIAAAADAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAO5rKAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAMAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAAAAAAAAAAAAA7msoAAAAAAAAAAAAAAAAAQAAAAMAAAADAAAAEgAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvi1AAAABAAAAADAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAO5rKAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAMAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAABAAAAALLQXgAAAAAA7msoAAAAAAAAAAAAAAAAAAAAABIAAAACAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAIAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAAAAAAAstBeAAAAAAEAAAABAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+M4AAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAASAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+LUAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{rzp9Fxxql+yoysglDjXm9+rsJeNX2GsSa7TOy3AzHOu4Y5Z8ICx52Q885gQGQWMtEP0w6yh83d6+o6kjC/WuAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('cdef45dd961d59375351ea7dd7ef6414ff49371a335723e84dafacea1e13665a', 17, 1, 'GDRW375MAYR46ODGF2WGANQC2RRZL7O246DYHHCGWTV2RE7IHE2QUQLD', 68719476737, 100, 1, '2019-06-03 16:35:33.898077', '2019-06-03 16:35:33.898077', 73014448128, 'AAAAAONt/6wGI884Zi6sYDYC1GOV/drnh4OcRrTrqJPoOTUKAAAAZAAAABAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsX//////////AAAAAAAAAAHoOTUKAAAAQIjLqcYXE8EAsH6Dx2hwPjiEfHGZ4jsMNZZc7PynNiJi9kFXjfvvLDlWizGAr2B9MFDrfDRDvjnBxKKhJifEcQM=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEQAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvjnAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEQAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvjnAAAABAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAARAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+OcAAAAEAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAARAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+OcAAAAEAAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAARAAAAAQAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAFVU0QAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAARAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+OcAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{iMupxhcTwQCwfoPHaHA+OIR8cZniOww1llzs/Kc2ImL2QVeN++8sOVaLMYCvYH0wUOt8NEO+OcHEoqEmJ8RxAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d1f593eb5e14f97027bc79821fa46628c107034fba9a5acef6a9da79e051ee73', 17, 2, 'GACAR2AEYEKITE2LKI5RMXF5MIVZ6Q7XILROGDT22O7JX4DSWFS7FDDP', 68719476737, 100, 1, '2019-06-03 16:35:33.898272', '2019-06-03 16:35:33.898273', 73014452224, 'AAAAAAQI6ATBFImTS1I7Fly9YiufQ/dC4uMOetO+m/BysWXyAAAAZAAAABAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABRVVSAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsX//////////AAAAAAAAAAFysWXyAAAAQI7hbwZc1+KWfheVnYAq5TXFX9ancHJmJq0wV0c9ONIfG6U8trhIVeVoiED2eUFFmhx+bBtF9TPSvifF/mfDlQk=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEQAAAAAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAACVAvjnAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEQAAAAAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAACVAvjnAAAABAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAARAAAAAAAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAJUC+OcAAAAEAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAARAAAAAAAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAJUC+OcAAAAEAAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAARAAAAAQAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAFFVVIAAAAAAC7C83M2T23Bu4kdQGqdfboZgjcxsJ2lBT23ifoRVFexAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAARAAAAAAAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAJUC+OcAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{juFvBlzX4pZ+F5WdgCrlNcVf1qdwcmYmrTBXRz040h8bpTy2uEhV5WiIQPZ5QUWaHH5sG0X1M9K+J8X+Z8OVCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('a5a9e3ca63e9cc155359c97337bcb14464cca56b230a4d0c7f27582644d16809', 16, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 8, 100, 1, '2019-06-03 16:35:33.910688', '2019-06-03 16:35:33.910688', 68719480832, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA423/rAYjzzhmLqxgNgLUY5X92ueHg5xGtOuok+g5NQoAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBhFD/bYaTZZJ3VJ9xJqXoW5eeLK0AeFaATBH92cRfx0WUTFqp6rXx47fMBUxkWYq8bAHMfYCS5XXPRg86sAGUK', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LajaV7cGAAAAAAAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LajaV7cGAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtwYAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqEVUvgYAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAADjbf+sBiPPOGYurGA2AtRjlf3a54eDnEa066iT6Dk1CgAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXt1EAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtzgAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{YRQ/22Gk2WSd1SfcSal6FuXniytAHhWgEwR/dnEX8dFlExaqeq18eO3zAVMZFmKvGwBzH2AkuV1z0YPOrABlCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6a056189b45760c607e331c90c5a8b4cd720961df8bc8cecfd4aa388b577a6cb', 16, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 9, 100, 1, '2019-06-03 16:35:33.910845', '2019-06-03 16:35:33.910845', 68719484928, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAJAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAABAjoBMEUiZNLUjsWXL1iK59D90Li4w56076b8HKxZfIAAAACVAvkAAAAAAAAAAABVvwF9wAAAEAxC5cl7tkjQI0cfFZTiIFDuo0SwyYnNqTUH2hxDBtm7h/vUkBG3cgwGXS87ninVkhmvdIpTWfeIeGiw7kgefUA', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LahFVL4GAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LahFVL4GAAAAAAAAAAJAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqEVUvgYAAAAAAAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtp7BRxQYAAAAAAAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAECOgEwRSJk0tSOxZcvWIrn0P3QuLjDnrTvpvwcrFl8gAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtzgAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtx8AAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{MQuXJe7ZI0CNHHxWU4iBQ7qNEsMmJzak1B9ocQwbZu4f71JARt3IMBl0vO54p1ZIZr3SKU1n3iHhosO5IHn1AA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('18bf6cce20cfbb0f9079c4b8783718949d13bd12d173a60363d2b8e3a07efead', 16, 3, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 10, 100, 1, '2019-06-03 16:35:33.910968', '2019-06-03 16:35:33.910968', 68719489024, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAKAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAALsLzczZPbcG7iR1Aap19uhmCNzGwnaUFPbeJ+hFUV7EAAAACVAvkAAAAAAAAAAABVvwF9wAAAEC/RVto6ytAqHpd6ZFWjwXQyXopKORz8QSvz0d8RoPrOEBgNEuAj8+kbyhS7QieOqwbiJrS0AU8YWaBQQ4zc+wL', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaewUcUGAAAAAAAAAAJAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAEAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaewUcUGAAAAAAAAAAKAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtp7BRxQYAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtpxtOzAYAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAuwvNzNk9twbuJHUBqnX26GYI3MbCdpQU9t4n6EVRXsQAAAAJUC+QAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtx8AAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAQAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXtwYAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{v0VbaOsrQKh6XemRVo8F0Ml6KSjkc/EEr89HfEaD6zhAYDRLgI/PpG8oUu0InjqsG4ia0tAFPGFmgUEOM3PsCw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('142c988b1f67984f74a1581de9caecf499e60f1e0eed496661aa2c559238764c', 15, 1, 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG', 55834574850, 100, 1, '2019-06-03 16:35:33.920804', '2019-06-03 16:35:33.920804', 64424513536, 'AAAAAI4uD0KXOiTFreQAeellyBKKdeD5+7Vurn+sRPr2O2TcAAAAZAAAAA0AAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAG5M9WO2kexu4dljTd0YSuyEnmDUsxKamzJWiv4FGkoQAAAABVVNEAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAAF9eEAAAAAAAAAAAH2O2TcAAAAQJBUx5tWfjAwXxab9U5HOjZvBRv3u95jXbyzuqeZ/kjsyMsU0jO/g03Rf1zgect1hj4hDYGN8mW4oEot0sSTZgw=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADwAAAAAAAAAAji4PQpc6JMWt5AB56WXIEop14Pn7tW6uf6xE+vY7ZNwAAAACThYCOAAAAA0AAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADwAAAAAAAAAAji4PQpc6JMWt5AB56WXIEop14Pn7tW6uf6xE+vY7ZNwAAAACThYCOAAAAA0AAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAOAAAAAQAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAFVU0QAAAAAAI4uD0KXOiTFreQAeellyBKKdeD5+7Vurn+sRPr2O2TcAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAPAAAAAQAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAFVU0QAAAAAAI4uD0KXOiTFreQAeellyBKKdeD5+7Vurn+sRPr2O2TcAAAAAAX14QB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAOAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJOFgKcAAAADQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAPAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJOFgI4AAAADQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{kFTHm1Z+MDBfFpv1Tkc6Nm8FG/e73mNdvLO6p5n+SOzIyxTSM7+DTdF/XOB5y3WGPiENgY3yZbigSi3SxJNmDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('30880dd42d8e402a30d8a3527b56c1e33e18e87c46e1332ea5cfc1721fd87cfb', 14, 1, 'GCHC4D2CS45CJRNN4QAHT2LFZAJIU5PA7H53K3VOP6WEJ6XWHNSNZKQG', 55834574849, 100, 1, '2019-06-03 16:35:33.933624', '2019-06-03 16:35:33.933624', 60129546240, 'AAAAAI4uD0KXOiTFreQAeellyBKKdeD5+7Vurn+sRPr2O2TcAAAAZAAAAA0AAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAG5M9WO2kexu4dljTd0YSuyEnmDUsxKamzJWiv4FGkoQAAAAAAAAAAAX14QAAAAAAAAAAAfY7ZNwAAABAieZSSuOZqlwtyjnj5d/S0GUSGiQvy0ipPLynpl4UvO8qc7CDz3vsLROlN2g50qXirydSOdao56hvRhrEfRsGCA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADgAAAAAAAAAAji4PQpc6JMWt5AB56WXIEop14Pn7tW6uf6xE+vY7ZNwAAAACVAvjnAAAAA0AAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADgAAAAAAAAAAji4PQpc6JMWt5AB56WXIEop14Pn7tW6uf6xE+vY7ZNwAAAACVAvjnAAAAA0AAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAOAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJUC+OcAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAOAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJaAcScAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAOAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJUC+OcAAAADQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAOAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJOFgKcAAAADQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAANAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJUC+QAAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAOAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJUC+OcAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ieZSSuOZqlwtyjnj5d/S0GUSGiQvy0ipPLynpl4UvO8qc7CDz3vsLROlN2g50qXirydSOdao56hvRhrEfRsGCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('dbd964fcfdb336a30f21c240fffdaf73d7c75880ed1b99375c62f84e3e592570', 14, 2, 'GANZGPKY5WSHWG5YOZMNG52GCK5SCJ4YGUWMJJVGZSK2FP4BI2JIJN2C', 55834574849, 100, 1, '2019-06-03 16:35:33.934033', '2019-06-03 16:35:33.934033', 60129550336, 'AAAAABuTPVjtpHsbuHZY03dGErshJ5g1LMSmpsyVor+BRpKEAAAAZAAAAA0AAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3H//////////AAAAAAAAAAGBRpKEAAAAQGDAV/5Op2DmFUP84dmyT5G/gxn1WzgdMrkSSU7wfpu39ycq36Sg+gs2ypRjw5hxxeMUj/GVEKipcDGndei38Aw=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADgAAAAAAAAAAG5M9WO2kexu4dljTd0YSuyEnmDUsxKamzJWiv4FGkoQAAAACWgHEnAAAAA0AAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADgAAAAAAAAAAG5M9WO2kexu4dljTd0YSuyEnmDUsxKamzJWiv4FGkoQAAAACWgHEnAAAAA0AAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAOAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJaAcScAAAADQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAOAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJaAcScAAAADQAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOAAAAAQAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAFVU0QAAAAAAI4uD0KXOiTFreQAeellyBKKdeD5+7Vurn+sRPr2O2TcAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAANAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJUC+QAAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAOAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJUC+OcAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{YMBX/k6nYOYVQ/zh2bJPkb+DGfVbOB0yuRJJTvB+m7f3JyrfpKD6CzbKlGPDmHHF4xSP8ZUQqKlwMad16LfwDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('cd8a8e9eb53fd268d1294e228995c27f422d90783c4054e44ab0028fc1da210a', 13, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 6, 100, 1, '2019-06-03 16:35:33.949136', '2019-06-03 16:35:33.949137', 55834578944, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAji4PQpc6JMWt5AB56WXIEop14Pn7tW6uf6xE+vY7ZNwAAAACVAvkAAAAAAAAAAABVvwF9wAAAEAUtdYWyr64yv/rKPr0/vV4vYyonfsWxpxHsiYLHKJ3bm6k+ypiAByc8t0K+7bzxSLPjmjKKN5Prw7AdenlC7MB', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaoEXalRAAAAAAAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaoEXalRAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqVEAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqW9asFEAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAANAAAAAAAAAACOLg9Clzokxa3kAHnpZcgSinXg+fu1bq5/rET69jtk3AAAAAJUC+QAAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAALAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqYMAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqWoAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{FLXWFsq+uMr/6yj69P71eL2MqJ37FsacR7ImCxyid25upPsqYgAcnPLdCvu288Uiz45oyijeT68OwHXp5QuzAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('bfbd5e9457d717bcf847291a6c24b7cd8db4ff784ecd4592be30d08146c0c264', 13, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 7, 100, 1, '2019-06-03 16:35:33.949341', '2019-06-03 16:35:33.949341', 55834583040, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAG5M9WO2kexu4dljTd0YSuyEnmDUsxKamzJWiv4FGkoQAAAACVAvkAAAAAAAAAAABVvwF9wAAAEDY1TiMj+qj8+zYb2Vb60h+qWxZtFfSGwb0kvKttSFAHQhGOjIddiVQopx9LDRO6UgPmLLxFvQpIzeGnagh3vQD', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LalvWrBRAAAAAAAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LalvWrBRAAAAAAAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqW9asFEAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqNpXt1EAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAANAAAAAAAAAAAbkz1Y7aR7G7h2WNN3RhK7ISeYNSzEpqbMlaK/gUaShAAAAAJUC+QAAAAADQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqWoAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAANAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqVEAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{2NU4jI/qo/Ps2G9lW+tIfqlsWbRX0hsG9JLyrbUhQB0IRjoyHXYlUKKcfSw0TulID5iy8Rb0KSM3hp2oId70Aw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('0e128647b2b93786b6b76e182dcda0173757066f8caf0523d1ba3b47fd6f720d', 12, 1, 'GDCVTBGSEEU7KLXUMHMSXBALXJ2T4I2KOPXW2S5TRLKDRIAXD5UDHAYO', 47244640257, 100, 1, '2019-06-03 16:35:33.960369', '2019-06-03 16:35:33.960369', 51539611648, 'AAAAAMVZhNIhKfUu9GHZK4QLunU+I0pz721Ls4rUOKAXH2gzAAAAZAAAAAsAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAg/K/Blr9FO/nVEGLdmCzChMYpmcQzxIhFm6NBzxznX0AAAAAHc1lAAAAAAAAAAABFx9oMwAAAEBwY9HQAR2SMPe3JPvmBBtBk2jfog0GFEFYkLNFzQNqvYl7iZitmO5FQmkKlv/NO5ZcaWBqXcHhOQpk0s2XSBQF', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAADAAAAAAAAAAAxVmE0iEp9S70YdkrhAu6dT4jSnPvbUuzitQ4oBcfaDMAAAACVAvjnAAAAAsAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAADAAAAAAAAAAAxVmE0iEp9S70YdkrhAu6dT4jSnPvbUuzitQ4oBcfaDMAAAACVAvjnAAAAAsAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAMAAAAAAAAAADFWYTSISn1LvRh2SuEC7p1PiNKc+9tS7OK1DigFx9oMwAAAAJUC+OcAAAACwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAMAAAAAAAAAADFWYTSISn1LvRh2SuEC7p1PiNKc+9tS7OK1DigFx9oMwAAAAI2Pn6cAAAACwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMAAAAAAAAAACD8r8GWv0U7+dUQYt2YLMKEximZxDPEiEWbo0HPHOdfQAAAAAdzWUAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAALAAAAAAAAAADFWYTSISn1LvRh2SuEC7p1PiNKc+9tS7OK1DigFx9oMwAAAAJUC+QAAAAACwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAMAAAAAAAAAADFWYTSISn1LvRh2SuEC7p1PiNKc+9tS7OK1DigFx9oMwAAAAJUC+OcAAAACwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{cGPR0AEdkjD3tyT75gQbQZNo36INBhRBWJCzRc0Dar2Je4mYrZjuRUJpCpb/zTuWXGlgal3B4TkKZNLNl0gUBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('67601a2ca212b84092a7d3c521172b67f4b93d72b726a06c540917d2ab83c1a1', 11, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 5, 100, 1, '2019-06-03 16:35:33.97177', '2019-06-03 16:35:33.97177', 47244644352, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAxVmE0iEp9S70YdkrhAu6dT4jSnPvbUuzitQ4oBcfaDMAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBHLko6/Tbv0v/5CWHkixXnbyoU6qQ6yewZGqPHFSzNxMfud86eYGkN0j4msMCXfLAou7iKOVn0MWyzlpvYRA0B', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaqZYKKDAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaqZYKKDAAAAAAAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAALAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqplgooMAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAALAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqgRdqYMAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAALAAAAAAAAAADFWYTSISn1LvRh2SuEC7p1PiNKc+9tS7OK1DigFx9oMwAAAAJUC+QAAAAACwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAKAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqplgopwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAALAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqplgooMAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Ry5KOv0279L/+Qlh5IsV528qFOqkOsnsGRqjxxUszcTH7nfOnmBpDdI+JrDAl3ywKLu4ijlZ9DFss5ab2EQNAQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('fdb696a797b769176cbaed3a50e4a6a8671119621f65a3f954a3bcf100c7ef0c', 10, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 38654705665, 200, 2, '2019-06-03 16:35:33.980859', '2019-06-03 16:35:33.980859', 42949677056, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAAAyAAAAAkAAAABAAAAAAAAAAAAAAACAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAABfXhAAAAAAAAAAAB2QViXgAAAEAxyl5gvCCDC7l0pq9b/Btd3cOUUcY9Rv0ALxVjul4EVSL1Vygr107GjDo11+YswdmlCuWf7KItU0chlogpns4L', 'AAAAAAAAAMgAAAAAAAAAAgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvjOAAAAAkAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvjOAAAAAkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAACAAAABAAAAAMAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqpZlshwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAKAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqpfjKlwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAKAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+M4AAAACQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAKAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJOFgI4AAAACQAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAQAAAADAAAACgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaqX4ypcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaqZYKKcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAADAAAACgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACThYCOAAAAAkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACSCAhOAAAAAkAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+QAAAAACQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAKAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+M4AAAACQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{McpeYLwggwu5dKavW/wbXd3DlFHGPUb9AC8VY7peBFUi9VcoK9dOxow6NdfmLMHZpQrln+yiLVNHIZaIKZ7OCw==}', 'none', NULL, NULL, true, 200);
INSERT INTO history_transactions VALUES ('66c28c0ccd5a2e47026aacafa2ecd3c501fe5de349ef376c0f8afb893c7bb55d', 9, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 4, 100, 1, '2019-06-03 16:35:33.989997', '2019-06-03 16:35:33.989997', 38654709760, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvkAAAAAAAAAAABVvwF9wAAAEARD6MVWgEASusfhr6JdF9K3Rie2XCRJKl/NoKyJcrd1kGs3ygpp55xu80YlFwgNVErZ/cEAHYOq06CwNfnE2sC', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LasraKscAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACQAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LasraKscAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytoqxwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqpZlshwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+QAAAAACQAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytoqzUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytoqxwAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{EQ+jFVoBAErrH4a+iXRfSt0YntlwkSSpfzaCsiXK3dZBrN8oKaeecbvNGJRcIDVRK2f3BAB2DqtOgsDX5xNrAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('dd74eee27a59843b28a05ad08abf65eaa231b7debe4d05550c0a7a424cca5929', 8, 1, 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB', 30064771073, 100, 1, '2019-06-03 16:35:33.999724', '2019-06-03 16:35:33.999725', 34359742464, 'AAAAADnqxUES0av4d13+BBKxKs+6A/cJNtM1/nf4eHV3Ow06AAAAZAAAAAcAAAABAAAAAAAAAAIAAAAAAAAAewAAAAEAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAJiWgAAAAAAAAAABdzsNOgAAAEBjk5EFqV8GiL9xU62OUCKeScXxGMTMqJoD7ryiGf5jLPZJRSphbWC3ZycHE+pDuu/6EKSqcNUri5AXzQmM+GYB', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACVAvicAAAAAcAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACVAvicAAAAAcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyrQFLUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyr2OlUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+JwAAAABwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJTc0vwAAAABwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+QAAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+OcAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Y5ORBalfBoi/cVOtjlAinknF8RjEzKiaA+68ohn+Yyz2SUUqYW1gt2cnBxPqQ7rv+hCkqnDVK4uQF80JjPhmAQ==}', 'id', '123', NULL, true, 100);
INSERT INTO history_transactions VALUES ('2551e76a3ce4881b7bc73fdfd89d670d511ea7d4e56156252b51777023202de7', 8, 2, 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB', 30064771074, 100, 1, '2019-06-03 16:35:34.000003', '2019-06-03 16:35:34.000004', 34359746560, 'AAAAADnqxUES0av4d13+BBKxKs+6A/cJNtM1/nf4eHV3Ow06AAAAZAAAAAcAAAACAAAAAAAAAAEAAAAFaGVsbG8AAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAACYloAAAAAAAAAAAXc7DToAAABAS2+MaPA79AjD0B7qjl0qEz0N6CkDmoS4kgnXjZfbvdc9IkqNm0S+vKBNgV80pSfixY147L+jvS/ganovqbLiAQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACU3NL8AAAAAcAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACU3NL8AAAAAcAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyr2OlUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyscX/UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJTc0vwAAAABwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJS2rVwAAAABwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+OcAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+M4AAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{S2+MaPA79AjD0B7qjl0qEz0N6CkDmoS4kgnXjZfbvdc9IkqNm0S+vKBNgV80pSfixY147L+jvS/ganovqbLiAQ==}', 'text', 'hello', NULL, true, 100);
INSERT INTO history_transactions VALUES ('3b36ecfbcc2adb0cfff08ae86199f64e12984f084bb03be9bb249611df82322b', 8, 3, 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB', 30064771075, 100, 1, '2019-06-03 16:35:34.000216', '2019-06-03 16:35:34.000216', 34359750656, 'AAAAADnqxUES0av4d13+BBKxKs+6A/cJNtM1/nf4eHV3Ow06AAAAZAAAAAcAAAADAAAAAAAAAAMBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQAAAAEAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAJiWgAAAAAAAAAABdzsNOgAAAEDC9hMtMYZ6hbx1iAdXngRcCYQmf8eu4zcB9SLH2998tVYca6QYig5Dsgy2oCMD1J7khIL9jz/VWjcPhvTVvC8L', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACUtq1cAAAAAcAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACUtq1cAAAAAcAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyscX/UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytChZUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJS2rVwAAAABwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJSQh7wAAAABwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+M4AAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+LUAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{wvYTLTGGeoW8dYgHV54EXAmEJn/HruM3AfUix9vffLVWHGukGIoOQ7IMtqAjA9Se5ISC/Y8/1Vo3D4b01bwvCw==}', 'hash', 'AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE=', NULL, true, 100);
INSERT INTO history_transactions VALUES ('e14885cb66af5f7f5e991b014eec475c61cc831292cf5526cdd0cda145300837', 8, 4, 'GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB', 30064771076, 100, 1, '2019-06-03 16:35:34.000411', '2019-06-03 16:35:34.000411', 34359754752, 'AAAAADnqxUES0av4d13+BBKxKs+6A/cJNtM1/nf4eHV3Ow06AAAAZAAAAAcAAAAEAAAAAAAAAAQCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgAAAAEAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAAAJiWgAAAAAAAAAABdzsNOgAAAEBOfq9PQ8EGcpjRWEaqGxvhBjSVuk6K5A2rthLYHnmAXmQ1JjJD3EddjiES3bPZUF5efGQvRjoEKgiB2dU3f2wF', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACUkIe8AAAAAcAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACUkIe8AAAAAcAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytChZUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqytoqzUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJSQh7wAAAABwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJRqYhwAAAABwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+LUAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+JwAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Tn6vT0PBBnKY0VhGqhsb4QY0lbpOiuQNq7YS2B55gF5kNSYyQ9xHXY4hEt2z2VBeXnxkL0Y6BCoIgdnVN39sBQ==}', 'return', 'AgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgI=', NULL, true, 100);
INSERT INTO history_transactions VALUES ('0fb9c2e20946222b23e1d1d660de9d74576c41cfd9b199f9d565a013c1ef89ca', 7, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 3, 100, 1, '2019-06-03 16:35:34.011682', '2019-06-03 16:35:34.011682', 30064775168, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAOerFQRLRq/h3Xf4EErEqz7oD9wk20zX+d/h4dXc7DToAAAACVAvkAAAAAAAAAAABVvwF9wAAAED8tIFyog9OeCqiaBNfxFdAlneNYTfjoNUMKi6FJCY5BqemnDBxGox3jKS/xx4zpxAToEFp3Y2M+NRJIU4g/H0J', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lau/0w21AAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lau/0w21AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDbUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyrQFLUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHAAAAAAAAAAA56sVBEtGr+Hdd/gQSsSrPugP3CTbTNf53+Hh1dzsNOgAAAAJUC+QAAAAABwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDc4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDbUAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{/LSBcqIPTngqomgTX8RXQJZ3jWE346DVDCouhSQmOQanppwwcRqMd4ykv8ceM6cQE6BBad2NjPjUSSFOIPx9CQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('a9085e13fbe9f84e07e320a0d445536de1afc2cfd8c7e4186687807edd2b4897', 6, 1, 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB', 17179869188, 100, 1, '2019-06-03 16:35:34.021962', '2019-06-03 16:35:34.021962', 25769807872, 'AAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAZAAAAAQAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAEAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAinuINUAAABA4PRAe0en/05ZH2leCeTOsxbT0cUu3wgUiWUcuDk4ya8G/gI90hlV6pzOYyAB6Zt5fN7pRrPRL/tTlnjgUAjaBvwNxUcAAABAFmdGR6JZukKJUC3Vr2YEJ/24G3tesqTv4cV5UcAozRhS2+w0PYVVqe7QTmOMNSGX/C3LxP1tSvpXdU/OhYsODw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvicAAAAAQAAAADAAAAAQAAAAAAAAAAAAAAAAECAgIAAAABAAAAAPZPnUyLZ+OYJjhn5Hkk43UuW6rOuemZPFQldOn8DcVHAAAAAQAAAAAAAAAAAAAAAQAAAAYAAAAAAAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAAlQL4nAAAAAEAAAABAAAAAEAAAAAAAAAAAAAAAABAgICAAAAAQAAAAD2T51Mi2fjmCY4Z+R5JON1LluqzrnpmTxUJXTp/A3FRwAAAAEAAAAAAAAAAAAAAAEAAAACAAAAAwAAAAYAAAAAAAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAAlQL4nAAAAAEAAAABAAAAAEAAAAAAAAAAAAAAAABAgICAAAAAQAAAAD2T51Mi2fjmCY4Z+R5JON1LluqzrnpmTxUJXTp/A3FRwAAAAEAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+JwAAAABAAAAAQAAAABAAAAAAAAAAAAAAAAAgICAgAAAAEAAAAA9k+dTItn45gmOGfkeSTjdS5bqs656Zk8VCV06fwNxUcAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAMAAAABAAAAAAAAAAAAAAAAAQICAgAAAAEAAAAA9k+dTItn45gmOGfkeSTjdS5bqs656Zk8VCV06fwNxUcAAAABAAAAAAAAAAAAAAABAAAABgAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvicAAAAAQAAAADAAAAAQAAAAAAAAAAAAAAAAECAgIAAAABAAAAAPZPnUyLZ+OYJjhn5Hkk43UuW6rOuemZPFQldOn8DcVHAAAAAQAAAAAAAAAA', '{4PRAe0en/05ZH2leCeTOsxbT0cUu3wgUiWUcuDk4ya8G/gI90hlV6pzOYyAB6Zt5fN7pRrPRL/tTlnjgUAjaBg==,FmdGR6JZukKJUC3Vr2YEJ/24G3tesqTv4cV5UcAozRhS2+w0PYVVqe7QTmOMNSGX/C3LxP1tSvpXdU/OhYsODw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('e9d1a3000aea36743142f2ede106d3cb37c3d7e88508e3f21b496370b5863858', 5, 1, 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB', 17179869185, 100, 1, '2019-06-03 16:35:34.03624', '2019-06-03 16:35:34.03624', 21474840576, 'AAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAZAAAAAQAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAEAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAASnuINUAAABASz0AtZeNzXSXkjPkKJfOE8aUTAuPR6pxMMbF337wxE3wzOTDaVcDQ2N5P3E9MKc+fbbFhZ9K+07+J0wMGltRBA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvi1AAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvi1AAAAAQAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+QAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+OcAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Sz0AtZeNzXSXkjPkKJfOE8aUTAuPR6pxMMbF337wxE3wzOTDaVcDQ2N5P3E9MKc+fbbFhZ9K+07+J0wMGltRBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('995b9269f9f9c4c1eace75501188766d6e8ae40c5413120811a50437683cb74c', 5, 2, 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB', 17179869186, 100, 1, '2019-06-03 16:35:34.036593', '2019-06-03 16:35:34.036593', 21474844672, 'AAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAZAAAAAQAAAACAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAD2T51Mi2fjmCY4Z+R5JON1LluqzrnpmTxUJXTp/A3FRwAAAAEAAAAAAAAAASnuINUAAABADpkMMc7kkkYjDoPwfUlOE9tLYvWHI/m+BBe/gCKN1cVvEF1UBVeCCuGBTjury4TqoxplKl4NZHJST5/Orr4XCA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvi1AAAAAQAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvi1AAAAAQAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAEAAAAA9k+dTItn45gmOGfkeSTjdS5bqs656Zk8VCV06fwNxUcAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+OcAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+M4AAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{DpkMMc7kkkYjDoPwfUlOE9tLYvWHI/m+BBe/gCKN1cVvEF1UBVeCCuGBTjury4TqoxplKl4NZHJST5/Orr4XCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('f78dca926455579b4a43009ffe35a0229a6da4bed32d3c999d7a06ad26605a25', 5, 3, 'GDXFAGJCSCI4CK2YHK6YRLA6TKEXFRX7BMGVMQOBMLIEUJRJ5YQNLMIB', 17179869187, 100, 1, '2019-06-03 16:35:34.036966', '2019-06-03 16:35:34.036966', 21474848768, 'AAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAZAAAAAQAAAADAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAABAAAAAgAAAAEAAAACAAAAAQAAAAIAAAAAAAAAAAAAAAAAAAABKe4g1QAAAEDglRRymtLjw+ImmGwTiBTKE7X7+2CywlHw8qed+t520SbAggcqboy5KXJaEP51/wRSMxtZUgDOFfaDn9Df04EA', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvi1AAAAAQAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAABAAAAAPZPnUyLZ+OYJjhn5Hkk43UuW6rOuemZPFQldOn8DcVHAAAAAQAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAAlQL4tQAAAAEAAAAAwAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAAD2T51Mi2fjmCY4Z+R5JON1LluqzrnpmTxUJXTp/A3FRwAAAAEAAAAAAAAAAAAAAAEAAAACAAAAAwAAAAUAAAAAAAAAAO5QGSKQkcErWDq9iKwemolyxv8LDVZBwWLQSiYp7iDVAAAAAlQL4tQAAAAEAAAAAwAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAAD2T51Mi2fjmCY4Z+R5JON1LluqzrnpmTxUJXTp/A3FRwAAAAEAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAMAAAABAAAAAAAAAAAAAAAAAQICAgAAAAEAAAAA9k+dTItn45gmOGfkeSTjdS5bqs656Zk8VCV06fwNxUcAAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+M4AAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+LUAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{4JUUcprS48PiJphsE4gUyhO1+/tgssJR8PKnnfredtEmwIIHKm6MuSlyWhD+df8EUjMbWVIAzhX2g5/Q39OBAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('66e27fb28870cb5256ea92764bcb222adbbaa5fec2d89a62a9aa8c9c8e2ee9e9', 4, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 2, 100, 1, '2019-06-03 16:35:34.052993', '2019-06-03 16:35:34.052994', 17179873280, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA7lAZIpCRwStYOr2IrB6aiXLG/wsNVkHBYtBKJinuINUAAAACVAvkAAAAAAAAAAABVvwF9wAAAECAUpO+hxiga/YgRsV3rFpBJydgOyn0TPImJCaQCMikkiG+sNXrQBsYXjJrlOiGjGsU3rk4uvGl85AriYD9PNYH', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaxU1gbOAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaxU1gbOAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBs4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDc4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAADuUBkikJHBK1g6vYisHpqJcsb/Cw1WQcFi0EomKe4g1QAAAAJUC+QAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBucAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBs4AAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{gFKTvocYoGv2IEbFd6xaQScnYDsp9EzyJiQmkAjIpJIhvrDV60AbGF4ya5TohoxrFN65OLrxpfOQK4mA/TzWBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4657f7ab2fd82ae203f04d209e6adec0e6bc4f0983b4fc3fa679820ed47e29d7', 3, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-03 16:35:34.065384', '2019-06-03 16:35:34.065384', 12884905984, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAQAAAAAAAABkAAAAAF4L0vAAAAAAAAAAAQAAAAAAAAAAAAAAAC6N7oJcJiUzTWRDL98Bj3fVrJUB19wFvCzEHh8nn/IOAAAAAlQL5AAAAAAAAAAAAVb8BfcAAABA8CyjzEXXVTMwnZTAbHfJeq2HCFzAWkU98ds2ZXFqjXR4EiN0YDSAb/pJwXc0TjMa//SiX83UvUFSqLa8hOXICQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBucAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAuje6CXCYlM01kQy/fAY931ayVAdfcBbwsxB4fJ5/yDgAAAAJUC+QAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{8CyjzEXXVTMwnZTAbHfJeq2HCFzAWkU98ds2ZXFqjXR4EiN0YDSAb/pJwXc0TjMa//SiX83UvUFSqLa8hOXICQ==}', 'none', NULL, '[100,1577833200)', true, 100);


--
-- Name: asset_stats asset_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY asset_stats
    ADD CONSTRAINT asset_stats_pkey PRIMARY KEY (id);


--
-- Name: exp_asset_stats exp_asset_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY exp_asset_stats
    ADD CONSTRAINT exp_asset_stats_pkey PRIMARY KEY (asset_code, asset_issuer, asset_type);


--
-- Name: gorp_migrations gorp_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY gorp_migrations
    ADD CONSTRAINT gorp_migrations_pkey PRIMARY KEY (id);


--
-- Name: history_assets history_assets_asset_code_asset_type_asset_issuer_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_assets
    ADD CONSTRAINT history_assets_asset_code_asset_type_asset_issuer_key UNIQUE (asset_code, asset_type, asset_issuer);


--
-- Name: history_assets history_assets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_assets
    ADD CONSTRAINT history_assets_pkey PRIMARY KEY (id);


--
-- Name: history_operation_participants history_operation_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_operation_participants
    ADD CONSTRAINT history_operation_participants_pkey PRIMARY KEY (id);


--
-- Name: history_transaction_participants history_transaction_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_transaction_participants
    ADD CONSTRAINT history_transaction_participants_pkey PRIMARY KEY (id);


--
-- Name: key_value_store key_value_store_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY key_value_store
    ADD CONSTRAINT key_value_store_pkey PRIMARY KEY (key);


--
-- Name: offers offers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY offers
    ADD CONSTRAINT offers_pkey PRIMARY KEY (offer_id);


--
-- Name: trust_lines trust_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY trust_lines
    ADD CONSTRAINT trust_lines_pkey PRIMARY KEY (ledger_key);


--
-- Name: accounts_data_account_id_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX accounts_data_account_id_name ON accounts_data USING btree (account_id, name);


--
-- Name: accounts_home_domain; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accounts_home_domain ON accounts USING btree (home_domain);


--
-- Name: accounts_inflation_destination; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accounts_inflation_destination ON accounts USING btree (inflation_destination);


--
-- Name: asset_by_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX asset_by_code ON history_assets USING btree (asset_code);


--
-- Name: asset_by_issuer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX asset_by_issuer ON history_assets USING btree (asset_issuer);


--
-- Name: by_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_account ON history_transactions USING btree (account, account_sequence);


--
-- Name: by_fee_account; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_fee_account ON history_transactions USING btree (fee_account) WHERE fee_account IS NOT NULL;

--
-- Name: by_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_hash ON public.history_transactions USING btree (transaction_hash);


--
-- Name: by_inner_hash; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX by_inner_hash ON history_transactions USING btree (inner_transaction_hash) WHERE inner_transaction_hash IS NOT NULL;

--
-- Name: by_ledger; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_ledger ON history_transactions USING btree (ledger_sequence, application_order);


--
-- Name: exp_asset_stats_by_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX exp_asset_stats_by_code ON exp_asset_stats USING btree (asset_code);


--
-- Name: exp_asset_stats_by_issuer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX exp_asset_stats_by_issuer ON exp_asset_stats USING btree (asset_issuer);


--
-- Name: hist_e_by_order; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_e_by_order ON history_effects USING btree (history_operation_id, "order");


--
-- Name: hist_e_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_e_id ON history_effects USING btree (history_account_id, history_operation_id, "order");


--
-- Name: hist_op_p_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_op_p_id ON history_operation_participants USING btree (history_account_id, history_operation_id);


--
-- Name: hist_tx_p_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_tx_p_id ON history_transaction_participants USING btree (history_account_id, history_transaction_id);


--
-- Name: hop_by_hoid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX hop_by_hoid ON history_operation_participants USING btree (history_operation_id);


--
-- Name: hs_ledger_by_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hs_ledger_by_id ON history_ledgers USING btree (id);


--
-- Name: hs_transaction_by_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hs_transaction_by_id ON history_transactions USING btree (id);


--
-- Name: htp_by_htid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htp_by_htid ON history_transaction_participants USING btree (history_transaction_id);


--
-- Name: htrd_by_base_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_base_account ON history_trades USING btree (base_account_id);


--
-- Name: htrd_by_base_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_base_offer ON history_trades USING btree (base_offer_id);


--
-- Name: htrd_by_counter_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_counter_account ON history_trades USING btree (counter_account_id);


--
-- Name: htrd_by_counter_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_counter_offer ON history_trades USING btree (counter_offer_id);


--
-- Name: htrd_by_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_offer ON history_trades USING btree (offer_id);


--
-- Name: htrd_counter_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_counter_lookup ON history_trades USING btree (counter_asset_id);


--
-- Name: htrd_pair_time_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_pair_time_lookup ON history_trades USING btree (base_asset_id, counter_asset_id, ledger_closed_at);


--
-- Name: htrd_pid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX htrd_pid ON history_trades USING btree (history_operation_id, "order");


--
-- Name: htrd_time_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_time_lookup ON history_trades USING btree (ledger_closed_at);


--
-- Name: index_history_accounts_on_address; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_accounts_on_address ON history_accounts USING btree (address);


--
-- Name: index_history_accounts_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_accounts_on_id ON history_accounts USING btree (id);


--
-- Name: index_history_effects_on_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_effects_on_type ON history_effects USING btree (type);


--
-- Name: index_history_ledgers_on_closed_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_ledgers_on_closed_at ON history_ledgers USING btree (closed_at);


--
-- Name: index_history_ledgers_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_id ON history_ledgers USING btree (id);


--
-- Name: index_history_ledgers_on_importer_version; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_ledgers_on_importer_version ON history_ledgers USING btree (importer_version);


--
-- Name: index_history_ledgers_on_ledger_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_ledger_hash ON history_ledgers USING btree (ledger_hash);


--
-- Name: index_history_ledgers_on_previous_ledger_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_previous_ledger_hash ON history_ledgers USING btree (previous_ledger_hash);


--
-- Name: index_history_ledgers_on_sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_sequence ON history_ledgers USING btree (sequence);


--
-- Name: index_history_operations_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_operations_on_id ON history_operations USING btree (id);


--
-- Name: index_history_operations_on_transaction_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_operations_on_transaction_id ON history_operations USING btree (transaction_id);


--
-- Name: index_history_operations_on_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_operations_on_type ON history_operations USING btree (type);


--
-- Name: index_history_transactions_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_transactions_on_id ON history_transactions USING btree (id);


--
-- Name: offers_by_buying_asset; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX offers_by_buying_asset ON offers USING btree (buying_asset);


--
-- Name: offers_by_last_modified_ledger; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX offers_by_last_modified_ledger ON offers USING btree (last_modified_ledger);


--
-- Name: offers_by_seller; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX offers_by_seller ON offers USING btree (seller_id);


--
-- Name: offers_by_selling_asset; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX offers_by_selling_asset ON offers USING btree (selling_asset);


--
-- Name: signers_by_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX signers_by_account ON accounts_signers USING btree (account_id);


--
-- Name: trade_effects_by_order_book; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trade_effects_by_order_book ON history_effects USING btree (((details ->> 'sold_asset_type'::text)), ((details ->> 'sold_asset_code'::text)), ((details ->> 'sold_asset_issuer'::text)), ((details ->> 'bought_asset_type'::text)), ((details ->> 'bought_asset_code'::text)), ((details ->> 'bought_asset_issuer'::text))) WHERE (type = 33);


--
-- Name: trust_lines_by_account_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trust_lines_by_account_id ON trust_lines USING btree (account_id);


--
-- Name: trust_lines_by_issuer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trust_lines_by_issuer ON trust_lines USING btree (asset_issuer);


--
-- Name: trust_lines_by_type_code_issuer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trust_lines_by_type_code_issuer ON trust_lines USING btree (asset_type, asset_code, asset_issuer);


--
-- Name: asset_stats asset_stats_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY asset_stats
    ADD CONSTRAINT asset_stats_id_fkey FOREIGN KEY (id) REFERENCES history_assets(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: history_trades history_trades_base_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_trades
    ADD CONSTRAINT history_trades_base_account_id_fkey FOREIGN KEY (base_account_id) REFERENCES history_accounts(id);


--
-- Name: history_trades history_trades_base_asset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_trades
    ADD CONSTRAINT history_trades_base_asset_id_fkey FOREIGN KEY (base_asset_id) REFERENCES history_assets(id);


--
-- Name: history_trades history_trades_counter_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_trades
    ADD CONSTRAINT history_trades_counter_account_id_fkey FOREIGN KEY (counter_account_id) REFERENCES history_accounts(id);


--
-- Name: history_trades history_trades_counter_asset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY history_trades
    ADD CONSTRAINT history_trades_counter_asset_id_fkey FOREIGN KEY (counter_asset_id) REFERENCES history_assets(id);


-- needed for migrations to work
ALTER TABLE accounts ADD sponsor TEXT;
CREATE INDEX accounts_by_sponsor ON accounts USING BTREE(sponsor);

ALTER TABLE accounts_data ADD sponsor TEXT;
CREATE INDEX accounts_data_by_sponsor ON accounts_data USING BTREE(sponsor);

ALTER TABLE trust_lines ADD sponsor TEXT;
CREATE INDEX trust_lines_by_sponsor ON trust_lines USING BTREE(sponsor);

ALTER TABLE offers ADD sponsor TEXT;
CREATE INDEX offers_by_sponsor ON offers USING BTREE(sponsor);

ALTER TABLE history_operation_participants
    DROP COLUMN id;

ALTER TABLE history_transaction_participants
    DROP COLUMN id;

DROP TABLE asset_stats cascade;

DROP INDEX exp_asset_stats_by_code;

DROP INDEX index_history_transactions_on_id;

DROP INDEX index_history_ledgers_on_id;

DROP INDEX asset_by_code;

ALTER TABLE history_transactions ADD account_muxed varchar(69) NULL, ADD fee_account_muxed varchar(69) NULL;
ALTER TABLE history_operations ADD source_account_muxed varchar(69) NULL;
ALTER TABLE history_effects ADD address_muxed varchar(69) NULL;


-- 49
CREATE INDEX IF NOT EXISTS htrd_agg_timestamp_brin ON history_trades_60000 USING brin(timestamp);

-- 50
CREATE TABLE liquidity_pools (
    id text NOT NULL, -- hex-encoded PoolID
    type smallint NOT NULL,
    fee integer NOT NULL,
    trustline_count bigint NOT NULL CHECK (trustline_count > 0),
    share_count bigint NOT NULL DEFAULT 0 CHECK(share_count >= 0),
    asset_reserves jsonb NOT NULL,
    last_modified_ledger integer NOT NULL,
    deleted boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);

CREATE INDEX liquidity_pools_by_asset_reserves ON liquidity_pools USING gin(asset_reserves jsonb_path_ops);
CREATE INDEX live_liquidity_pools ON liquidity_pools USING BTREE (deleted, last_modified_ledger);

CREATE SEQUENCE history_liquidity_pools_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE history_liquidity_pools (
    id bigint NOT NULL DEFAULT nextval('history_liquidity_pools_id_seq'::regclass),
    liquidity_pool_id text NOT NULL
);

CREATE UNIQUE INDEX index_history_liquidity_pools_on_id ON history_liquidity_pools USING btree (id);
CREATE UNIQUE INDEX index_history_liquidity_pools_on_liquidity_pool_id ON history_liquidity_pools USING btree (liquidity_pool_id);

CREATE TABLE history_operation_liquidity_pools (
    history_operation_id bigint NOT NULL,
    history_liquidity_pool_id bigint NOT NULL
);

CREATE UNIQUE INDEX index_history_operation_liquidity_pools_on_ids ON history_operation_liquidity_pools USING btree (history_operation_id , history_liquidity_pool_id);
CREATE INDEX index_history_operation_liquidity_pools_on_operation_id ON history_operation_liquidity_pools USING btree (history_operation_id);

CREATE TABLE history_transaction_liquidity_pools (
    history_transaction_id bigint NOT NULL,
    history_liquidity_pool_id bigint NOT NULL
);

CREATE UNIQUE INDEX index_history_transaction_liquidity_pools_on_ids ON history_transaction_liquidity_pools USING btree (history_transaction_id , history_liquidity_pool_id);
CREATE INDEX index_history_transaction_liquidity_pools_on_transaction_id ON history_transaction_liquidity_pools USING btree (history_transaction_id);

ALTER TABLE trust_lines ADD liquidity_pool_id text;
CREATE INDEX trust_lines_by_liquidity_pool_id ON trust_lines USING BTREE(liquidity_pool_id);

DROP INDEX htrd_by_offer;
DROP INDEX htrd_counter_lookup;

ALTER TABLE history_trades DROP offer_id,
                           ALTER base_account_id DROP NOT NULL,
                           ALTER counter_account_id DROP NOT NULL,
                           ADD base_liquidity_pool_id bigint,
                           ADD counter_liquidity_pool_id bigint,
                           ADD liquidity_pool_fee int;

CREATE INDEX htrd_by_base_liquidity_pool_id ON history_trades USING BTREE(base_liquidity_pool_id);
CREATE INDEX htrd_by_counter_liquidity_pool_id ON history_trades USING BTREE(counter_liquidity_pool_id);

-- 51
DROP INDEX IF EXISTS by_account;
DROP INDEX IF EXISTS by_fee_account;

--
-- PostgreSQL database dump complete
--
