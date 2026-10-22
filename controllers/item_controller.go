package controllers

import (
	"net/http"

	"github.com/linn221/expenses/models"
)

type ItemController struct {
	MyController
}

func (c *ItemController) Create() http.HandlerFunc {
	return CreateHandler[models.Item, models.NewItem](c.Container, c.Services.ItemService)
}

func (c *ItemController) Update() http.HandlerFunc {
	return UpdateHandler[models.Item, models.NewItem](c.Container, c.Services.ItemService)
}

func (c *ItemController) Delete() http.HandlerFunc {
	return DeleteHandler(c.Container, c.Services.ItemService)
}

func (c *ItemController) Get() http.HandlerFunc {
	return GetHandler(c.Services.ItemService)
}
