package services

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type profileService struct {
	*app.Container
}

func (s *profileService) Get(session *app.Session) (*dtos.Profile, error) {
	var user models.User
	if err := s.DB.WithContext(session.Ctx).First(&user, session.UserId).Error; err != nil {
		return nil, err
	}

	var wallets []models.Wallet
	err := s.DB.WithContext(session.Ctx).
		Where("user_id = ?", session.UserId).
		Order("id").
		Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	balances := make([]dtos.ProfileWallet, 0, len(wallets))
	for _, wallet := range wallets {
		balances = append(balances, dtos.ProfileWallet{
			Id:       wallet.Id,
			Name:     wallet.Name,
			Currency: wallet.Currency,
			Current:  wallet.Current,
		})
	}

	return &dtos.Profile{
		Name:    user.Name,
		Email:   user.Email,
		Wallets: balances,
	}, nil
}
