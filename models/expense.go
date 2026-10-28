package models

import (
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Expense struct {
	HasWalletId
	HasUserId
	Id            int             `json:"id" gorm:"primaryKey"`
	ExpenseDate   MyDate          `json:"expense_date" gorm:"index"`
	ExpenseAmount decimal.Decimal `json:"expense_amount"`
	Note          string          `json:"note"`

	ItemId    int             `json:"item_id" gorm:"index"`
	Qty       int             `json:"qty"`
	ItemPrice decimal.Decimal `json:"item_price"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// type DailyExpense struct {
// 	HasUserId
// 	HasWalletId
// 	Date   MyDate          `json:"date" gorm:"index"`
// 	Amount decimal.Decimal `json:"amount"`
// }

func (existing *Expense) UpdateValues(input *Expense) map[string]any {
	m := map[string]any{
		"ExpenseDate":   input.ExpenseDate,
		"ExpenseAmount": input.ExpenseAmount,
		"Note":          input.Note,
		"ItemId":        input.ItemId,
		"Qty":           input.Qty,
		"ItemPrice":     input.ItemPrice,
	}
	return m
}

func (input *Expense) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		if !input.ItemPrice.GreaterThan(decimal.Zero) {
			af(validate.Failure("item price must be greater than 0"))
			return
		}
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("wallets", input.WalletId, "wallet not found", userFilter),
			validate.ExistsRule("items", input.ItemId, "item not found", userFilter),
		)
	})
}
