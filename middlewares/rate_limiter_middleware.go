package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/linn221/expenses/utils"
)

// using closure for dependencies
func NewRateLimiter(counter utils.CacheService, window time.Duration, limit int64, label string, getLimitBy func(r *http.Request) (string, error)) func(http.HandlerFunc) http.HandlerFunc {
	// Get the IP address or user identifier from the request.
	if label == "" {
		panic("label must not be an empty string")
	}

	return func(next http.HandlerFunc) http.HandlerFunc {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limitBy, err := getLimitBy(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			key := label + ":" + limitBy // using userId for convenience
			// Check if the key exists in the cache.
			exists, err := counter.Exists(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If the key doesn't exist, create it and set expiry.
			if exists == 0 {
				err := counter.SetValue(key, "1", window)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// If the key exists, get the current count.
			// number of requests within the current window
			count, err := counter.Incr(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// If the count exceeds the limit, return an error response.
			if count > limit {
				http.Error(w, fmt.Sprintf("Rate limit exceeded. Try again in %v seconds", window.Seconds()), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

}
