package controllers

import (
	"net/http"

	"github.com/linn221/expenses/models"
)

type CategoryController struct {
	MyController
}

func (c *CategoryController) Create() http.HandlerFunc {
	return CreateHandler[models.Category, models.NewCategory](c.Container, c.Services.CategoryService)
}

func (c *CategoryController) Update() http.HandlerFunc {
	return UpdateHandler[models.Category, models.NewCategory](c.Container, c.Services.CategoryService)
}

func (c *CategoryController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Category](c.Container, c.Services.CategoryService)
}

func (c *CategoryController) Get() http.HandlerFunc {
	return GetHandler[models.Category](c.Services.CategoryService)
}

func (c *CategoryController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.Category, models.CategoryFilter](c.Services.CategoryService)
}

func (c *CategoryController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.CategoryService)
}
