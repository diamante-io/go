-- +migrate Up

CREATE TABLE auth_methods_audit (
  audit_id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  audit_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  audit_user TEXT NOT NULL DEFAULT USER,
  audit_op audit_op NOT NULL,
  LIKE auth_methods
);

-- +migrate StatementBegin
CREATE FUNCTION record_auth_methods_audit() RETURNS TRIGGER AS $BODY$
  BEGIN
    IF (TG_OP = 'INSERT') THEN
      INSERT INTO auth_methods_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
      INSERT INTO auth_methods_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
      INSERT INTO auth_methods_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, OLD.*);
      RETURN OLD;
    END IF;
  END;
$BODY$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER record_auth_methods_audit
AFTER INSERT OR UPDATE OR DELETE ON auth_methods
  FOR EACH ROW EXECUTE PROCEDURE record_auth_methods_audit();

-- +migrate Down

DROP TRIGGER record_auth_methods_audit ON auth_methods;
DROP FUNCTION record_auth_methods_audit;
DROP TABLE auth_methods_audit;
