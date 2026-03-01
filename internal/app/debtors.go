package app

import (
	"econ-stats/internal/models"
)

func (a *App) GetDebtors() ([]models.DebtorDetail, error) {
	return a.debtorSvc.GetAll()
}

func (a *App) ToggleDebtorPinned(id int64, pinned bool) error {
	return a.debtorSvc.TogglePinned(id, pinned)
}

func (a *App) GetPinnedDebtors() ([]models.DebtorDetail, error) {
	return a.debtorSvc.GetPinned()
}

func (a *App) CreateDebtor(name, icon, color string) (*models.Debtor, error) {
	return a.debtorSvc.Create(name, icon, color)
}

func (a *App) UpdateDebtor(id int64, name, icon, color string) error {
	return a.debtorSvc.Update(id, name, icon, color)
}

func (a *App) DeleteDebtor(id int64) error {
	return a.debtorSvc.Delete(id)
}

func (a *App) AddDebtorMerchantKey(debtorID int64, merchantKey string) error {
	return a.debtorSvc.AddMerchantKey(debtorID, merchantKey)
}

func (a *App) RemoveDebtorMerchantKey(debtorID int64, merchantKey string) error {
	return a.debtorSvc.RemoveMerchantKey(debtorID, merchantKey)
}

func (a *App) GetDebtorTransactions(debtorID int64) ([]models.Transaction, error) {
	return a.debtorSvc.GetTransactions(debtorID)
}

func (a *App) GetAllMerchantKeys() ([]string, error) {
	return a.debtorSvc.GetAllMerchantKeys()
}

func (a *App) AddDebtorTransaction(debtorID int64, description string, amount float64, date string) error {
	return a.debtorSvc.AddManualTransaction(debtorID, description, amount, date)
}

func (a *App) UpdateDebtorTransaction(id int64, description string, amount float64, date string) error {
	return a.debtorSvc.UpdateManualTransaction(id, description, amount, date)
}

func (a *App) DeleteDebtorTransaction(id int64) error {
	return a.debtorSvc.DeleteManualTransaction(id)
}
