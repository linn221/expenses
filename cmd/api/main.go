package main

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
	"github.com/linn221/expenses/routing"
)

func main() {
	db := config.ConnectDB()
	cache := config.NewInMemoryCache()
	a := app.Container{
		DB:    db,
		Cache: cache,
	}
	routing.Serve(&a)
}
