package models

type MonthlyStats struct {
	Month             string          `json:"month"` // YYYY-MM
	TotalExpenses     float64         `json:"totalExpenses"`
	TotalIncome       float64         `json:"totalIncome"`
	NetSavings        float64         `json:"netSavings"`
	SavingsRate       float64         `json:"savingsRate"`
	AvgDailySpend     float64         `json:"avgDailySpend"`
	MonthOverMonth    float64         `json:"monthOverMonth"` // % change
	CategoryBreakdown []CategorySpend `json:"categoryBreakdown"`
	TopMerchants      []MerchantSpend `json:"topMerchants"`
	LargestExpenses   []Transaction   `json:"largestExpenses"`
	DailySpending     []DailySpend    `json:"dailySpending"`
}

type CategorySpend struct {
	CategoryID   int64   `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	CategoryColor string `json:"categoryColor"`
	CategoryIcon string  `json:"categoryIcon"`
	Total        float64 `json:"total"`
	Count        int     `json:"count"`
}

type MerchantSpend struct {
	MerchantKey string  `json:"merchantKey"`
	Total       float64 `json:"total"`
	Count       int     `json:"count"`
}

type DailySpend struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type SpendingTrend struct {
	Month    string  `json:"month"`
	Expenses float64 `json:"expenses"`
	Income   float64 `json:"income"`
}

type AvailableMonth struct {
	Month string `json:"month"` // YYYY-MM
	Label string `json:"label"` // "February 2026"
}
