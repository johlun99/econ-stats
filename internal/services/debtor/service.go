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

		if len(keys) > 0 {
			balance, err := s.computeBalance(keys)
			if err != nil {
				return nil, err
			}
			debtors[i].Balance = balance
		}
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
	if len(keys) == 0 {
		return []models.Transaction{}, nil
	}

	placeholders := make([]string, len(keys))
	args := make([]any, len(keys))
	for i, k := range keys {
		placeholders[i] = "?"
		args[i] = k
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
		       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.merchant_key IN (%s)
		ORDER BY t.booking_date DESC
	`, strings.Join(placeholders, ","))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var isTransfer int
		var catID sql.NullInt64
		var catName, catColor, catIcon sql.NullString
		if err := rows.Scan(&t.ID, &t.BookingDate, &t.TransactionDate, &t.Description,
			&t.Amount, &t.Balance, &catID, &catName, &catColor,
			&catIcon, &t.MerchantKey, &isTransfer); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
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

func (s *Service) computeBalance(keys []string) (float64, error) {
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

	var balance float64
	if err := s.db.QueryRow(query, args...).Scan(&balance); err != nil {
		return 0, fmt.Errorf("compute balance: %w", err)
	}
	return balance, nil
}
