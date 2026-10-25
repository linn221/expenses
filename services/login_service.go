package services

import (
	"context"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
	"golang.org/x/crypto/bcrypt"
)

type loginService struct {
	*app.Container
}

func (s *loginService) Login(ctx context.Context, input *dtos.LoginInfo) (*models.User, error) {
	var user models.User
	if err := s.DB.WithContext(ctx).Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, utils.NewBadRequest("email or password do not match")
	}
	err := utils.ComparePassword(user.Password, input.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, utils.NewBadRequest("email or password do not match")
	}
	if err != nil {
		return nil, utils.NewInternalServerError(err.Error())
	}
	return &user, nil
}

func (s *loginService) Register(ctx context.Context, input *dtos.NewUser) (*models.User, error) {
	bs, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(bs),
	}

	err = validate.Validate(s.Container, &app.Session{Ctx: ctx}, &user, validate.ActionCreate)
	if err != nil {
		return nil, err
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	err = tx.Create(&user).Error
	if err != nil {
		return nil, err
	}
	defaultWallet := models.Wallet{
		Name:     "Wallet",
		Currency: "MMK",
	}
	defaultWallet.UserId = user.Id
	err = tx.Create(&defaultWallet).Error
	if err != nil {
		return nil, err
	}
	return &user, tx.Commit().Error
}
