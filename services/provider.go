package services

import (
	"github.com/linn221/expenses/app"
)

type Provider struct {
	CategoryService  *categoryService
	LoginService     *loginService
	SessionIdService *sessionIdService
	BalanceService   *balanceService
	ItemService      *itemService
	IncomeService    *incomeService
	ExpenseService   *expenseService
}

// for accessing of app container without the dependency injection
var Default *defaultService

func New(c *app.Container) *Provider {
	Default = &defaultService{container: c}
	return &Provider{
		CategoryService:  &categoryService{Container: c},
		LoginService:     &loginService{Container: c},
		SessionIdService: &sessionIdService{Container: c},
		BalanceService:   &balanceService{Container: c},
		ItemService:      &itemService{Container: c},
		IncomeService:    &incomeService{Container: c},
		ExpenseService:   &expenseService{Container: c},
	}
}
