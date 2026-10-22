package models

import (
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	v "github.com/linn221/expenses/validate"
)

type Category struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"index"`
	HasIsActive
	HasUserId
}

type NewCategory struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type CategoryFilter struct {
	Name     string
	IsActive bool
}

func ParseCategoryFilter(r *http.Request) (filter CategoryFilter) {
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = name
	}
	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		filter.IsActive = isActive == "1"
	}
	return
}

func (i NewCategory) ToModel(s *app.Session) Category {
	cat := Category{
		Name:      i.Name,
		HasUserId: UserId(s),
	}

	return cat
}

func (existing *Category) UpdateValues(input *Category) map[string]any {
	m := make(map[string]any)
	if existing.Name != input.Name {
		m["Name"] = input.Name
	}
	return m
}

// will be invoked by service
// rules are generated based on state
func (i *NewCategory) ValidationRules(vc *v.ValidationContext) {

	vc.OnMutate(func(s *app.Session, appendFunc func(...validate.Rule)) {
		userFilter := UserFilter(s)
		appendFunc(
			validate.ExistsRule("categories", s.ResId, "category not found", userFilter).
				When(s.ResId > 0),
			validate.UniqueRule("categories", "name", i.Name, s.ResId, userFilter),
		)
	})
}

func (existing *Category) ValidationRules(vc *v.ValidationContext) {
	vc.OnDelete(func(s *app.Session, appendFunc func(...v.Rule)) {
		userFilter := UserFilter(s)
		appendFunc(validate.NoResultRule("items", "category is used in item", userFilter))
	})
	vc.OnActiveStatusChange(func(s *app.Session, af v.AppendFunc) {
		userFilter := UserFilter(s)
		af(validate.NoResultRule("items", "category is used in item", userFilter).When(existing.IsActive == false))
	})
}
