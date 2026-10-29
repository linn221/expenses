package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type itemService struct {
	*app.Container
}

func (s *itemService) Create(session *app.Session, input models.Item) (*models.Item, error) {
	return createModel(s.Container, session, input)
}

func (s *itemService) Update(session *app.Session, existing *models.Item, input *models.Item) (*models.Item, error) {
	return updateModel(s.Container, session, existing, input)
}

func (s *itemService) Get(session *app.Session) (*models.Item, error) {
	return selectResource[models.Item](s.Container, session, "Category")
}

func (s *itemService) Delete(session *app.Session, record *models.Item) (*models.Item, error) {
	return deleteModel(s.Container, session, record)
}

func (s *itemService) List(session *app.Session, filter *dtos.ItemFilter) ([]models.Item, error) {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Item{}).Preload("Category").Where("user_id = ?", session.UserId)

	if filter != nil {
		if filter.CategoryId > 0 {
			db.Where("category_id = ?", filter.CategoryId)
		}
		if filter.Search != "" {
			db.Where("name LIKE ? OR code LIKE ?",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
			)
		}
	}

	var results []models.Item
	err := db.Find(&results).Error
	return results, err
}

func (s *itemService) HttpResponse() dtos.HttpResponseFromModel[models.Item] {
	return dtos.RespondItem
}

func (s *itemService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.Item] {
	return dtos.RespondItemCollection
}
