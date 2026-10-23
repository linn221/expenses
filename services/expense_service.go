package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
)

type expenseService struct {
	*app.Container
}

func (s *expenseService) Create(session *app.Session, input models.Expense) (*models.Expense, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	//get balance
	var balance models.Balance
	err := session.SelectForUpdate(&balance, "id = ?", input.BalanceId)
	if err != nil {
		return nil, err
	}
	if balance.Current.LessThan(input.Amount) {
		return nil, utils.NewBadRequest("expense amount must not be greater than wallet's balance")
	}

	delta := input.Amount.Neg()
	if err := incrementBalance(session, input.BalanceId, delta); err != nil {
		return nil, err
	}

	updateItemPrice := app.RetrieveFromSession[bool](session, app.SessionValueUpdateItemPrice)
	if updateItemPrice {
		err := session.Tx.Exec(`UPDATE items SET item_price = ? WHERE id = ?`, input.ItemPrice, input.ItemId).Error
		if err != nil {
			return nil, err
		}
	}

	result, err := createModel(s.Container, session, input)
	if err != nil {
		return nil, err
	}

	return result, session.Commit()
}

func (s *expenseService) Update(session *app.Session, existing *models.Expense, input *models.Expense) (*models.Expense, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	//get balance
	var balance models.Balance
	err := session.SelectForUpdate(&balance, "id = ?", input.BalanceId)
	if err != nil {
		return nil, err
	}
	delta := input.Amount.Sub(existing.Amount).Neg()
	if balance.Current.LessThan(delta.Abs()) {
		return nil, utils.NewBadRequest("expense amount must not be greater than wallet's balance")
	}

	if err := incrementBalance(session, input.BalanceId, delta); err != nil {
		return nil, err
	}

	updateItemPrice := app.RetrieveFromSession[bool](session, app.SessionValueUpdateItemPrice)
	if updateItemPrice {
		err := session.Tx.Exec(`UPDATE items SET item_price = ? WHERE id = ?`, input.ItemPrice, input.ItemId).Error
		if err != nil {
			return nil, err
		}
	}

	result, err := updateModel(s.Container, session, existing, input)
	if err != nil {
		return nil, err
	}

	return result, session.Commit()
}

func (s *expenseService) Delete(session *app.Session, record *models.Expense) (*models.Expense, error) {

	session.StartTransaction(s.Container)
	defer session.Rollback()

	err := session.Tx.Exec(`UPDATE balances SET current = current + ? WHERE id = ?`, record.Amount.Abs(), record.BalanceId).Error
	if err != nil {
		return nil, err
	}

	result, err := deleteModel(s.Container, session, record)
	if err != nil {
		return nil, err
	}
	return result, session.Commit()
}

func (s *expenseService) Get(session *app.Session) (*models.Expense, error) {
	return selectResource[models.Expense](s.Container, session)
}

func (s *expenseService) List(session *app.Session, filter *models.ExpenseFilter) ([]models.Expense, error) {
	var results []models.Expense
	return results, nil
}
