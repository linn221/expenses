package models

import (
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

func (existing *Category) UpdateValues(input *Category) map[string]any {
	m := make(map[string]any)
	if existing.Name != input.Name {
		m["Name"] = input.Name
	}
	return m
}

// will be invoked by service
// rules are generated based on state
func (i *Category) ValidationRules(vc *v.ValidationContext) {

	vc.OnMutate(func(s *app.Session, appendFunc func(...validate.Rule)) {
		userFilter := UserFilter(s)
		appendFunc(
			validate.ExistsRule("categories", s.ResId, "category not found", userFilter).
				When(s.ResId > 0),
			validate.UniqueRule("categories", "name", i.Name, s.ResId, userFilter),
		)
	})
	vc.OnDelete(func(s *app.Session, appendFunc func(...v.Rule)) {
		appendFunc(validate.NoResultRule("items", "category is used in item", validate.NewFilter("category_id = ?", s.ResId)))
	})
	vc.OnActiveStatusChange(func(s *app.Session, af v.AppendFunc) {
		af(validate.NoResultRule("items", "category is used in item", validate.NewFilter("category_id = ?", s.ResId)).When(i.IsActive == false))
	})
}
