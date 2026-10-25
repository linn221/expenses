package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

// when a resource reaches a service, it is assumed by default that all validations have been passed
// what you need is a straight crud
type categoryService struct {
	*app.Container
}

func (s *categoryService) Create(session *app.Session, model models.Category) (*models.Category, error) {
	return createModel(s.Container, session, model)
}

func (s *categoryService) Update(session *app.Session, existing *models.Category, input *models.Category) (*models.Category, error) {

	return updateModel(s.Container, session, existing, input)
}

func (s *categoryService) Delete(session *app.Session, record *models.Category) (*models.Category, error) {
	return deleteModel(s.Container, session, record)
}

func (s *categoryService) Get(session *app.Session) (*models.Category, error) {
	return selectResource[models.Category](s.Container, session)
}

func (s *categoryService) List(session *app.Session, filter *dtos.CategoryFilter) ([]models.Category, error) {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Category{}).Where("user_id = ?", session.UserId)
	if filter.Name != "" {
		db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.IsActive != nil {
		db.Where("is_active = ?", filter.IsActive)
	}

	var results []models.Category
	err := db.Find(&results).Error
	return results, err
}

// optional
func (s *categoryService) HttpResponse() dtos.HttpResponseFromModel[models.Category] {
	return dtos.RespondCategory
}

func (s *categoryService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.Category] {
	return dtos.RespondCategoryCollection
}

// func (s *categoryService) HttpResponseCollection(results []models.Category) any {
// 	var res dtos.Cate
// 	return
// }
// func (s *categoryService) CollectionResource() *models.CategoryCollection {
// 	return &models.CategoryCollection{}
// }
