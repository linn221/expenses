package dtos

import (
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
)

type NewCostPerspective struct {
	Name         string                          `json:"name" validate:"required,min=3,max=255"`
	Value        decimal.Decimal                 `json:"value" validate:"required"`
	Type         string                          `json:"type" validate:"required,oneof=percentage ratio rate_hour rate_day rate_month"`
	Applications []NewCostPerspectiveApplication `json:"applications" validate:"required,dive"`
}

type NewCostPerspectiveApplication struct {
	MoneyInsightCategory []models.MoneyInsightCategory `json:"money_insight_category" validate:"required,min=1,dive,oneof=expense income daily_expense"`
	StartingValue        decimal.Decimal               `json:"starting_value" validate:"required"`
}

func (input NewCostPerspective) ToModel(session *app.Session) models.CostPerspective {
	var applications []models.CostPerspectiveApplication
	for _, application := range input.Applications {
		for _, insight := range application.MoneyInsightCategory {
			applications = append(applications, models.CostPerspectiveApplication{
				MoneyInsightCategory: insight,
				StartingValue:        application.StartingValue,
			})
		}
	}

	model := models.CostPerspective{
		Name:         input.Name,
		Value:        input.Value,
		Type:         models.CostPerspectiveType(input.Type),
		Applications: applications,
	}
	model.UserId = session.UserId
	return model
}

type CostPerspectiveFilter struct {
	Name     string
	Type     string
	IsActive *bool
}

func (filter CostPerspectiveFilter) Parse(r *http.Request) (*CostPerspectiveFilter, error) {
	q := r.URL.Query()
	if name := q.Get("name"); name != "" {
		filter.Name = name
	}
	if perspectiveType := q.Get("type"); perspectiveType != "" {
		filter.Type = perspectiveType
	}
	if isActiveValue := q.Get("is_active"); isActiveValue != "" {
		isActive := isActiveValue == "1"
		filter.IsActive = &isActive
	}
	return &filter, nil
}

type costPerspectiveApplicationResponse struct {
	Id                   int                         `json:"id"`
	MoneyInsightCategory models.MoneyInsightCategory `json:"money_insight_category"`
	StartingValue        decimal.Decimal             `json:"starting_value"`
}

func respondApplications(applications []models.CostPerspectiveApplication) []costPerspectiveApplicationResponse {
	return addToCollection(applications, func(application models.CostPerspectiveApplication) costPerspectiveApplicationResponse {
		return costPerspectiveApplicationResponse{
			Id:                   application.Id,
			MoneyInsightCategory: application.MoneyInsightCategory,
			StartingValue:        application.StartingValue,
		}
	})
}

type CostPerspectiveResponse struct {
	Id           int                                  `json:"id"`
	Name         string                               `json:"name"`
	Value        decimal.Decimal                      `json:"value"`
	Type         models.CostPerspectiveType           `json:"type"`
	Applications []costPerspectiveApplicationResponse `json:"applications"`
}

func RespondCostPerspective(cp *models.CostPerspective) any {
	return &CostPerspectiveResponse{
		Id:           cp.Id,
		Name:         cp.Name,
		Value:        cp.Value,
		Type:         cp.Type,
		Applications: respondApplications(cp.Applications),
	}
}

type costPerspectiveCollectionInstance struct {
	Id           int                                  `json:"id"`
	Name         string                               `json:"name"`
	Value        decimal.Decimal                      `json:"value"`
	Type         models.CostPerspectiveType           `json:"type"`
	Applications []costPerspectiveApplicationResponse `json:"applications"`
}

type CostPerspectiveCollectionResponse struct {
	Results []costPerspectiveCollectionInstance `json:"results"`
}

func RespondCostPerspectiveCollection(results []models.CostPerspective) any {
	instances := addToCollection(results, func(cp models.CostPerspective) costPerspectiveCollectionInstance {
		return costPerspectiveCollectionInstance{
			Id:           cp.Id,
			Name:         cp.Name,
			Value:        cp.Value,
			Type:         cp.Type,
			Applications: respondApplications(cp.Applications),
		}
	})
	return &CostPerspectiveCollectionResponse{Results: instances}
}
