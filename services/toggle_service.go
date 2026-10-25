package services

import (
	"github.com/linn221/expenses/app"
)

func ToggleModel[T Togglable](session *app.Session, model T, activeStatus bool) error {
	var err error
	Default.AccessContainer(func(c *app.Container) {
		// change status
		err = c.DB.WithContext(session.Ctx).Model(model).Updates(map[string]any{
			"is_active": activeStatus,
		}).Error
	})
	return err
}
