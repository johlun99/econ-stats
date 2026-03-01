package categorizer

import (
	"econ-stats/internal/services/category"
	"econ-stats/internal/services/transaction"
)

type Engine struct {
	categorySvc    *category.Service
	transactionSvc *transaction.Service
}

func NewEngine(categorySvc *category.Service, transactionSvc *transaction.Service) *Engine {
	return &Engine{
		categorySvc:    categorySvc,
		transactionSvc: transactionSvc,
	}
}

// AutoCategorize applies all existing rules to uncategorized transactions.
// Returns the number of newly categorized transactions.
func (e *Engine) AutoCategorize() (int, error) {
	rules, err := e.categorySvc.GetRulesMap()
	if err != nil {
		return 0, err
	}
	return e.transactionSvc.AutoCategorize(rules)
}

// CategorizeByMerchant creates a rule and applies it retroactively.
// Returns the number of transactions updated.
func (e *Engine) CategorizeByMerchant(merchantKey string, categoryID int64) (int64, error) {
	// Create or update the rule
	if err := e.categorySvc.CreateRule(merchantKey, categoryID); err != nil {
		return 0, err
	}

	// Apply retroactively to all transactions with this merchant
	return e.transactionSvc.SetCategoryByMerchant(merchantKey, categoryID)
}
