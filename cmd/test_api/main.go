package main

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
	"github.com/linn221/expenses/routing"
)

func main() {
	db := config.ConnectSQLite()
	cache := config.NewInMemoryCache()
	routing.Serve(&app.Container{
		DB:      db,
		Cache:   cache,
		Locker:  cache,
		Counter: cache,
	})
}
