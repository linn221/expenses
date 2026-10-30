package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type incomeController struct {
	myController
}

func (c *incomeController) Create() http.HandlerFunc {
	return CreateHandler[models.Income, dtos.NewIncome](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Update() http.HandlerFunc {
	return UpdateHandler[models.Income, dtos.NewIncome](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Delete() http.HandlerFunc {
	return DeleteHandler[models.Income](c.Container, c.Services.IncomeService)
}

func (c *incomeController) Get() http.HandlerFunc {
	return GetHandler[models.Income](c.Services.IncomeService)
}

func (c *incomeController) List() http.HandlerFunc {
	return nil
}

func (c *incomeController) ListByWalletId() http.HandlerFunc {
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
		idValue := r.PathValue("walletId")
		if idValue == "" {
			http.Error(w, "resource id not found", http.StatusBadRequest)
			return
		}
		walletId, err := strconv.Atoi(idValue)
		if err != nil {
			http.Error(w, "resource id not found", http.StatusBadRequest)
			return
		}
		var filterValue dtos.IncomeFilter
		filter, err := filterValue.Parse(r)
		if err != nil {
			writeError(w, err)
			return
		}
		results, err := c.Services.IncomeService.List(&session, walletId, filter)
		if err != nil {
			writeError(w, err)
			return
		}
		resp := getCollectionResource(c.Services.IncomeService, results)
		json.NewEncoder(w).Encode(resp)
	}

}
