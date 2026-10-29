package config

import (
	"github.com/linn221/expenses/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Wallet{},
		&models.Item{},
		&models.Income{},
		&models.Expense{},
		&models.CostPerspective{},
		&models.CostPerspectiveApplication{},
	)
	if err != nil {
		panic(err)
	}
}
