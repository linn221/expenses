package controllers

import (
	"net/http"

	"github.com/linn221/expenses/models"
)

type balanceController struct {
	myController
}

func (c *balanceController) Create() http.HandlerFunc {
	return CreateHandler[models.Balance, models.NewBalance](c.Container, c.Services.BalanceService)
}

func (c *balanceController) Update() http.HandlerFunc {
	return UpdateHandler[models.Balance, models.NewBalance](c.Container, c.Services.BalanceService)
}

func (c *balanceController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Balance](c.Container, c.Services.BalanceService)
}

func (c *balanceController) Get() http.HandlerFunc {
	return GetHandler[models.Balance](c.Services.BalanceService)
}

func (c *balanceController) List() http.HandlerFunc {
	return ListAllHandler[models.Balance]
}

func (c *balanceController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.BalanceService)
}
