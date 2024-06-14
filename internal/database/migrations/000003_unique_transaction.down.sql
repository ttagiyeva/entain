BEGIN;

    ALTER TABLE transactions DROP CONSTRAINT unique_transaction_id;

COMMIT;   