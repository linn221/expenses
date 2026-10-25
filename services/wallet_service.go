package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
)

type walletService struct {
	*app.Container
}

func incrementWallet(session *app.Session, walletId int, amount decimal.Decimal) error {

	if session.Tx == nil {
		panic("not allowed to increase wallet without a transaction")
	}

	//using an atomic expression
	raw := `UPDATE wallets SET current = current + ? WHERE id = ?`
	err := session.Tx.Exec(raw, amount, walletId).Error
	return err
}

func (s *walletService) Create(session *app.Session, model models.Wallet) (*models.Wallet, error) {
	return createModel(s.Container, session, model)
}

func (s *walletService) Update(session *app.Session, existing *models.Wallet, input *models.Wallet) (*models.Wallet, error) {
	return updateModel(s.Container, session, existing, input)
}

func (s *walletService) Delete(session *app.Session, record *models.Wallet) (*models.Wallet, error) {
	return deleteModel(s.Container, session, record)
}

func (s *walletService) Get(session *app.Session) (*models.Wallet, error) {
	return selectResource[models.Wallet](s.Container, session)
}
