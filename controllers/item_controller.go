package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type itemController struct {
	myController
}

func (c *itemController) Create() http.HandlerFunc {
	return CreateHandler[models.Item, dtos.NewItem](c.Container, c.Services.ItemService)
}

func (c *itemController) Update() http.HandlerFunc {
	return UpdateHandler[models.Item, dtos.NewItem](c.Container, c.Services.ItemService)
}

func (c *itemController) Delete() http.HandlerFunc {
	return DeleteHandler(c.Container, c.Services.ItemService)
}

func (c *itemController) Get() http.HandlerFunc {
	return GetHandler(c.Services.ItemService)
}

func (c *itemController) List() http.HandlerFunc {
	return ListWithFilterHandler(c.Services.ItemService)
}

func (c *itemController) Paginate() http.HandlerFunc {
	return PaginateHandler[models.Item, dtos.ItemFilter](c.Services.ItemService)
}

func (c *itemController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.ItemService)
}
