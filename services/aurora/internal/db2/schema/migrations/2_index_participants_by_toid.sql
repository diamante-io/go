-- +migrate Up

CREATE INDEX hop_by_hoid ON history_operation_participants USING btree (history_operation_id);
CREATE INDEX htp_by_htid ON history_transaction_participants USING btree (history_transaction_id);

-- +migrate Down

DROP INDEX hop_by_hoid;
DROP INDEX htp_by_htid;
