package routing

import (
	"net/http"

	"github.com/linn221/expenses/config"
)

type RestRouter struct {
	routes     []*ResourceRoute
	middleware Middleware
}

func NewRestRouter(middlewares ...Middleware) *RestRouter {
	return &RestRouter{
		middleware: chain(middlewares...),
	}
}

func (r *RestRouter) NewResource(pattern string, res ResourceController, middlewares ...Middleware) *ResourceRoute {
	resourceInstance := &ResourceRoute{
		pattern:    pattern,
		controller: res,
		middleware: chain(middlewares...),
	}
	r.routes = append(r.routes, resourceInstance)
	return resourceInstance
}

func (r *RestRouter) Register(mux *http.ServeMux) {
	for _, res := range r.routes {
		res.applyRoutes(mux, r.middleware)
	}
}

type ResourceRoute struct {
	pattern    string
	controller ResourceController
	middleware Middleware
}

func (res *ResourceRoute) applyRoutes(mux *http.ServeMux, middleware Middleware) {

	registerRoute := func(pattern string, f http.HandlerFunc) {
		mux.HandleFunc(pattern, f)
	}

	if config.Cfg.DEV_MODE() {
		m := myMux{ServeMux: mux}
		registerRoute = func(pattern string, f http.HandlerFunc) {
			m.HandleFunc(pattern, f)
		}
	}
	if res.middleware != nil {
		parent := middleware
		middleware = func(hf http.HandlerFunc) http.HandlerFunc {
			return res.middleware(use(parent, hf))
		}
	}

	registerRoute("POST "+res.pattern, use(middleware, res.controller.Create()))
	registerRoute("PUT "+res.pattern+"/{id}", use(middleware, res.controller.Update()))
	registerRoute("DELETE "+res.pattern+"/{id}", use(middleware, res.controller.Delete()))
	registerRoute("GET "+res.pattern+"/{id}", use(middleware, res.controller.Get()))

	if patchController, ok := res.controller.(PatchController); ok {
		registerRoute("PATCH "+res.pattern+"/{id}", use(middleware, patchController.Patch()))
	}

	if listController, ok := res.controller.(ListController); ok {
		registerRoute("GET "+res.pattern, use(middleware, listController.List()))
	}

	if paginateController, ok := res.controller.(PaginateController); ok {
		registerRoute("GET "+res.pattern+"/paginate", use(middleware, paginateController.Paginate()))
	}

	if activeToggleController, ok := res.controller.(IsActiveToggleController); ok {
		registerRoute("POST /is-active"+res.pattern+"/{id}", activeToggleController.ToggleIsActive())
	}
}
