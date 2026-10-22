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

func (u *User) Resource() any {
	return struct {
		Id    int `json:"id"`
		Name  string
		Email string
	}{
		u.Id, u.Name, u.Email,
	}
}

type NewUser struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (input *NewUser) ValidationRules(vc *validate.ValidationContext) {
	vc.OnMutate(func(s *app.Session, af validate.AppendFunc) {
		af(
			validate.UniqueRule("users", "email", input.Email, s.ResId, nil),
		)
	})
}

// func (input *NewUser) Validate(db *gorm.DB, id int) error {

// }

type LoginInfo struct {
	Email    string
	Password string
}
