package models

import (
	"github.com/shopspring/decimal"
)

// for informing the user how much is a certain expense
// eg, 5 % of his salary, 5 breakfasts, 7 hours of your time worth, etc
type CostPerspective struct {
	Id           int                          `json:"id" gorm:"primaryKey"`
	Name         string                       `json:"name" gorm:"index"`
	Value        decimal.Decimal              `json:"value"`
	Type         CostPerspectiveType          `json:"type" gorm:"index"`
	Applications []CostPerspectiveApplication `json:"applications" gorm:"foreignKey:CostPerspectiveId"`
	HasUserId
	HasIsActive
}

// the certain perspective will be applied when starting value is reached on a certain insight
type CostPerspectiveApplication struct {
	Id                   int                  `json:"id" gorm:"primaryKey"`
	CostPerspectiveId    int                  `json:"cost_perspective_id" gorm:"index"`
	MoneyInsightCategory MoneyInsightCategory `json:"money_insight_category"`
	StartingValue        decimal.Decimal      `json:"starting_value"`
}

type MoneyInsightCategory string

const (
	MoneyInsightCategoryExpense      = "expense"
	MoneyInsightCategoryIncome       = "income"
	MoneyInsightCategoryDailyExpense = "daily_expense"
)

func (existing *CostPerspective) UpdateValues(new *CostPerspective) map[string]any {
	changes := make(map[string]any)
	if existing.Name != new.Name {
		changes["name"] = new.Name
	}
	if existing.Value.String() != new.Value.String() {
		changes["value"] = new.Value
	}
	if existing.Type != new.Type {
		changes["type"] = new.Type
	}
	return changes
}
