package dtos

import "github.com/shopspring/decimal"

type ProfileWallet struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Currency string          `json:"currency"`
	Current  decimal.Decimal `json:"current"`
}

type Profile struct {
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Wallets []ProfileWallet `json:"wallets"`
}
