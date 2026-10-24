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
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
)

type Getter[T any] interface {
	Get(*app.Session) (*T, error)
}
type Lister[T, F any] interface {
	List(*app.Session, F) ([]T, error)
}
type Creater[T any] interface {
	Create(*app.Session, T) (*T, error)
}
type Updater[T any] interface {
	Getter[T]
	Update(*app.Session, *T, *T) (*T, error)
}
type Deleter[T any] interface {
	Getter[T]
	Delete(*app.Session, *T) (*T, error)
}

type QueryParser[T any] interface {
	Parse(*http.Request) (*T, error)
}

func CreateHandler[T any, I models.NewModel[T]](c *app.Container, service Creater[T]) http.HandlerFunc {
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

		// Convert input to model
		record := (*input).ToModel(&session)
		if validateInput, ok := any(input).(validate.Input); ok {
			err := validate.Validate(c, &session, validateInput, validate.ActionCreate)
			if err != nil {
				writeError(w, err)
				return
			}

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

func UpdateHandler[T any, I models.NewModel[T]](c *app.Container, service Updater[T]) http.HandlerFunc {
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

func DeleteHandler[T any](c *app.Container, service Deleter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		record, err := service.Get(&session)
		if err != nil {
			writeError(w, err)
			return
		}

		err = validateIfRulesExist(c, &session, record, validate.ActionDelete)
		if err != nil {
			writeError(w, err)
			return
		}

		// Call service to delete record
		_, err = service.Delete(&session, record)
		if err != nil {
			writeError(w, err)
			return
		}

		// Send response
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetHandler[T any](service Getter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// Call service to get record
		result, err := service.Get(&session)
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

func ListWithFilterHandler[T any, F QueryParser[F]](service Lister[T, *F]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		var filterZero F
		filter, err := filterZero.Parse(r)
		if err != nil {
			writeError(w, err)
			return
		}

		results, err := service.List(&session, filter)
		if err != nil {
			writeError(w, err)
			return
		}
		response := getCollectionResource(service, results)
		utils.Must(json.NewEncoder(w).Encode(response))
	}
}

func ListAllHandler[T any](w http.ResponseWriter, r *http.Request) {

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

	services.Default.AccessContainer(func(c *app.Container) {
		results, err := services.ListAllResource[T](c, &session)
		if err != nil {
			writeError(w, err)
			return
		}
		utils.Must(json.NewEncoder(w).Encode(results))
	})
}

func ToggleIsActiveHandler[T services.Togglable](getService Getter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// get resource id
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

		// get is_active status that the user wants to change
		value := r.URL.Query().Get("is_active")
		isActive := value == "1"

		model, err := getService.Get(&session)
		if err != nil {
			writeError(w, err)
			return
		}
		// nothing to change
		if (*model).GetIsActive() == isActive {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// run the validations only when we are sure the request will modify the model's is active status
		// validation if valiationRules exist on the resource model
		services.Default.AccessContainer(func(c *app.Container) {
			err = validateIfRulesExist(c, &session, *model, validate.ActionActiveStatusChange)
			if err != nil {
				writeError(w, err)
				return
			}
		})

		err = services.ToggleModel(&session, *model, isActive)
		if err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
