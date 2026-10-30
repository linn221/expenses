package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

func ListAllHandler[T any](w http.ResponseWriter, r *http.Request) {

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

	services.Default.AccessContainer(func(c *app.Container) {
		results, err := services.ListAllResource[T](c, &session)
		if err != nil {
			writeError(w, err)
			return
		}
		utils.Must(json.NewEncoder(w).Encode(results))
	})
}
