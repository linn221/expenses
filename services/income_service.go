package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"gorm.io/gorm"
)

type incomeService struct {
	*app.Container
}

func (s *incomeService) Create(session *app.Session, model models.Income) (*models.Income, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	if err := session.Tx.Create(&model).Error; err != nil {
		return nil, err
	}

	if err := incrementWallet(session, model.WalletId, model.Amount); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	session.ResId = model.Id
	return s.Get(session)
}

func (s *incomeService) Update(session *app.Session, existing *models.Income, input *models.Income) (*models.Income, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	delta := input.Amount.Sub(existing.Amount)
	result, err := updateModel(s.Container, session, existing, input)
	if err != nil {
		return nil, err
	}

	if err := incrementWallet(session, existing.WalletId, delta); err != nil {
		return nil, err
	}

	return result, session.Commit()
}

func (s *incomeService) Delete(session *app.Session, record *models.Income) (*models.Income, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	subtracted := record.Amount
	result, err := deleteModel(s.Container, session, record)
	if err != nil {
		return nil, err
	}

	if err := incrementWallet(session, record.WalletId, subtracted.Neg()); err != nil {
		return nil, err
	}
	return result, session.Commit()
}

func (s *incomeService) Get(session *app.Session) (*models.Income, error) {
	return selectResource[models.Income](s.Container, session, "Wallet")
}

func (s *incomeService) ApplyFilters(session *app.Session, filter *dtos.IncomeFilter) *gorm.DB {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Income{})
	if filter == nil {
		return db
	}
	db = db.Where("wallet_id = ?", filter.WalletId)
	if filter.IncomeDate != (models.MyDate{}) {
		db = db.Where("income_date = ?", filter.IncomeDate)
	}
	if filter.Search != "" {
		db = db.Where("title LIKE ?", "%"+filter.Search+"%")
	}
	if !filter.MaxAmount.IsZero() {
		db = db.Where("amount <= ?", filter.MaxAmount)
	}
	if !filter.MinAmount.IsZero() {
		db = db.Where("amount >= ?", filter.MinAmount)
	}
	return db
}

func (s *incomeService) List(session *app.Session, filter *dtos.IncomeFilter) ([]models.Income, error) {
	var incomes []models.Income
	if err := s.ApplyFilters(session, filter).Find(&incomes).Error; err != nil {
		return nil, err
	}
	return incomes, nil
}

func (s *incomeService) HttpResponse() dtos.HttpResponseFromModel[models.Income] {
	return dtos.RespondIncome
}

func (s *incomeService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.Income] {
	return dtos.RespondIncomeCollection
}
