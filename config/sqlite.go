package config

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectSQLite() *gorm.DB {
	dsn := os.Getenv("SQLITE_PATH")
	if dsn == "" {
		dsn = "expenses.db?_foreign_keys=on"
	}

	db, err := gorm.Open(sqlite.Open(dsn), initConfig())
	if err != nil {
		panic("Fail To Connect Database")
	}
	migrate(db)
	return db
}
