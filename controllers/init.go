package controllers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var myValidate *validator.Validate

func init() {
	myValidate = validator.New()
	myValidate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})
}
