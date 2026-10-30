package services

import (
	"strings"

	"github.com/linn221/expenses/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Pagination is the page number and the column to sort by, with asc or desc.
// Collection, when set, turns the page of models into the HTTP collection body.
type Pagination[T any] struct {
	Page       int
	PerPage    int
	Column     string
	Direction  string
	Collection func([]T) any
}

type PageResult[T any] struct {
	Results []T
	Data    any
	Total   int64
	Page    int
	PerPage int
}

// Paginate counts and loads one page from a query that already has its conditions.
func Paginate[T any](db *gorm.DB, page Pagination[T]) (PageResult[T], error) {
	if page.Page < 1 {
		page.Page = 1
	}
	if page.PerPage < 1 {
		page.PerPage = config.Cfg.PageSize()
	}

	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return PageResult[T]{}, err
	}

	query := db
	if page.Column != "" {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: page.Column},
			Desc:   strings.EqualFold(page.Direction, "desc"),
		})
	}

	results := make([]T, 0)
	err := query.Limit(page.PerPage).Offset((page.Page - 1) * page.PerPage).Find(&results).Error
	if err != nil {
		return PageResult[T]{}, err
	}

	data := any(results)
	if page.Collection != nil {
		data = page.Collection(results)
	}

	return PageResult[T]{
		Results: results,
		Data:    data,
		Total:   total,
		Page:    page.Page,
		PerPage: page.PerPage,
	}, nil
}
