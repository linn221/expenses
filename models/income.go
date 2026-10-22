package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Income struct {
	HasUserId
	HasBalanceId
	Id         int             `json:"id" gorm:"primaryKey"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	IncomeDate MyDate          `json:"income_date" gorm:"index"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type NewIncome struct {
	IncomeDate MyDate          `json:"income_date" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Title      string          `json:"title" validate:"required,min=3,max=255"`
}
