package services

import (
	"github.com/linn221/expenses/app"
)

type Togglable interface {
	GetIsActive() bool
}

func ToggleModel[T Togglable](session *app.Session, model T, activeStatus bool) error {
	c := Default.Container
	// change status
	err := c.DB.WithContext(session.Ctx).Model(model).Updates(map[string]any{
		"is_active": activeStatus,
	}).Error
	return err
}
