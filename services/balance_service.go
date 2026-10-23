package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
)

type balanceService struct {
	*app.Container
}

func incrementBalance(session *app.Session, balanceId int, amount decimal.Decimal) error {

	if session.Tx == nil {
		panic("not allowed to increase balance without a transaction")
	}

	//using an atomic expression
	raw := `UPDATE balances SET current = current + ? WHERE id = ?`
	err := session.Tx.Exec(raw, amount, balanceId).Error
	return err
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
