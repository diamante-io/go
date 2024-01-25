--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.1
-- Dumped by pg_dump version 9.6.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

DROP INDEX IF EXISTS public.upgradehistbyseq;
DROP INDEX IF EXISTS public.scpquorumsbyseq;
DROP INDEX IF EXISTS public.scpenvsbyseq;
DROP INDEX IF EXISTS public.ledgersbyseq;
DROP INDEX IF EXISTS public.histfeebyseq;
DROP INDEX IF EXISTS public.histbyseq;
DROP INDEX IF EXISTS public.bestofferindex;
DROP INDEX IF EXISTS public.accountbalances;
ALTER TABLE IF EXISTS ONLY public.upgradehistory DROP CONSTRAINT IF EXISTS upgradehistory_pkey;
ALTER TABLE IF EXISTS ONLY public.txhistory DROP CONSTRAINT IF EXISTS txhistory_pkey;
ALTER TABLE IF EXISTS ONLY public.txfeehistory DROP CONSTRAINT IF EXISTS txfeehistory_pkey;
ALTER TABLE IF EXISTS ONLY public.trustlines DROP CONSTRAINT IF EXISTS trustlines_pkey;
ALTER TABLE IF EXISTS ONLY public.storestate DROP CONSTRAINT IF EXISTS storestate_pkey;
ALTER TABLE IF EXISTS ONLY public.scpquorums DROP CONSTRAINT IF EXISTS scpquorums_pkey;
ALTER TABLE IF EXISTS ONLY public.quoruminfo DROP CONSTRAINT IF EXISTS quoruminfo_pkey;
ALTER TABLE IF EXISTS ONLY public.pubsub DROP CONSTRAINT IF EXISTS pubsub_pkey;
ALTER TABLE IF EXISTS ONLY public.publishqueue DROP CONSTRAINT IF EXISTS publishqueue_pkey;
ALTER TABLE IF EXISTS ONLY public.peers DROP CONSTRAINT IF EXISTS peers_pkey;
ALTER TABLE IF EXISTS ONLY public.offers DROP CONSTRAINT IF EXISTS offers_pkey;
ALTER TABLE IF EXISTS ONLY public.ledgerheaders DROP CONSTRAINT IF EXISTS ledgerheaders_pkey;
ALTER TABLE IF EXISTS ONLY public.ledgerheaders DROP CONSTRAINT IF EXISTS ledgerheaders_ledgerseq_key;
ALTER TABLE IF EXISTS ONLY public.ban DROP CONSTRAINT IF EXISTS ban_pkey;
ALTER TABLE IF EXISTS ONLY public.accounts DROP CONSTRAINT IF EXISTS accounts_pkey;
ALTER TABLE IF EXISTS ONLY public.accountdata DROP CONSTRAINT IF EXISTS accountdata_pkey;
DROP TABLE IF EXISTS public.upgradehistory;
DROP TABLE IF EXISTS public.txhistory;
DROP TABLE IF EXISTS public.txfeehistory;
DROP TABLE IF EXISTS public.trustlines;
DROP TABLE IF EXISTS public.storestate;
DROP TABLE IF EXISTS public.scpquorums;
DROP TABLE IF EXISTS public.scphistory;
DROP TABLE IF EXISTS public.quoruminfo;
DROP TABLE IF EXISTS public.pubsub;
DROP TABLE IF EXISTS public.publishqueue;
DROP TABLE IF EXISTS public.peers;
DROP TABLE IF EXISTS public.offers;
DROP TABLE IF EXISTS public.ledgerheaders;
DROP TABLE IF EXISTS public.ban;
DROP TABLE IF EXISTS public.accounts;
DROP TABLE IF EXISTS public.accountdata;
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

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accountdata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accountdata (
    accountid character varying(56) NOT NULL,
    dataname character varying(88) NOT NULL,
    datavalue character varying(112) NOT NULL,
    lastmodified integer NOT NULL
);


