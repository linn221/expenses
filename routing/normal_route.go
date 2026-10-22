package routing

import (
	"net/http"

	"github.com/linn221/expenses/config"
)

type NormalRouter struct {
	middleware Middleware
	instances  []*NormalRouteInstance
}

type NormalRouteInstance struct {
	pattern    string
	handle     http.HandlerFunc
	middleware Middleware
}

func NewNormalRouter(middlewares ...Middleware) *NormalRouter {
	return &NormalRouter{
		middleware: chain(middlewares...),
	}
}

func (r *NormalRouter) NewInstance(pattern string, handle http.HandlerFunc, middlewares ...Middleware) *NormalRouteInstance {
	instance := &NormalRouteInstance{
		pattern:    pattern,
		handle:     handle,
		middleware: chain(middlewares...),
	}
	r.instances = append(r.instances, instance)
	return instance
}

func (r *NormalRouter) Register(mux *http.ServeMux) {
	registerRoute := func(pattern string, f http.HandlerFunc) {
		mux.HandleFunc(pattern, f)
	}

	if config.Cfg.DEV_MODE() {
		m := myMux{ServeMux: mux}
		registerRoute = func(pattern string, f http.HandlerFunc) {
			m.HandleFunc(pattern, f)
		}
	}
	for _, inst := range r.instances {
		registerRoute(inst.pattern, use(chain(r.middleware, inst.middleware), inst.handle))
	}
}
