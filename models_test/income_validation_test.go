package models_test

import (
	"context"
	"errors"
	"testing"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

func TestIncomeAmountMustBePositive(t *testing.T) {
	session := &app.Session{UserId: 1, Ctx: context.Background()}
	income := &models.Income{
		Amount:      decimal.NewFromInt(-1),
		HasUserId:   models.HasUserId{UserId: 1},
		HasWalletId: models.HasWalletId{WalletId: 1},
	}

	err := validate.Validate(nil, session, income, validate.ActionCreate)
	if !errors.Is(err, validate.ErrClient) {
		t.Fatalf("negative amount: got %v", err)
	}
}
