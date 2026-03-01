package category

import (
	"database/sql"
	"fmt"

	"econ-stats/internal/models"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll() ([]models.Category, error) {
	rows, err := s.db.Query("SELECT id, name, color, icon, is_income, is_expense, sort_order FROM categories ORDER BY sort_order, name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var isIncome, isExpense int
		if err := rows.Scan(&c.ID, &c.Name, &c.Color, &c.Icon, &isIncome, &isExpense, &c.SortOrder); err != nil {
			return nil, err
		}
		c.IsIncome = isIncome == 1
		c.IsExpense = isExpense == 1
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (s *Service) Create(name, color, icon string, isIncome, isExpense bool) (*models.Category, error) {
	incomeVal := boolToInt(isIncome)
	expenseVal := boolToInt(isExpense)

	// Get max sort order
	var maxOrder int
	s.db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM categories").Scan(&maxOrder)

	res, err := s.db.Exec(
		"INSERT INTO categories (name, color, icon, is_income, is_expense, sort_order) VALUES (?, ?, ?, ?, ?, ?)",
		name, color, icon, incomeVal, expenseVal, maxOrder+1,
	)
	if err != nil {
		return nil, fmt.Errorf("insert category: %w", err)
	}

	id, _ := res.LastInsertId()
	return &models.Category{
		ID: id, Name: name, Color: color, Icon: icon, IsIncome: isIncome, IsExpense: isExpense, SortOrder: maxOrder + 1,
	}, nil
}

func (s *Service) Update(id int64, name, color, icon string, isIncome, isExpense bool) error {
	_, err := s.db.Exec(
		"UPDATE categories SET name=?, color=?, icon=?, is_income=?, is_expense=? WHERE id=?",
		name, color, icon, boolToInt(isIncome), boolToInt(isExpense), id,
	)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (s *Service) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM categories WHERE id=?", id)
	return err
}

// Rules

func (s *Service) GetRules() ([]models.CategoryRule, error) {
	rows, err := s.db.Query(`
		SELECT cr.id, cr.merchant_key, cr.category_id, c.name
		FROM category_rules cr
		JOIN categories c ON cr.category_id = c.id
		ORDER BY cr.merchant_key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.CategoryRule
	for rows.Next() {
		var r models.CategoryRule
		if err := rows.Scan(&r.ID, &r.MerchantKey, &r.CategoryID, &r.CategoryName); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, rows.Err()
}

func (s *Service) CreateRule(merchantKey string, categoryID int64) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO category_rules (merchant_key, category_id) VALUES (?, ?)",
		merchantKey, categoryID,
	)
	return err
}

func (s *Service) DeleteRule(id int64) error {
	_, err := s.db.Exec("DELETE FROM category_rules WHERE id=?", id)
	return err
}

func (s *Service) GetRulesMap() (map[string]int64, error) {
	rows, err := s.db.Query("SELECT merchant_key, category_id FROM category_rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make(map[string]int64)
	for rows.Next() {
		var key string
		var catID int64
		if err := rows.Scan(&key, &catID); err != nil {
			return nil, err
		}
		rules[key] = catID
	}
	return rules, rows.Err()
}
