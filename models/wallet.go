package models

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	Id int `json:"id" gorm:"primaryKey"`
	HasUserId
	HasIsActive
	Name     string          `json:"name"`
	Current  decimal.Decimal `json:"current"`
	Currency string          `json:"currency"`
}

func (existing *Wallet) UpdateValues(input *Wallet) map[string]any {
	return map[string]any{
		"Name":     input.Name,
		"Currency": input.Currency,
	}
}

func (i *Wallet) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("wallets", s.ResId, "wallet not found", userFilter).When(s.ResId > 0),
			validate.UniqueRule("wallets", "name", i.Name, s.ResId, userFilter),
			validate.UniqueRule("wallets", "currency", i.Currency, s.ResId, userFilter),
		)
	})
	vc.OnDelete(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("wallets", s.ResId, "wallet not found", userFilter),
			validate.NoResultRule("expenses", "wallet is used", validate.NewFilter(
				"wallet_id = ?", s.ResId,
			)),
		)
	})
	vc.OnActiveStatusChange(func(s *app.Session, af validate.AppendFunc) {
		af(
			// no need to check if wallet exists since the model is already fetched by default handler
			validate.NoResultRule("expenses", "wallet is used", validate.NewFilter(
				"wallet_id = ?", s.ResId,
			)).When(i.IsActive == false), // validate the wallet is not used if it is to be deactivated
		)

	})
}
