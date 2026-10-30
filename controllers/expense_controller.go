package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
)

type expenseController struct {
	myController
}

func (c *expenseController) Create() http.HandlerFunc {
	return CreateHandler[models.Expense, dtos.NewExpense](c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Update() http.HandlerFunc {
	return UpdateHandler[models.Expense, dtos.NewExpense](c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Delete() http.HandlerFunc {
	return DeleteHandler(c.Container, c.Services.ExpenseService)
}

func (c *expenseController) Get() http.HandlerFunc {
	return GetHandler(c.Services.ExpenseService)
}

func (c *expenseController) List() http.HandlerFunc {
	return nil
}

func (c *expenseController) ListByWalletId() http.HandlerFunc {
	// return c.Services.ExpenseService.

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
		var filter dtos.ExpenseFilter
		f, err := filter.Parse(r)
		if err != nil {
			writeError(w, err)
			return
		}
		results, err := c.Services.ExpenseService.List(&session, walletId, f)
		if err != nil {
			writeError(w, err)
			return
		}

		response := getCollectionResource(c.Services.ExpenseService, results)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			writeError(w, err)
			return
		}
	}
}
