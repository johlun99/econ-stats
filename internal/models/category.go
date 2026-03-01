package models

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
	IsIncome  bool   `json:"isIncome"`
	IsExpense bool   `json:"isExpense"`
	SortOrder int    `json:"sortOrder"`
}

type CategoryRule struct {
	ID           int64  `json:"id"`
	MerchantKey  string `json:"merchantKey"`
	CategoryID   int64  `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}
