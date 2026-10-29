package controllers

import (
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type costPerspectiveController struct {
	myController
}

func (c *costPerspectiveController) Create() http.HandlerFunc {
	return CreateHandler[models.CostPerspective, dtos.NewCostPerspective](c.Container, c.Services.CostPerspectiveService)
}

func (c *costPerspectiveController) Update() http.HandlerFunc {
	return UpdateHandler[models.CostPerspective, dtos.NewCostPerspective](c.Container, c.Services.CostPerspectiveService)
}

func (c *costPerspectiveController) Delete() http.HandlerFunc {
	return DeleteHandler[models.CostPerspective](c.Container, c.Services.CostPerspectiveService)
}

func (c *costPerspectiveController) Get() http.HandlerFunc {
	return GetHandler[models.CostPerspective](c.Services.CostPerspectiveService)
}

func (c *costPerspectiveController) List() http.HandlerFunc {
	return ListWithFilterHandler[models.CostPerspective, dtos.CostPerspectiveFilter](c.Services.CostPerspectiveService)
}

func (c *costPerspectiveController) ToggleIsActive() http.HandlerFunc {
	return ToggleIsActiveHandler(c.Services.CostPerspectiveService)
}
