package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"gorm.io/gorm"
)

type expenseService struct {
	*app.Container
}

func (s *expenseService) Create(session *app.Session, input models.Expense) (*models.Expense, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	//get wallet
	var wallet models.Wallet
	err := session.SelectForUpdate(&wallet, "id = ?", input.WalletId)
	if err != nil {
		return nil, err
	}
	if wallet.Current.LessThan(input.ExpenseAmount) {
		return nil, utils.NewBadRequest("expense amount must not be greater than wallet's balance")
	}

	delta := input.ExpenseAmount.Neg()
	if err := incrementWallet(session, input.WalletId, delta); err != nil {
		return nil, err
	}

	updateItemPrice := app.RetrieveFromSession[bool](session, app.SessionValueUpdateItemPrice)
	if updateItemPrice {
		err := session.Tx.Exec(`UPDATE items SET price = ? WHERE id = ?`, input.ItemPrice, input.ItemId).Error
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
	//get wallet
	var wallet models.Wallet
	err := session.SelectForUpdate(&wallet, "id = ?", input.WalletId)
	if err != nil {
		return nil, err
	}
	delta := input.ExpenseAmount.Sub(existing.ExpenseAmount).Neg()
	// delta is what the wallet balance changes by: negative when the expense grows
	if delta.IsNegative() && wallet.Current.LessThan(delta.Abs()) {
		return nil, utils.NewBadRequest("expense amount must not be greater than wallet's balance")
	}

	if err := incrementWallet(session, input.WalletId, delta); err != nil {
		return nil, err
	}

	updateItemPrice := app.RetrieveFromSession[bool](session, app.SessionValueUpdateItemPrice)
	if updateItemPrice {
		err := session.Tx.Exec(`UPDATE items SET price = ? WHERE id = ?`, input.ItemPrice, input.ItemId).Error
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

	err := session.Tx.Exec(`UPDATE wallets SET current = current + ? WHERE id = ?`, record.ExpenseAmount.Abs(), record.WalletId).Error
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
	return selectResource[models.Expense](s.Container, session, "Wallet", "Item.Category")
}

func (s *expenseService) ApplyFilters(session *app.Session, filter *dtos.ExpenseFilter) *gorm.DB {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.Expense{}).Preload("Wallet").Preload("Item.Category")
	if filter == nil {
		return db
	}
	db = db.Where("wallet_id = ?", filter.WalletId)
	if filter.Search != "" {
		db = db.Where("note LIKE ?", "%"+filter.Search+"%")
	}
	if !filter.MinAmount.IsZero() {
		db = db.Where("expense_amount >= ?", filter.MinAmount)
	}
	if !filter.MaxAmount.IsZero() {
		db = db.Where("expense_amount <= ?", filter.MaxAmount)
	}
	if !filter.ExpenseDate.IsZero() {
		db = db.Where("expense_date = ?", filter.ExpenseDate.Time)
	}
	return db
}

func (s *expenseService) List(session *app.Session, filter *dtos.ExpenseFilter) ([]models.Expense, error) {
	var results []models.Expense
	err := s.ApplyFilters(session, filter).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *expenseService) HttpResponse() dtos.HttpResponseFromModel[models.Expense] {
	return dtos.RespondExpense
}

func (s *expenseService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.Expense] {
	return dtos.RespondExpenseCollection
}
