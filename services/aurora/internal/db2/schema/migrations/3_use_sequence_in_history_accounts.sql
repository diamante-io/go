-- +migrate Up
CREATE SEQUENCE history_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
SELECT setval('history_accounts_id_seq', (SELECT MAX(id) FROM history_accounts));
ALTER TABLE ONLY history_accounts ALTER COLUMN id SET DEFAULT nextval('history_accounts_id_seq'::regclass);

-- +migrate Down
ALTER TABLE ONLY history_accounts ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE history_accounts_id_seq;
