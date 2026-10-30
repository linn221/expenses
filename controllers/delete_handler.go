package controllers

import (
	"net/http"
	"strconv"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
)

func DeleteHandler[T any](c *app.Container, service Deleter[T]) http.HandlerFunc {
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
