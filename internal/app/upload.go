package app

import (
	"econ-stats/internal/models"
	"econ-stats/internal/services/parser"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SelectAndImportFile() (*models.ImportResult, error) {
	filepath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Välj Handelsbanken-fil",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.xlsx"},
		},
	})
	if err != nil {
		return nil, err
	}
	if filepath == "" {
		return nil, nil // User cancelled
	}

	return a.ImportFile(filepath)
}

func (a *App) ImportFile(filepath string) (*models.ImportResult, error) {
	transactions, err := parser.ParseHandelsbanken(filepath)
	if err != nil {
		return nil, err
	}

	result, err := a.transactionSvc.Import(transactions)
	if err != nil {
		return nil, err
	}

	// Auto-categorize with existing rules
	autoCategorized, err := a.categorizerEng.AutoCategorize()
	if err != nil {
		return nil, err
	}
	result.AutoCategorized = autoCategorized

	uncategorized, err := a.transactionSvc.GetUncategorizedCount()
	if err != nil {
		return nil, err
	}
	result.Uncategorized = uncategorized

	return result, nil
}
