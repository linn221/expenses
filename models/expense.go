package models

import (
	"net/http"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Expense struct {
	HasBalanceId
	HasUserId
	Id          int             `json:"id" gorm:"primaryKey"`
	ExpenseDate MyDate          `json:"expense_date" gorm:"index"`
	Amount      decimal.Decimal `json:"amount"`
	Note        string          `json:"note"`
	BalanceId   int

	ItemId    int             `json:"item_id" gorm:"index"`
	Qty       int             `json:"qty"`
	ItemPrice decimal.Decimal `json:"item_price"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// type DailyExpense struct {
// 	HasUserId
// 	HasBalanceId
// 	Date   MyDate          `json:"date" gorm:"index"`
// 	Amount decimal.Decimal `json:"amount"`
// }

type NewExpense struct {
	ExpenseDate MyDate          `json:"expense_date" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Note        string          `json:"note" validate:"required,min=3,max=255"`
	BalanceId   int

	ItemId          int             `json:"item_id" validate:"required,gte=1"`
	Qty             int             `json:"qty" validate:"required,gte=1"`
	ItemPrice       decimal.Decimal `json:"item_price" validate:"required"`
	UpdateItemPrice bool
}

type ExpenseFilter struct {
	ExpenseDate MyDate
	Search      string
	MaxAmount   decimal.Decimal
	MinAmount   decimal.Decimal
}

func (e ExpenseFilter) Parse(r *http.Request) (*ExpenseFilter, error) {
	q := r.URL.Query()

	if search := q.Get("search"); search != "" {
		e.Search = search
	}

	if dateq := q.Get("expense_date"); dateq != "" {
		date, err := time.Parse("2006-01-02", dateq)
		if err != nil {
			return nil, utils.NewBadRequest("expense_date is not a valid date")
		}
		e.ExpenseDate = MyDate{Time: date}
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

func (input NewExpense) ToModel(session *app.Session) Expense {

	session.PutValue(app.SessionValueUpdateItemPrice, input.ItemPrice)
	e := Expense{
		ItemId:      input.ItemId,
		Qty:         input.Qty,
		ItemPrice:   input.ItemPrice,
		Amount:      input.ItemPrice.Mul(decimal.NewFromInt(int64(input.Qty))),
		ExpenseDate: input.ExpenseDate,
		Note:        input.Note,
	}

	e.UserId = session.UserId
	e.BalanceId = input.BalanceId
	return e
}

func (existing *Expense) UpdateValues(input *Expense) map[string]any {
	m := map[string]any{
		"ExpenseDate": input.ExpenseDate,
		"Note":        input.Note,
		"ItemId":      input.ItemId,
		"Qty":         input.Qty,
		"ItemPrice":   input.ItemPrice,
	}
	return m
}

func (input *NewExpense) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		if !input.ItemPrice.GreaterThan(decimal.Zero) {
			af(validate.Failure("item price must be greater than 0"))
			return
		}
		userFilter := UserFilter(s)
		af(
			validate.ExistsRule("balances", input.BalanceId, "balance not found", userFilter),
			validate.ExistsRule("items", input.ItemId, "item not found", userFilter),
		)
	})
}
