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

INSERT INTO asset_stats VALUES (1, '5000000000', 2, 0, '');
INSERT INTO asset_stats VALUES (2, '5000000000', 2, 0, '');


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
INSERT INTO gorp_migrations VALUES ('18_account_for_signers.sql', '2019-10-31 14:19:49.123835+01');


--
-- Data for Name: history_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_accounts VALUES (1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_accounts VALUES (2, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_accounts VALUES (3, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_accounts VALUES (4, 'GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG');
INSERT INTO history_accounts VALUES (5, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_accounts_id_seq', 5, true);


--
-- Data for Name: history_assets; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_assets VALUES (1, 'credit_alphanum4', 'USD', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_assets VALUES (2, 'credit_alphanum4', 'EUR', 'GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG');


--
-- Name: history_assets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_assets_id_seq', 2, true);


--
-- Data for Name: history_effects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_effects VALUES (2, 38654709761, 1, 33, '{"seller": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "offer_id": 2, "sold_amount": "49.9999999", "bought_amount": "45.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "bought_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (1, 38654709761, 2, 33, '{"seller": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "offer_id": 2, "sold_amount": "45.0000000", "bought_amount": "49.9999999", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "bought_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 38654709761, 3, 33, '{"seller": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "offer_id": 3, "sold_amount": "100.0000000", "bought_amount": "80.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "bought_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (1, 38654709761, 4, 33, '{"seller": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "offer_id": 3, "sold_amount": "80.0000000", "bought_amount": "100.0000000", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "bought_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 34359742465, 1, 33, '{"seller": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "offer_id": 1, "sold_amount": "100.0000000", "bought_amount": "100.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "bought_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (1, 34359742465, 2, 33, '{"seller": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "offer_id": 1, "sold_amount": "100.0000000", "bought_amount": "100.0000000", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "bought_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 34359742465, 3, 33, '{"seller": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "offer_id": 2, "sold_amount": "50.0000000", "bought_amount": "45.0000000", "sold_asset_code": "USD", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "EUR", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "bought_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (1, 34359742465, 4, 33, '{"seller": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "offer_id": 2, "sold_amount": "45.0000000", "bought_amount": "50.0000000", "sold_asset_code": "EUR", "sold_asset_type": "credit_alphanum4", "bought_asset_code": "USD", "bought_asset_type": "credit_alphanum4", "sold_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "bought_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 17179873281, 1, 2, '{"amount": "500.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (3, 17179873281, 2, 3, '{"amount": "500.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 17179877377, 1, 2, '{"amount": "500.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (4, 17179877377, 2, 3, '{"amount": "500.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (2, 12884905985, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (1, 12884910081, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (2, 12884914177, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (1, 12884918273, 1, 20, '{"limit": "922337203685.4775807", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');
INSERT INTO history_effects VALUES (2, 8589938689, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (5, 8589938689, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 8589938689, 3, 10, '{"weight": 1, "public_key": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (1, 8589942785, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (5, 8589942785, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 8589942785, 3, 10, '{"weight": 1, "public_key": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"}');
INSERT INTO history_effects VALUES (3, 8589946881, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (5, 8589946881, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (3, 8589946881, 3, 10, '{"weight": 1, "public_key": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}');
INSERT INTO history_effects VALUES (4, 8589950977, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (5, 8589950977, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (4, 8589950977, 3, 10, '{"weight": 1, "public_key": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_ledgers VALUES (9, '486cb2fea3a46891d301a6dd1d06f09d3cfe6144df09581888a0f51299eca4c9', 'beaea632702de50e50e1828cc4f6f2fa5177750a52ffb7f980e0eebc791e69da', 1, 1, '2019-06-03 16:35:18', '2019-06-03 16:35:14.561384', '2019-06-03 16:35:14.561384', 38654705664, 16, 1000000000000000000, 1500, 100, 100000000, 1000000, 11, 'AAAAC76upjJwLeUOUOGCjMT28vpRd3UKUv+3+YDg7rx5Hmnak9zC8GPJ6ZzGP5+PUo/5nHalJyPDK1xoG7op7rz8DI4AAAAAXPVMRgAAAAAAAAAA5u7HOB8cX0WnqixnvzOlY3D3Y3oP0dYzzW9yWnVbm0aw5yPUnzjewy4+r4wev9UkS7PVQ/EqDoIUgz/Be3VFmgAAAAkN4Lazp2QAAAAAAAAAAAXcAAAAAAAAAAAAAAAEAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (8, 'beaea632702de50e50e1828cc4f6f2fa5177750a52ffb7f980e0eebc791e69da', '0d3f343b2f9f0caebd2e9b5ad5233d384ddef567e01fdcd8e79b3a6f010e14cd', 1, 1, '2019-06-03 16:35:17', '2019-06-03 16:35:14.606009', '2019-06-03 16:35:14.606009', 34359738368, 16, 1000000000000000000, 1400, 100, 100000000, 1000000, 11, 'AAAACw0/NDsvnwyuvS6bWtUjPThN3vVn4B/c2OebOm8BDhTNF+e1/1KxUzO74nHG3kqkGpSmFdahZJnhcECgqI4Wul8AAAAAXPVMRQAAAAAAAAAAAOMGyexOuevk4QIYCK/e2TLQXecE/TBuu1dvGQilTsNvK1rbv4HEl5bCA6pVvDMw5mE6JcfcRSMCkb3u4w2XMQAAAAgN4Lazp2QAAAAAAAAAAAV4AAAAAAAAAAAAAAADAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (7, '0d3f343b2f9f0caebd2e9b5ad5233d384ddef567e01fdcd8e79b3a6f010e14cd', 'c3a9cd326f60a94427dcadeb699d3f6ecd39353cc932acc7d27451273326bfc8', 1, 1, '2019-06-03 16:35:16', '2019-06-03 16:35:14.625706', '2019-06-03 16:35:14.625706', 30064771072, 16, 1000000000000000000, 1300, 100, 100000000, 1000000, 11, 'AAAAC8OpzTJvYKlEJ9yt62mdP27NOTU8yTKsx9J0USczJr/IuByoQOqmTAMUTx6gy0jhoGG1FiLn3uIQDEu52Ts13f4AAAAAXPVMRAAAAAAAAAAA8011I8cweiVULRk7g+I3eWM/k+U+svLLvJL8X9ENB1pwLIQgQIFbTzzBFKq7Ll+k2W2vdjGQsQwexfujYVC00AAAAAcN4Lazp2QAAAAAAAAAAAUUAAAAAAAAAAAAAAADAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (6, 'c3a9cd326f60a94427dcadeb699d3f6ecd39353cc932acc7d27451273326bfc8', 'bcea057e1cd211d2cfb7e3fb995fd45cd6a44221f38acc0522b45b8c81d96524', 1, 1, '2019-06-03 16:35:15', '2019-06-03 16:35:14.638052', '2019-06-03 16:35:14.638052', 25769803776, 16, 1000000000000000000, 1200, 100, 100000000, 1000000, 11, 'AAAAC7zqBX4c0hHSz7fj+5lf1FzWpEIh84rMBSK0W4yB2WUkjmMqt+2EvvBEZaLG6tbmWrlOaBGNpqS9R+fV7D1C1bEAAAAAXPVMQwAAAAAAAAAAuQrPQr2qlNVpCIqhQaeQMsD9JWWNaxISOIwMsxvilIUACd+JYsOyVEMrgEvOK6iALcFeWbKnOiuPJcluEtVi9gAAAAYN4Lazp2QAAAAAAAAAAASwAAAAAAAAAAAAAAACAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (5, 'bcea057e1cd211d2cfb7e3fb995fd45cd6a44221f38acc0522b45b8c81d96524', '7cbd2c56a6b343fa00f87e24264de55ff07e32f3b149084d99e565ef0cd4f02a', 1, 1, '2019-06-03 16:35:14', '2019-06-03 16:35:14.647891', '2019-06-03 16:35:14.647891', 21474836480, 16, 1000000000000000000, 1100, 100, 100000000, 1000000, 11, 'AAAAC3y9LFams0P6APh+JCZN5V/wfjLzsUkITZnlZe8M1PAqI/JFnzikQwlAWB2pC99Ya/VbQ2JSlXrpxN9Ed9L/rDMAAAAAXPVMQgAAAAAAAAAA5tFak/fh7MEOithjlydErBU0/M0ilaefu1Ub0Cw1cbQzN/qMhlghXJshjEg7qm8V0NxHBj4ZOikOX0S6QnMzRAAAAAUN4Lazp2QAAAAAAAAAAARMAAAAAAAAAAAAAAABAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (4, '7cbd2c56a6b343fa00f87e24264de55ff07e32f3b149084d99e565ef0cd4f02a', 'aa733b6e19d9d1df08ab8ac1bb70ce941adb9c319716dbbed0b83a4f9717a427', 2, 2, '2019-06-03 16:35:13', '2019-06-03 16:35:14.660033', '2019-06-03 16:35:14.660033', 17179869184, 16, 1000000000000000000, 1000, 100, 100000000, 1000000, 11, 'AAAAC6pzO24Z2dHfCKuKwbtwzpQa25wxlxbbvtC4Ok+XF6QnjkhXuTkDWtl4LNb/YNQTjpJEIGD+9shjsz3Jzsqv4/YAAAAAXPVMQQAAAAAAAAAAN5B5W9KsxoOXOpscNNdDinauNxC/Pw2kk7DDUo0SXOXZntDAxTfXn66m9IhE0XmliaKvkeIpNJfz9KPboAIZOwAAAAQN4Lazp2QAAAAAAAAAAAPoAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO history_ledgers VALUES (3, 'aa733b6e19d9d1df08ab8ac1bb70ce941adb9c319716dbbed0b83a4f9717a427', '5b6634cc9e9199fd3fbbaf0655c54ca5e7b812d562e63eca0fb2b9a89b180c4f', 4, 4, '2019-06-03 16:35:12', '2019-06-03 16:35:14.672241', '2019-06-03 16:35:14.672241', 12884901888, 16, 1000000000000000000, 800, 100, 100000000, 1000000, 11, 'AAAAC1tmNMyekZn9P7uvBlXFTKXnuBLVYuY+yg+yuaibGAxPE5AOt+233THPtTNbUnjSwPvmo6g3zyFcFPoc8VLPGGAAAAAAXPVMQAAAAAAAAAAAShzqI9aLY+7Kp5Zw36Ye703iZxnHA4Dw1m5KF4RHKmONRLiXB5aHfKAGWr3HxolQQxgyxXCdjPya6ZVzaVSFTwAAAAMN4Lazp2QAAAAAAAAAAAMgAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0);
INSERT INTO history_ledgers VALUES (2, '5b6634cc9e9199fd3fbbaf0655c54ca5e7b812d562e63eca0fb2b9a89b180c4f', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 4, 4, '2019-06-03 16:35:11', '2019-06-03 16:35:14.6861', '2019-06-03 16:35:14.6861', 8589934592, 16, 1000000000000000000, 400, 100, 100000000, 1000000, 11, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76Zh82QgsuASmS0xgFp1fjiT2wNgvVWRClSW5uMP344SMoAAAAAXPVMPwAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAAh7rHH3ehqvUsWEPfAlvVspAA2/ZR0JjPAFIdHbQ2sSiuJxnU7ly+zexZusXx0NLXyoK5vVQIvDzDn+FmYFOehgAAAAIN4Lazp2QAAAAAAAAAAAGQAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0);
INSERT INTO history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-03 16:35:14.698632', '2019-06-03 16:35:14.698632', 4294967296, 16, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operation_participants VALUES (1, 38654709761, 2);
INSERT INTO history_operation_participants VALUES (2, 34359742465, 2);
INSERT INTO history_operation_participants VALUES (3, 30064775169, 1);
INSERT INTO history_operation_participants VALUES (4, 25769807873, 1);
INSERT INTO history_operation_participants VALUES (5, 21474840577, 1);
INSERT INTO history_operation_participants VALUES (6, 17179873281, 3);
INSERT INTO history_operation_participants VALUES (7, 17179873281, 2);
INSERT INTO history_operation_participants VALUES (8, 17179877377, 4);
INSERT INTO history_operation_participants VALUES (9, 17179877377, 1);
INSERT INTO history_operation_participants VALUES (10, 12884905985, 2);
INSERT INTO history_operation_participants VALUES (11, 12884910081, 1);
INSERT INTO history_operation_participants VALUES (12, 12884914177, 2);
INSERT INTO history_operation_participants VALUES (13, 12884918273, 1);
INSERT INTO history_operation_participants VALUES (14, 8589938689, 5);
INSERT INTO history_operation_participants VALUES (15, 8589938689, 2);
INSERT INTO history_operation_participants VALUES (16, 8589942785, 5);
INSERT INTO history_operation_participants VALUES (17, 8589942785, 1);
INSERT INTO history_operation_participants VALUES (18, 8589946881, 5);
INSERT INTO history_operation_participants VALUES (19, 8589946881, 3);
INSERT INTO history_operation_participants VALUES (20, 8589950977, 5);
INSERT INTO history_operation_participants VALUES (21, 8589950977, 4);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_operation_participants_id_seq', 21, true);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operations VALUES (38654709761, 38654709760, 1, 3, '{"price": "0.8000000", "amount": "151.0000000", "price_r": {"d": 5, "n": 4}, "offer_id": 0, "buying_asset_code": "EUR", "buying_asset_type": "credit_alphanum4", "selling_asset_code": "USD", "selling_asset_type": "credit_alphanum4", "buying_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "selling_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (34359742465, 34359742464, 1, 3, '{"price": "0.9000000", "amount": "150.0000000", "price_r": {"d": 10, "n": 9}, "offer_id": 0, "buying_asset_code": "EUR", "buying_asset_type": "credit_alphanum4", "selling_asset_code": "USD", "selling_asset_type": "credit_alphanum4", "buying_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "selling_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (30064775169, 30064775168, 1, 3, '{"price": "1.2500000", "amount": "80.0000000", "price_r": {"d": 4, "n": 5}, "offer_id": 0, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_code": "EUR", "selling_asset_type": "credit_alphanum4", "buying_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "selling_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (25769807873, 25769807872, 1, 3, '{"price": "1.1111111", "amount": "90.0000000", "price_r": {"d": 1000000000, "n": 1111111111}, "offer_id": 0, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_code": "EUR", "selling_asset_type": "credit_alphanum4", "buying_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "selling_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (21474840577, 21474840576, 1, 3, '{"price": "1.0000000", "amount": "100.0000000", "price_r": {"d": 1, "n": 1}, "offer_id": 0, "buying_asset_code": "USD", "buying_asset_type": "credit_alphanum4", "selling_asset_code": "EUR", "selling_asset_type": "credit_alphanum4", "buying_asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "selling_asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (17179873281, 17179873280, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "amount": "500.0000000", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4');
INSERT INTO history_operations VALUES (17179877377, 17179877376, 1, 1, '{"to": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "from": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "amount": "500.0000000", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG');
INSERT INTO history_operations VALUES (12884905985, 12884905984, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (12884910081, 12884910080, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "USD", "asset_type": "credit_alphanum4", "asset_issuer": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (12884914177, 12884914176, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "trustor": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (12884918273, 12884918272, 1, 6, '{"limit": "922337203685.4775807", "trustee": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "trustor": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "asset_code": "EUR", "asset_type": "credit_alphanum4", "asset_issuer": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG"}', 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');
INSERT INTO history_operations VALUES (8589938689, 8589938688, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589942785, 8589942784, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589946881, 8589946880, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589950977, 8589950976, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_trades VALUES (38654709761, 0, '2019-06-03 16:35:18', 2, 2, 1, 499999999, 1, 2, 450000000, false, 1000000000, 1111111111, 4, 2);
INSERT INTO history_trades VALUES (38654709761, 1, '2019-06-03 16:35:18', 3, 2, 1, 1000000000, 1, 2, 800000000, false, 4, 5, 4, 3);
INSERT INTO history_trades VALUES (34359742465, 0, '2019-06-03 16:35:17', 1, 2, 1, 1000000000, 1, 2, 1000000000, false, 1, 1, 4611686052787130369, 1);
INSERT INTO history_trades VALUES (34359742465, 1, '2019-06-03 16:35:17', 2, 2, 1, 500000000, 1, 2, 450000000, false, 1000000000, 1111111111, 4611686052787130369, 2);


--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transaction_participants VALUES (1, 38654709760, 2);
INSERT INTO history_transaction_participants VALUES (2, 34359742464, 2);
INSERT INTO history_transaction_participants VALUES (3, 30064775168, 1);
INSERT INTO history_transaction_participants VALUES (4, 25769807872, 1);
INSERT INTO history_transaction_participants VALUES (5, 21474840576, 1);
INSERT INTO history_transaction_participants VALUES (6, 17179873280, 3);
INSERT INTO history_transaction_participants VALUES (7, 17179873280, 2);
INSERT INTO history_transaction_participants VALUES (8, 17179877376, 4);
INSERT INTO history_transaction_participants VALUES (9, 17179877376, 1);
INSERT INTO history_transaction_participants VALUES (10, 12884905984, 2);
INSERT INTO history_transaction_participants VALUES (11, 12884910080, 1);
INSERT INTO history_transaction_participants VALUES (12, 12884914176, 2);
INSERT INTO history_transaction_participants VALUES (13, 12884918272, 1);
INSERT INTO history_transaction_participants VALUES (14, 8589938688, 5);
INSERT INTO history_transaction_participants VALUES (15, 8589938688, 2);
INSERT INTO history_transaction_participants VALUES (16, 8589942784, 5);
INSERT INTO history_transaction_participants VALUES (17, 8589942784, 1);
INSERT INTO history_transaction_participants VALUES (18, 8589946880, 5);
INSERT INTO history_transaction_participants VALUES (19, 8589946880, 3);
INSERT INTO history_transaction_participants VALUES (20, 8589950976, 5);
INSERT INTO history_transaction_participants VALUES (21, 8589950976, 4);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_transaction_participants_id_seq', 21, true);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transactions VALUES ('54e27e3d7c84fdb79d6709922c52aaa752fd5080dcd433af901cd17860a957b2', 9, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934596, 100, 1, '2019-06-03 16:35:14.561744', '2019-06-03 16:35:14.561745', 38654709760, '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', '{1q6EaaT1LdjuTrlh1ErZEAeMfBiK2+R64z7ec8I/eEaBSTP1YZrg9jzgJ4yCaSXXmJyT6VgWtd9vgugOqG5NCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('1fa96abeeb0e2b71fa947d8f7003ca5d3e95734652c38b328ae62a8c1b1b3d73', 8, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934595, 100, 1, '2019-06-03 16:35:14.606344', '2019-06-03 16:35:14.606345', 34359742464, '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', '{jjPwd3EIxz1iPsgxjjf6CdEvmZ5QxMY7kw2oez7V9qhYKcL7QgLetFeNhaiKA+cQEOZM3FBh80D0tBA0flV3DA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4a0ac16e948710199659b0e8131be06df5944e57806c0db722859e4e8aad1e44', 7, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934597, 100, 1, '2019-06-03 16:35:14.625995', '2019-06-03 16:35:14.625996', 30064775168, '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', '{YK8tE7Ae0YbqbY8nQgjPcVEcZTH1bKCo1cc1A83RimT+hM2W0kC4+fnYVn47t8ia9wd98wtUwIUM0NOAWeshAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2d0d264d9e2e3364c29faf6a6ec815a8a8e6a46e23096da1411f292309d78712', 6, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934596, 100, 1, '2019-06-03 16:35:14.638427', '2019-06-03 16:35:14.638427', 25769807872, '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', '{eh97duB2MF3rrgCuqsouPb+JrhstrliUb8wvu4VVEpZNunSWlsoJrRsVHkAScIcP8XImAqHdOI06krccVGRZDw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c7ab3843f0a0c4fe79dece0ff1b8391f7a9d34c47cbd7e35034f7b28f60dfc00', 3, 3, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934594, 100, 1, '2019-06-03 16:35:14.672808', '2019-06-03 16:35:14.672808', 12884914176, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABRVVSAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSX//////////AAAAAAAAAAGu5L5MAAAAQO/nblo8KAkSOf8cQOOiADXygx+I0ZdWoM4Vg4EKPAAJXFntctjCIyQ4csVUywaW32J/keQWYby52BjiNhT/6Qo=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rJOAAAAAIAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rJOAAAAAIAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFFVVIAAAAAAKD8GPxXf3CMjVKv3qhzs6au2HaMHH0s5B32RVUQithJAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{7+duWjwoCRI5/xxA46IANfKDH4jRl1agzhWDgQo8AAlcWe1y2MIjJDhyxVTLBpbfYn+R5BZhvLnYGOI2FP/pCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('356eec1f0432bbbba09ddcc1937e28253a4dc8e98ff77d5e9dee140be5093a70', 5, 1, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934595, 100, 1, '2019-06-03 16:35:14.648216', '2019-06-03 16:35:14.648216', 21474840576, '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', '{ToCTIP5fc+hzoWPqVqVi+YdR1s3fy66foJfkxz9wmcVBqq1kE6gDAOe1OXUbP9/ut5GldwlhJ1D6kRxSf/Y8AQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2837a1b3def2c2ddfdcde5b44bf08d7a11a9328d870df17fa2bb66d4c83260c7', 4, 1, 'GC23QF2HUE52AMXUFUH3AYJAXXGXXV2VHXYYR6EYXETPKDXZSAW67XO4', 8589934593, 100, 1, '2019-06-03 16:35:14.660189', '2019-06-03 16:35:14.660189', 17179873280, 'AAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAEqBfIAAAAAAAAAAAH5kC3vAAAAQGoQPiD0HEBi0U8cHN6nlZ3okEfdmt7mqkQHIt2tuRLaZZ1iMQwU43M8v+ntJQsA4c2eBXt9GYp/29FLjnba1Qw=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAAAO5rJnAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAAAO5rJnAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAADAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAEAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAASoF8gB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ahA+IPQcQGLRTxwc3qeVneiQR92a3uaqRAci3a25EtplnWIxDBTjczy/6e0lCwDhzZ4Fe30Zin/b0UuOdtrVDA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('e4febb6aab5f83df4c19b4f1f6b8687f11886169e996e16baa75aac2d8a09c68', 4, 2, 'GCQPYGH4K57XBDENKKX55KDTWOTK5WDWRQOH2LHEDX3EKVIQRLMESGBG', 8589934593, 100, 1, '2019-06-03 16:35:14.660364', '2019-06-03 16:35:14.660365', 17179877376, 'AAAAAKD8GPxXf3CMjVKv3qhzs6au2HaMHH0s5B32RVUQithJAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAABRVVSAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSQAAAAEqBfIAAAAAAAAAAAEQithJAAAAQEghWBLDjmNLzcPF6o8dqUHMsI0WhttWE/ABSKaHNc+0FqsF+ui5+eky4ERyu99YR6BEHF4NlgyLnSq3zcDr9Qo=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAAoPwY/Fd/cIyNUq/eqHOzpq7YdowcfSzkHfZFVRCK2EkAAAAAO5rJnAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAoPwY/Fd/cIyNUq/eqHOzpq7YdowcfSzkHfZFVRCK2EkAAAAAO5rJnAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAgAAAAMAAAADAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFFVVIAAAAAAKD8GPxXf3CMjVKv3qhzs6au2HaMHH0s5B32RVUQithJAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAAAAAAEAAAAEAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFFVVIAAAAAAKD8GPxXf3CMjVKv3qhzs6au2HaMHH0s5B32RVUQithJAAAAASoF8gB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSQAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSQAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{SCFYEsOOY0vNw8Xqjx2pQcywjRaG21YT8AFIpoc1z7QWqwX66Ln56TLgRHK731hHoEQcXg2WDIudKrfNwOv1Cg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('bd486dbdd02d460817671c4a5a7e9d6e865ca29cb41e62d7aaf70a2fee5b36de', 3, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934593, 100, 1, '2019-06-03 16:35:14.672439', '2019-06-03 16:35:14.67244', 12884905984, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt73//////////AAAAAAAAAAGu5L5MAAAAQB9kmKW2q3v7Qfy8PMekEb1TTI5ixqkI0BogXrOt7gO162Qbkh2dSTUfeDovc0PAafhDXxthVAlsLujlBmyjBAY=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rJOAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rJOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAQAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{H2SYpbare/tB/Lw8x6QRvVNMjmLGqQjQGiBes63uA7XrZBuSHZ1JNR94Oi9zQ8Bp+ENfG2FUCWwu6OUGbKMEBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('00ab9cfce2b4c4141d8bb6768dd094bdbb1c7406710dbb3ba0ef98870f63a344', 3, 2, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934593, 100, 1, '2019-06-03 16:35:14.67266', '2019-06-03 16:35:14.67266', 12884910080, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABVVNEAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt73//////////AAAAAAAAAAEnciLVAAAAQLVbII+1LeizxgncDI46KHyBt05+H92n1+R328J9zNl2fgJW2nfn3FIoLVs2qV1+CUpr121a2B7AM6HKr4nBLAI=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rJOAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rJOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msk4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msk4AAAAAgAAAAEAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFVU0QAAAAAALW4F0ehO6Ay9C0PsGEgvc1711U98Yj4mLkm9Q75kC3vAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{tVsgj7Ut6LPGCdwMjjoofIG3Tn4f3afX5Hfbwn3M2XZ+Albad+fcUigtWzapXX4JSmvXbVrYHsAzocqvicEsAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('5ef9c06bb625d2da2281118bdf14d808768353ee2fca7457c8e506fbeb91fc55', 3, 4, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2', 8589934594, 100, 1, '2019-06-03 16:35:14.672941', '2019-06-03 16:35:14.672942', 12884918272, 'AAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAYAAAABRVVSAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSX//////////AAAAAAAAAAEnciLVAAAAQOG2GKO9i60cM1QK2UN1gXrHEYjeGLRXFT5snCqO5FnPET5cVs30N7ITPZ6HH6QcZ2IdC1c66wLge4GyR8vpww0=', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAGAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rJOAAAAAIAAAABAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rJOAAAAAIAAAACAAAAAQAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msk4AAAAAgAAAAIAAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msk4AAAAAgAAAAIAAAACAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAQAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAFFVVIAAAAAAKD8GPxXf3CMjVKv3qhzs6au2HaMHH0s5B32RVUQithJAAAAAAAAAAB//////////wAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msk4AAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{4bYYo72LrRwzVArZQ3WBescRiN4YtFcVPmycKo7kWc8RPlxWzfQ3shM9nocfpBxnYh0LVzrrAuB7gbJHy+nDDQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d', 2, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-03 16:35:14.686222', '2019-06-03 16:35:14.686223', 8589938688, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rKAAAAAAAAAAABVvwF9wAAAECDzqvkQBQoNAJifPRXDoLhvtycT3lFPCQ51gkdsFHaBNWw05S/VhW0Xgkr0CBPE4NaFV2Kmcs3ZwLmib4TRrML', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/5wAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrNryTRwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{g86r5EAUKDQCYnz0Vw6C4b7cnE95RTwkOdYJHbBR2gTVsNOUv1YVtF4JK9AgTxODWhVdipnLN2cC5om+E0azCw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('164a5064eba64f2cdbadb856bf3448485fc626247ada3ed39cddf0f6902133b6', 2, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 2, 100, 1, '2019-06-03 16:35:14.686364', '2019-06-03 16:35:14.686364', 8589942784, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEASEZiZbeFwCsrKBnKIus/05VtJDBrgosuhLQ/U6XUj4twWyhs7UtS4CMexOM6JqcfqJK10WlBkkwn4g8PIfjIG', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrNryTRwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrMwLmpwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAA7YL8A7jlgEPe0dUU7VHcDQx6Q/wlHqc3UD15aJ3Ii1QAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{EhGYmW3hcArKygZyiLrP9OVbSQwa4KLLoS0P1Ol1I+LcFsobO1LUuAjHsTjOianH6iStdFpQZJMJ+IPDyH4yBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2d365c3c49f376570df856bca62503966f0e269a2f51cdb68ce2ee19a7f8245a', 2, 3, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 3, 100, 1, '2019-06-03 16:35:14.686469', '2019-06-03 16:35:14.686469', 8589946880, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAtbgXR6E7oDL0LQ+wYSC9zXvXVT3xiPiYuSb1DvmQLe8AAAAAO5rKAAAAAAAAAAABVvwF9wAAAEB/JBgvIM71gLBIh0TON9b+l+ApZz1CKDQiUFSV0scRguB1anyMwMR6s5SiaCwtDnxsPna12RdUQKlH2aeMAy8H', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrMwLmpwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrL0k6BwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAC1uBdHoTugMvQtD7BhIL3Ne9dVPfGI+Ji5JvUO+ZAt7wAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{fyQYLyDO9YCwSIdEzjfW/pfgKWc9Qig0IlBUldLHEYLgdWp8jMDEerOUomgsLQ58bD52tdkXVECpR9mnjAMvBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('003e91101d19aabb429491953806886d777c260233c6478f1c928a79ec4e2743', 2, 4, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 4, 100, 1, '2019-06-03 16:35:14.686579', '2019-06-03 16:35:14.686579', 8589950976, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAoPwY/Fd/cIyNUq/eqHOzpq7YdowcfSzkHfZFVRCK2EkAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEDIOudzujfo+dSIJXXb06SjLBLLXsFxnVnR1HJejfq2NgFUtLuX2KrVNSZyRBG+WvfdoXwCPcp85hDRbCmjbPYM', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrL0k6BwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrK4+NZwAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAACg/Bj8V39wjI1Sr96oc7Omrth2jBx9LOQd9kVVEIrYSQAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/5wAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{yDrnc7o36PnUiCV129OkoywSy17BcZ1Z0dRyXo36tjYBVLS7l9iq1TUmckQRvlr33aF8Aj3KfOYQ0Wwpo2z2DA==}', 'none', NULL, NULL, true, 100);


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

-- added manually
ALTER TABLE IF EXISTS ONLY public.key_value_store DROP CONSTRAINT IF EXISTS key_value_store_pkey;
DROP TABLE IF EXISTS public.key_value_store;
CREATE TABLE key_value_store (
                                 key character varying(255) NOT NULL,
                                 value character varying(255) NOT NULL
);
INSERT INTO key_value_store VALUES ('exp_ingest_last_ledger', '0');
ALTER TABLE ONLY key_value_store
    ADD CONSTRAINT key_value_store_pkey PRIMARY KEY (key);

CREATE TABLE accounts_signers (
                                  account character varying(64),
                                  signer character varying(64),
                                  weight integer NOT NULL,
    -- we will query by signer so that is why signer is the first item in the composite key
                                  PRIMARY KEY (signer, account)
);

--
-- PostgreSQL database dump complete
--

