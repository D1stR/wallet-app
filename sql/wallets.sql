CREATE TABLE IF NOT EXISTS wallets
(
    id      UUID PRIMARY KEY,
    balance FLOAT DEFAULT 0 CHECK (balance >= 0)
);
