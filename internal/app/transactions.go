package app

import (
	"econ-stats/internal/models"
)

func (a *App) GetTransactions(month string) ([]models.Transaction, error) {
	return a.transactionSvc.GetAll(month)
}

func (a *App) SearchTransactions(term string, month string) ([]models.Transaction, error) {
	return a.transactionSvc.Search(term, month)
}

func (a *App) GetUncategorizedMerchants() ([]models.MerchantGroup, error) {
	return a.transactionSvc.GetUncategorizedMerchants()
}

func (a *App) GetMerchantTransactions(merchantKey string) ([]models.Transaction, error) {
	return a.transactionSvc.GetMerchantTransactions(merchantKey)
}

func (a *App) UpdateTransactionCategory(id int64, categoryID *int64) error {
	return a.transactionSvc.UpdateCategory(id, categoryID)
}

func (a *App) GetUncategorizedCount() (int, error) {
	return a.transactionSvc.GetUncategorizedCount()
}
