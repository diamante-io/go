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
DROP INDEX IF EXISTS public.trade_effects_by_order_book;
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
DROP INDEX IF EXISTS public.by_ledger;
DROP INDEX IF EXISTS public.by_hash;
DROP INDEX IF EXISTS public.by_account;
DROP INDEX IF EXISTS public.asset_by_issuer;
DROP INDEX IF EXISTS public.asset_by_code;
ALTER TABLE IF EXISTS ONLY public.history_transaction_participants DROP CONSTRAINT IF EXISTS history_transaction_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_operation_participants DROP CONSTRAINT IF EXISTS history_operation_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_asset_code_asset_type_asset_issuer_key;
ALTER TABLE IF EXISTS ONLY public.gorp_migrations DROP CONSTRAINT IF EXISTS gorp_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.asset_stats DROP CONSTRAINT IF EXISTS asset_stats_pkey;
ALTER TABLE IF EXISTS public.history_transaction_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_operation_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_assets ALTER COLUMN id DROP DEFAULT;
DROP TABLE IF EXISTS public.history_transactions;
DROP SEQUENCE IF EXISTS public.history_transaction_participants_id_seq;
DROP TABLE IF EXISTS public.history_transaction_participants;
DROP TABLE IF EXISTS public.history_trades;
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
DROP TABLE IF EXISTS public.asset_stats;
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
    failed_transaction_count integer
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
    CONSTRAINT history_trades_base_amount_check CHECK ((base_amount > 0)),
    CONSTRAINT history_trades_check CHECK ((base_asset_id < counter_asset_id)),
    CONSTRAINT history_trades_counter_amount_check CHECK ((counter_amount > 0))
);


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

