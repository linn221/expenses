package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

var orderColumnPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type pageBody struct {
	Data    any   `json:"data"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

func PaginateHandler[T any, F QueryParser[F]](service ApplyFilters[*F]) http.HandlerFunc {
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

		var filterZero F
		filter, err := filterZero.Parse(r)
		if err != nil {
			writeError(w, err)
			return
		}

		page, err := parsePagination[T](r)
		if err != nil {
			writeError(w, err)
			return
		}
		if collection, ok := service.(CollectionResponseGenerator[T]); ok {
			page.Collection = collection.HttpResponseCollection()
		}

		db := service.ApplyFilters(&session, filter)
		result, err := services.Paginate(db, page)
		if err != nil {
			writeError(w, err)
			return
		}

		body := pageBody{
			Data:    result.Data,
			Page:    result.Page,
			PerPage: result.PerPage,
			Total:   result.Total,
		}
		utils.Must(json.NewEncoder(w).Encode(body))
	}
}

func parsePagination[T any](r *http.Request) (services.Pagination[T], error) {
	q := r.URL.Query()
	page := 1
	if raw := q.Get("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return services.Pagination[T]{}, utils.NewBadRequest("page must be a positive integer")
		}
		page = n
	}

	column := q.Get("order")
	direction := strings.ToLower(q.Get("direction"))
	if column != "" && !orderColumnPattern.MatchString(column) {
		return services.Pagination[T]{}, utils.NewBadRequest("order is not a valid column")
	}
	if direction != "" && direction != "asc" && direction != "desc" {
		return services.Pagination[T]{}, utils.NewBadRequest("direction must be asc or desc")
	}
	if column == "" && direction != "" {
		return services.Pagination[T]{}, utils.NewBadRequest("order is required when direction is set")
	}

	return services.Pagination[T]{
		Page:      page,
		Column:    column,
		Direction: direction,
	}, nil
}