--
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accounts (
    accountid character varying(56) NOT NULL,
    balance bigint NOT NULL,
    seqnum bigint NOT NULL,
    numsubentries integer NOT NULL,
    inflationdest character varying(56),
    homedomain character varying(44) NOT NULL,
    thresholds text NOT NULL,
    flags integer NOT NULL,
    lastmodified integer NOT NULL,
    buyingliabilities bigint,
    sellingliabilities bigint,
    signers text,
    CONSTRAINT accounts_balance_check CHECK ((balance >= 0)),
    CONSTRAINT accounts_buyingliabilities_check CHECK ((buyingliabilities >= 0)),
    CONSTRAINT accounts_numsubentries_check CHECK ((numsubentries >= 0)),
    CONSTRAINT accounts_sellingliabilities_check CHECK ((sellingliabilities >= 0))
);


--
-- Name: ban; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE ban (
    nodeid character(56) NOT NULL
);


--
-- Name: ledgerheaders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE ledgerheaders (
    ledgerhash character(64) NOT NULL,
    prevhash character(64) NOT NULL,
    bucketlisthash character(64) NOT NULL,
    ledgerseq integer,
    closetime bigint NOT NULL,
    data text NOT NULL,
    CONSTRAINT ledgerheaders_closetime_check CHECK ((closetime >= 0)),
    CONSTRAINT ledgerheaders_ledgerseq_check CHECK ((ledgerseq >= 0))
);


--
-- Name: offers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE offers (
    sellerid character varying(56) NOT NULL,
    offerid bigint NOT NULL,
    sellingasset text NOT NULL,
    buyingasset text NOT NULL,
    amount bigint NOT NULL,
    pricen integer NOT NULL,
    priced integer NOT NULL,
    price double precision NOT NULL,
    flags integer NOT NULL,
    lastmodified integer NOT NULL,
    CONSTRAINT offers_amount_check CHECK ((amount >= 0)),
    CONSTRAINT offers_offerid_check CHECK ((offerid >= 0))
);


--
-- Name: peers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE peers (
    ip character varying(15) NOT NULL,
    port integer DEFAULT 0 NOT NULL,
    nextattempt timestamp without time zone NOT NULL,
    numfailures integer DEFAULT 0 NOT NULL,
    type integer NOT NULL,
    CONSTRAINT peers_numfailures_check CHECK ((numfailures >= 0)),
    CONSTRAINT peers_port_check CHECK (((port > 0) AND (port <= 65535)))
);


--
-- Name: publishqueue; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE publishqueue (
    ledger integer NOT NULL,
    state text
);


--
-- Name: pubsub; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE pubsub (
    resid character(32) NOT NULL,
    lastread integer
);


--
-- Name: quoruminfo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE quoruminfo (
    nodeid character(56) NOT NULL,
    qsethash character(64) NOT NULL
);


--
-- Name: scphistory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE scphistory (
    nodeid character(56) NOT NULL,
    ledgerseq integer NOT NULL,
    envelope text NOT NULL,
    CONSTRAINT scphistory_ledgerseq_check CHECK ((ledgerseq >= 0))
);


--
-- Name: scpquorums; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE scpquorums (
    qsethash character(64) NOT NULL,
    lastledgerseq integer NOT NULL,
    qset text NOT NULL,
    CONSTRAINT scpquorums_lastledgerseq_check CHECK ((lastledgerseq >= 0))
);


--
-- Name: storestate; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE storestate (
    statename character(32) NOT NULL,
    state text
);


--
-- Name: trustlines; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE trustlines (
    accountid character varying(56) NOT NULL,
    assettype integer NOT NULL,
    issuer character varying(56) NOT NULL,
    assetcode character varying(12) NOT NULL,
    tlimit bigint NOT NULL,
    balance bigint NOT NULL,
    flags integer NOT NULL,
    lastmodified integer NOT NULL,
    buyingliabilities bigint,
    sellingliabilities bigint,
    CONSTRAINT trustlines_balance_check CHECK ((balance >= 0)),
    CONSTRAINT trustlines_buyingliabilities_check CHECK ((buyingliabilities >= 0)),
    CONSTRAINT trustlines_sellingliabilities_check CHECK ((sellingliabilities >= 0)),
    CONSTRAINT trustlines_tlimit_check CHECK ((tlimit > 0))
);


--
-- Name: txfeehistory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE txfeehistory (
    txid character(64) NOT NULL,
    ledgerseq integer NOT NULL,
    txindex integer NOT NULL,
    txchanges text NOT NULL,
    CONSTRAINT txfeehistory_ledgerseq_check CHECK ((ledgerseq >= 0))
);


