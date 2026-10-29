package models

import (
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"index"`
	Email    string `json:"email" gorm:"index"`
	Password string `json:"password"`
	HasIsActive
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) HttpResponse() any {
	return struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		u.Id, u.Name, u.Email,
	}
}

func (input *User) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		af(
			validate.UniqueRule("users", "email", input.Email, s.ResId, nil),
		)
	})
}
