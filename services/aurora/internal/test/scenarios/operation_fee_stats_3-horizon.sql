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

INSERT INTO history_accounts VALUES (1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO history_accounts VALUES (2, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_accounts_id_seq', 2, true);


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

INSERT INTO history_effects VALUES (1, 38654709761, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 38654709761, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 38654713857, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 38654713857, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 38654717953, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 38654717953, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 34359742465, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 34359742465, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 30064775169, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 30064775169, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 25769807873, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 25769807873, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 21474840577, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 21474840577, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 17179873281, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 17179873281, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (1, 17179873282, 1, 2, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 17179873282, 2, 3, '{"amount": "10.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 12884905985, 1, 0, '{"starting_balance": "1000.0000000"}');
INSERT INTO history_effects VALUES (1, 12884905985, 2, 3, '{"amount": "1000.0000000", "asset_type": "native"}');
INSERT INTO history_effects VALUES (2, 12884905985, 3, 10, '{"weight": 1, "public_key": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_ledgers VALUES (9, '2bf34a98b988141b851bcef629b7ef31ed7e74ca39000cdf92d002a91d1c6935', 'da47b612dc8d5b2fa7bf9bfa09db11f3dbc3badf2497083a519978698880a78f', 3, 3, '2019-06-03 16:37:31', '2019-06-03 16:37:27.870696', '2019-06-03 16:37:27.870696', 38654705664, 16, 1000000000000000000, 1000, 100, 100000000, 1000000, 11, 'AAAAC9pHthLcjVsvp7+b+gnbEfPbw7rfJJcIOlGZeGmIgKePkJ5YHJBvZSO3FWn9fV3rAIQ0bm9gsHzXfU7tEazHg/8AAAAAXPVMywAAAAAAAAAAdWzSHhfCzXGLIkOOjoC5Ym8OBprpH5cEOKDNFjF7mH1f3SyVVZeHY50rSogB2TlMZKc8GcCSomNW+7rp7NdslwAAAAkN4Lazp2QAAAAAAAAAAAPoAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO history_ledgers VALUES (8, 'da47b612dc8d5b2fa7bf9bfa09db11f3dbc3badf2497083a519978698880a78f', '7c5e9b5ec8ee15e5fbde0406499973a85a0d90b243addce0c1ff16a8510970e0', 1, 1, '2019-06-03 16:37:30', '2019-06-03 16:37:27.889132', '2019-06-03 16:37:27.889132', 34359738368, 16, 1000000000000000000, 700, 100, 100000000, 1000000, 11, 'AAAAC3xem17I7hXl+94EBkmZc6haDZCyQ63c4MH/FqhRCXDgKSvRyi2Ch/rUCgC7mvyEgkoSjwLop2En+fKX4u77S9sAAAAAXPVMygAAAAAAAAAA1WZy5HNzQG6v2/yQfdo/18UP2Mi3LOX+HsIFwjtPY3dqPyAyz0bMBb26v8R9NyFd8BQr6IfTXpUfoqWipxdbCAAAAAgN4Lazp2QAAAAAAAAAAAK8AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (7, '7c5e9b5ec8ee15e5fbde0406499973a85a0d90b243addce0c1ff16a8510970e0', 'f5d5ec2689d187facad52c2ef940f0b3a86822c24e8fdaeacf6f7317d36134e6', 1, 1, '2019-06-03 16:37:29', '2019-06-03 16:37:27.900016', '2019-06-03 16:37:27.900017', 30064771072, 16, 1000000000000000000, 600, 100, 100000000, 1000000, 11, 'AAAAC/XV7CaJ0Yf6ytUsLvlA8LOoaCLCTo/a6s9vcxfTYTTmi3GZssPukt5B4pW9nYeAi9c4fBJ/G2009FXT+aofPn0AAAAAXPVMyQAAAAAAAAAAOgSmIClOx9JX9c2gFj6mbAWu2y1mlO0kGdlw5hFcYz1uNWrkAVJl78mTro7woogMHVmygB4fJpIjuQr0kHVoUAAAAAcN4Lazp2QAAAAAAAAAAAJYAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (6, 'f5d5ec2689d187facad52c2ef940f0b3a86822c24e8fdaeacf6f7317d36134e6', '5ae88507efcf646686d3ed03321124529f7b9928bf0531fde31aa55baa3b0fec', 1, 1, '2019-06-03 16:37:28', '2019-06-03 16:37:27.910025', '2019-06-03 16:37:27.910025', 25769803776, 16, 1000000000000000000, 500, 100, 100000000, 1000000, 11, 'AAAAC1rohQfvz2RmhtPtAzIRJFKfe5kovwUx/eMapVuqOw/sn/+p+bTn7uolyYuy5JK2LC1rXH1wYlbg1Er/DS2nlQ0AAAAAXPVMyAAAAAAAAAAAqqhBNhS+oY0Fq/RSpEEko+knV66p0YlC+DiVf0GYlCGg9YIn6PSEfLVsngLGVPIAMbBYSS/o7J4XE21595K6KwAAAAYN4Lazp2QAAAAAAAAAAAH0AAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (5, '5ae88507efcf646686d3ed03321124529f7b9928bf0531fde31aa55baa3b0fec', '10b6161843eb9d652e093486fa8174f93fc44a538870829f7f0ac47da8faa651', 1, 1, '2019-06-03 16:37:27', '2019-06-03 16:37:27.920568', '2019-06-03 16:37:27.920568', 21474836480, 16, 1000000000000000000, 400, 100, 100000000, 1000000, 11, 'AAAACxC2FhhD651lLgk0hvqBdPk/xEpTiHCCn38KxH2o+qZRD70AR0B0Om/2/JDgto1Z3HdONB5TdVAW8+uTkr5a4g0AAAAAXPVMxwAAAAAAAAAAfZ5r/iAId0HYP/J85oWovpWvupf3zLdy0YhEdBPnW26KeNNzkVWivC8sNo3dpubpXqq5+Y4rtQrnFnp9TEn+ngAAAAUN4Lazp2QAAAAAAAAAAAGQAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (4, '10b6161843eb9d652e093486fa8174f93fc44a538870829f7f0ac47da8faa651', '647dbe45977be72bc6fe73b13b5476a93b405cd7a562af189628c7cbd41ec891', 1, 2, '2019-06-03 16:37:26', '2019-06-03 16:37:27.931569', '2019-06-03 16:37:27.931569', 17179869184, 16, 1000000000000000000, 300, 100, 100000000, 1000000, 11, 'AAAAC2R9vkWXe+crxv5zsTtUdqk7QFzXpWKvGJYox8vUHsiRiP4WGa5E15yHbJ/Fjs7aniDcPvHvQvRRTagpeuy2LZcAAAAAXPVMxgAAAAAAAAAAvUi8p6JAmXRap5H687O+OXcAKN2//CsLZbyxWLVlZb9ZUMyJYVoYrjI5mh3fV0Ba3qZVHr/lKgv/ck1FfNB45gAAAAQN4Lazp2QAAAAAAAAAAAEsAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (3, '647dbe45977be72bc6fe73b13b5476a93b405cd7a562af189628c7cbd41ec891', '87d37dcbbad4bced6acfa589f9869dea83939ff0f1dbfef99727c9b8ecb0e3f1', 1, 1, '2019-06-03 16:37:25', '2019-06-03 16:37:27.942838', '2019-06-03 16:37:27.942839', 12884901888, 16, 1000000000000000000, 100, 100, 100000000, 1000000, 11, 'AAAAC4fTfcu61Lztas+lifmGneqDk5/w8dv++ZcnybjssOPx6p7KmzoLq6THlwtzdRXidN9Q1s+IoUveox2wGSAVnYkAAAAAXPVMxQAAAAAAAAAAfFnUZMxpcaRFgW684JniUG/dzZ5jn4eP2mZ8LIGonSrIQlCO7YXSBRlsJC2BmSJC6FqFvhzU66vHkiNcsFs6hAAAAAMN4Lazp2QAAAAAAAAAAABkAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 1, 0);
INSERT INTO history_ledgers VALUES (2, '87d37dcbbad4bced6acfa589f9869dea83939ff0f1dbfef99727c9b8ecb0e3f1', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 0, 0, '2019-06-03 16:37:24', '2019-06-03 16:37:27.952867', '2019-06-03 16:37:27.952867', 8589934592, 16, 1000000000000000000, 0, 100, 100000000, 1000000, 11, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76ZuZRHr9UdXKbTKiclfOjy72YZFJUkJPVcKT5htvorm1QAAAAAXPVMxAAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERnv+Ltt/ycz/x8/+lFB80rnVx7j2Mrm29EpusUR+gv9ZAAAAAIN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);
INSERT INTO history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-03 16:37:27.958898', '2019-06-03 16:37:27.958898', 4294967296, 16, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operation_participants VALUES (1, 38654709761, 2);
INSERT INTO history_operation_participants VALUES (2, 38654709761, 1);
INSERT INTO history_operation_participants VALUES (3, 38654713857, 2);
INSERT INTO history_operation_participants VALUES (4, 38654713857, 1);
INSERT INTO history_operation_participants VALUES (5, 38654717953, 2);
INSERT INTO history_operation_participants VALUES (6, 38654717953, 1);
INSERT INTO history_operation_participants VALUES (7, 34359742465, 2);
INSERT INTO history_operation_participants VALUES (8, 34359742465, 1);
INSERT INTO history_operation_participants VALUES (9, 30064775169, 2);
INSERT INTO history_operation_participants VALUES (10, 30064775169, 1);
INSERT INTO history_operation_participants VALUES (11, 25769807873, 2);
INSERT INTO history_operation_participants VALUES (12, 25769807873, 1);
INSERT INTO history_operation_participants VALUES (13, 21474840577, 2);
INSERT INTO history_operation_participants VALUES (14, 21474840577, 1);
INSERT INTO history_operation_participants VALUES (15, 17179873281, 2);
INSERT INTO history_operation_participants VALUES (16, 17179873281, 1);
INSERT INTO history_operation_participants VALUES (17, 17179873282, 2);
INSERT INTO history_operation_participants VALUES (18, 17179873282, 1);
INSERT INTO history_operation_participants VALUES (19, 12884905985, 1);
INSERT INTO history_operation_participants VALUES (20, 12884905985, 2);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_operation_participants_id_seq', 20, true);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_operations VALUES (38654709761, 38654709760, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (38654713857, 38654713856, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (38654717953, 38654717952, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (34359742465, 34359742464, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (30064775169, 30064775168, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (25769807873, 25769807872, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (21474840577, 21474840576, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (17179873281, 17179873280, 1, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (17179873282, 17179873280, 2, 1, '{"to": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "from": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "amount": "10.0000000", "asset_type": "native"}', 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY');
INSERT INTO history_operations VALUES (12884905985, 12884905984, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY", "starting_balance": "1000.0000000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transaction_participants VALUES (1, 38654709760, 2);
INSERT INTO history_transaction_participants VALUES (2, 38654709760, 1);
INSERT INTO history_transaction_participants VALUES (3, 38654713856, 2);
INSERT INTO history_transaction_participants VALUES (4, 38654713856, 1);
INSERT INTO history_transaction_participants VALUES (5, 38654717952, 2);
INSERT INTO history_transaction_participants VALUES (6, 38654717952, 1);
INSERT INTO history_transaction_participants VALUES (7, 34359742464, 2);
INSERT INTO history_transaction_participants VALUES (8, 34359742464, 1);
INSERT INTO history_transaction_participants VALUES (9, 30064775168, 2);
INSERT INTO history_transaction_participants VALUES (10, 30064775168, 1);
INSERT INTO history_transaction_participants VALUES (11, 25769807872, 2);
INSERT INTO history_transaction_participants VALUES (12, 25769807872, 1);
INSERT INTO history_transaction_participants VALUES (13, 21474840576, 2);
INSERT INTO history_transaction_participants VALUES (14, 21474840576, 1);
INSERT INTO history_transaction_participants VALUES (15, 17179873280, 2);
INSERT INTO history_transaction_participants VALUES (16, 17179873280, 1);
INSERT INTO history_transaction_participants VALUES (17, 12884905984, 1);
INSERT INTO history_transaction_participants VALUES (18, 12884905984, 2);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('history_transaction_participants_id_seq', 18, true);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO history_transactions VALUES ('6a349e7331e93a251367287e274fb1699abaf723bde37aebe96248c76fd3071a', 9, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901894, 400, 1, '2019-06-03 16:37:27.870885', '2019-06-03 16:37:27.870885', 38654709760, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABkAAAAAMAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAvG2IEoAgIDgfSZC0D4ClAMlvU8rCmn1JtgrmtA9HShVsqoMPeyC8rbXu+Dizq74y9TSl1/9P37YY9kWfU09oBw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACMEiafAAAAAMAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACMEiafAAAAAMAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrF3G2GcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrF9EUKcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJp8AAAAAwAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIqUrl8AAAAAwAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAIAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJuoAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJtEAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{vG2IEoAgIDgfSZC0D4ClAMlvU8rCmn1JtgrmtA9HShVsqoMPeyC8rbXu+Dizq74y9TSl1/9P37YY9kWfU09oBw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('9a719ea0bc6fd18082cbaec8d1f06c074e6c6aa784fa9ee9f0b015cf8a398bd5', 9, 2, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901895, 400, 1, '2019-06-03 16:37:27.87116', '2019-06-03 16:37:27.871161', 38654713856, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABkAAAAAMAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAxG3ZbC4djlBXwWQidTeJb/7Q2fr0GPD1mx/2bF++HE+eBPrKP0ol1VSNUQVaW7mMcdFjQcTHSb+uBoq+kd3dCg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACKlK5fAAAAAMAAAAGAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACKlK5fAAAAAMAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrF9EUKcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrGDByOcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIqUrl8AAAAAwAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIkXNh8AAAAAwAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJtEAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJrgAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{xG3ZbC4djlBXwWQidTeJb/7Q2fr0GPD1mx/2bF++HE+eBPrKP0ol1VSNUQVaW7mMcdFjQcTHSb+uBoq+kd3dCg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('25ded52d9314195e638c758b6eeef7cd07c0cf4c896697f6d5cb228c44dacdd8', 9, 3, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901896, 400, 1, '2019-06-03 16:37:27.871322', '2019-06-03 16:37:27.871323', 38654717952, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABkAAAAAMAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAa2qrw54P1lv9IGMKjXGfCNlcdCRXl33v57V+uAmZYf1UvGMsakdNbZFHENg75vdnxM4aHyAcrTMoSTqyvMc7CQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACJFzYfAAAAAMAAAAHAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACJFzYfAAAAAMAAAAIAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrGDByOcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrGI/QScAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIkXNh8AAAAAwAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIeZvd8AAAAAwAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJrgAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAJAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJp8AAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{a2qrw54P1lv9IGMKjXGfCNlcdCRXl33v57V+uAmZYf1UvGMsakdNbZFHENg75vdnxM4aHyAcrTMoSTqyvMc7CQ==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('fbeb854b57c7ea853028f23ebe71de61c1ecbd8a64f6437da735ee37883ce558', 8, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901893, 300, 1, '2019-06-03 16:37:27.889321', '2019-06-03 16:37:27.889321', 34359742464, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABLAAAAAMAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABArAAIYpB4GOYOqjJiwKvRsZ+V3AZXshTLQb5MRvOuue/lSawV12iNSTEBIpPOqYUc0hfVudWfmLd2aWZ5UQd9AA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAACAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACNj58qAAAAAMAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAACAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACNj58qAAAAAMAAAAFAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFxJYCcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrF3G2GcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAIAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI2PnyoAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAIwSJuoAAAAAwAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAHAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI2Pn0MAAAAAwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAIAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI2PnyoAAAAAwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{rAAIYpB4GOYOqjJiwKvRsZ+V3AZXshTLQb5MRvOuue/lSawV12iNSTEBIpPOqYUc0hfVudWfmLd2aWZ5UQd9AA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('d2a62bf7b9e118b182c33b2fd93b2cc2013dbe9a8d77f35a239b70c8a667e5e5', 7, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901892, 400, 1, '2019-06-03 16:37:27.900178', '2019-06-03 16:37:27.900178', 30064775168, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABkAAAAAMAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAcKnXL1cr7aTkY83f55Oh0M/PNjPSTaZooDIfmoZz16BgDN94hqraJ73vmRdHmqtJaKYdwtcNgovdEvVxFYaIBg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABwAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACPDReDAAAAAMAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABwAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACPDReDAAAAAMAAAAEAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAGAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFrL5+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFxJYCcAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAHAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI8NF4MAAAAAwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI2Pn0MAAAAAwAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAGAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI8NF5wAAAAAwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAHAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI8NF4MAAAAAwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{cKnXL1cr7aTkY83f55Oh0M/PNjPSTaZooDIfmoZz16BgDN94hqraJ73vmRdHmqtJaKYdwtcNgovdEvVxFYaIBg==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('b4499cd4bc782623f9ac9654040d49c154fab6ab8d83b2110002c620a5eb7407', 6, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901891, 400, 1, '2019-06-03 16:37:27.910253', '2019-06-03 16:37:27.910254', 25769807872, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAABkAAAAAMAAAADAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAABfxa1tvLDgKKRnsVwm97GeZmHtvBJee12Q49wseNvKHjwb0amqXGJVYFN7PGH5ZZ56Se9GvyiL99zLLTz29Dw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACQio/cAAAAAMAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABgAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACQio/cAAAAAMAAAADAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAFAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFlOb6cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFrL5+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAGAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJCKj9wAAAAAwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAI8NF5wAAAAAwAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAFAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJCKj/UAAAAAwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAGAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJCKj9wAAAAAwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ABfxa1tvLDgKKRnsVwm97GeZmHtvBJee12Q49wseNvKHjwb0amqXGJVYFN7PGH5ZZ56Se9GvyiL99zLLTz29Dw==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('b8fd5e6ed3d2658aa66040319e076e30006f7950e18e9a03e1eddeedfccbb418', 5, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901890, 200, 1, '2019-06-03 16:37:27.920777', '2019-06-03 16:37:27.920777', 21474840576, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAAAyAAAAAMAAAACAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAdkFYl4AAABAY8zQeTlk6qu1feh/23t9EMxnoOW+6moGmjXKum57BkkQq6zoV/VciJ7IVIpi+jPVZSk+KSrCQdAm6EV4jBbvBA==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACSCAg1AAAAAMAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABQAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACSCAg1AAAAAMAAAACAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAMAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFfQ92cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFlOb6cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAFAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJIICDUAAAAAwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJCKj/UAAAAAwAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAAEAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJIICE4AAAAAwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJIICDUAAAAAwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Y8zQeTlk6qu1feh/23t9EMxnoOW+6moGmjXKum57BkkQq6zoV/VciJ7IVIpi+jPVZSk+KSrCQdAm6EV4jBbvBA==}', 'none', NULL, NULL, true, 100);
INSERT INTO history_transactions VALUES ('ba38e7c204b3f8ab8907a4b9618417854bccb54a7fa494a36c3d185bb45d07d6', 4, 1, 'GAG52TW6QAB6TGNMOTL32Y4M3UQQLNNNHPEHYAIYRP6SFF6ZAVRF5ZQY', 12884901889, 200, 2, '2019-06-03 16:37:27.931785', '2019-06-03 16:37:27.931785', 17179873280, 'AAAAAA3dTt6AA+mZrHTXvWOM3SEFta07yHwBGIv9IpfZBWJeAAAAyAAAAAMAAAABAAAAAAAAAAAAAAACAAAAAAAAAAEAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcAAAAAAAAAAAX14QAAAAAAAAAAAQAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9wAAAAAAAAAABfXhAAAAAAAAAAAB2QViXgAAAEBzT3nPm0xtu6CkU5jiXuBFFlZ9yTXnlEKy5HLcoVo9ym4phM8ja3knZbLZ4zJiNklsNl99mmSVkJKz7XXgOXEH', 'AAAAAAAAAMgAAAAAAAAAAgAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAABAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvjOAAAAAMAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvjOAAAAAMAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAACAAAABAAAAAMAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBucAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFZTfycAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAMAAAAEAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+M4AAAAAwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJOFgI4AAAAAwAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAQAAAADAAAABAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaxWU38nAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4LaxX0PdnAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAADAAAABAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACThYCOAAAAAMAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAABAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACSCAhOAAAAAMAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAA=', 'AAAAAgAAAAMAAAADAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+QAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+M4AAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{c095z5tMbbugpFOY4l7gRRZWfck155RCsuRy3KFaPcpuKYTPI2t5J2Wy2eMyYjZJbDZffZpklZCSs+114DlxBw==}', 'none', NULL, NULL, true, 200);
INSERT INTO history_transactions VALUES ('f1d63c0b88a1ab68a44bcd02e7c9dd7c7da818ac1ff87762e922acac9958766e', 3, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-03 16:37:27.943187', '2019-06-03 16:37:27.943188', 12884905984, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAADd1O3oAD6ZmsdNe9Y4zdIQW1rTvIfAEYi/0il9kFYl4AAAACVAvkAAAAAAAAAAABVvwF9wAAAEDUWAnn6bBg8wR8y/D76fh6M+FmmxKaCQL33EyRWWYFxlFN4w2rpaZ3uW69gVg3ooM8LCkF+P8AWaxcKBMjrBMC', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAQAAAAIAAAADAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAAAAAAAYvwdC9CRsrYcDdZWNGsqaNfTR8bywsjubQRHAlb8BfcN4Lazp2P/nAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAABAAAAAwAAAAMAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrFTWBucAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAN3U7egAPpmax0171jjN0hBbWtO8h8ARiL/SKX2QViXgAAAAJUC+QAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{1FgJ5+mwYPMEfMvw++n4ejPhZpsSmgkC99xMkVlmBcZRTeMNq6Wmd7luvYFYN6KDPCwpBfj/AFmsXCgTI6wTAg==}', 'none', NULL, NULL, true, 100);


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

