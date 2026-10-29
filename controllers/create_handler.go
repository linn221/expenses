package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/myctx"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
)

func CreateHandler[T any, I models.ModelConverter[T]](c *app.Container, service Creater[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Parsing request body
		input, err := parseJson[I](r)
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
			// w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{"errors": validationErrors})
			return
		}

		// Construct user session
		userId, err := myctx.GetUserId(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		session := app.Session{
			UserId: userId,
			Ctx:    r.Context(),
		}

		// Convert input to model. Rules live on the model, not the request DTO.
		record := (*input).ToModel(&session)
		err = validateIfRulesExist(c, &session, &record, validate.ActionCreate)
		if err != nil {
			writeError(w, err)
			return
		}

		// Call service
		result, err := service.Create(&session, record)
		if err != nil {
			writeError(w, err)
			return
		}

		// Get presentation for http response
		res := getResource(service, result)

		// Send response
		w.WriteHeader(http.StatusCreated)
		utils.Must(json.NewEncoder(w).Encode(res))
	}
}
