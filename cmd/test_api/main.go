package main

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
	"github.com/linn221/expenses/routing"
)

func main() {
	routing.Serve(newContainer())
}

func newContainer() *app.Container {
	db := config.ConnectSQLite()
	cache := config.NewInMemoryCache()
	return &app.Container{
		DB:      db,
		Cache:   cache,
		Locker:  cache,
		Counter: cache,
	}
}
