package models

import (
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"github.com/shopspring/decimal"
)

type Item struct {
	Id         int             `json:"id" gorm:"primaryKey"`
	Name       string          `json:"name"`
	Code       string          `json:"code" gorm:"index"`
	Price      decimal.Decimal `json:"price"`
	CategoryId int             `json:"category_id" gorm:"index"`
	Category   Category        `json:"category"`
	HasUserId
	HasIsActive
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (item *Item) Resource() any {
	return struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Price    string `json:"price"`
		Category string `json:"category"`
	}{item.Id, item.Name, item.Code, item.Price.String(), item.Category.Name}
}

func (existing *Item) UpdateValues(input *Item) map[string]any {
	return map[string]any{
		"Name":       input.Name,
		"Price":      input.Price,
		"Code":       input.Code,
		"CategoryId": input.CategoryId,
	}
}

func (input *Item) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		if !input.Price.GreaterThan(decimal.Zero) {
			af(validate.Failure("input price must be greater than zero"))
			return
		}

		userFiler := UserFilter(s)
		af(
			validate.UniqueRule("items", "code", input.Code, s.ResId, userFiler),
			validate.ExistsRule("categories", input.CategoryId, "category not found", userFiler),
		)
	})
}
