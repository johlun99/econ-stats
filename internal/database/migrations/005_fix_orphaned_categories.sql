-- Null out category_id for transactions referencing deleted categories
UPDATE transactions SET category_id = NULL
WHERE category_id IS NOT NULL
  AND category_id NOT IN (SELECT id FROM categories);
