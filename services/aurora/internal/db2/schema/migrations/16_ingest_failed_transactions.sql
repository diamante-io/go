-- +migrate Up

-- Check db2/history.Transaction.Successful field comment for more information.
ALTER TABLE history_transactions ADD successful boolean;

-- +migrate Down

-- Remove failed transactions and operations from failed transactions!
DELETE FROM history_operations USING history_transactions
  WHERE history_transactions.id = history_operations.transaction_id AND successful = false;
DELETE FROM history_transactions WHERE successful = false;
ALTER TABLE history_transactions DROP COLUMN successful;
