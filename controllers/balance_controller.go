package controllers

import (
	"net/http"

	"github.com/linn221/expenses/models"
)

type BalanceController struct {
	MyController
}

func (c *BalanceController) Create() http.HandlerFunc {
	return CreateHandler[models.Balance, models.NewBalance](c.Container, c.Services.BalanceService)
}

func (c *BalanceController) Update() http.HandlerFunc {
	return UpdateHandler[models.Balance, models.NewBalance](c.Container, c.Services.BalanceService)
}

func (c *BalanceController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Balance](c.Container, c.Services.BalanceService)
}

func (c *BalanceController) Get() http.HandlerFunc {
	return GetHandler[models.Balance](c.Services.BalanceService)
}

func (c *BalanceController) List() http.HandlerFunc {
	return ListAllHandler[models.Balance]
}

func (c *BalanceController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.BalanceService)
}
