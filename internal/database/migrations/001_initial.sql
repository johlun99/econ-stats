-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    color TEXT NOT NULL DEFAULT '#6B7280',
    icon TEXT NOT NULL DEFAULT '📦',
    is_income INTEGER NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
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
    UNIQUE(booking_date, transaction_date, description, amount)
);

CREATE INDEX IF NOT EXISTS idx_transactions_merchant_key ON transactions(merchant_key);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_booking_date ON transactions(booking_date);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_date ON transactions(transaction_date);

-- Category rules: merchant_key -> category mapping
CREATE TABLE IF NOT EXISTS category_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    merchant_key TEXT NOT NULL UNIQUE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed default categories
INSERT OR IGNORE INTO categories (name, color, icon, is_income, sort_order) VALUES
    ('Boende', '#EF4444', '🏠', 0, 1),
    ('Mat', '#F97316', '🛒', 0, 2),
    ('Restaurang', '#F59E0B', '🍽️', 0, 3),
    ('Transport', '#3B82F6', '🚗', 0, 4),
    ('Nöje', '#8B5CF6', '🎮', 0, 5),
    ('Hälsa', '#10B981', '💊', 0, 6),
    ('Prenumerationer', '#EC4899', '📱', 0, 7),
    ('Shopping', '#6366F1', '🛍️', 0, 8),
    ('Sparande', '#14B8A6', '💰', 0, 9),
    ('Försäkring', '#64748B', '🛡️', 0, 10),
    ('Övrigt', '#6B7280', '📦', 0, 11),
    ('Lön', '#22C55E', '💵', 1, 12),
    ('Övrig inkomst', '#84CC16', '💸', 1, 13);
