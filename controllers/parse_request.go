package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/linn221/expenses/utils"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// parse the request body to json
func parseJson[T any](r *http.Request) (*T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return &v, nil
}

func WithInputParsing[T any](f func(w http.ResponseWriter, r *http.Request, input *T) error) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		input, err := parseJson[T](r)
		if err != nil {
			writeError(w, err)
			return
		}

		// validate the input's data shape
		var validationErrors []ValidationError
		err = myValidate.Struct(input)
		if err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				for _, verr := range ve {
					validationErrors = append(validationErrors, ValidationError{
						Field:   verr.Field(),
						Message: fmt.Sprintf("Field validation for '%s' failed on the '%s' constraint", verr.Field(), verr.Tag()),
					})
				}
			}

			w.WriteHeader(http.StatusBadRequest)
			utils.Must(json.NewEncoder(w).Encode(validationErrors))
			return
		}

		err = f(w, r, input)
		if err != nil {
			writeError(w, err)
			return
		}
	}
}
