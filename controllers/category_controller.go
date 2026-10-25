package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type categoryController struct {
	myController
}

func (c *categoryController) Create() http.HandlerFunc {
	return CreateHandler[models.Category, dtos.NewCategory](c.Container, c.Services.CategoryService)
}

func (c *categoryController) Update() http.HandlerFunc {
	return UpdateHandler[models.Category, dtos.NewCategory](c.Container, c.Services.CategoryService)
}

func (c *categoryController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Category](c.Container, c.Services.CategoryService)
}

func (c *categoryController) Get() http.HandlerFunc {
	return GetHandler[models.Category](c.Services.CategoryService)
}

func (c *categoryController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.Category, dtos.CategoryFilter](c.Services.CategoryService)
}

func (c *categoryController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.CategoryService)
}
