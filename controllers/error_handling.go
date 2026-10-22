package controllers

import (
	"errors"
	"net/http"

	"github.com/linn221/expenses/utils"
	"github.com/linn221/expenses/validate"
	"gorm.io/gorm"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		http.Error(w, "record not found", http.StatusNotFound)
	case errors.Is(err, validate.ErrClient):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, validate.ErrServer):
		http.Error(w, err.Error(), http.StatusInternalServerError)
	case errors.Is(err, utils.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, utils.ErrInternalServerError):
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
