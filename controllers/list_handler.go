package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/utils"
)

func ListWithFilterHandler[T any, F QueryParser[F]](service Lister[T, *F]) http.HandlerFunc {
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
