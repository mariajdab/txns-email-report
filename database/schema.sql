CREATE TABLE IF NOT EXISTS register_transactions
(
    id serial PRIMARY KEY,
    date_at TIMESTAMP,
    txn NUMERIC
);