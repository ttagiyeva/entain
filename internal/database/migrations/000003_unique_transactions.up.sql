BEGIN;

    ALTER TABLE transactions ADD CONSTRAINT unique_transaction_id UNIQUE (transaction_id);

COMMIT;    