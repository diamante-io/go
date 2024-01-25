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
    fee_charged integer,
    new_max_fee integer
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
INSERT INTO history_accounts VALUES (2, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_accounts VALUES (3, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_accounts_id_seq', 3, true);


--
-- Data for Name: history_assets; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Name: history_assets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_assets_id_seq', 1, false);


--
-- Data for Name: history_effects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_effects VALUES (1, 30064775169, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 30064775169, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 30064779265, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 30064779265, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 30064783361, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 30064783361, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 25769807873, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 25769807873, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 25769811969, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 25769811969, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 21474840577, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 21474840577, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 21474844673, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 21474844673, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 21474848769, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 21474848769, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 21474852865, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 21474852865, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 17179873281, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 17179873281, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 12884905985, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 12884905985, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 8589938689, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (3, 8589938689, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 8589938689, 3, 10, '{"weight": 1, "public_key": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO history_effects VALUES (2, 8589942785, 1, 0, '{"starting_balance": "100.0000000"}');
INSERT INTO history_effects VALUES (3, 8589942785, 2, 3, '{"amount": "100.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 8589942785, 3, 10, '{"weight": 1, "public_key": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_ledgers VALUES (7, '16f6fdb5c204980e34ec5334c409b452b5906550d6e2bc5e65b780335e128f25', '12f18c34bb467618af7d20936de6ca4a3d1caa881c36ac9a7b89ce713454ed42', 3, 3, '2019-06-03 16:36:46', '2019-06-03 16:36:44.314788', '2019-06-03 16:36:44.314789', 30064771072, 16, 1000000000000000000, 1300, 100, 100000000, 1000000, 11, 'AAAACxLxjDS7RnYYr30gk23myko9HKqIHDasmnuJznE0VO1CJAEap9zN9ulrCqs7yeKqKBciTLmHONm+AASqGkGhtEUAAAAAXPVMngAAAAAAAAAAL150/k9Xfs0TJQCfIu9DCsOEki1Kn+hL2IFyyLVMdXLSqRMksZg/1s/TN2MT0PTALvZwe0hrL/vO9qy8DXIyhAAAAAcN4Lazp2QAAAAAAAAAAAUUAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO history_ledgers VALUES (6, '12f18c34bb467618af7d20936de6ca4a3d1caa881c36ac9a7b89ce713454ed42', '8c22cb8c9c0b37cf98f7fe67fc3a860765c3e68cc0e3ca734a2245d666a24daf', 2, 2, '2019-06-03 16:36:45', '2019-06-03 16:36:44.343267', '2019-06-03 16:36:44.343267', 25769803776, 16, 1000000000000000000, 1000, 100, 100000000, 1000000, 11, 'AAAAC4wiy4ycCzfPmPf+Z/w6hgdlw+aMwOPKc0oiRdZmok2vimN7pygrRl39wo4rGR3ojemz1j9vGRpGXgvKvFhV0TIAAAAAXPVMnQAAAAAAAAAA+GFZHqUIknS3p5sZb4czc8Kj0abGSMVfNaq+9zzlt2RJy5twyGiNy+/nAx1r7Fd+9Gwx1fWOSMkGfXyWCtxQqgAAAAYN4Lazp2QAAAAAAAAAAAPoAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO history_ledgers VALUES (5, '8c22cb8c9c0b37cf98f7fe67fc3a860765c3e68cc0e3ca734a2245d666a24daf', '161c1bd6b25ae1a13742a87c6cb92c865a36c388fb88cd7998b861aeba6696d9', 4, 4, '2019-06-03 16:36:44', '2019-06-03 16:36:44.357815', '2019-06-03 16:36:44.357815', 21474836480, 16, 1000000000000000000, 800, 100, 100000000, 1000000, 11, 'AAAACxYcG9ayWuGhN0KofGy5LIZaNsOI+4jNeZi4Ya66ZpbZv9yVLXLZ14BAjbXRKWKqMHgP4d/yiTLYGoLhZYpcsGQAAAAAXPVMnAAAAAAAAAAAK7ZFHpUv5xgPNV5PVN8eW9yyrJcBUZDK2CSKHdvPgtg5bkfXmLDAtVbwJmXDeKBhcjvQCVwGdkqsINH6rRuEGwAAAAUN4Lazp2QAAAAAAAAAAAMgAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 4, 0);
INSERT INTO history_ledgers VALUES (4, '161c1bd6b25ae1a13742a87c6cb92c865a36c388fb88cd7998b861aeba6696d9', 'e17d4c547e8c25da958839a24d65fad4245c061b18053c0b3355e480abbaa181', 1, 1, '2019-06-03 16:36:43', '2019-06-03 16:36:44.378676', '2019-06-03 16:36:44.378676', 17179869184, 16, 1000000000000000000, 400, 100, 100000000, 1000000, 11, 'AAAAC+F9TFR+jCXalYg5ok1l+tQkXAYbGAU8CzNV5ICruqGBv7gH+MirXcGlvfrSyn6jsW5bxOF4aGsLXUWY6b0GVDcAAAAAXPVMmwAAAAAAAAAA9Nf75AfzKHcfz2nf3wc7i1MpbTisZhsly7MED4o1H3nnYKwcsxtlShiwrJAjUtVbfOPY9454F0BOVuh7a5ArjQAAAAQN4Lazp2QAAAAAAAAAAAGQAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (3, 'e17d4c547e8c25da958839a24d65fad4245c061b18053c0b3355e480abbaa181', 'b296f06555c13ef529da4922b5f8092f7b39ec8364ce93ac8a00b7ac131ca8a3', 1, 1, '2019-06-03 16:36:42', '2019-06-03 16:36:44.395408', '2019-06-03 16:36:44.395408', 12884901888, 16, 1000000000000000000, 300, 100, 100000000, 1000000, 11, 'AAAAC7KW8GVVwT71KdpJIrX4CS97OeyDZM6TrIoAt6wTHKijYjXpTzYAT3UVd5+80mub7DyftrMo9PFH/aZNNRh0SvMAAAAAXPVMmgAAAAAAAAAAzpWxfUdkwpI+mahMUGVw0hP2fK7oFHG9i01PTevvzr/JhKl/HdOjHaSKDl370ucnP4E3gPIT4qzvwbdPBSjEJwAAAAMN4Lazp2QAAAAAAAAAAAEsAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (2, 'b296f06555c13ef529da4922b5f8092f7b39ec8364ce93ac8a00b7ac131ca8a3', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 2, 2, '2019-06-03 16:36:41', '2019-06-03 16:36:44.405856', '2019-06-03 16:36:44.405856', 8589934592, 16, 1000000000000000000, 200, 100, 100000000, 1000000, 11, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76ZRZSEryEgGUi8Runb12uW8dKmZaWm1LdHQ2+tzinFtaAAAAAAXPVMmQAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAA+9VGePGCiconnKwGWQgSrBOnhFsIJOmYJO/8HociG0TDBUH5mFPGH330CygWgjI9Tns2ILOmNsQ1bQ4sqzooEwAAAAIN4Lazp2QAAAAAAAAAAADIAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-03 16:36:44.422671', '2019-06-03 16:36:44.422672', 4294967296, 16, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operation_participants VALUES (1, 30064775169, 2);
INSERT INTO history_operation_participants VALUES (2, 30064775169, 1);
INSERT INTO history_operation_participants VALUES (3, 30064779265, 1);
INSERT INTO history_operation_participants VALUES (4, 30064779265, 2);
INSERT INTO history_operation_participants VALUES (5, 30064783361, 2);
INSERT INTO history_operation_participants VALUES (6, 30064783361, 1);
INSERT INTO history_operation_participants VALUES (7, 25769807873, 1);
INSERT INTO history_operation_participants VALUES (8, 25769807873, 2);
INSERT INTO history_operation_participants VALUES (9, 25769811969, 2);
INSERT INTO history_operation_participants VALUES (10, 25769811969, 1);
INSERT INTO history_operation_participants VALUES (11, 21474840577, 1);
INSERT INTO history_operation_participants VALUES (12, 21474840577, 2);
INSERT INTO history_operation_participants VALUES (13, 21474844673, 2);
INSERT INTO history_operation_participants VALUES (14, 21474844673, 1);
INSERT INTO history_operation_participants VALUES (15, 21474848769, 2);
INSERT INTO history_operation_participants VALUES (16, 21474848769, 1);
INSERT INTO history_operation_participants VALUES (17, 21474852865, 2);
INSERT INTO history_operation_participants VALUES (18, 21474852865, 1);
INSERT INTO history_operation_participants VALUES (19, 17179873281, 1);
INSERT INTO history_operation_participants VALUES (20, 17179873281, 2);
INSERT INTO history_operation_participants VALUES (21, 12884905985, 2);
INSERT INTO history_operation_participants VALUES (22, 12884905985, 1);
INSERT INTO history_operation_participants VALUES (23, 8589938689, 3);
INSERT INTO history_operation_participants VALUES (24, 8589938689, 1);
INSERT INTO history_operation_participants VALUES (25, 8589942785, 3);
INSERT INTO history_operation_participants VALUES (26, 8589942785, 2);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_operation_participants_id_seq', 26, true);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operations VALUES (30064775169, 30064775168, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (30064779265, 30064779264, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (30064783361, 30064783360, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (25769807873, 25769807872, 1, 1, '{"to": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "10.0000000", "asset_type": "native"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (25769811969, 25769811968, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (21474840577, 21474840576, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (21474844673, 21474844672, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (21474848769, 21474848768, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (21474852865, 21474852864, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (17179873281, 17179873280, 1, 1, '{"to": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "10.0000000", "asset_type": "native"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO history_operations VALUES (12884905985, 12884905984, 1, 1, '{"to": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "from": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "amount": "10.0000000", "asset_type": "native"}', 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO history_operations VALUES (8589938689, 8589938688, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_operations VALUES (8589942785, 8589942784, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "starting_balance": "100.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transaction_participants VALUES (1, 30064775168, 1);
INSERT INTO history_transaction_participants VALUES (2, 30064775168, 2);
INSERT INTO history_transaction_participants VALUES (3, 30064779264, 2);
INSERT INTO history_transaction_participants VALUES (4, 30064779264, 1);
INSERT INTO history_transaction_participants VALUES (5, 30064783360, 2);
INSERT INTO history_transaction_participants VALUES (6, 30064783360, 1);
INSERT INTO history_transaction_participants VALUES (7, 25769807872, 2);
INSERT INTO history_transaction_participants VALUES (8, 25769807872, 1);
INSERT INTO history_transaction_participants VALUES (9, 25769811968, 1);
INSERT INTO history_transaction_participants VALUES (10, 25769811968, 2);
INSERT INTO history_transaction_participants VALUES (11, 21474840576, 2);
INSERT INTO history_transaction_participants VALUES (12, 21474840576, 1);
INSERT INTO history_transaction_participants VALUES (13, 21474844672, 2);
INSERT INTO history_transaction_participants VALUES (14, 21474844672, 1);
INSERT INTO history_transaction_participants VALUES (15, 21474848768, 2);
INSERT INTO history_transaction_participants VALUES (16, 21474848768, 1);
INSERT INTO history_transaction_participants VALUES (17, 21474852864, 2);
INSERT INTO history_transaction_participants VALUES (18, 21474852864, 1);
INSERT INTO history_transaction_participants VALUES (19, 17179873280, 1);
INSERT INTO history_transaction_participants VALUES (20, 17179873280, 2);
INSERT INTO history_transaction_participants VALUES (21, 12884905984, 2);
INSERT INTO history_transaction_participants VALUES (22, 12884905984, 1);
INSERT INTO history_transaction_participants VALUES (23, 8589938688, 3);
INSERT INTO history_transaction_participants VALUES (24, 8589938688, 1);
INSERT INTO history_transaction_participants VALUES (25, 8589942784, 3);
INSERT INTO history_transaction_participants VALUES (26, 8589942784, 2);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_transaction_participants_id_seq', 26, true);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transactions VALUES ('720754335df899ed59dac9b35c89a1de60a7054c006db48ce12dc5081a6bbc5f', 7, 1, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934599, 100, 1, '2019-06-03 16:36:44.314976', '2019-06-03 16:36:44.314976', 30064775168, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABA5VdEXTh7kOlnlAmZykb462dFL6+URfv7kn222WD7uoXBt4zp0JTSPtB3DGyhjfAtX25rFJcc7YdXNuhQbSNDAw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAI8NCfAAAAAIAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAI8NCfAAAAAIAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck04AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABZaC44AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0J8AAAAAgAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAdzWF8AAAAAgAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0OoAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0NEAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{5VdEXTh7kOlnlAmZykb462dFL6+URfv7kn222WD7uoXBt4zp0JTSPtB3DGyhjfAtX25rFJcc7YdXNuhQbSNDAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('dfb03490f1faae8720598126bd770af8e7f081ac8e0683cea55a8aa35a6ba60a', 7, 2, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934600, 100, 1, '2019-06-03 16:36:44.315227', '2019-06-03 16:36:44.315227', 30064779264, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABANTM8JSXXbrK4BBvy/6y8Q6Uxw7u5HeV4ZfjafiHnXdepdTLsX1vObQgHDKwFSQ1bVtDORqDrhJ9ljRe4HgiRBQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAHc1hfAAAAAIAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAHc1hfAAAAAIAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABZaC44AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABfXg84AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAdzWF8AAAAAgAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAX14B8AAAAAgAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0NEAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0LgAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{NTM8JSXXbrK4BBvy/6y8Q6Uxw7u5HeV4ZfjafiHnXdepdTLsX1vObQgHDKwFSQ1bVtDORqDrhJ9ljRe4HgiRBQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('c9579846d3165e8f4271222caa4ce10d4ba1244011a8204dc68b810f055840fa', 7, 3, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934601, 100, 1, '2019-06-03 16:36:44.315378', '2019-06-03 16:36:44.315378', 30064783360, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAJAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAzc3x+vst1GnJX/uxRYwVtT877yIiHEiZQyAgYG+P+4pnqEM6h/+9NNgotuSXCVb8dfbGanBDQVE/qrmUInYiBw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAF9eAfAAAAAIAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAF9eAfAAAAAIAAAAJAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABfXg84AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABlU/A4AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAX14B8AAAAAgAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAR4Z98AAAAAgAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0LgAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0J8AAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{zc3x+vst1GnJX/uxRYwVtT877yIiHEiZQyAgYG+P+4pnqEM6h/+9NNgotuSXCVb8dfbGanBDQVE/qrmUInYiBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('f60bf7761d5459f1087262ca47486c901808d239486096c92541efa445cc4fe9', 6, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934594, 100, 1, '2019-06-03 16:36:44.343487', '2019-06-03 16:36:44.343488', 25769807872, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAAAAAAAX14QAAAAAAAAAAAa7kvkwAAABA+lv7NIE3yrIXlVPXxn1pYF38xMsqkaa42kprQQwAQlAdG8ICI4t+ZLX4pel6cAZFGYx73fZyXKBHruV0RvNGCA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAU3JNOAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABgAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAU3JNOAAAAAIAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0OoAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAApuSSoAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck04AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABNfGw4AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck2cAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck04AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{+lv7NIE3yrIXlVPXxn1pYF38xMsqkaa42kprQQwAQlAdG8ICI4t+ZLX4pel6cAZFGYx73fZyXKBHruV0RvNGCA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('9e56352113669f286d4950e55d0d9706fd41a0b3a0b74d5e51268bd0298e7e47', 6, 2, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934598, 100, 1, '2019-06-03 16:36:44.343733', '2019-06-03 16:36:44.343733', 25769811968, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAFeeBKecfo5v061dy3QfbF8zgO6gEUR8ildkKig42N0Yl4Z437Kpj4M0LWfDibhvKd6+voXM5rEFLVMpYFr/5AA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAKbkkqAAAAAIAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABgAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAKbkkqAAAAAIAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABNfGw4AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck04AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAApuSSoAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0OoAAAAAgAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0QMAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0OoAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{FeeBKecfo5v061dy3QfbF8zgO6gEUR8ildkKig42N0Yl4Z437Kpj4M0LWfDibhvKd6+voXM5rEFLVMpYFr/5AA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('f04617a26c4212f15a73578f59ea9913500fcb4818f828c17190f1454a04186c', 5, 1, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934594, 100, 1, '2019-06-03 16:36:44.358147', '2019-06-03 16:36:44.358147', 21474840576, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAG5HCDK2pf/77Ppffgv5hal7Q0yyfubULLN9szm3nJYL9YT60pLsuIC4YSwxAyVvsUHyQ3iJ48EQ+3VS/uIiiAg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rIDAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rIDAAAAAIAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKqcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msgMAAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA1pOcMAAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msk4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{G5HCDK2pf/77Ppffgv5hal7Q0yyfubULLN9szm3nJYL9YT60pLsuIC4YSwxAyVvsUHyQ3iJ48EQ+3VS/uIiiAg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('b31dce4295a5ef0e8f19fb83816ed880d573e5a2669f85be5cc602a90f9250c5', 5, 2, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934595, 100, 1, '2019-06-03 16:36:44.3585', '2019-06-03 16:36:44.3585', 21474844672, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAADAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAgzQF38kzgeHKf4Y3rZKYXturoU3n2LXyyuISFdK6/D5seTrjOXHU+m4kiIVeWUNtHx7ep3MSD1wIXuKjT0ReAA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAANaTnDAAAAAIAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAANaTnDAAAAAIAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKqcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABHhoucAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA1pOcMAAAAAgAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAvrwYMAAAAAgAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msk4AAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msjUAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{gzQF38kzgeHKf4Y3rZKYXturoU3n2LXyyuISFdK6/D5seTrjOXHU+m4kiIVeWUNtHx7ep3MSD1wIXuKjT0ReAA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('b0db02f48666502f7bb227049e48f7805380f07aee6e704b6a12b25d48c7b591', 5, 3, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934596, 100, 1, '2019-06-03 16:36:44.359124', '2019-06-03 16:36:44.359124', 21474848768, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABABqMU3WJT6ur297rvjulCqylVeC1bNKyQbClqyad+ou+x8u7GtYDf6o+aP/sLKitYYGnlDUvpTgdIyuMqSncYAw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAL68GDAAAAAIAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAL68GDAAAAAIAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABHhoucAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABNfGycAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAvrwYMAAAAAgAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAApuSUMAAAAAgAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msjUAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7mshwAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{BqMU3WJT6ur297rvjulCqylVeC1bNKyQbClqyad+ou+x8u7GtYDf6o+aP/sLKitYYGnlDUvpTgdIyuMqSncYAw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('4d1e89bd9835caded110ada3ace56bd1a919b8e53e3bb458aeb51330068269c6', 5, 4, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934597, 100, 1, '2019-06-03 16:36:44.359529', '2019-06-03 16:36:44.35953', 21474852864, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAKDnpoJwX4M7e52BdtZadIq1SC7dJxAjJDiXzMAK6ysLY2VKGVvXWs/RWmZYiXIkDO0ECyKfIov+1y4stQypZDw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAKbklDAAAAAIAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAKbklDAAAAAIAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABNfGycAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABTck2cAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAApuSUMAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAAjw0QMAAAAAgAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7mshwAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msgMAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{KDnpoJwX4M7e52BdtZadIq1SC7dJxAjJDiXzMAK6ysLY2VKGVvXWs/RWmZYiXIkDO0ECyKfIov+1y4stQypZDw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d178e262b8b2aac66a75f69e70ce3cb7fdf01a6060433636c7b4a3a178236429', 4, 1, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934593, 100, 1, '2019-06-03 16:36:44.378888', '2019-06-03 16:36:44.378888', 17179873280, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAAAAAAAX14QAAAAAAAAAAAa7kvkwAAABA9qncr+3eaHaYqpDspvoIbiENnY3te9dqrCYtGbiT13CWh/b+cm+CUe9//x0NDxiU/ptY0QlY/z54IF7jF0H7CQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAQZCqnAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAQZCqnAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAADAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA1pOicAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKqcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKsAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKqcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{9qncr+3eaHaYqpDspvoIbiENnY3te9dqrCYtGbiT13CWh/b+cm+CUe9//x0NDxiU/ptY0QlY/z54IF7jF0H7CQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('aeb26130d1715c5e9b0c85de46d454200cec26513e8ce06ab05628069bea0793', 3, 1, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON', 8589934593, 100, 1, '2019-06-03 16:36:44.395647', '2019-06-03 16:36:44.395647', 12884905984, 'AAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAAAAAAAX14QAAAAAAAAAAAW8UiFoAAABAJ1AToo3dEH9+7//OjpIHtWCDsL/0MUQlbjUSQC2+I3TVEl9chqrpqx5GG6yjN8INl3IZ7/HSA0EfRB2xZ9VMCg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rJnAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rJnAAAAAIAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAABBkKsAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAADAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA1pOicAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msmcAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{J1AToo3dEH9+7//OjpIHtWCDsL/0MUQlbjUSQC2+I3TVEl9chqrpqx5GG6yjN8INl3IZ7/HSA0EfRB2xZ9VMCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d', 2, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-03 16:36:44.406035', '2019-06-03 16:36:44.406036', 8589938688, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rKAAAAAAAAAAABVvwF9wAAAECDzqvkQBQoNAJifPRXDoLhvtycT3lFPCQ51gkdsFHaBNWw05S/VhW0Xgkr0CBPE4NaFV2Kmcs3ZwLmib4TRrML', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrNryTU4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{g86r5EAUKDQCYnz0Vw6C4b7cnE95RTwkOdYJHbBR2gTVsNOUv1YVtF4JK9AgTxODWhVdipnLN2cC5om+E0azCw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('1ed82961fb013d39f96aa7e428c4174caa4a5a43dbc65713a37d46d96ee5c314', 2, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 2, 100, 1, '2019-06-03 16:36:44.40629', '2019-06-03 16:36:44.40629', 8589942784, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rKAAAAAAAAAAABVvwF9wAAAECuHh3Q7Zpbulw+xzp7NeMPdfYErNzJOrvQi8GOkN7WgfwSPgzHcPE/E/s8CL/AQrjBtw067aUZAvoaVf12oCQB', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAAAAAABAAAAAwAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrNryTU4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrMwLms4AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{rh4d0O2aW7pcPsc6ezXjD3X2BKzcyTq70IvBjpDe1oH8Ej4Mx3DxPxP7PAi/wEK4wbcNOu2lGQL6GlX9dqAkAQ==}', 'none', NULL, NULL, true, 100);


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

