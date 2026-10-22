package models

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Balance struct {
	Id int `json:"id" gorm:"primaryKey"`
	HasUserId
	HasIsActive
	Name     string
	Current  decimal.Decimal `json:"current"`
	Currency string          `json:"currency"`
}

type NewBalance struct {
	Name     string
	Currency string
}

func (i NewBalance) ToModel(session *app.Session) Balance {
	return Balance{
		Name:      i.Name,
		Currency:  i.Currency,
		HasUserId: UserId(session),
	}
}

func (existing *Balance) UpdateValues(input *Balance) map[string]any {
	return map[string]any{
		"Name":     input.Name,
		"Currency": input.Currency,
	}
}

func (i *NewBalance) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("balances", s.ResId, "balance not found", userFilter).When(s.ResId > 0),
			validate.UniqueRule("balances", "name", i.Name, s.ResId, userFilter),
			validate.UniqueRule("balances", "currency", i.Currency, s.ResId, userFilter),
		)
	})
}

func (i *Balance) ValidationRules(vc *validate.ValidationContext) {
	vc.OnDelete(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("balances", s.ResId, "balance not found", userFilter),
			validate.NoResultRule("expenses", "balance is used", validate.NewFilter(
				"balance_id = ?", s.ResId,
			)),
		)
	})
	vc.OnActiveStatusChange(func(s *app.Session, af validate.AppendFunc) {
		af(
			// no need to check if balance exists since the model is already fetched by default handler
			validate.NoResultRule("expenses", "balance is used", validate.NewFilter(
				"balance_id = ?", s.ResId,
			)).When(i.IsActive == false), // validate the balance is not used if it is to be deactivated
		)

	})
}
