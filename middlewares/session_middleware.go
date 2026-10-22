package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/linn221/expenses/myctx"
	"github.com/linn221/expenses/utils"
)

func Session(sessionStore utils.CacheService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies, err := r.Cookie("session_id")
			if err != nil && err != http.ErrNoCookie {
				http.Error(w, err.Error(), 400)
				return
			}
			if cookies != nil && cookies.Value != "" {
				userId, err := lookUpUser(sessionStore, cookies.Value)
				if err != nil {
					utils.RemoveCookies(w, "session_id")
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := r.Context()
				ctx = myctx.SetUserId(ctx, userId)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func lookUpUser(store utils.CacheService, sessionId string) (int, error) {
	value, ok, err := store.GetValue("Session:" + sessionId)
	if err != nil || !ok {
		return 0, errors.New("user does not exist")
	}
	userId, err := strconv.Atoi(value)
	if err != nil || userId <= 0 {
		return 0, errors.New("user does not exist")
	}

	return userId, nil
}
