package controllers

import (
	"net/http"
	"strconv"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/validate"
)

func ToggleIsActiveHandler[T services.Togglable](getService Getter[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Construct user session
		userId, err := app.GetUserId(r.Context())
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
