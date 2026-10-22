package services

import (
	"github.com/linn221/expenses/app"
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
	return selectResource[models.Item](s.Container, session)
}

func (s *itemService) Delete(session *app.Session, record *models.Item) (*models.Item, error) {
	return deleteModel(s.Container, session, record)
}
