package app

import (
	"context"
	"database/sql"

	"econ-stats/internal/database"
	"econ-stats/internal/services/categorizer"
	"econ-stats/internal/services/category"
	"econ-stats/internal/services/stats"
	"econ-stats/internal/services/transaction"
)

type App struct {
	ctx            context.Context
	db             *sql.DB
	transactionSvc *transaction.Service
	categorySvc    *category.Service
	categorizerEng *categorizer.Engine
	statsSvc       *stats.Service
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	db, err := database.Open()
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	a.db = db

	a.transactionSvc = transaction.NewService(db)
	a.categorySvc = category.NewService(db)
	a.categorizerEng = categorizer.NewEngine(a.categorySvc, a.transactionSvc)
	a.statsSvc = stats.NewService(db)
}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}
