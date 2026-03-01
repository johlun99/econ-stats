package stats

import (
	"database/sql"
	"fmt"
	"time"

	"econ-stats/internal/models"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetMonthlyStats(month string) (*models.MonthlyStats, error) {
	stats := &models.MonthlyStats{Month: month}

	// Total expenses (negative amounts + returns in expense categories, excluding transfers)
	s.db.QueryRow(`
		SELECT COALESCE(SUM(-t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
	`, month+"%").Scan(&stats.TotalExpenses)

	// Total income (positive amounts, excluding transfers and returns in expense categories)
	s.db.QueryRow(`
		SELECT COALESCE(SUM(t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.amount > 0 AND t.is_transfer = 0
		  AND (c.is_expense = 0 OR c.is_expense IS NULL)
	`, month+"%").Scan(&stats.TotalIncome)

	stats.NetSavings = stats.TotalIncome - stats.TotalExpenses

	if stats.TotalIncome > 0 {
		stats.SavingsRate = stats.NetSavings / stats.TotalIncome * 100
	}

	// Average daily spend
	t, err := time.Parse("2006-01", month)
	if err == nil {
		daysInMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		// If current month, use days elapsed so far
		now := time.Now()
		if t.Year() == now.Year() && t.Month() == now.Month() {
			daysInMonth = now.Day()
		}
		if daysInMonth > 0 {
			stats.AvgDailySpend = stats.TotalExpenses / float64(daysInMonth)
		}
	}

	// Month-over-month change
	prevMonth := prevMonthStr(month)
	var prevExpenses float64
	s.db.QueryRow(`
		SELECT COALESCE(SUM(-t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
	`, prevMonth+"%").Scan(&prevExpenses)

	if prevExpenses > 0 {
		stats.MonthOverMonth = (stats.TotalExpenses - prevExpenses) / prevExpenses * 100
	}

	// Category breakdown
	stats.CategoryBreakdown, _ = s.getCategoryBreakdown(month)

	// Top merchants
	stats.TopMerchants, _ = s.getTopMerchants(month)

	// Largest expenses
	stats.LargestExpenses, _ = s.getLargestExpenses(month)

	// Daily spending
	stats.DailySpending, _ = s.getDailySpending(month)

	return stats, nil
}

func (s *Service) getCategoryBreakdown(month string) ([]models.CategorySpend, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(t.category_id, 0), COALESCE(c.name, 'Okategoriserad'),
		       COALESCE(c.color, '#9CA3AF'), COALESCE(c.icon, '❓'),
		       SUM(-t.amount), COUNT(*)
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
		GROUP BY t.category_id
		HAVING SUM(-t.amount) > 0
		ORDER BY SUM(-t.amount) DESC
	`, month+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CategorySpend
	for rows.Next() {
		var cs models.CategorySpend
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.CategoryColor, &cs.CategoryIcon, &cs.Total, &cs.Count); err != nil {
			return nil, err
		}
		result = append(result, cs)
	}
	return result, rows.Err()
}

func (s *Service) getTopMerchants(month string) ([]models.MerchantSpend, error) {
	rows, err := s.db.Query(`
		SELECT t.merchant_key, SUM(-t.amount), COUNT(*)
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
		GROUP BY t.merchant_key
		HAVING SUM(-t.amount) > 0
		ORDER BY SUM(-t.amount) DESC
		LIMIT 10
	`, month+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MerchantSpend
	for rows.Next() {
		var ms models.MerchantSpend
		if err := rows.Scan(&ms.MerchantKey, &ms.Total, &ms.Count); err != nil {
			return nil, err
		}
		result = append(result, ms)
	}
	return result, rows.Err()
}

func (s *Service) getLargestExpenses(month string) ([]models.Transaction, error) {
	rows, err := s.db.Query(`
		SELECT t.id, t.booking_date, t.transaction_date, t.description, t.amount, t.balance,
		       t.category_id, c.name, c.color, c.icon, t.merchant_key, t.is_transfer
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.amount < 0 AND t.is_transfer = 0
		ORDER BY t.amount ASC
		LIMIT 5
	`, month+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Transaction
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
		result = append(result, t)
	}
	return result, rows.Err()
}

func (s *Service) getDailySpending(month string) ([]models.DailySpend, error) {
	rows, err := s.db.Query(`
		SELECT t.transaction_date, SUM(-t.amount)
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
		GROUP BY t.transaction_date
		HAVING SUM(-t.amount) > 0
		ORDER BY t.transaction_date
	`, month+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.DailySpend
	for rows.Next() {
		var ds models.DailySpend
		if err := rows.Scan(&ds.Date, &ds.Total); err != nil {
			return nil, err
		}
		result = append(result, ds)
	}
	return result, rows.Err()
}

func (s *Service) GetSpendingTrend(months int) ([]models.SpendingTrend, error) {
	query := `
		SELECT substr(t.transaction_date, 1, 7) as month,
		       COALESCE(SUM(CASE WHEN t.is_transfer = 0 AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0)) THEN -t.amount ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN t.amount > 0 AND t.is_transfer = 0 AND (c.is_expense = 0 OR c.is_expense IS NULL) THEN t.amount ELSE 0 END), 0)
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		GROUP BY month
		ORDER BY month DESC
	`
	var args []interface{}
	if months > 0 {
		query += " LIMIT ?"
		args = append(args, months)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.SpendingTrend
	for rows.Next() {
		var st models.SpendingTrend
		if err := rows.Scan(&st.Month, &st.Expenses, &st.Income); err != nil {
			return nil, err
		}
		result = append(result, st)
	}

	// Reverse to chronological order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, rows.Err()
}

func (s *Service) GetAvailableMonths() ([]models.AvailableMonth, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT substr(transaction_date, 1, 7) as month
		FROM transactions
		ORDER BY month DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthNames := map[string]string{
		"01": "Januari", "02": "Februari", "03": "Mars", "04": "April",
		"05": "Maj", "06": "Juni", "07": "Juli", "08": "Augusti",
		"09": "September", "10": "Oktober", "11": "November", "12": "December",
	}

	var result []models.AvailableMonth
	for rows.Next() {
		var am models.AvailableMonth
		if err := rows.Scan(&am.Month); err != nil {
			return nil, err
		}
		if len(am.Month) >= 7 {
			am.Label = fmt.Sprintf("%s %s", monthNames[am.Month[5:7]], am.Month[:4])
		}
		result = append(result, am)
	}
	return result, rows.Err()
}

func (s *Service) GetYearlyStats(year string) (*models.YearlyStats, error) {
	stats := &models.YearlyStats{Year: year}

	// Total expenses (negative amounts + returns in expense categories, excluding transfers)
	s.db.QueryRow(`
		SELECT COALESCE(SUM(-t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
	`, year+"%").Scan(&stats.TotalExpenses)

	// Total income (positive amounts, excluding transfers and returns in expense categories)
	s.db.QueryRow(`
		SELECT COALESCE(SUM(t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.amount > 0 AND t.is_transfer = 0
		  AND (c.is_expense = 0 OR c.is_expense IS NULL)
	`, year+"%").Scan(&stats.TotalIncome)

	stats.NetSavings = stats.TotalIncome - stats.TotalExpenses

	if stats.TotalIncome > 0 {
		stats.SavingsRate = stats.NetSavings / stats.TotalIncome * 100
	}

	// Average monthly spend
	yearInt := 0
	fmt.Sscanf(year, "%d", &yearInt)
	now := time.Now()
	monthsElapsed := 12
	if yearInt == now.Year() {
		monthsElapsed = int(now.Month())
	}
	if monthsElapsed > 0 {
		stats.AvgMonthlySpend = stats.TotalExpenses / float64(monthsElapsed)
	}

	// Year-over-year change
	prevYear := fmt.Sprintf("%d", yearInt-1)
	var prevExpenses float64
	s.db.QueryRow(`
		SELECT COALESCE(SUM(-t.amount), 0) FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ? AND t.is_transfer = 0
		  AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0))
	`, prevYear+"%").Scan(&prevExpenses)

	if prevExpenses > 0 {
		stats.YearOverYear = (stats.TotalExpenses - prevExpenses) / prevExpenses * 100
	}

	// Category breakdown
	stats.CategoryBreakdown, _ = s.getCategoryBreakdown(year)

	// Top merchants
	stats.TopMerchants, _ = s.getTopMerchants(year)

	// Largest expenses
	stats.LargestExpenses, _ = s.getLargestExpenses(year)

	// Monthly spending within the year
	stats.MonthlySpending, _ = s.getMonthlySpending(year)

	return stats, nil
}

func (s *Service) getMonthlySpending(year string) ([]models.MonthlySpend, error) {
	rows, err := s.db.Query(`
		SELECT substr(t.transaction_date, 1, 7) as month,
		       COALESCE(SUM(CASE WHEN t.is_transfer = 0 AND (t.amount < 0 OR (c.is_expense = 1 AND t.amount > 0)) THEN -t.amount ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN t.amount > 0 AND t.is_transfer = 0 AND (c.is_expense = 0 OR c.is_expense IS NULL) THEN t.amount ELSE 0 END), 0)
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.transaction_date LIKE ?
		GROUP BY month
		ORDER BY month
	`, year+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MonthlySpend
	for rows.Next() {
		var ms models.MonthlySpend
		if err := rows.Scan(&ms.Month, &ms.Total, &ms.Income); err != nil {
			return nil, err
		}
		result = append(result, ms)
	}
	return result, rows.Err()
}

func (s *Service) GetAvailableYears() ([]models.AvailableYear, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT substr(transaction_date, 1, 4) as year
		FROM transactions
		ORDER BY year DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.AvailableYear
	for rows.Next() {
		var ay models.AvailableYear
		if err := rows.Scan(&ay.Year); err != nil {
			return nil, err
		}
		ay.Label = ay.Year
		result = append(result, ay)
	}
	return result, rows.Err()
}

func prevMonthStr(month string) string {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return month
	}
	prev := t.AddDate(0, -1, 0)
	return prev.Format("2006-01")
}
