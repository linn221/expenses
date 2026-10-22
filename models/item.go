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

type NewItem struct {
	Name       string          `json:"name" validate:"required,min=3,max=255"`
	Code       string          `json:"code" validate:"required,min=3,max=255"`
	Price      decimal.Decimal `json:"price" validate:"required,gt=0"`
	CategoryId int             `json:"category_id" validate:"required"`
}

func (input NewItem) ToModel(session *app.Session) Item {
	return Item{
		HasUserId:  UserId(session),
		Name:       input.Name,
		Code:       input.Code,
		Price:      input.Price,
		CategoryId: input.CategoryId,
	}
}

func (existing *Item) UpdateValues(input *Item) map[string]any {
	return map[string]any{
		"Name":       input.Name,
		"Price":      input.Price,
		"Code":       input.Code,
		"CategoryId": input.CategoryId,
	}
}

func (input *NewItem) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		if !input.Price.GreaterThan(decimal.Zero) {
			af(validate.Failure("input price must be greater than zero"))
			return
		}

		userFiler := UserFilter(s)
		af(
			validate.UniqueRule("items", "code", input.Code, s.ResId, userFiler),
			validate.ExistsRule("items", input.CategoryId, "category not found", userFiler),
		)
	})
}
