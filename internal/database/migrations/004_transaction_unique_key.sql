CREATE TABLE transactions_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    booking_date TEXT NOT NULL,
    transaction_date TEXT NOT NULL,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    balance REAL NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    merchant_key TEXT NOT NULL,
    is_transfer INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(transaction_date, description, amount)
);

INSERT INTO transactions_new (id, booking_date, transaction_date, description, amount, balance, category_id, merchant_key, is_transfer, created_at)
    SELECT id, booking_date, transaction_date, description, amount, balance, category_id, merchant_key, is_transfer, created_at
    FROM transactions
    GROUP BY transaction_date, description, amount;

DROP TABLE transactions;
ALTER TABLE transactions_new RENAME TO transactions;

CREATE INDEX IF NOT EXISTS idx_transactions_merchant_key ON transactions(merchant_key);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_booking_date ON transactions(booking_date);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_date ON transactions(transaction_date);
