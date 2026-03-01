package models

type Debtor struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type DebtorDetail struct {
	Debtor
	MerchantKeys []string `json:"merchantKeys"`
	Balance      float64  `json:"balance"`
}

type DebtorTransaction struct {
	ID              int64   `json:"id"`
	DebtorID        int64   `json:"debtorId"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transactionDate"`
}
