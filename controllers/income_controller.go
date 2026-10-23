package controllers

import (
	"net/http"

	"github.com/linn221/expenses/models"
)

type incomeController struct {
	myController
}

func (c *incomeController) Create() http.HandlerFunc {
	return CreateHandler[models.Income, models.NewIncome](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Update() http.HandlerFunc {
	return UpdateHandler[models.Income, models.NewIncome](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Income](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Get() http.HandlerFunc {
	return GetHandler[models.Income](c.Services.IncomeService)
}

func (c *incomeController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.Income, models.IncomeFilter](c.Services.IncomeService)
}
