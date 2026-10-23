package models

import (
	"net/http"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
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
	BalanceId  int             `json:"balance_id" validate:"required,gte=1"`
}

type IncomeFilter struct {
	IncomeDate MyDate
	Search     string
	MaxAmount  decimal.Decimal
	MinAmount  decimal.Decimal
}

func (f IncomeFilter) Parse(r *http.Request) (*IncomeFilter, error) {
	q := r.URL.Query()

	if search := q.Get("search"); search != "" {
		f.Search = search
	}

	if dateq := q.Get("income_date"); dateq != "" {
		date, err := time.Parse("2006-01-02", dateq)
		if err != nil {
			return nil, utils.NewBadRequest("income_date is not a valid date")
		}
		f.IncomeDate = MyDate{Time: date}
	}

	if maxq := q.Get("max_amount"); maxq != "" {
		maxAmount, err := decimal.NewFromString(maxq)
		if err != nil {
			return nil, utils.NewBadRequest("max_amount is not a valid decimal")
		}
		f.MaxAmount = maxAmount
	}

	if minq := q.Get("min_amount"); minq != "" {
		minAmount, err := decimal.NewFromString(minq)
		if err != nil {
			return nil, utils.NewBadRequest("min_amount is not a valid decimal")
		}
		f.MinAmount = minAmount
	}

	return &f, nil
}

func (input NewIncome) ToModel(session *app.Session) Income {
	model := Income{
		Title:      input.Title,
		Amount:     input.Amount,
		IncomeDate: input.IncomeDate,
	}
	model.UserId = session.UserId
	model.BalanceId = input.BalanceId
	return model
}

func (existing *Income) UpdateValues(input *Income) map[string]any {
	return map[string]any{
		"Title":     input.Title,
		"Amount":    input.Amount,
		"IcomeDate": input.IncomeDate,
		// don't let balance_id change
	}
}

func (input *NewIncome) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		userFilter := UserFilter(s)
		af(validate.ExistsRule("balances", input.BalanceId, "balance not found", userFilter))
	})
}
