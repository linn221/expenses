package app

import (
	"context"

	"github.com/linn221/expenses/utils"
	"gorm.io/gorm"
)

type Container struct {
	DB    *gorm.DB
	Cache utils.CacheService
}

type Session struct {
	UserId int
	Ctx    context.Context
	ResId  int
	Tx     *gorm.DB
}
