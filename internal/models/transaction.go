package models

type Transaction struct {
	ID              int64   `json:"id"`
	BookingDate     string  `json:"bookingDate"`
	TransactionDate string  `json:"transactionDate"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	Balance         float64 `json:"balance"`
	CategoryID      *int64  `json:"categoryId"`
	CategoryName    *string `json:"categoryName"`
	CategoryColor   *string `json:"categoryColor"`
	CategoryIcon    *string `json:"categoryIcon"`
	MerchantKey     string  `json:"merchantKey"`
	IsTransfer      bool    `json:"isTransfer"`
}

type ImportResult struct {
	TotalRows        int `json:"totalRows"`
	NewTransactions  int `json:"newTransactions"`
	DuplicatesSkipped int `json:"duplicatesSkipped"`
	AutoCategorized  int `json:"autoCategorized"`
	Uncategorized    int `json:"uncategorized"`
}

type MerchantGroup struct {
	MerchantKey  string        `json:"merchantKey"`
	Count        int           `json:"count"`
	TotalAmount  float64       `json:"totalAmount"`
	IncomeTotal  float64       `json:"incomeTotal"`
	ExpenseTotal float64       `json:"expenseTotal"`
	FirstDate    string        `json:"firstDate"`
	LastDate     string        `json:"lastDate"`
	Transactions []Transaction `json:"transactions"`
}
