package controllers

import "github.com/go-playground/validator/v10"

var myValidate *validator.Validate

func init() {
	myValidate = validator.New()
}
