package dtos

import (
	"net/http"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
)

type NewIncome struct {
	IncomeDate models.MyDate   `json:"income_date" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required"`
	Title      string          `json:"title" validate:"required,min=3,max=255"`
	WalletId   int             `json:"wallet_id" validate:"required,gte=1"`
}

type IncomeFilter struct {
	WalletId   int
	IncomeDate models.MyDate
	Search     string
	MaxAmount  decimal.Decimal
	MinAmount  decimal.Decimal
}

func (f IncomeFilter) Parse(r *http.Request) (*IncomeFilter, error) {
	walletID, err := walletIDFromPath(r)
	if err != nil {
		return nil, err
	}
	f.WalletId = walletID

	q := r.URL.Query()

	if search := q.Get("search"); search != "" {
		f.Search = search
	}

	if dateq := q.Get("income_date"); dateq != "" {
		date, err := time.Parse("2006-01-02", dateq)
		if err != nil {
			return nil, utils.NewBadRequest("income_date is not a valid date")
		}
		f.IncomeDate = models.MyDate{Time: date}
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

func (input NewIncome) ToModel(session *app.Session) models.Income {
	model := models.Income{
		Title:      input.Title,
		Amount:     input.Amount,
		IncomeDate: input.IncomeDate,
	}
	model.UserId = session.UserId
	model.WalletId = input.WalletId
	return model
}

type IncomeResource struct {
	Id         int             `json:"id"`
	Title      string          `json:"title"`
	Amount     decimal.Decimal `json:"amount"`
	IncomeDate time.Time       `json:"income_date"`
	Wallet     WalletResource  `json:"wallet"`
}

type IncomeCollection struct {
	Incomes []IncomeCollectionInstance `json:"incomes"`
}

type IncomeCollectionInstance struct {
	Id     int             `json:"id"`
	Title  string          `json:"title"`
	Amount decimal.Decimal `json:"amount"`
}

func RespondIncome(income *models.Income) any {
	return &IncomeResource{
		Id:         income.Id,
		Title:      income.Title,
		Amount:     income.Amount,
		IncomeDate: income.IncomeDate.Time,
		Wallet: WalletResource{
			Id:   income.WalletId,
			Name: income.Wallet.Name,
		},
	}
}

func RespondIncomeCollection(results []models.Income) any {
	incomes := addToCollection(results, func(r models.Income) IncomeCollectionInstance {
		return IncomeCollectionInstance{
			Id:     r.Id,
			Title:  r.Title,
			Amount: r.Amount,
		}
	})
	return &IncomeCollection{Incomes: incomes}
}
