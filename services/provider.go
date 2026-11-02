package services

import (
	"github.com/linn221/expenses/app"
)

type Provider struct {
	CategoryService        *categoryService
	LoginService           *loginService
	ProfileService         *profileService
	SessionIdService       *sessionIdService
	WalletService          *walletService
	WalletTransferService  *walletTransferService
	ItemService            *itemService
	IncomeService          *incomeService
	ExpenseService         *expenseService
	CostPerspectiveService *costPerspectiveService
}

// for accessing of app container without the dependency injection
var Default *defaultService

func New(c *app.Container) *Provider {
	Default = &defaultService{container: c}
	return &Provider{
		CategoryService:        &categoryService{Container: c},
		LoginService:           &loginService{Container: c},
		ProfileService:         &profileService{Container: c},
		SessionIdService:       &sessionIdService{Container: c},
		WalletService:          &walletService{Container: c},
		WalletTransferService:  &walletTransferService{Container: c},
		ItemService:            &itemService{Container: c},
		IncomeService:          &incomeService{Container: c},
		ExpenseService:         &expenseService{Container: c},
		CostPerspectiveService: &costPerspectiveService{Container: c},
	}
}
