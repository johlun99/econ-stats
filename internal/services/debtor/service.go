package debtor

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

func (s *Service) GetAll() ([]models.DebtorDetail, error) {
	rows, err := s.db.Query("SELECT id, name, icon, color FROM debtors ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("query debtors: %w", err)
	}
	defer rows.Close()

	var debtors []models.DebtorDetail
	for rows.Next() {
		var d models.DebtorDetail
		if err := rows.Scan(&d.ID, &d.Name, &d.Icon, &d.Color); err != nil {
			return nil, fmt.Errorf("scan debtor: %w", err)
		}
		debtors = append(debtors, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range debtors {
		keys, err := s.getMerchantKeys(debtors[i].ID)
		if err != nil {
			return nil, err
		}
		debtors[i].MerchantKeys = keys

		balance, err := s.computeBalance(debtors[i].ID, keys)
		if err != nil {
			return nil, err
		}
		debtors[i].Balance = balance
	}

	return debtors, nil
}

func (s *Service) Create(name, icon, color string) (*models.Debtor, error) {
	res, err := s.db.Exec(
		"INSERT INTO debtors (name, icon, color) VALUES (?, ?, ?)",
		name, icon, color,
	)
	if err != nil {
		return nil, fmt.Errorf("insert debtor: %w", err)
	}

	id, _ := res.LastInsertId()
	return &models.Debtor{ID: id, Name: name, Icon: icon, Color: color}, nil
}

func (s *Service) Update(id int64, name, icon, color string) error {
	_, err := s.db.Exec(
		"UPDATE debtors SET name=?, icon=?, color=? WHERE id=?",
		name, icon, color, id,
	)
	return err
}

func (s *Service) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM debtors WHERE id=?", id)
	return err
}

func (s *Service) AddMerchantKey(debtorID int64, merchantKey string) error {
	_, err := s.db.Exec(
		"INSERT INTO debtor_merchant_keys (debtor_id, merchant_key) VALUES (?, ?)",
		debtorID, merchantKey,
	)
	if err != nil {
		return fmt.Errorf("add merchant key: %w", err)
	}
	return nil
}

func (s *Service) RemoveMerchantKey(debtorID int64, merchantKey string) error {
	_, err := s.db.Exec(
		"DELETE FROM debtor_merchant_keys WHERE debtor_id=? AND merchant_key=?",
		debtorID, merchantKey,
	)
	return err
}

func (s *Service) GetTransactions(debtorID int64) ([]models.Transaction, error) {
	keys, err := s.getMerchantKeys(debtorID)
	if err != nil {
		return nil, err
	}

	var parts []string
	var args []any

	if len(keys) > 0 {
		placeholders := make([]string, len(keys))
		for i, k := range keys {
			placeholders[i] = "?"
			args = append(args, k)
		}
		parts = append(parts, fmt.Sprintf(`
			SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
			       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer, 0 AS is_manual
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			WHERE t.merchant_key IN (%s)
		`, strings.Join(placeholders, ",")))
	}

	parts = append(parts, `
		SELECT dt.id, dt.transaction_date, dt.transaction_date, dt.description, dt.amount, 0,
		       NULL, NULL, NULL, NULL, '', 0, 1 AS is_manual
		FROM debtor_transactions dt
		WHERE dt.debtor_id = ?
	`)
	args = append(args, debtorID)

	query := strings.Join(parts, " UNION ALL ") + " ORDER BY transaction_date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var isTransfer, isManual int
		var catID sql.NullInt64
		var catName, catColor, catIcon sql.NullString
		if err := rows.Scan(&t.ID, &t.BookingDate, &t.TransactionDate, &t.Description,
			&t.Amount, &t.Balance, &catID, &catName, &catColor,
			&catIcon, &t.MerchantKey, &isTransfer, &isManual); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		t.IsTransfer = isTransfer == 1
		t.IsManual = isManual == 1
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

func (s *Service) AddManualTransaction(debtorID int64, description string, amount float64, date string) error {
	_, err := s.db.Exec(
		"INSERT INTO debtor_transactions (debtor_id, description, amount, transaction_date) VALUES (?, ?, ?, ?)",
		debtorID, description, amount, date,
	)
	if err != nil {
		return fmt.Errorf("add manual transaction: %w", err)
	}
	return nil
}

func (s *Service) UpdateManualTransaction(id int64, description string, amount float64, date string) error {
	_, err := s.db.Exec(
		"UPDATE debtor_transactions SET description=?, amount=?, transaction_date=? WHERE id=?",
		description, amount, date, id,
	)
	if err != nil {
		return fmt.Errorf("update manual transaction: %w", err)
	}
	return nil
}

func (s *Service) DeleteManualTransaction(id int64) error {
	_, err := s.db.Exec("DELETE FROM debtor_transactions WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("delete manual transaction: %w", err)
	}
	return nil
}

func (s *Service) GetAllMerchantKeys() ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT merchant_key FROM transactions ORDER BY merchant_key")
	if err != nil {
		return nil, fmt.Errorf("query all merchant keys: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (s *Service) getMerchantKeys(debtorID int64) ([]string, error) {
	rows, err := s.db.Query(
		"SELECT merchant_key FROM debtor_merchant_keys WHERE debtor_id=? ORDER BY merchant_key",
		debtorID,
	)
	if err != nil {
		return nil, fmt.Errorf("query merchant keys: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (s *Service) computeBalance(debtorID int64, keys []string) (float64, error) {
	var balance float64

	if len(keys) > 0 {
		placeholders := make([]string, len(keys))
		args := make([]any, len(keys))
		for i, k := range keys {
			placeholders[i] = "?"
			args[i] = k
		}

		query := fmt.Sprintf(
			"SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE merchant_key IN (%s)",
			strings.Join(placeholders, ","),
		)

		if err := s.db.QueryRow(query, args...).Scan(&balance); err != nil {
			return 0, fmt.Errorf("compute balance: %w", err)
		}
	}

	var manualBalance float64
	if err := s.db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM debtor_transactions WHERE debtor_id=?",
		debtorID,
	).Scan(&manualBalance); err != nil {
		return 0, fmt.Errorf("compute manual balance: %w", err)
	}

	return balance + manualBalance, nil
}
