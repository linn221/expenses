package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type walletTransferController struct {
	myController
}

func (c *walletTransferController) Create() http.HandlerFunc {
	return CreateHandler[models.WalletTransfer, dtos.NewWalletTransfer](c.Container, c.Services.WalletTransferService)
}

func (c *walletTransferController) Update() http.HandlerFunc {
	return UpdateHandler[models.WalletTransfer, dtos.NewWalletTransfer](c.Container, c.Services.WalletTransferService)
}

func (c *walletTransferController) Delete() http.HandlerFunc {
	return DeleteHandler[models.WalletTransfer](c.Container, c.Services.WalletTransferService)
}

func (c *walletTransferController) Get() http.HandlerFunc {
	return GetHandler[models.WalletTransfer](c.Services.WalletTransferService)
}

func (c *walletTransferController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.WalletTransfer, dtos.WalletTransferFilter](c.Services.WalletTransferService)
}

func (c *walletTransferController) Paginate() http.HandlerFunc {
	return PaginateHandler[models.WalletTransfer, dtos.WalletTransferFilter](c.Services.WalletTransferService)
}
