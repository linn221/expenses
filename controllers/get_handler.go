package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/myctx"
	"github.com/linn221/expenses/utils"
)

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
