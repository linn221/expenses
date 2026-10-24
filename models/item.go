package models

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/utils"
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
	Code       string          `json:"code" validate:"omitempty,min=3,max=255"`
	Price      decimal.Decimal `json:"price" validate:"required"`
	CategoryId int             `json:"category_id" validate:"required"`
}

func (input NewItem) ToModel(session *app.Session) Item {

	code := strings.ToLower(input.Name)[:min(6, len(input.Name))]
	return Item{
		HasUserId:  UserId(session),
		Name:       input.Name,
		Code:       code,
		Price:      input.Price,
		CategoryId: input.CategoryId,
	}
}

func (item *Item) Resource() any {
	return struct {
		Id       int
		Name     string
		Code     string
		Price    string
		Category string
	}{item.Id, item.Name, item.Code, item.Price.String(), item.Category.Name}
}

type itemCollectionInstance struct {
	Id       int
	Name     string
	Price    decimal.Decimal
	Category string
}

type ItemCollection struct {
	Items []itemCollectionInstance
}

func (collection ItemCollection) AddInstances(results []Item) *ItemCollection {
	addInstancesToCollection(results, func(m Item) itemCollectionInstance {
		return itemCollectionInstance{
			Id:       m.Id,
			Name:     m.Name,
			Price:    m.Price,
			Category: m.Category.Name,
		}
	})
	return &collection
}

type ItemFilter struct {
	Search string
	// Price      decimal.Decimal
	CategoryId int
}

func (f ItemFilter) Parse(r *http.Request) (*ItemFilter, error) {
	if search := r.URL.Query().Get("search"); search != "" {
		f.Search = search
	}

	// if priceq := r.URL.Query().Get("price"); priceq != "" {
	// 	price, err := decimal.NewFromString(priceq)
	// 	if err != nil {
	// 		return nil, utils.NewBadRequest("price is not a valid decimal")
	// 	}
	// 	f.Price = price
	// }

	if categoryIdq := r.URL.Query().Get("category_id"); categoryIdq != "" {
		categoryId, err := strconv.Atoi(categoryIdq)
		if err != nil {
			return nil, utils.NewBadRequest("category_id is not a valid int")
		}
		f.CategoryId = categoryId
	}
	return &f, nil
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
			validate.ExistsRule("categories", input.CategoryId, "category not found", userFiler),
		)
	})
}
