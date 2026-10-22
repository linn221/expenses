package config

import "gorm.io/gorm"

type Dbh struct {
	db *gorm.DB
}
