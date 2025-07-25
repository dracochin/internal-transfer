CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE accounts (
  account_id BIGINT PRIMARY KEY,
  balance NUMERIC(20, 8) NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  source_account_id BIGINT NOT NULL,
  destination_account_id BIGINT NOT NULL,
  amount NUMERIC(20, 8) NOT NULL CHECK (amount > 0),
  idempotency_key TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE entries (
  id SERIAL PRIMARY KEY,
  transaction_id UUID NOT NULL,
  account_id BIGINT NOT NULL,
  amount NUMERIC(20, 8) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (transaction_id) REFERENCES transactions(id),
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);
