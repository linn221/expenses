package models

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"gorm.io/gorm"
)

func UserFilter(s *app.Session) validate.FilterFunc {
	return func(d *gorm.DB) {
		d.Where("user_id = ?", s.UserId)
	}
}

func UserId(s *app.Session) HasUserId {
	return HasUserId{
		UserId: s.UserId,
	}
}