CREATE TABLE history_transactions (
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
    fee_charged integer
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
-- Data for Name: asset_stats; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO asset_stats VALUES (1, '1009876000', 1, 1, 'https://test.com/.well-known/diamcircle.toml');
INSERT INTO asset_stats VALUES (2, '3000010434000', 2, 1, 'https://test.com/.well-known/diamcircle.toml');
INSERT INTO asset_stats VALUES (3, '10000000000', 1, 2, '');


--
-- Data for Name: gorp_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO gorp_migrations VALUES ('1_initial_schema.sql', '2019-06-03 18:28:47.032496+02');
INSERT INTO gorp_migrations VALUES ('2_index_participants_by_toid.sql', '2019-06-03 18:28:47.039657+02');
INSERT INTO gorp_migrations VALUES ('3_use_sequence_in_history_accounts.sql', '2019-06-03 18:28:47.044048+02');
INSERT INTO gorp_migrations VALUES ('4_add_protocol_version.sql', '2019-06-03 18:28:47.054532+02');
INSERT INTO gorp_migrations VALUES ('5_create_trades_table.sql', '2019-06-03 18:28:47.063028+02');
INSERT INTO gorp_migrations VALUES ('6_create_assets_table.sql', '2019-06-03 18:28:47.068415+02');
INSERT INTO gorp_migrations VALUES ('7_modify_trades_table.sql', '2019-06-03 18:28:47.081625+02');
INSERT INTO gorp_migrations VALUES ('8_create_asset_stats_table.sql', '2019-06-03 18:28:47.087463+02');
INSERT INTO gorp_migrations VALUES ('8_add_aggregators.sql', '2019-06-03 18:28:47.090109+02');
INSERT INTO gorp_migrations VALUES ('9_add_header_xdr.sql', '2019-06-03 18:28:47.092718+02');
INSERT INTO gorp_migrations VALUES ('10_add_trades_price.sql', '2019-06-03 18:28:47.095973+02');
INSERT INTO gorp_migrations VALUES ('11_add_trades_account_index.sql', '2019-06-03 18:28:47.099698+02');
INSERT INTO gorp_migrations VALUES ('12_asset_stats_amount_string.sql', '2019-06-03 18:28:47.107549+02');
INSERT INTO gorp_migrations VALUES ('13_trade_offer_ids.sql', '2019-06-03 18:28:47.112768+02');
INSERT INTO gorp_migrations VALUES ('14_fix_asset_toml_field.sql', '2019-06-03 18:28:47.115116+02');
INSERT INTO gorp_migrations VALUES ('15_ledger_failed_txs.sql', '2019-06-03 18:28:47.116796+02');
INSERT INTO gorp_migrations VALUES ('16_ingest_failed_transactions.sql', '2019-06-03 18:28:47.117989+02');
INSERT INTO gorp_migrations VALUES ('17_transaction_fee_paid.sql', '2019-06-03 18:28:47.120034+02');


--
-- Data for Name: history_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_accounts VALUES (1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_accounts VALUES (2, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_accounts VALUES (3, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_accounts VALUES (4, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_accounts_id_seq', 4, true);


--
-- Data for Name: history_assets; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_assets VALUES (1, 'credit_alphanum4', 'BTC', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_assets VALUES (2, 'credit_alphanum4', 'USD', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_assets VALUES (3, 'credit_alphanum4', 'SCOT', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');


--
-- Name: history_assets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_assets_id_seq', 3, true);


--
-- Data for Name: history_effects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_effects VALUES (1, 34359742465, 1, 2, '{"amount": "1.3623000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 34359742465, 2, 3, '{"amount": "1.3623000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 34359746561, 1, 2, '{"amount": "31.6577680", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 34359746561, 2, 3, '{"amount": "31.6577680", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 30064775169, 1, 22, '{"limit": "1000000000.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 30064779265, 1, 2, '{"amount": "89.9500000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 30064779265, 2, 3, '{"amount": "89.9500000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 25769807873, 1, 2, '{"amount": "1000.0000000", "asset_code": "SCOT", "asset_type": "credit_alphanum4", "asset_issuer": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (1, 25769807873, 2, 3, '{"amount": "1000.0000000", "asset_code": "SCOT", "asset_type": "credit_alphanum4", "asset_issuer": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (1, 25769811969, 1, 2, '{"amount": "100000.1200000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 25769811969, 2, 3, '{"amount": "100000.1200000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 25769816065, 1, 2, '{"amount": "100.9876000", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 25769816065, 2, 3, '{"amount": "100.9876000", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 25769820161, 1, 2, '{"amount": "200000.9234000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 25769820161, 2, 3, '{"amount": "200000.9234000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 21474840577, 1, 23, '{"trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 21474844673, 1, 23, '{"trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 21474848769, 1, 23, '{"trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 17179873281, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 17179877377, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 17179881473, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 17179885569, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "SCOT", "asset_type": "credit_alphanum4", "asset_issuer": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (1, 12884905985, 1, 6, '{"auth_revocable_flag": true}');
INSERT INTO history_effects VALUES (3, 12884910081, 1, 5, '{"home_domain": "test.com"}');
INSERT INTO history_effects VALUES (3, 12884914177, 1, 6, '{"auth_required_flag": true}');
INSERT INTO history_effects VALUES (3, 8589938689, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (4, 8589938689, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (3, 8589938689, 3, 10, '{"weight": 1, "public_key": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 8589942785, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (4, 8589942785, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 8589942785, 3, 10, '{"weight": 1, "public_key": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (2, 8589946881, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (4, 8589946881, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 8589946881, 3, 10, '{"weight": 1, "public_key": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_ledgers VALUES (8, 'b126571afcaf19157b639f5eba881e25930d94a37a52815db11819ccd5340b94', '1491ff6cbfb702b38010befb6f0a64d62293bc00539c24a4cfe3038bf7174339', 2, 2, '2019-06-03 16:37:47', '2019-06-03 16:37:45.007637', '2019-06-03 16:37:45.007637', 34359738368, 16, 1000000000000000000, 2100, 100, 100000000, 1000000, 11, 'AAAACxSR/2y/twKzgBC++28KZNYik7wAU5wkpM/jA4v3F0M5SrnXL2P8lSzo6IVDmtY1ZEbJBT/DkKVVD6JAdZ0hcmsAAAAAXPVM2wAAAAAAAAAAZxdQWfDM1MxUCuefAme/bjqOchzkHE89OsTko56c4I4iAccXBST+f0Cn0xFhOxpER1nlhqjiIIsMXna2VGlwLAAAAAgN4Lazp2QAAAAAAAAAAAg0AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO history_ledgers VALUES (7, '1491ff6cbfb702b38010befb6f0a64d62293bc00539c24a4cfe3038bf7174339', 'e48cf2fed5e3b8fbf715a418f92f81b8362c8804fc334e73d59b9ec99ee4224b', 2, 2, '2019-06-03 16:37:46', '2019-06-03 16:37:45.030701', '2019-06-03 16:37:45.030702', 30064771072, 16, 1000000000000000000, 1900, 100, 100000000, 1000000, 11, 'AAAAC+SM8v7V47j79xWkGPkvgbg2LIgE/DNOc9Wbnsme5CJLT03OgiB/+MODRPQfafM9wx2p409pQvaZdj1O67SYKJ0AAAAAXPVM2gAAAAAAAAAAp6X4a1fQWRN1fT4Oha1VAb23rxjrafhg/8wjyOPQ4elbWMK8Kf5rrB0MVt6oD2Tl6w8vCi+f4g0Og6SINTYEDQAAAAcN4Lazp2QAAAAAAAAAAAdsAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO history_ledgers VALUES (6, 'e48cf2fed5e3b8fbf715a418f92f81b8362c8804fc334e73d59b9ec99ee4224b', '333ec5e9ed90bbebb591729290b241ba6530ddf9f6102e0244f2e5b79721c1e4', 4, 4, '2019-06-03 16:37:45', '2019-06-03 16:37:45.043381', '2019-06-03 16:37:45.043381', 25769803776, 16, 1000000000000000000, 1700, 100, 100000000, 1000000, 11, 'AAAACzM+xentkLvrtZFykpCyQbplMN359hAuAkTy5beXIcHku8Zf/iDsmpUtzhUhvCkZxGvlhB58WJOSt+WZKn6exCUAAAAAXPVM2QAAAAAAAAAAXbLrzjyhc09OMR0rGIcNlJyvUas57e4z2TpUhlLwGBrLfhuK3ISMsCuVJlDnaAm7yhR0gRbaia/WaeLEQ9kV8gAAAAYN4Lazp2QAAAAAAAAAAAakAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0);
INSERT INTO history_ledgers VALUES (5, '333ec5e9ed90bbebb591729290b241ba6530ddf9f6102e0244f2e5b79721c1e4', 'c06da07d9df310bcff40b5fd1f635316526a35aa2170720d7e7af6a1231bc889', 3, 3, '2019-06-03 16:37:44', '2019-06-03 16:37:45.06041', '2019-06-03 16:37:45.06041', 21474836480, 16, 1000000000000000000, 1300, 100, 100000000, 1000000, 11, 'AAAAC8BtoH2d8xC8/0C1/R9jUxZSajWqIXByDX569qEjG8iJbKMKb88bKwjF/tf5B2srqlB1Wt1wpqMZDxD5nC0zt2YAAAAAXPVM2AAAAAAAAAAAAvTWoAZ8sb3n9zINDA1hqoanXiN/pcjfEtKLDrzVvDNudBo4h9aHdtfWXTsgJUG48KFTs4PTWToRXu/VLVdQ5wAAAAUN4Lazp2QAAAAAAAAAAAUUAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO history_ledgers VALUES (4, 'c06da07d9df310bcff40b5fd1f635316526a35aa2170720d7e7af6a1231bc889', '73d45d07ebf473cd3a6ae2d4aa6f56d01f6ffd71d7b3567afd5060b6fc921b56', 4, 4, '2019-06-03 16:37:43', '2019-06-03 16:37:45.072877', '2019-06-03 16:37:45.072877', 17179869184, 16, 1000000000000000000, 1000, 100, 100000000, 1000000, 11, 'AAAAC3PUXQfr9HPNOmri1KpvVtAfb/1x17NWev1QYLb8khtWr+ch0bnpNLHQ9oKhg8ODAYI15MnucRB9qXMxtHv/+0cAAAAAXPVM1wAAAAAAAAAAg4bDLmPB7NuPhcFNu33PjSojcPo2TH0x5yLLubhqMjzJvBY5PAIvHSKyu1oe3+8JuK/LnHyRVafV32eOU81flAAAAAQN4Lazp2QAAAAAAAAAAAPoAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0);
INSERT INTO history_ledgers VALUES (3, '73d45d07ebf473cd3a6ae2d4aa6f56d01f6ffd71d7b3567afd5060b6fc921b56', '4a4b1ff214c4817fd203885fdf554c8d86271d6788616565150d7de7566fdf98', 3, 3, '2019-06-03 16:37:42', '2019-06-03 16:37:45.085791', '2019-06-03 16:37:45.085792', 12884901888, 16, 1000000000000000000, 600, 100, 100000000, 1000000, 11, 'AAAAC0pLH/IUxIF/0gOIX99VTI2GJx1niGFlZRUNfedWb9+YispBRYjJkjprlnEXQeQh2xZxp4jvWl+jqXzxvvNnVZkAAAAAXPVM1gAAAAAAAAAAz4e0cNoSrSyJbiU9ww4Xz/C47nFVIBEHkR11kgme+DzF1UX/tKW/NixJgxyLo68LswymO6lok5PKGtG9vmaZrQAAAAMN4Lazp2QAAAAAAAAAAAJYAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO history_ledgers VALUES (2, '4a4b1ff214c4817fd203885fdf554c8d86271d6788616565150d7de7566fdf98', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 3, 3, '2019-06-03 16:37:41', '2019-06-03 16:37:45.10023', '2019-06-03 16:37:45.10023', 8589934592, 16, 1000000000000000000, 300, 100, 100000000, 1000000, 11, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76Z9GwX6q8KzxXUDk3+fdZshszCWX2SbuAnJNMiSoulCykAAAAAXPVM1QAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAACyk/eWGAYEOhla9SxZHjkxGIQ61pIijmPF9hVp1Qv5rtOXsmkRJcyXmle2ZKFjzhTk0u31dZoIkX51uoncfP4wAAAAIN4Lazp2QAAAAAAAAAAAEsAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-03 16:37:45.111028', '2019-06-03 16:37:45.111029', 4294967296, 16, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operation_participants VALUES (1, 34359742465, 2);
INSERT INTO history_operation_participants VALUES (2, 34359742465, 1);
INSERT INTO history_operation_participants VALUES (3, 34359746561, 1);
INSERT INTO history_operation_participants VALUES (4, 34359746561, 2);
INSERT INTO history_operation_participants VALUES (5, 30064775169, 2);
INSERT INTO history_operation_participants VALUES (6, 30064779265, 2);
INSERT INTO history_operation_participants VALUES (7, 30064779265, 1);
INSERT INTO history_operation_participants VALUES (8, 25769807873, 1);
INSERT INTO history_operation_participants VALUES (9, 25769807873, 2);
INSERT INTO history_operation_participants VALUES (10, 25769811969, 3);
INSERT INTO history_operation_participants VALUES (11, 25769811969, 1);
INSERT INTO history_operation_participants VALUES (12, 25769816065, 3);
INSERT INTO history_operation_participants VALUES (13, 25769816065, 1);
INSERT INTO history_operation_participants VALUES (14, 25769820161, 3);
INSERT INTO history_operation_participants VALUES (15, 25769820161, 2);
INSERT INTO history_operation_participants VALUES (16, 21474840577, 3);
INSERT INTO history_operation_participants VALUES (17, 21474840577, 1);
INSERT INTO history_operation_participants VALUES (18, 21474844673, 3);
INSERT INTO history_operation_participants VALUES (19, 21474844673, 1);
INSERT INTO history_operation_participants VALUES (20, 21474848769, 3);
INSERT INTO history_operation_participants VALUES (21, 21474848769, 2);
INSERT INTO history_operation_participants VALUES (22, 17179873281, 2);
INSERT INTO history_operation_participants VALUES (23, 17179877377, 1);
INSERT INTO history_operation_participants VALUES (24, 17179881473, 1);
INSERT INTO history_operation_participants VALUES (25, 17179885569, 2);
INSERT INTO history_operation_participants VALUES (26, 12884905985, 1);
INSERT INTO history_operation_participants VALUES (27, 12884910081, 3);
INSERT INTO history_operation_participants VALUES (28, 12884914177, 3);
INSERT INTO history_operation_participants VALUES (29, 8589938689, 4);
INSERT INTO history_operation_participants VALUES (30, 8589938689, 3);
INSERT INTO history_operation_participants VALUES (31, 8589942785, 4);
INSERT INTO history_operation_participants VALUES (32, 8589942785, 1);
INSERT INTO history_operation_participants VALUES (33, 8589946881, 4);
INSERT INTO history_operation_participants VALUES (34, 8589946881, 2);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_operation_participants_id_seq', 34, true);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operations VALUES (34359742465, 34359742464, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "amount": "1.3623000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (34359746561, 34359746560, 1, 1, '{"to": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "31.6577680", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (30064775169, 30064775168, 1, 6, '{"limit": "1000000000.0000000", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (30064779265, 30064779264, 1, 1, '{"to": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "89.9500000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (25769807873, 25769807872, 1, 1, '{"to": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "1000.0000000", "asset_code": "SCOT", "asset_type": "credit_alphanum4", "asset_issuer": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (25769811969, 25769811968, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "amount": "100000.1200000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (25769816065, 25769816064, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "amount": "100.9876000", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (25769820161, 25769820160, 1, 1, '{"to": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "from": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "amount": "200000.9234000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (21474840577, 21474840576, 1, 7, '{"trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "authorize": true, "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (21474844673, 21474844672, 1, 7, '{"trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "authorize": true, "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (21474848769, 21474848768, 1, 7, '{"trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "authorize": true, "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (17179873281, 17179873280, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (17179877377, 17179877376, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (17179881473, 17179881472, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "BTC", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (17179885569, 17179885568, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "SCOT", "asset_type": "credit_alphanum4", "asset_issuer": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (12884905985, 12884905984, 1, 5, '{"set_flags": [2], "set_flags_s": ["auth_revocable"]}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (12884910081, 12884910080, 1, 5, '{"home_domain": "test.com"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (12884914177, 12884914176, 1, 5, '{"set_flags": [1], "set_flags_s": ["auth_required"]}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (8589938689, 8589938688, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589942785, 8589942784, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589946881, 8589946880, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transaction_participants VALUES (1, 34359742464, 2);
INSERT INTO history_transaction_participants VALUES (2, 34359742464, 1);
INSERT INTO history_transaction_participants VALUES (3, 34359746560, 1);
INSERT INTO history_transaction_participants VALUES (4, 34359746560, 2);
INSERT INTO history_transaction_participants VALUES (5, 30064775168, 2);
INSERT INTO history_transaction_participants VALUES (6, 30064779264, 1);
INSERT INTO history_transaction_participants VALUES (7, 30064779264, 2);
INSERT INTO history_transaction_participants VALUES (8, 25769807872, 1);
INSERT INTO history_transaction_participants VALUES (9, 25769807872, 2);
INSERT INTO history_transaction_participants VALUES (10, 25769811968, 3);
INSERT INTO history_transaction_participants VALUES (11, 25769811968, 1);
INSERT INTO history_transaction_participants VALUES (12, 25769816064, 1);
INSERT INTO history_transaction_participants VALUES (13, 25769816064, 3);
INSERT INTO history_transaction_participants VALUES (14, 25769820160, 3);
INSERT INTO history_transaction_participants VALUES (15, 25769820160, 2);
INSERT INTO history_transaction_participants VALUES (16, 21474840576, 3);
INSERT INTO history_transaction_participants VALUES (17, 21474840576, 1);
INSERT INTO history_transaction_participants VALUES (18, 21474844672, 3);
INSERT INTO history_transaction_participants VALUES (19, 21474844672, 1);
INSERT INTO history_transaction_participants VALUES (20, 21474848768, 2);
INSERT INTO history_transaction_participants VALUES (21, 21474848768, 3);
INSERT INTO history_transaction_participants VALUES (22, 17179873280, 2);
INSERT INTO history_transaction_participants VALUES (23, 17179877376, 1);
INSERT INTO history_transaction_participants VALUES (24, 17179881472, 1);
INSERT INTO history_transaction_participants VALUES (25, 17179885568, 2);
INSERT INTO history_transaction_participants VALUES (26, 12884905984, 1);
INSERT INTO history_transaction_participants VALUES (27, 12884910080, 3);
INSERT INTO history_transaction_participants VALUES (28, 12884914176, 3);
INSERT INTO history_transaction_participants VALUES (29, 8589938688, 4);
INSERT INTO history_transaction_participants VALUES (30, 8589938688, 3);
INSERT INTO history_transaction_participants VALUES (31, 8589942784, 4);
INSERT INTO history_transaction_participants VALUES (32, 8589942784, 1);
INSERT INTO history_transaction_participants VALUES (33, 8589946880, 2);
INSERT INTO history_transaction_participants VALUES (34, 8589946880, 4);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_transaction_participants_id_seq', 34, true);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transactions VALUES ('142d3dbe5948eb39db1fd62d912ce67131b1b300adb015acf0f17d91a057429d', 8, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934596, 100, 1, '2019-06-03 16:37:45.007926', '2019-06-03 16:37:45.007926', 34359742464, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAz97YAAAAAAAAAAEnciLVAAAAQHbmlPqVcxoIqzJFayddJwGRM8Vxm0BYlui3LVu9d/nB2hb/tsUWgUZLCUnNv/CPjsMTAN2LmVkYOMtCdYc+NQ8=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvicAAAAAIAAAADAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvicAAAAAIAAAAEAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6J8aF6B//////////wAAAAEAAAAAAAAAAAAAAAEAAAAIAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6J/p9nh//////////wAAAAEAAAAAAAAAAAAAAAMAAAAHAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0d90TjAAI4byb8EAAAAAAAEAAAAAAAAAAAAAAAEAAAAIAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0d6kb1gAI4byb8EAAAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+LUAAAAAgAAAAMAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+JwAAAAAgAAAAMAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{duaU+pVzGgirMkVrJ10nAZEzxXGbQFiW6LctW713+cHaFv+2xRaBRksJSc2/8I+OwxMA3YuZWRg4y0J1hz41Dw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('3362c9b76d85a844c739b338dbef4213ce64eca1ceb6c0d70e878975ab1477b1', 8, 2, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934598, 100, 1, '2019-06-03 16:37:45.008333', '2019-06-03 16:37:45.008334', 34359746560, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAS3peQAAAAAAAAAAGu5L5MAAAAQKnjaWS6Rk617nkw1/KuCffaeN1Mymuz8m9Brm0RJ1IYNKdnudV+72HsCM1Vnfnz/+iB6ERFxOsEp1mBHpUMQwk=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvhqAAAAAIAAAAFAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvhqAAAAAIAAAAGAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAIAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0d6kb1gAI4byb8EAAAAAAAEAAAAAAAAAAAAAAAEAAAAIAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0fGDBugAI4byb8EAAAAAAAEAAAAAAAAAAAAAAAMAAAAIAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6J/p9nh//////////wAAAAEAAAAAAAAAAAAAAAEAAAAIAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6I0LXuh//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+IMAAAAAgAAAAUAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+GoAAAAAgAAAAUAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{qeNpZLpGTrXueTDX8q4J99p43UzKa7Pyb0GubREnUhg0p2e51X7vYewIzVWd+fP/6IHoREXE6wSnWYEelQxDCQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('fa17f7c083fddc53e8e28885be934e19bf637e287c1951be581dd05c0be93b56', 7, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934595, 100, 1, '2019-06-03 16:37:45.03086', '2019-06-03 16:37:45.03086', 30064775168, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAADAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAjhvJvwQAAAAAAAAAAAAEnciLVAAAAQNI8SXbUBWJi/xf8bWtBBKonww9YpbLck1/295qxZOYN5vjFDYQLaG3b1aGWqzWZqa9FMHkJ2tAEDPjEHIMkzAw=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvi1AAAAAIAAAACAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvi1AAAAAIAAAADAAAAAgAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAGAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0anXBlB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAHAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0anXBlAAI4byb8EAAAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+LUAAAAAgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{0jxJdtQFYmL/F/xta0EEqifDD1ilstyTX/b3mrFk5g3m+MUNhAtobdvVoZarNZmpr0UweQna0AQM+MQcgyTMDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c42c988a72ac8aed3bb9a7b7dfb96b905e33d1506f4e663360135e6c6e115078', 7, 2, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934597, 100, 1, '2019-06-03 16:37:45.031078', '2019-06-03 16:37:45.031078', 30064779264, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA1nUfgAAAAAAAAAAGu5L5MAAAAQAGYynFy2CKfKZyhmWMLfgmhdJtJHXW7ogTdyZ7aviECOHYJSQKPkcnMoG4N76ipkuVH6hjuxDHBJ83+HnyhbAQ=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAviDAAAAAIAAAAEAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAviDAAAAAIAAAAFAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0anXBlAAI4byb8EAAAAAAAEAAAAAAAAAAAAAAAEAAAAHAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAB0d90TjAAI4byb8EAAAAAAAEAAAAAAAAAAAAAAAMAAAAGAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6NS3X4B//////////wAAAAEAAAAAAAAAAAAAAAEAAAAHAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAA6J8aF6B//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+JwAAAAAgAAAAQAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+IMAAAAAgAAAAQAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{AZjKcXLYIp8pnKGZYwt+CaF0m0kddbuiBN3Jntq+IQI4dglJAo+Rycygbg3vqKmS5UfqGO7EMcEnzf4efKFsBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('cf0f5fcd46881458ba623f9e6e7c52489d4bd3979a4196819882bb6240b4e855', 6, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934596, 100, 1, '2019-06-03 16:37:45.043543', '2019-06-03 16:37:45.043543', 25769807872, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABU0NPVAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+QAAAAAAAAAAAGu5L5MAAAAQLSYQCC1+DGQ8srHLxi6SfnN/dn8t7mAcXlDniU3J+d6Ezg1U6lg9i0jWOsfamioYVbJ9dAiQBZyIsn7TB5cLww=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvicAAAAAIAAAADAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABgAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvicAAAAAIAAAAEAAAAAgAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAAEAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFTQ09UAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAGAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFTQ09UAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAAlQL5AB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAMAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+JwAAAAAgAAAAMAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{tJhAILX4MZDyyscvGLpJ+c392fy3uYBxeUOeJTcn53oTODVTqWD2LSNY6x9qaKhhVsn10CJAFnIiyftMHlwvDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2300600248f841cd5f50276fc18eb16bc88a734e7a290f287ab3a2aa92684826', 6, 2, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934598, 100, 1, '2019-06-03 16:37:45.043709', '2019-06-03 16:37:45.04371', 25769811968, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAOjUt1+AAAAAAAAAAAH5kC3vAAAAQIqp3RfP1ueB0TRJRYXnao+kmde4BDh8q0Ep7q14Q8oRNx1R9utncfpoXr7JOcqiwtgarT9k6KmMyjda97H5RgM=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvg4AAAAAIAAAAFAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+DgAAAAAgAAAAYAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABQAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAAAAAABAAAABgAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAOjUt1+Af/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+IMAAAAAgAAAAUAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAYAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4agAAAACAAAABQAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{iqndF8/W54HRNElFhedqj6SZ17gEOHyrQSnurXhDyhE3HVH262dx+mhevsk5yqLC2BqtP2ToqYzKN1r3sflGAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('5a48a811ec874fc9c5d77c7caeb8abcea076c1baa51b755b2a878391a089c7d1', 6, 3, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934599, 100, 1, '2019-06-03 16:37:45.043846', '2019-06-03 16:37:45.043846', 25769816064, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA8MXwgAAAAAAAAAAH5kC3vAAAAQGEXqpE9OKOxah6oBhR955A4BYmO+yuLNMMtcALlLsKj2M1e9QTlBvAzuwkgECvg2iw8qXZB2kHteYw8qoozcQA=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvg4AAAAAIAAAAGAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+DgAAAAAgAAAAcAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABQAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAAAAAABAAAABgAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA8MXwgf/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAGAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+GoAAAAAgAAAAUAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAYAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4UQAAAACAAAABQAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{YReqkT04o7FqHqgGFH3nkDgFiY77K4s0wy1wAuUuwqPYzV71BOUG8DO7CSAQK+DaLDypdkHaQe15jDyqijNxAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4e4779f0d69db51ec4f7b73387b60c239433804d2747def21b7771e9b71d75be', 6, 4, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934600, 100, 1, '2019-06-03 16:37:45.043977', '2019-06-03 16:37:45.043977', 25769820160, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAdGp1wZQAAAAAAAAAAH5kC3vAAAAQHQFhOcK6JMPYxfRWB+xO13EkPDqkvvPG/Hp8EWDTIMTpHHi4Mqr3/SreJLUxOi3qGSqYFJHiAoK65rFYQaPEAQ=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvg4AAAAAIAAAAHAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+DgAAAAAgAAAAgAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABQAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAAAAAABAAAABgAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAdGp1wZQf/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAGAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+FEAAAAAgAAAAUAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAYAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4OAAAAACAAAABQAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{dAWE5wrokw9jF9FYH7E7XcSQ8OqS+88b8enwRYNMgxOkceLgyqvf9Kt4ktTE6LeoZKpgUkeICgrrmsVhBo8QBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('3b666a253313fc7a0d241ee28064eec78aaa5ebd0a7c0ae7f85259e80fad029f', 5, 1, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934595, 100, 1, '2019-06-03 16:37:45.060663', '2019-06-03 16:37:45.060663', 21474840576, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAADAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAEAAAAAAAAAAfmQLe8AAABAL6czYFvSBhdVeD4fbXOHuXFa2CDqLpFfc+QJnoiPLt/23YViURGLyfg388FKMKsbNJEgmFsCJjtgl3fj7wr/Aw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAviDAAAAAIAAAACAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+IMAAAAAgAAAAMAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAAAAAAAAAAAAAAAAAABAAAABQAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAIAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4tQAAAACAAAAAgAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{L6czYFvSBhdVeD4fbXOHuXFa2CDqLpFfc+QJnoiPLt/23YViURGLyfg388FKMKsbNJEgmFsCJjtgl3fj7wr/Aw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d9d6b816a0a3c640637d48fe33fa00f9ef116103c204834a1c18a9765803fd5d', 5, 2, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934596, 100, 1, '2019-06-03 16:37:45.061037', '2019-06-03 16:37:45.061037', 21474844672, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAEAAAAAAAAAAfmQLe8AAABAMIB8sKelxTqFOLPILjB0nItcfrGrCwursIhshVeKHSw2IC4pmCeg7KGDOLpfUCLc23n5HeTsxJsb/CrHJF/XDQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAviDAAAAAIAAAADAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+IMAAAAAgAAAAQAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAAAAAAAAAAAAAAAAAABAAAABQAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+LUAAAAAgAAAAIAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4nAAAAACAAAAAgAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{MIB8sKelxTqFOLPILjB0nItcfrGrCwursIhshVeKHSw2IC4pmCeg7KGDOLpfUCLc23n5HeTsxJsb/CrHJF/XDQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('6ab66668ea2801de6a7239c94d44e5d41f361812607748125da372b27b66cd3c', 5, 3, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934597, 100, 1, '2019-06-03 16:37:45.061357', '2019-06-03 16:37:45.061357', 21474848768, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAcAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAEAAAAAAAAAAfmQLe8AAABA78VZpv8Z9a3XM9gv6hyMLt2bBrZ5sKsFRU4GKXYtxY2MkAt9J9ENrSRZn1M0jlx9FFGtCvtFFZi8DhxvqDyaBQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAHAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAviDAAAAAIAAAAEAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+IMAAAAAgAAAAUAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAABAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAAAAAAAAAAAAAAAAAABAAAABQAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAAAAAAAf/////////8AAAABAAAAAAAAAAA=', 'AAAAAgAAAAMAAAAFAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+JwAAAAAgAAAAIAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAlQL4gwAAAACAAAAAgAAAAAAAAAAAAAAAQAAAAh0ZXN0LmNvbQEAAAAAAAAAAAAAAAAAAAA=', '{78VZpv8Z9a3XM9gv6hyMLt2bBrZ5sKsFRU4GKXYtxY2MkAt9J9ENrSRZn1M0jlx9FFGtCvtFFZi8DhxvqDyaBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('00ab9cfce2b4c4141d8bb6768dd094bdbb1c7406710dbb3ba0ef98870f63a344', 4, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934593, 100, 1, '2019-06-03 16:37:45.073022', '2019-06-03 16:37:45.073022', 17179873280, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt73//////////AAAAAAAAAAEnciLVAAAAQLVbII+1LeizxgncDI46KHyBt05+H92n1+R328J9zNl2fgJW2nfn3FIoLVs2qV1+CUpr121a2B7AM6HKr4nBLAI=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvjOAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvjOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+OcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{tVsgj7Ut6LPGCdwMjjoofIG3Tn4f3afX5Hfbwn3M2XZ+Albad+fcUigtWzapXX4JSmvXbVrYHsAzocqvicEsAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4486298e04ffb1f3620c521f81adb5207f5d12c21b08a076589d2be3d8dae543', 4, 2, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934594, 100, 1, '2019-06-03 16:37:45.073203', '2019-06-03 16:37:45.073203', 17179877376, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt73//////////AAAAAAAAAAGu5L5MAAAAQFp8rsD4Au1oeZkBT1RHIJRyxWayau3f5UjeA0w4+0LzjLEyi9nGMs8elAH4lDhhDJxCJ8HhxbG+XT/cmQsu1QA=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvi1AAAAAIAAAABAAAAAAAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvi1AAAAAIAAAACAAAAAAAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAIAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAIAAAABAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+OcAAAAAgAAAAEAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+M4AAAAAgAAAAEAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{WnyuwPgC7Wh5mQFPVEcglHLFZrJq7d/lSN4DTDj7QvOMsTKL2cYyzx6UAfiUOGEMnEInweHFsb5dP9yZCy7VAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('1d6308dc6e9617bee39a69f68176cf6f3abcf4d3617db3c766647bd198a5e442', 4, 3, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934595, 100, 1, '2019-06-03 16:37:45.07334', '2019-06-03 16:37:45.07334', 17179881472, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAADAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABQlRDAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt73//////////AAAAAAAAAAGu5L5MAAAAQLEyHlSQ5gb4aQ7evOl4mZ6lSTIF7kShyso/iyP0uz3ipHocd38/dLiu7lVvMGXwo6ymJ7mixdDuNLIWiI9TbQI=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvi1AAAAAIAAAACAAAAAQAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvi1AAAAAIAAAADAAAAAQAAAAAAAAACAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAMAAAABAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAMAAAACAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFCVEMAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+M4AAAAAgAAAAEAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+LUAAAAAgAAAAEAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{sTIeVJDmBvhpDt686XiZnqVJMgXuRKHKyj+LI/S7PeKkehx3fz90uK7uVW8wZfCjrKYnuaLF0O40shaIj1NtAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('647eaba7f3bc5726dc1041553fe4741542ed0a2af2d098d93b0bac5b6f3c624c', 4, 4, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934594, 100, 1, '2019-06-03 16:37:45.073472', '2019-06-03 16:37:45.073472', 17179885568, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABU0NPVAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TH//////////AAAAAAAAAAEnciLVAAAAQHTUKeZaZX/yonQdzrGY0klZqwhUZd7ontUbjpQmLk+XRY8uYos+AI2Z3qqU3QF27EV4VRsVcUUvvn57fqFdzgQ=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvjOAAAAAIAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvjOAAAAAIAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFTQ09UAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+OcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+M4AAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{dNQp5lplf/KidB3OsZjSSVmrCFRl3uie1RuOlCYuT5dFjy5iiz4AjZneqpTdAXbsRXhVGxVxRS++fnt+oV3OBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('bd60680a1378ffec739e1ffa2db4cd51f58babfb714e04a52bd2b65bf8a31b4f', 3, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934593, 100, 1, '2019-06-03 16:37:45.086095', '2019-06-03 16:37:45.086096', 12884905984, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABruS+TAAAAEBkz5uRgU5FxqOu8Yak7Bbdc0BtgvEJ0FjurZz/LgGwT2EX91Y81YrdSVu2NPR0lbhSAotGQlvSPYEy5vN67p4C', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvjnAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvjnAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+OcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+OcAAAAAgAAAAEAAAAAAAAAAAAAAAIAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+OcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ZM+bkYFORcajrvGGpOwW3XNAbYLxCdBY7q2c/y4BsE9hF/dWPNWK3UlbtjT0dJW4UgKLRkJb0j2BMubzeu6eAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c780569c402c298b7b5f3f1a6a20ac1219a06df39a78fb3ac6d93ca53ad4e5ed', 3, 2, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934593, 100, 1, '2019-06-03 16:37:45.086312', '2019-06-03 16:37:45.086312', 12884910080, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAACHRlc3QuY29tAAAAAAAAAAAAAAAB+ZAt7wAAAEBHwkZcyIWmaPvEtDlR8Ed4dD1Mep2juLtHF3n5RG0jurJhKq/3MB1zR6bDHr+wow35ijK92ihjHWqTxjzKDhkO', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvjOAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvjOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+OcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{R8JGXMiFpmj7xLQ5UfBHeHQ9THqdo7i7Rxd5+URtI7qyYSqv9zAdc0emwx6/sKMN+YoyvdooYx1qk8Y8yg4ZDg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2d317dcef8626e639bcaab4a4b1ca1e8e6647eb46d65ca8d98137cd98eb10ae7', 3, 3, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934594, 100, 1, '2019-06-03 16:37:45.086491', '2019-06-03 16:37:45.086491', 12884914176, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB+ZAt7wAAAEB8q5Of+GA0eadw+hTrTCIAoedKyFge/Kv+RUNsq7sv7pSoLAQFWqwFIvxCGBul0XhSxOomG/gWgmIiwj6a1goM', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvjOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvjOAAAAAIAAAACAAAAAAAAAAAAAAAAAAAACHRlc3QuY29tAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAIAAAAAAAAAAAAAAAEAAAAIdGVzdC5jb20BAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+OcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+M4AAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{fKuTn/hgNHmncPoU60wiAKHnSshYHvyr/kVDbKu7L+6UqCwEBVqsBSL8QhgbpdF4UsTqJhv4FoJiIsI+mtYKDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('db398eb4ae89756325643cad21c94e13bfc074b323ee83e141bf701a5d904f1b', 2, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-03 16:37:45.10044', '2019-06-03 16:37:45.10044', 8589938688, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAACVAvkAAAAAAAAAAABVvwF9wAAAEAYjQcPT2G5hqnBmgGGeg9J8l4c1EnUlxklElH9sqZr0971F6OLWfe/m4kpFtI+sI0i1qLit5A0JyWnbhYLW5oD', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBrUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{GI0HD09huYapwZoBhnoPSfJeHNRJ1JcZJRJR/bKma9Pe9Reji1n3v5uJKRbSPrCNItai4reQNCclp24WC1uaAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('f97caffab8c16023a37884165cb0b3ff1aa2daf4000fef49d21efc847ddbfbea', 2, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 2, 100, 1, '2019-06-03 16:37:45.10063', '2019-06-03 16:37:45.10063', 8589942784, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBmKpSgvrwKO20XCOfYfXsGEEUtwYaaEfqSu6ymJmlDma+IX6I7IggbUZMocQdZ94IMAfKdQANqXbIO7ysweeMC', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBrUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDbUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ZiqUoL68CjttFwjn2H17BhBFLcGGmhH6kruspiZpQ5mviF+iOyIIG1GTKHEHWfeCDAHynUADal2yDu8rMHnjAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('725756b1fbdf83b08127f385efedf0909cc820b6cce71f1c0897d15427cb5add', 2, 3, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 3, 100, 1, '2019-06-03 16:37:45.100739', '2019-06-03 16:37:45.10074', 8589946880, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAACVAvkAAAAAAAAAAABVvwF9wAAAEBj4gBQ/BAbgqf7qOotatgZUHjDlsOtDNdp7alZR5/Fk9fGj+lxEygAZWzY7/LY1Z3SF6c0qs172LhAkkvV8p0M', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtq7/TDbUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtqyrQFLUAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAJUC+QAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Y+IAUPwQG4Kn+6jqLWrYGVB4w5bDrQzXae2pWUefxZPXxo/pcRMoAGVs2O/y2NWd0henNKrNe9i4QJJL1fKdDA==}', 'none', NULL, NULL, true, 100);


--
-- Name: asset_stats asset_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY asset_stats
    ADD CONSTRAINT asset_stats_pkey PRIMARY KEY (id);


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
-- Name: by_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_hash ON history_transactions USING btree (transaction_hash);


--
-- Name: by_ledger; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_ledger ON history_transactions USING btree (ledger_sequence, application_order);


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
-- Name: trade_effects_by_order_book; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trade_effects_by_order_book ON history_effects USING btree (((details ->> 'sold_asset_type'::text)), ((details ->> 'sold_asset_code'::text)), ((details ->> 'sold_asset_issuer'::text)), ((details ->> 'bought_asset_type'::text)), ((details ->> 'bought_asset_code'::text)), ((details ->> 'bought_asset_issuer'::text))) WHERE (type = 33);


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


--
-- PostgreSQL database dump complete
--

