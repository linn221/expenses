package controllers

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/validate"
)

type myController struct {
	*app.Container
	Services *services.Provider
}

func validateIfRulesExist(c *app.Container, s *app.Session, inputModel any, action validate.Action) error {

	vInput, ok := inputModel.(validate.Input)
	if ok {
		return validate.Validate(c, s, vInput, action)
	}
	return nil
}
