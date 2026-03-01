package transaction

import (
	"database/sql"
	"fmt"
	"strings"

	"econ-stats/internal/models"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Import(transactions []models.Transaction) (*models.ImportResult, error) {
	result := &models.ImportResult{
		TotalRows: len(transactions),
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	existsStmt, err := tx.Prepare(`SELECT 1 FROM transactions WHERE transaction_date = ? AND description = ? AND amount = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare exists: %w", err)
	}
	defer existsStmt.Close()

	upsertStmt, err := tx.Prepare(`
		INSERT INTO transactions
			(booking_date, transaction_date, description, amount, balance, merchant_key, is_transfer)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(transaction_date, description, amount) DO UPDATE SET
			booking_date = excluded.booking_date,
			balance = excluded.balance,
			is_transfer = excluded.is_transfer
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare upsert: %w", err)
	}
	defer upsertStmt.Close()

	for _, t := range transactions {
		isTransfer := 0
		if t.IsTransfer {
			isTransfer = 1
		}

		var exists int
		err := existsStmt.QueryRow(t.TransactionDate, t.Description, t.Amount).Scan(&exists)
		alreadyExists := err == nil

		if _, err := upsertStmt.Exec(t.BookingDate, t.TransactionDate, t.Description, t.Amount, t.Balance, t.MerchantKey, isTransfer); err != nil {
			return nil, fmt.Errorf("upsert: %w", err)
		}

		if alreadyExists {
			result.Updated++
		} else {
			result.NewTransactions++
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return result, nil
}

func (s *Service) GetAll(month string) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
		       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
	`
	args := []interface{}{}

	if month != "" {
		query += " WHERE t.transaction_date LIKE ?"
		args = append(args, month+"%")
	}

	query += " ORDER BY t.transaction_date DESC, t.id DESC"

	return s.queryTransactions(query, args...)
}

func (s *Service) Search(term string, month string) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
		       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.description LIKE ?
	`
	args := []interface{}{"%" + term + "%"}

	if month != "" {
		query += " AND t.transaction_date LIKE ?"
		args = append(args, month+"%")
	}

	query += " ORDER BY t.transaction_date DESC, t.id DESC"

	return s.queryTransactions(query, args...)
}

func (s *Service) GetUncategorizedMerchants() ([]models.MerchantGroup, error) {
	rows, err := s.db.Query(`
		SELECT merchant_key, COUNT(*) as cnt, SUM(amount) as total,
		       COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) as income_total,
		       COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) as expense_total,
		       MIN(transaction_date) as first_date, MAX(transaction_date) as last_date
		FROM transactions
		WHERE category_id IS NULL AND is_transfer = 0
		GROUP BY merchant_key
		ORDER BY cnt DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.MerchantGroup
	for rows.Next() {
		var g models.MerchantGroup
		if err := rows.Scan(&g.MerchantKey, &g.Count, &g.TotalAmount, &g.IncomeTotal, &g.ExpenseTotal, &g.FirstDate, &g.LastDate); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (s *Service) GetMerchantTransactions(merchantKey string) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
		       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.merchant_key = ?
		ORDER BY t.transaction_date DESC
	`
	return s.queryTransactions(query, merchantKey)
}

func (s *Service) UpdateCategory(id int64, categoryID *int64) error {
	_, err := s.db.Exec("UPDATE transactions SET category_id = ? WHERE id = ?", categoryID, id)
	return err
}

func (s *Service) SetCategoryByMerchant(merchantKey string, categoryID int64) (int64, error) {
	res, err := s.db.Exec("UPDATE transactions SET category_id = ? WHERE merchant_key = ?", categoryID, merchantKey)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Service) GetUncategorizedCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE category_id IS NULL AND is_transfer = 0").Scan(&count)
	return count, err
}

func (s *Service) AutoCategorize(rules map[string]int64) (int, error) {
	if len(rules) == 0 {
		return 0, nil
	}

	total := 0
	for merchantKey, categoryID := range rules {
		res, err := s.db.Exec(
			"UPDATE transactions SET category_id = ? WHERE merchant_key = ? AND category_id IS NULL",
			categoryID, merchantKey,
		)
		if err != nil {
			return total, err
		}
		n, _ := res.RowsAffected()
		total += int(n)
	}
	return total, nil
}

func (s *Service) DeleteAll() error {
	_, err := s.db.Exec("DELETE FROM transactions")
	return err
}

func (s *Service) queryTransactions(query string, args ...interface{}) ([]models.Transaction, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var isTransfer int
		var catName, catColor, catIcon sql.NullString
		var catID sql.NullInt64
		if err := rows.Scan(
			&t.ID, &t.BookingDate, &t.TransactionDate, &t.Description, &t.Amount, &t.Balance,
			&catID, &catName, &catColor, &catIcon, &t.MerchantKey, &isTransfer,
		); err != nil {
			return nil, err
		}
		t.IsTransfer = isTransfer == 1
		if catID.Valid {
			t.CategoryID = &catID.Int64
			t.CategoryName = &catName.String
			t.CategoryColor = &catColor.String
			t.CategoryIcon = &catIcon.String
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}

// UncategorizedCount returns count of uncategorized non-transfer transactions for a search
func (s *Service) CountByFilter(month string, categoryFilter string) (int, error) {
	query := "SELECT COUNT(*) FROM transactions WHERE 1=1"
	args := []interface{}{}

	if month != "" {
		query += " AND transaction_date LIKE ?"
		args = append(args, month+"%")
	}

	switch strings.ToLower(categoryFilter) {
	case "uncategorized":
		query += " AND category_id IS NULL AND is_transfer = 0"
	case "transfer":
		query += " AND is_transfer = 1"
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count, err
}
