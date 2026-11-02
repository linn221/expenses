package dtos_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
)

func TestDecimalAmountsDoNotPanic(t *testing.T) {
	v := validator.New()
	income := dtos.NewIncome{
		Title:      "salary",
		Amount:     dtos.NaturalDecimal{Decimal: decimal.NewFromInt(500000)},
		IncomeDate: models.MyDate{Time: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)},
		WalletId:   1,
	}
	if err := v.Struct(income); err != nil {
		t.Fatal(err)
	}

	expense := dtos.NewExpense{
		ExpenseDate: models.MyDate{Time: time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)},
		Amount:      decimal.NewFromInt(4000),
		Note:        "latte",
		ItemId:      1,
		Qty:         1,
		ItemPrice:   dtos.NaturalDecimal{Decimal: decimal.NewFromInt(4000)},
	}
	if err := v.Struct(expense); err != nil {
		t.Fatal(err)
	}
}

func TestNaturalDecimalUnmarshal(t *testing.T) {
	var amount dtos.NaturalDecimal
	if err := json.Unmarshal([]byte(`"12.50"`), &amount); err != nil {
		t.Fatal(err)
	}
	if !amount.Equal(decimal.RequireFromString("12.50")) {
		t.Fatalf("amount = %s", amount)
	}

	if err := json.Unmarshal([]byte(`3`), &amount); err != nil {
		t.Fatal(err)
	}
	if !amount.Equal(decimal.NewFromInt(3)) {
		t.Fatalf("amount = %s", amount)
	}

	for _, raw := range []string{`"0"`, `0`, `"-1"`, `-1`, `null`, `"nope"`} {
		err := json.Unmarshal([]byte(raw), &amount)
		if !errors.Is(err, utils.ErrBadRequest) {
			t.Fatalf("%s: got %v", raw, err)
		}
	}
}
