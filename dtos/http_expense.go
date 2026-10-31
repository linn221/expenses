package dtos

import (
	"net/http"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
)

type NewExpense struct {
	ExpenseDate models.MyDate   `json:"expense_date" validate:"required"`
	Amount      decimal.Decimal `json:"amount"`
	Note        string          `json:"note" validate:"omitempty,min=3,max=255"`
	WalletId    int

	ItemId          int             `json:"item_id" validate:"required,gte=1"`
	Qty             int             `json:"qty" validate:"required,gte=1"`
	ItemPrice       decimal.Decimal `json:"item_price" validate:"required"`
	UpdateItemPrice bool
}

type ExpenseFilter struct {
	WalletId    int
	ExpenseDate models.MyDate
	Search      string
	MaxAmount   decimal.Decimal
	MinAmount   decimal.Decimal
}

func (e ExpenseFilter) Parse(r *http.Request) (*ExpenseFilter, error) {
	walletID, err := walletIDFromPath(r)
	if err != nil {
		return nil, err
	}
	e.WalletId = walletID

	q := r.URL.Query()

	if search := q.Get("search"); search != "" {
		e.Search = search
	}

	if dateq := q.Get("expense_date"); dateq != "" {
		date, err := time.Parse("2006-01-02", dateq)
		if err != nil {
			return nil, utils.NewBadRequest("expense_date is not a valid date")
		}
		e.ExpenseDate = models.MyDate{Time: date}
	}

	if maxq := q.Get("max_amount"); maxq != "" {
		maxAmount, err := decimal.NewFromString(maxq)
		if err != nil {
			return nil, utils.NewBadRequest("max_amount is not a valid decimal")
		}
		e.MaxAmount = maxAmount
	}

	if minq := q.Get("min_amount"); minq != "" {
		minAmount, err := decimal.NewFromString(minq)
		if err != nil {
			return nil, utils.NewBadRequest("min_amount is not a valid decimal")
		}
		e.MinAmount = minAmount
	}

	return &e, nil
}

func (input NewExpense) ToModel(session *app.Session) models.Expense {

	session.PutValue(app.SessionValueUpdateItemPrice, input.UpdateItemPrice)
	e := models.Expense{
		ItemId:        input.ItemId,
		Qty:           input.Qty,
		ItemPrice:     input.ItemPrice,
		ExpenseAmount: input.ItemPrice.Mul(decimal.NewFromInt(int64(input.Qty))),
		ExpenseDate:   input.ExpenseDate,
		Note:          input.Note,
	}

	e.UserId = session.UserId
	e.WalletId = input.WalletId
	return e
}

// use the same struct for collection response
type expenseResponse struct {
	Id          int             `json:"id"`
	ExpenseDate models.MyDate   `json:"expense_date"`
	Amount      decimal.Decimal `json:"amount"`
	Note        string          `json:"note"`
	Wallet      WalletResource  `json:"wallet"`
	Item        ItemResource    `json:"item"`
	Insights    []string        `json:"insights"`
}

type ItemResource struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
	Category string          `json:"category"`
}

func toExpenseResponse(expense *models.Expense) expenseResponse {
	return expenseResponse{
		Id:          expense.Id,
		ExpenseDate: expense.ExpenseDate,
		Amount:      expense.ExpenseAmount,
		Note:        expense.Note,
		Wallet: WalletResource{
			Id:   expense.WalletId,
			Name: expense.Wallet.Name,
		},
		Item: ItemResource{
			Id:       expense.ItemId,
			Name:     expense.Item.Name,
			Price:    expense.Item.Price,
			Category: expense.Item.Category.Name,
		},
		Insights: []string{},
	}
}

func RespondExpense(expense *models.Expense) any {
	response := toExpenseResponse(expense)
	return &response
}

type expenseCollection struct {
	Expenses []expenseResponse `json:"expenses"`
}

func RespondExpenseCollection(results []models.Expense) any {
	expenses := addToCollection(results, func(expense models.Expense) expenseResponse {
		return toExpenseResponse(&expense)
	})
	return &expenseCollection{Expenses: expenses}
}
