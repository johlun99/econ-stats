-- Add is_expense column (default 1 = true for all existing categories)
ALTER TABLE categories ADD COLUMN is_expense INTEGER NOT NULL DEFAULT 1;

-- Income-only categories: Lön, Övrig inkomst
UPDATE categories SET is_expense = 0 WHERE name IN ('Lön', 'Övrig inkomst');
