package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
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

	if err := incrementBalance(session, model.BalanceId, model.Amount); err != nil {
		return nil, err
	}

	return &model, session.Commit()
}

func (s *incomeService) Update(session *app.Session, existing *models.Income, input *models.Income) (*models.Income, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()

	delta := input.Amount.Sub(existing.Amount)
	result, err := updateModel(s.Container, session, existing, input)
	if err != nil {
		return nil, err
	}

	if err := incrementBalance(session, existing.BalanceId, delta); err != nil {
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

	if err := incrementBalance(session, record.BalanceId, subtracted.Neg()); err != nil {
		return nil, err
	}
	return result, session.Commit()
}

func (s *incomeService) Get(session *app.Session) (*models.Income, error) {
	return selectResource[models.Income](s.Container, session)
}

func (s *incomeService) List(session *app.Session, filter *models.IncomeFilter) ([]models.Income, error) {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Income{})
	if filter != nil {
		if filter.IncomeDate != (models.MyDate{}) {
			db = db.Where("income_date = ?", filter.IncomeDate)
		}
		if filter.Search != "" {
			db = db.Where("title LIKE ?", "%"+filter.Search+"%")
		}
		if filter.MaxAmount.GreaterThan(decimal.Zero) {
			db = db.Where("amount <= ?", filter.MaxAmount)
		}
		if filter.MinAmount.GreaterThan(decimal.Zero) {
			db = db.Where("amount >= ?", filter.MinAmount)
		}
	}
	var incomes []models.Income
	if err := db.Find(&incomes).Error; err != nil {
		return nil, err
	}
	return incomes, nil
}
