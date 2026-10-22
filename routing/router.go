package routing

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

// chain composes middlewares so the first one runs first.
// An empty list returns nil.
func chain(middlewares ...Middleware) Middleware {
	var ms []Middleware
	for _, m := range middlewares {
		if m != nil {
			ms = append(ms, m)
		}
	}
	if len(ms) == 0 {
		return nil
	}
	if len(ms) == 1 {
		return ms[0]
	}
	return func(hf http.HandlerFunc) http.HandlerFunc {
		for i := len(ms) - 1; i >= 0; i-- {
			hf = ms[i](hf)
		}
		return hf
	}
}

func use(middleware Middleware, hf http.HandlerFunc) http.HandlerFunc {
	if middleware == nil {
		return hf
	}
	return middleware(hf)
}
