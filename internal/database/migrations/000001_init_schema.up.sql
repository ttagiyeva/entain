BEGIN;

	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TYPE source_type AS ENUM ('game', 'server', 'payment');
	CREATE TYPE state AS ENUM ('win', 'lost');

	CREATE TABLE IF NOT EXISTS
		users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			balance NUMERIC(18,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

	CREATE TABLE IF NOT EXISTS
		transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users (id),
			transaction_id VARCHAR(100) NOT NULL,
			source_type source_type NOT NULL,
			state state NOT NULL,
			amount NUMERIC(18,2) NOT NULL CHECK (amount >= 0),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			cancelled BOOLEAN NOT NULL DEFAULT FALSE,
			cancelled_at TIMESTAMP WITH TIME ZONE 
		);
		
COMMIT;