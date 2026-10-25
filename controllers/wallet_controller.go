package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type walletController struct {
	myController
}

func (c *walletController) Create() http.HandlerFunc {
	return CreateHandler[models.Wallet, dtos.NewWallet](c.Container, c.Services.WalletService)
}

func (c *walletController) Update() http.HandlerFunc {
	return UpdateHandler[models.Wallet, dtos.NewWallet](c.Container, c.Services.WalletService)
}

func (c *walletController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Wallet](c.Container, c.Services.WalletService)
}

func (c *walletController) Get() http.HandlerFunc {
	return GetHandler[models.Wallet](c.Services.WalletService)
}

func (c *walletController) List() http.HandlerFunc {
	return ListAllHandler[models.Wallet]
}

func (c *walletController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.WalletService)
}