--
-- Name: txhistory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE txhistory (
    txid character(64) NOT NULL,
    ledgerseq integer NOT NULL,
    txindex integer NOT NULL,
    txbody text NOT NULL,
    txresult text NOT NULL,
    txmeta text NOT NULL,
    CONSTRAINT txhistory_ledgerseq_check CHECK ((ledgerseq >= 0))
);


--
-- Name: upgradehistory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE upgradehistory (
    ledgerseq integer NOT NULL,
    upgradeindex integer NOT NULL,
    upgrade text NOT NULL,
    changes text NOT NULL,
    CONSTRAINT upgradehistory_ledgerseq_check CHECK ((ledgerseq >= 0))
);


--
-- Data for Name: accountdata; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO accounts VALUES ('GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1000000000000000000, 0, 0, NULL, '', 'AQAAAA==', 0, 1, NULL, NULL, NULL);


--
-- Data for Name: ban; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: ledgerheaders; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO ledgerheaders VALUES ('63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', '0000000000000000000000000000000000000000000000000000000000000000', '572a2e32ff248a07b0e70fd1f6d318c1facd20b6cc08c33d5775259868125a16', 1, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('fddd6d06d74e0d5677eb020b0ff83c321a12e1ec8c447225dccec603b7f2a20b', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 'eff8bb6dff2733ff1f3ffa5141f34ae7571ee3d8cae6dbd129bac511fa0bfd64', 2, 1559579827, 'AAAAC2PZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76ZuZRHr9UdXKbTKiclfOjy72YZFJUkJPVcKT5htvorm1QAAAAAXPVMswAAAAIAAAAIAAAAAQAAAAsAAAAIAAAAAwAPQkAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERnv+Ltt/ycz/x8/+lFB80rnVx7j2Mrm29EpusUR+gv9ZAAAAAIN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('3cdae71be1076874ec5a3955a3f27c5a8524ed514a76c76724254c9a9a16deb6', 'fddd6d06d74e0d5677eb020b0ff83c321a12e1ec8c447225dccec603b7f2a20b', 'eff8bb6dff2733ff1f3ffa5141f34ae7571ee3d8cae6dbd129bac511fa0bfd64', 3, 1559579828, 'AAAAC/3dbQbXTg1Wd+sCCw/4PDIaEuHsjERyJdzOxgO38qILONQ6wuACH6xfhjuX9wreJJSyXtsFr1EbF3uAJuDtKZAAAAAAXPVMtAAAAAAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERnv+Ltt/ycz/x8/+lFB80rnVx7j2Mrm29EpusUR+gv9ZAAAAAMN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('1e5ddd4a413af2671918d67cbc187a1cdc6039a0307430d915ca7456c16892e4', '3cdae71be1076874ec5a3955a3f27c5a8524ed514a76c76724254c9a9a16deb6', 'f79fa0bf4f0941a78d93bbee679f206d87fb0da208857e8fac6ce60968444614', 4, 1559579829, 'AAAACzza5xvhB2h07Fo5VaPyfFqFJO1RSnbHZyQlTJqaFt62YJ1PffsPE4Bgb2MGQtV4n4R1RhaPlEnNsv2tOFqFsDgAAAAAXPVMtQAAAAAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERn3n6C/TwlBp42Tu+5nnyBth/sNogiFfo+sbOYJaERGFAAAAAQN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('2f0ba640058008a15d8fba4436f4c24d7993994f9f73a426c7e23879ca31a93f', '1e5ddd4a413af2671918d67cbc187a1cdc6039a0307430d915ca7456c16892e4', 'f79fa0bf4f0941a78d93bbee679f206d87fb0da208857e8fac6ce60968444614', 5, 1559579830, 'AAAACx5d3UpBOvJnGRjWfLwYehzcYDmgMHQw2RXKdFbBaJLk6Q9O9qon6XPoND6LH8FnSyN245R2LczccTxE8t48/TAAAAAAXPVMtgAAAAAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERn3n6C/TwlBp42Tu+5nnyBth/sNogiFfo+sbOYJaERGFAAAAAUN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('a01cbeeaeb18cbed0ee913a13ba4f760660f0a6541be8f19021816601f96a7d8', '2f0ba640058008a15d8fba4436f4c24d7993994f9f73a426c7e23879ca31a93f', '6a6ce2f01ea7c9b517e0fe337cc3df702d30f312fc7389eceb9dc49db9e7785c', 6, 1559579831, 'AAAACy8LpkAFgAihXY+6RDb0wk15k5lPn3OkJsfiOHnKMak/UNAhrzjlDDebLKtCXYZeDsKoDrfnW+Amjcfeo9DYiesAAAAAXPVMtwAAAAAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERlqbOLwHqfJtRfg/jN8w99wLTDzEvxziezrncSdued4XAAAAAYN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');
INSERT INTO ledgerheaders VALUES ('1f60b8c6e0bd7b4f76b500ed04d0ce93f8f31853a8c2bdfd46faec8bbac9a2a1', 'a01cbeeaeb18cbed0ee913a13ba4f760660f0a6541be8f19021816601f96a7d8', '6a6ce2f01ea7c9b517e0fe337cc3df702d30f312fc7389eceb9dc49db9e7785c', 7, 1559579832, 'AAAAC6AcvurrGMvtDukToTuk92BmDwplQb6PGQIYFmAflqfYB0IOinnjwk9l13AyWCV0ITfGiJwVvpobnJxxD0agw6oAAAAAXPVMuAAAAAAAAAAA3z9hmASpL9tAVxktxD3XSOp3itxSvEmM6AUkwBS4ERlqbOLwHqfJtRfg/jN8w99wLTDzEvxziezrncSdued4XAAAAAcN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA');


--
-- Data for Name: offers; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: peers; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: publishqueue; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: pubsub; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: quoruminfo; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: scphistory; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 2, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAIAAAACAAAAAQAAAEi5lEev1R1cptMqJyV86PLvZhkUlSQk9VwpPmG2+iubVAAAAABc9UyzAAAAAgAAAAgAAAABAAAACwAAAAgAAAADAA9CQAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABAjprck2DJswhyf49SIq7uzv9rt6lVYn6CxHMzSrEEUPbjelNVA0x7AP27WbMXlslFSwH70fRAAzCIVjslFFBqDA==');
INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 3, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAMAAAACAAAAAQAAADA41DrC4AIfrF+GO5f3Ct4klLJe2wWvURsXe4Am4O0pkAAAAABc9Uy0AAAAAAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABACgGkamC+/IlixaMNuP88/CyreSLBJ2CZ6L7fpu0KPc4vtF1HuKsG5Y5FBb3uT6KLtlCaL+5ARmJna8GiYtxBDw==');
INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 4, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAQAAAACAAAAAQAAADBgnU99+w8TgGBvYwZC1XifhHVGFo+USc2y/a04WoWwOAAAAABc9Uy1AAAAAAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABAHDSNcqVsa3J0bPPcv9Sse1QpDhTFgbncF2a3B9QyQe3w1KmsdvRLwWUmVrX1EkxjTJMVydxj9+o1p4FqbZ6wAg==');
INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 5, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAUAAAACAAAAAQAAADDpD072qifpc+g0PosfwWdLI3bjlHYtzNxxPETy3jz9MAAAAABc9Uy2AAAAAAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABAkYYwCvwPVdU96BJIkmXkRt+Fm9cOrKWKavuPUBG9CC7g4F1/yRTa1P9H95BNqwMkn6d7v/Z+f7v9aZ+JEnrVDA==');
INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 6, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAYAAAACAAAAAQAAADBQ0CGvOOUMN5ssq0Jdhl4OwqgOt+db4CaNx96j0NiJ6wAAAABc9Uy3AAAAAAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABAHKiZd0L7g+0I8TBX57oCPANAo/usLpsHc2XIxtIm+vgU7HHTivJeA1mj8YCkojtt82DIr1gQWIMpYEupWNDRBA==');
INSERT INTO scphistory VALUES ('GD5L6QUAEOOV2Q5UQI64NMNFISONLPXQH4HBAEEFHSVMPKYUSXH3VSGG', 7, 'AAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAAAAAAcAAAACAAAAAQAAADAHQg6KeePCT2XXcDJYJXQhN8aInBW+mhucnHEPRqDDqgAAAABc9Uy4AAAAAAAAAAAAAAABjSHRVUr6VMlEcau3H59kKM/insurCW81XEu/NG3CuHsAAABAzZtPgd3/YNHsLyLpY6woBhg1YNAjPJXNZ01y8D4WC7vCBttYMnpvHRdZ/iZC9NtOkBxN3VaXUNrLaxPOr6NHCw==');


--
-- Data for Name: scpquorums; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO scpquorums VALUES ('8d21d1554afa54c94471abb71f9f6428cfe29ecbab096f355c4bbf346dc2b87b', 7, 'AAAAAQAAAAEAAAAA+r9CgCOdXUO0gj3GsaVEnNW+8D8OEBCFPKrHqxSVz7oAAAAA');


--
-- Data for Name: storestate; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO storestate VALUES ('lastscpdata                     ', 'AAAAAgAAAAD6v0KAI51dQ7SCPcaxpUSc1b7wPw4QEIU8qserFJXPugAAAAAAAAAHAAAAA40h0VVK+lTJRHGrtx+fZCjP4p7LqwlvNVxLvzRtwrh7AAAAAQAAAJgHQg6KeePCT2XXcDJYJXQhN8aInBW+mhucnHEPRqDDqgAAAABc9Uy4AAAAAAAAAAEAAAAA+r9CgCOdXUO0gj3GsaVEnNW+8D8OEBCFPKrHqxSVz7oAAABAEEBMIBpEICnNLCE2wCF4kwS2DjXGrcXONqPTWQnXCGQxr9Bcp+BF0y36nSLhVCFvQ0d84eXpc8YCboGtQMpCAAAAAAEAAACYB0IOinnjwk9l13AyWCV0ITfGiJwVvpobnJxxD0agw6oAAAAAXPVMuAAAAAAAAAABAAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAQBBATCAaRCApzSwhNsAheJMEtg41xq3Fzjaj01kJ1whkMa/QXKfgRdMt+p0i4VQhb0NHfOHl6XPGAm6BrUDKQgAAAABAoX0v8Uqbdi5xvjI7M+SEk/LCHzjOBX71zF9emwTX6d46Z9b0RIG3+gZdp0Pl370F1Q4EaVELh2j9iWa7n4jYBAAAAAD6v0KAI51dQ7SCPcaxpUSc1b7wPw4QEIU8qserFJXPugAAAAAAAAAHAAAAAgAAAAEAAAAwB0IOinnjwk9l13AyWCV0ITfGiJwVvpobnJxxD0agw6oAAAAAXPVMuAAAAAAAAAAAAAAAAY0h0VVK+lTJRHGrtx+fZCjP4p7LqwlvNVxLvzRtwrh7AAAAQM2bT4Hd/2DR7C8i6WOsKAYYNWDQIzyVzWdNcvA+Fgu7wgbbWDJ6bx0XWf4mQvTbTpAcTd1Wl1Day2sTzq+jRwsAAAABoBy+6usYy+0O6ROhO6T3YGYPCmVBvo8ZAhgWYB+Wp9gAAAAAAAAAAQAAAAEAAAABAAAAAPq/QoAjnV1DtII9xrGlRJzVvvA/DhAQhTyqx6sUlc+6AAAAAA==');
INSERT INTO storestate VALUES ('databaseschema                  ', '10');
INSERT INTO storestate VALUES ('networkpassphrase               ', 'Test SDF Network ; September 2015');
INSERT INTO storestate VALUES ('forcescponnextlaunch            ', 'false');
INSERT INTO storestate VALUES ('ledgerupgrades                  ', '{
    "time": 0,
    "version": {
        "has": false
    },
    "fee": {
        "has": false
    },
    "maxtxsize": {
        "has": false
    },
    "reserve": {
        "has": false
    }
}');
INSERT INTO storestate VALUES ('lastclosedledger                ', '1f60b8c6e0bd7b4f76b500ed04d0ce93f8f31853a8c2bdfd46faec8bbac9a2a1');
INSERT INTO storestate VALUES ('historyarchivestate             ', '{
    "version": 1,
    "server": "v11.1.0",
    "currentLedger": 7,
    "currentBuckets": [
        {
            "curr": "056cd09fa286b2de1ff329b224259b858f94ec1dcc93f0e0512b2f6819348951",
            "next": {
                "state": 0
            },
            "snap": "056cd09fa286b2de1ff329b224259b858f94ec1dcc93f0e0512b2f6819348951"
        },
        {
            "curr": "e3bc65ded24f6c7ee979f9350aaecf9856f52e75490a401f1a66bb3d77d767dd",
            "next": {
                "state": 1,
                "output": "056cd09fa286b2de1ff329b224259b858f94ec1dcc93f0e0512b2f6819348951"
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        },
        {
            "curr": "0000000000000000000000000000000000000000000000000000000000000000",
            "next": {
                "state": 0
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        }
    ]
}');


--
-- Data for Name: trustlines; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: txfeehistory; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: txhistory; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: upgradehistory; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO upgradehistory VALUES (2, 1, 'AAAAAQAAAAs=', 'AAAAAA==');
INSERT INTO upgradehistory VALUES (2, 2, 'AAAAAwAPQkA=', 'AAAAAA==');


--
-- Name: accountdata accountdata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY accountdata
    ADD CONSTRAINT accountdata_pkey PRIMARY KEY (accountid, dataname);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (accountid);


--
-- Name: ban ban_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY ban
    ADD CONSTRAINT ban_pkey PRIMARY KEY (nodeid);


--
-- Name: ledgerheaders ledgerheaders_ledgerseq_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY ledgerheaders
    ADD CONSTRAINT ledgerheaders_ledgerseq_key UNIQUE (ledgerseq);


--
-- Name: ledgerheaders ledgerheaders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY ledgerheaders
    ADD CONSTRAINT ledgerheaders_pkey PRIMARY KEY (ledgerhash);


--
-- Name: offers offers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY offers
    ADD CONSTRAINT offers_pkey PRIMARY KEY (offerid);


--
-- Name: peers peers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY peers
    ADD CONSTRAINT peers_pkey PRIMARY KEY (ip, port);


--
-- Name: publishqueue publishqueue_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY publishqueue
    ADD CONSTRAINT publishqueue_pkey PRIMARY KEY (ledger);


--
-- Name: pubsub pubsub_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY pubsub
    ADD CONSTRAINT pubsub_pkey PRIMARY KEY (resid);


--
-- Name: quoruminfo quoruminfo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY quoruminfo
    ADD CONSTRAINT quoruminfo_pkey PRIMARY KEY (nodeid);


--
-- Name: scpquorums scpquorums_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY scpquorums
    ADD CONSTRAINT scpquorums_pkey PRIMARY KEY (qsethash);


--
-- Name: storestate storestate_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY storestate
    ADD CONSTRAINT storestate_pkey PRIMARY KEY (statename);


--
-- Name: trustlines trustlines_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY trustlines
    ADD CONSTRAINT trustlines_pkey PRIMARY KEY (accountid, issuer, assetcode);


--
-- Name: txfeehistory txfeehistory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY txfeehistory
    ADD CONSTRAINT txfeehistory_pkey PRIMARY KEY (ledgerseq, txindex);


--
-- Name: txhistory txhistory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY txhistory
    ADD CONSTRAINT txhistory_pkey PRIMARY KEY (ledgerseq, txindex);


--
-- Name: upgradehistory upgradehistory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY upgradehistory
    ADD CONSTRAINT upgradehistory_pkey PRIMARY KEY (ledgerseq, upgradeindex);


--
-- Name: accountbalances; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accountbalances ON accounts USING btree (balance) WHERE (balance >= 1000000000);


--
-- Name: bestofferindex; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX bestofferindex ON offers USING btree (sellingasset, buyingasset, price);


--
-- Name: histbyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX histbyseq ON txhistory USING btree (ledgerseq);


--
-- Name: histfeebyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX histfeebyseq ON txfeehistory USING btree (ledgerseq);


--
-- Name: ledgersbyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX ledgersbyseq ON ledgerheaders USING btree (ledgerseq);


--
-- Name: scpenvsbyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX scpenvsbyseq ON scphistory USING btree (ledgerseq);


--
-- Name: scpquorumsbyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX scpquorumsbyseq ON scpquorums USING btree (lastledgerseq);


--
-- Name: upgradehistbyseq; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX upgradehistbyseq ON upgradehistory USING btree (ledgerseq);


--
-- PostgreSQL database dump complete
--

