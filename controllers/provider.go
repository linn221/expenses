package controllers

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
)

type Provider struct {
	CategoryController        *categoryController
	ItemController            *itemController
	LoginController           *loginController
	ProfileController         *profileController
	WalletController          *walletController
	WalletTransferController  *walletTransferController
	IncomeController          *incomeController
	ExpenseController         *expenseController
	CostPerspectiveController *costPerspectiveController
}

func New(c *app.Container, provider *services.Provider) *Provider {
	base := myController{Container: c, Services: provider}
	return &Provider{
		CategoryController:        &categoryController{myController: base},
		ItemController:            &itemController{myController: base},
		LoginController:           &loginController{Services: provider},
		ProfileController:         &profileController{myController: base},
		WalletController:          &walletController{myController: base},
		WalletTransferController:  &walletTransferController{myController: base},
		IncomeController:          &incomeController{myController: base},
		ExpenseController:         &expenseController{myController: base},
		CostPerspectiveController: &costPerspectiveController{myController: base},
	}
}
