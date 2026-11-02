package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"gorm.io/gorm"
)

type walletTransferService struct {
	*app.Container
}

func (s *walletTransferService) Create(session *app.Session, input models.WalletTransfer) (*models.WalletTransfer, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	var fromWallet, toWallet models.Wallet
	err := session.SelectForUpdate(&fromWallet, "id = ?", input.FromWalletId)
	if err != nil {
		return nil, err
	}
	if err := session.SelectForUpdate(&toWallet, "id = ?", input.ToWalletId); err != nil {
		return nil, err
	}

	if input.TransferAmount.GreaterThan(fromWallet.Current) {
		return nil, utils.NewBadRequest("source wallet does not have enough balance")
	}

	if err := incrementWallet(session, fromWallet.Id, input.TransferAmount.Neg()); err != nil {
		return nil, err
	}
	if err := incrementWallet(session, toWallet.Id, input.TransferAmount); err != nil {
		return nil, err
	}

	created, err := createModel(s.Container, session, input)
	if err != nil {
		return nil, err
	}

	return created, session.Commit()
}

func (s *walletTransferService) Update(session *app.Session, existing *models.WalletTransfer, input *models.WalletTransfer) (*models.WalletTransfer, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	var fromWallet, toWallet models.Wallet

	// from/to wallet cannot be changed
	err := session.SelectForUpdate(&fromWallet, "id = ?", existing.FromWalletId)
	if err != nil {
		return nil, err
	}
	if err := session.SelectForUpdate(&toWallet, "id = ?", existing.ToWalletId); err != nil {
		return nil, err
	}

	delta := input.TransferAmount.Sub(existing.TransferAmount)
	if delta.IsPositive() {
		if delta.GreaterThan(fromWallet.Current) {
			return nil, utils.NewBadRequest("source wallet does not have enough balance")
		}
	} else {
		if delta.GreaterThan(toWallet.Current) {
			return nil, utils.NewBadRequest("destination wallet does not have enough balance")
		}
	}

	if err := incrementWallet(session, fromWallet.Id, delta.Neg()); err != nil {
		return nil, err
	}
	if err := incrementWallet(session, toWallet.Id, delta); err != nil {
		return nil, err
	}

	updated, err := updateModel(s.Container, session, existing, input)
	if err != nil {
		return nil, err
	}

	return updated, session.Commit()
}

func (s *walletTransferService) Delete(session *app.Session, record *models.WalletTransfer) (*models.WalletTransfer, error) {
	session.StartTransaction(s.Container)
	defer session.Rollback()
	var fromWallet, toWallet models.Wallet

	err := session.SelectForUpdate(&fromWallet, "id = ?", record.FromWalletId)
	if err != nil {
		return nil, err
	}
	if err := session.SelectForUpdate(&toWallet, "id = ?", record.ToWalletId); err != nil {
		return nil, err
	}

	if record.TransferAmount.GreaterThan(toWallet.Current) {
		return nil, utils.NewBadRequest("destination wallet does not have enough balance")
	}

	if err := incrementWallet(session, fromWallet.Id, record.TransferAmount); err != nil {
		return nil, err
	}
	if err := incrementWallet(session, toWallet.Id, record.TransferAmount.Neg()); err != nil {
		return nil, err
	}

	deleted, err := deleteModel(s.Container, session, record)
	if err != nil {
		return nil, err
	}

	return deleted, session.Commit()

}

func (s *walletTransferService) Get(session *app.Session) (*models.WalletTransfer, error) {
	return selectResource[models.WalletTransfer](s.Container, session, "FromWallet", "ToWallet")
}

func (s *walletTransferService) ApplyFilters(session *app.Session, filter *dtos.WalletTransferFilter) *gorm.DB {
	db := s.Container.DB.WithContext(session.Ctx).Model(&models.WalletTransfer{}).
		Preload("FromWallet").Preload("ToWallet").
		Where("user_id = ?", session.UserId)
	if filter == nil {
		return db
	}
	if filter.Search != "" {
		db = db.Where("description LIKE ?", "%"+filter.Search+"%")
	}
	if !filter.TransferDate.IsZero() {
		db = db.Where("transfer_date = ?", filter.TransferDate)
	}
	if filter.FromWalletId != 0 {
		db = db.Where("from_wallet_id = ?", filter.FromWalletId)
	}
	if filter.ToWalletId != 0 {
		db = db.Where("to_wallet_id = ?", filter.ToWalletId)
	}
	return db
}

func (s *walletTransferService) List(session *app.Session, filter *dtos.WalletTransferFilter) ([]models.WalletTransfer, error) {
	var results []models.WalletTransfer
	err := s.ApplyFilters(session, filter).Find(&results).Error
	return results, err
}

func (s *walletTransferService) HttpResponse() dtos.HttpResponseFromModel[models.WalletTransfer] {
	return dtos.RespondWalletTransfer
}

func (s *walletTransferService) HttpResponseCollection() dtos.HttpResponseFromModelCollection[models.WalletTransfer] {
	return dtos.RespondWalletTransferCollection
}
