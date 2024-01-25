-- +migrate Up

CREATE TYPE audit_op AS ENUM (
  'INSERT',
  'UPDATE',
  'DELETE'
);

CREATE TABLE accounts_audit (
  audit_id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  audit_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  audit_user TEXT NOT NULL DEFAULT USER,
  audit_op audit_op NOT NULL,
  LIKE accounts
);

-- +migrate StatementBegin
CREATE FUNCTION record_accounts_audit() RETURNS TRIGGER AS $BODY$
  BEGIN
    IF (TG_OP = 'INSERT') THEN
      INSERT INTO accounts_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
      INSERT INTO accounts_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
      INSERT INTO accounts_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, OLD.*);
      RETURN OLD;
    END IF;
  END;
$BODY$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER record_accounts_audit
AFTER INSERT OR UPDATE OR DELETE ON accounts
  FOR EACH ROW EXECUTE PROCEDURE record_accounts_audit();

-- +migrate Down

DROP TRIGGER record_accounts_audit ON accounts;
DROP FUNCTION record_accounts_audit;
DROP TABLE accounts_audit;
DROP TYPE audit_op;
