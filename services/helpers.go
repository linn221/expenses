package services

import (
	"github.com/linn221/expenses/app"
	"gorm.io/gorm"
)

func startTransaction(c *app.Container, s *app.Session) {
	tx := c.DB.WithContext(s.Ctx).Begin()
	s.Tx = tx
}

func transaction(c *app.Container, s *app.Session, f func(tx *gorm.DB) error) error {
	return c.DB.WithContext(s.Ctx).Transaction(f)
}

func selectResource[T any](c *app.Container, s *app.Session, preloads ...string) (*T, error) {
	var v T
	dbctx := c.DB.WithContext(s.Ctx).Where("user_id = ?", s.UserId)
	for _, p := range preloads {
		dbctx = dbctx.Preload(p)
	}
	err := dbctx.First(&v, s.ResId).Error
	return &v, err
}

func listResource[T any](c *app.Container, s *app.Session) ([]T, error) {
	var v []T
	err := c.DB.WithContext(s.Ctx).Where("user_id = ?", s.UserId).Find(&v).Error
	return v, err
}

func ListAllResource[T any](c *app.Container, s *app.Session) ([]T, error) {
	var results []T
	err := c.DB.WithContext(s.Ctx).Where("user_id = ?", s.UserId).Find(&results).Error
	return results, err
}

func createModel[T any](c *app.Container, s *app.Session, model T) (*T, error) {
	var err error
	if s.Tx == nil {
		err = c.DB.WithContext(s.Ctx).Create(&model).Error
	} else {
		err = s.Tx.Create(&model).Error
	}
	return &model, err
}

func updateModel[T UpdatableModel[T]](c *app.Container, s *app.Session, existing T, input T) (T, error) {
	updates := existing.UpdateValues(input)
	var err error
	if s.Tx == nil {
		err = c.DB.WithContext(s.Ctx).Model(existing).Updates(updates).Error
	} else {
		err = s.Tx.Model(existing).Updates(updates).Error
	}
	return existing, err
}

func deleteModel[T any](c *app.Container, s *app.Session, model T) (T, error) {
	var err error
	if s.Tx == nil {
		err = c.DB.WithContext(s.Ctx).Delete(&model).Error
	} else {
		err = s.Tx.Delete(&model).Error
	}
	return model, err
}
