package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
)

type balanceService struct {
	*app.Container
}

func (s *balanceService) Create(session *app.Session, model models.Balance) (*models.Balance, error) {
	return createModel(s.Container, session, model)
}

func (s *balanceService) Update(session *app.Session, existing *models.Balance, input *models.Balance) (*models.Balance, error) {
	return updateModel(s.Container, session, existing, input)
}

func (s *balanceService) Delete(session *app.Session, record *models.Balance) (*models.Balance, error) {
	return deleteModel(s.Container, session, record)
}

func (s *balanceService) Get(session *app.Session) (*models.Balance, error) {
	return selectResource[models.Balance](s.Container, session)
}
