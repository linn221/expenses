package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/myctx"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
)

func UpdateHandler[T any, I models.ModelConverter[T]](c *app.Container, service Updater[T]) http.HandlerFunc {
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
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "resource id not found", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(idValue)
		if err != nil {
			http.Error(w, "resource id not found", http.StatusBadRequest)
			return
		}
		session.ResId = resId

		// Call service to get existing record
		existing, err := service.Get(&session)
		if err != nil {
			writeError(w, err)
			return
		}

		// Convert input to model
		inputModel := (*input).ToModel(&session)
		// Validate if rules exist
		err = validateIfRulesExist(c, &session, input, validate.ActionUpdate)
		if err != nil {
			writeError(w, err)
			return
		}

		// Call service to update record
		result, err := service.Update(&session, existing, &inputModel)
		if err != nil {
			writeError(w, err)
			return
		}

		// Get presentation for http response
		res := getResource(service, result)

		// Send response
		w.WriteHeader(http.StatusOK)
		utils.Must(json.NewEncoder(w).Encode(res))
	}
}
