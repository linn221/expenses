package dtos

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
)

func TestDecimalAmountsDoNotPanic(t *testing.T) {
	v := validator.New()
	income := NewIncome{
		Title:      "salary",
		Amount:     decimal.NewFromInt(500000),
		IncomeDate: models.MyDate{Time: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)},
		WalletId:   1,
	}
	if err := v.Struct(income); err != nil {
		t.Fatal(err)
	}

	expense := NewExpense{
		ExpenseDate: models.MyDate{Time: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)},
		Amount:      decimal.NewFromInt(4000),
		Note:        "latte",
		ItemId:      1,
		Qty:         1,
		ItemPrice:   decimal.NewFromInt(4000),
	}
	if err := v.Struct(expense); err != nil {
		t.Fatal(err)
	}
}
