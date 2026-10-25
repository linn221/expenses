package models

import (
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Income struct {
	HasUserId
	HasWalletId
	Id         int             `json:"id" gorm:"primaryKey"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	IncomeDate MyDate          `json:"income_date" gorm:"index"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

func (existing *Income) UpdateValues(input *Income) map[string]any {
	return map[string]any{
		"Title":     input.Title,
		"Amount":    input.Amount,
		"IcomeDate": input.IncomeDate,
		// don't let wallet_id change
	}
}

func (input *Income) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(validate.ExistsRule("wallets", input.WalletId, "wallet not found", userFilter))
	})
}
