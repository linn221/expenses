package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type expenseController struct {
	myController
}

func (c *expenseController) Create() http.HandlerFunc {
	return CreateHandler[models.Expense, dtos.NewExpense](c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Update() http.HandlerFunc {
	return UpdateHandler[models.Expense, dtos.NewExpense](c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Delete() http.HandlerFunc {
	return DeleteHandler(c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Get() http.HandlerFunc {
	return GetHandler(c.Services.ExpenseService)
}

func (c *expenseController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.Expense, dtos.ExpenseFilter](c.Services.ExpenseService)
}
