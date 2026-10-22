package utils

import "gorm.io/gorm"

func First[T any](tx *gorm.DB, id any) (*T, error) {
	var result T
	err := tx.First(&result, id).Error
	return &result, err
}
