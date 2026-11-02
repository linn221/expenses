package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/utils"
)

type profileController struct {
	myController
}

func (c *profileController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := app.GetUserId(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		session := app.Session{
			UserId: userId,
			Ctx:    r.Context(),
		}

		profile, err := c.Services.ProfileService.Get(&session)
		if err != nil {
			writeError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		utils.Must(json.NewEncoder(w).Encode(profile))
	}
}
