package app

import (
	"github.com/linn221/expenses/utils"
	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Cache  utils.CacheService
	Locker utils.Locker
}
