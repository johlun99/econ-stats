package app

import (
	"econ-stats/internal/models"
)

func (a *App) GetCategories() ([]models.Category, error) {
	return a.categorySvc.GetAll()
}

func (a *App) CreateCategory(name, color, icon string, isIncome, isExpense bool) (*models.Category, error) {
	return a.categorySvc.Create(name, color, icon, isIncome, isExpense)
}

func (a *App) UpdateCategory(id int64, name, color, icon string, isIncome, isExpense bool) error {
	return a.categorySvc.Update(id, name, color, icon, isIncome, isExpense)
}

func (a *App) DeleteCategory(id int64) error {
	return a.categorySvc.Delete(id)
}

func (a *App) GetCategoryRules() ([]models.CategoryRule, error) {
	return a.categorySvc.GetRules()
}

func (a *App) CategorizeByMerchant(merchantKey string, categoryID int64) (int64, error) {
	return a.categorizerEng.CategorizeByMerchant(merchantKey, categoryID)
}

func (a *App) DeleteCategoryRule(id int64) error {
	return a.categorySvc.DeleteRule(id)
}
