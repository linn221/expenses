package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"gorm.io/gorm"
)

type itemService struct {
	*app.Container
}

func (s *itemService) Create(session *app.Session, input models.Item) (*models.Item, error) {
	return createModel(s.Container, session, input)
}

func (s *itemService) Update(session *app.Session, existing *models.Item, input *models.Item) (*models.Item, error) {
	if _, err := updateModel(s.Container, session, existing, input); err != nil {
		return nil, err
	}
	return s.Get(session)
}

func (s *itemService) Get(session *app.Session) (*models.Item, error) {
	return selectResource[models.Item](s.Container, session, "Category")
}

func (s *itemService) Delete(session *app.Session, record *models.Item) (*models.Item, error) {
	return deleteModel(s.Container, session, record)
}

func (s *itemService) ApplyFilters(session *app.Session, filter *dtos.ItemFilter) *gorm.DB {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Item{}).Preload("Category").Where("user_id = ?", session.UserId)
	if filter == nil {
		return db
	}
	if filter.CategoryId > 0 {
		db = db.Where("category_id = ?", filter.CategoryId)
	}
	if filter.Search != "" {
		db = db.Where("name LIKE ? OR code LIKE ?",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
	}
	return db
}

func (s *itemService) List(session *app.Session, filter *dtos.ItemFilter) ([]models.Item, error) {
	var results []models.Item
	err := s.ApplyFilters(session, filter).Find(&results).Error
	return results, err
}

func (s *itemService) HttpResponse() dtos.HttpResponseFromModel[models.Item] {
	return dtos.RespondItem
}

func (s *itemService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.Item] {
	return dtos.RespondItemCollection
}
