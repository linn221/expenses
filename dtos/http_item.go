package dtos

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
)

type NewItem struct {
	Name       string          `json:"name" validate:"required,min=3,max=255"`
	Code       string          `json:"code" validate:"omitempty,min=3,max=255"`
	Price      decimal.Decimal `json:"price" validate:"required"`
	CategoryId int             `json:"category_id" validate:"required"`
}

func (input NewItem) ToModel(session *app.Session) models.Item {

	code := strings.ToLower(input.Name)[:min(6, len(input.Name))]
	item := models.Item{

		Name:       input.Name,
		Code:       code,
		Price:      input.Price,
		CategoryId: input.CategoryId,
	}
	item.UserId = session.UserId
	return item
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

type itemResponse struct {
	Id       int
	Name     string
	Code     string
	Price    decimal.Decimal
	Category string
}

func RespondItem(i *models.Item) any {
	return &itemResponse{
		Id:       i.Id,
		Name:     i.Name,
		Code:     i.Code,
		Price:    i.Price,
		Category: i.Category.Name,
	}
}

type itemCollectionInstance struct {
	Id       int
	Name     string
	Category string
}

type itemCollection struct {
	Items []itemCollectionInstance
}

func RespondItemCollection(items []models.Item) any {
	results := addToCollection(items, func(i models.Item) itemCollectionInstance {
		return itemCollectionInstance{
			Id:       i.Id,
			Name:     i.Name,
			Category: i.Category.Name,
		}
	})
	return &itemCollection{Items: results}
}
