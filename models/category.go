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
	IsActive *bool
}

// CategoryCollection for http presentation
type CategoryCollection struct {
	Instances []categoryCollectionInstance
}

type categoryCollectionInstance struct {
	Id   int
	Name string
}

func addInstancesToCollection[Model, Instance any](results []Model, toCollectionInstance func(m Model) Instance) []Instance {
	instances := make([]Instance, 0, len(results))
	for _, r := range results {
		instances = append(instances, toCollectionInstance(r))
	}

	return instances
}

func (list CategoryCollection) AddInstances(results []Category) *CategoryCollection {
	instances := addInstancesToCollection(results, func(m Category) categoryCollectionInstance {
		return categoryCollectionInstance{
			Id:   m.Id,
			Name: m.Name,
		}
	})
	list.Instances = instances
	return &list
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
		appendFunc(validate.NoResultRule("items", "category is used in item", validate.NewFilter("category_id = ?", s.ResId)))
	})
	vc.OnActiveStatusChange(func(s *app.Session, af v.AppendFunc) {
		af(validate.NoResultRule("items", "category is used in item", validate.NewFilter("category_id = ?", s.ResId)).When(existing.IsActive == false))
	})
}
