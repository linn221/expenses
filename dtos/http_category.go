package dtos

import (
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
)

type NewCategory struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type CategoryFilter struct {
	Name     string
	IsActive *bool
}

func (filter CategoryFilter) Parse(r *http.Request) (*CategoryFilter, error) {
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = name
	}
	if isActiveValue := r.URL.Query().Get("is_active"); isActiveValue != "" {
		isActive := isActiveValue == "1"
		filter.IsActive = &isActive
	}
	return &filter, nil
}

func (i NewCategory) ToModel(s *app.Session) models.Category {
	cat := models.Category{
		Name: i.Name,
	}
	cat.UserId = s.UserId

	return cat
}

type CategoryResponse struct {
	Id   int
	Name string
}

func RespondCategory(category *models.Category) any {

	return &CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

type categoryCollectionResponseInstance struct {
	Id   int
	Name string
}

type CategoryCollectionResponse struct {
	Results []categoryCollectionResponseInstance
}

func RespondCategoryCollection(categories []models.Category) any {
	results := addToCollection(categories, func(c models.Category) categoryCollectionResponseInstance {
		return categoryCollectionResponseInstance{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	return &CategoryCollectionResponse{Results: results}
}
