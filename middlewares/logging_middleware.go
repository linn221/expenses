package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/linn221/expenses/utils"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := newResponseWriter(w)
		clientIP := utils.GetClientIP(r)

		next.ServeHTTP(wrappedWriter, r)

		log.Printf("[%s] %s %s %s %d %v",
			clientIP,
			r.Method,
			r.URL.Path,
			r.Proto,
			wrappedWriter.StatusCode(),
			time.Since(start),
		)
	})
}
