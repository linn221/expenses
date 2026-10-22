package middlewares

import (
	"net/http"

	"github.com/linn221/expenses/myctx"
	"github.com/linn221/expenses/utils"
)

// func Auth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := myctx.GetUserId(r.Context())
// 		if err != nil {
// 			utils.RemoveCookies(w, "session_id")
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func Auth() func(http.HandlerFunc) http.HandlerFunc {
	return func(hf http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			_, err := myctx.GetUserId(r.Context())
			if err != nil {
				utils.RemoveCookies(w, "session_id")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			hf.ServeHTTP(w, r)
		}
	}
}
