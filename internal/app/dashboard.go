package app

import (
	"econ-stats/internal/models"
)

func (a *App) GetMonthlyStats(month string) (*models.MonthlyStats, error) {
	return a.statsSvc.GetMonthlyStats(month)
}

func (a *App) GetSpendingTrend(months int) ([]models.SpendingTrend, error) {
	return a.statsSvc.GetSpendingTrend(months)
}

func (a *App) GetAvailableMonths() ([]models.AvailableMonth, error) {
	return a.statsSvc.GetAvailableMonths()
}

func (a *App) GetYearlyStats(year string) (*models.YearlyStats, error) {
	return a.statsSvc.GetYearlyStats(year)
}

func (a *App) GetAvailableYears() ([]models.AvailableYear, error) {
	return a.statsSvc.GetAvailableYears()
}
