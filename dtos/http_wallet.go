package dtos

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
)

type NewWallet struct {
	Name     string
	Currency string
}

func (i NewWallet) ToModel(session *app.Session) models.Wallet {
	return models.Wallet{
		Name:      i.Name,
		Currency:  i.Currency,
		HasUserId: models.UserId(session),
	}
}

type WalletResource struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
