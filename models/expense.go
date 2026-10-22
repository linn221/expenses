package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	HasBalanceId
	HasUserId
	Id          int             `json:"id" gorm:"primaryKey"`
	ExpenseDate MyDate          `json:"expense_date" gorm:"index"`
	Amount      decimal.Decimal `json:"amount"`
	Note        string          `json:"note"`

	ItemId       int             `json:"item_id" gorm:"index"`
	Qty          int             `json:"qty"`
	ExpensePrice decimal.Decimal `json:"expense_price"`
	CreatedAt    time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type DailyExpense struct {
	HasUserId
	HasBalanceId
	Date   MyDate          `json:"date" gorm:"index"`
	Amount decimal.Decimal `json:"amount"`
}

type NewExpense struct {
	ExpenseDate MyDate          `json:"expense_date" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Note        string          `json:"note" validate:"required,min=3,max=255"`

	ItemId       int             `json:"item_id" validate:"required"`
	Qty          int             `json:"qty" validate:"required,gt=0"`
	ExpensePrice decimal.Decimal `json:"expense_price" validate:"required,gt=0"`
}
