package services

import (
	"github.com/linn221/expenses/app"
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

func (s *categoryService) List(session *app.Session, filter models.CategoryFilter) ([]models.Category, error) {
	return listResource[models.Category](s.Container, session)
}

// optional
func (s *categoryService) Resource(result *models.Category) (any, error) {

	return struct {
		Id   int
		Name string
	}{Id: result.Id,
		Name: result.Name}, nil

}
