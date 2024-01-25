-- +migrate Up

CREATE TABLE identities_audit (
  audit_id BIGINT NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  audit_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  audit_user TEXT NOT NULL DEFAULT USER,
  audit_op audit_op NOT NULL,
  LIKE identities
);

-- +migrate StatementBegin
CREATE FUNCTION record_identities_audit() RETURNS TRIGGER AS $BODY$
  BEGIN
    IF (TG_OP = 'INSERT') THEN
      INSERT INTO identities_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
      INSERT INTO identities_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, NEW.*);
      RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
      INSERT INTO identities_audit VALUES (DEFAULT, DEFAULT, DEFAULT, TG_OP::audit_op, OLD.*);
      RETURN OLD;
    END IF;
  END;
$BODY$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER record_identities_audit
AFTER INSERT OR UPDATE OR DELETE ON identities
  FOR EACH ROW EXECUTE PROCEDURE record_identities_audit();

-- +migrate Down

DROP TRIGGER record_identities_audit ON identities;
DROP FUNCTION record_identities_audit;
DROP TABLE identities_audit;
