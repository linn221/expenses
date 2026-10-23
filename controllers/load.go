package controllers

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
)

type Controllers struct {
	CategoryController *categoryController
	ItemController     *itemController
	LoginController    *loginController
	BalanceController  *balanceController
	IncomeController   *incomeController
	ExpenseController  *expenseController
}

func LoadControllers(c *app.Container, provider *services.Provider) Controllers {
	base := myController{Container: c, Services: provider}
	return Controllers{
		CategoryController: &categoryController{myController: base},
		ItemController:     &itemController{myController: base},
		LoginController:    &loginController{Services: provider},
		BalanceController:  &balanceController{myController: base},
		IncomeController:   &incomeController{myController: base},
		ExpenseController:  &expenseController{myController: base},
	}
}
