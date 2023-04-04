BEGIN;

    DROP TABLE IF EXISTS transactions;
    DROP TABLE IF EXISTS users;
    DROP TYPE IF EXISTS source_type;
    DROP TYPE IF EXISTS state;
    DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;