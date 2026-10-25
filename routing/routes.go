package routing

import (
	"fmt"
	"net/http"
	"time"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
	"github.com/linn221/expenses/controllers"
	"github.com/linn221/expenses/middlewares"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

func Serve(c *app.Container) {
	mux := http.NewServeMux()
	services := services.New(c)
	controllers := controllers.New(c, services)
	sessionMiddleware := middlewares.Session(c.Cache)
	var rateLimiter = middlewares.NewRateLimiter(c.Cache, time.Minute*2, 50, "login", func(r *http.Request) (string, error) {
		return utils.GetClientIP(r), nil
	})

	r := NewNormalRouter()
	r.NewInstance("POST /login", controllers.LoginController.Login(), rateLimiter)
	r.NewInstance("POST /register", controllers.LoginController.Register(), rateLimiter)
	r.Register(mux)

	rest := NewRestRouter(middlewares.Auth())
	rest.NewResource("/categories", controllers.CategoryController)
	rest.NewResource("/wallets", controllers.WalletController)
	rest.NewResource("/items", controllers.ItemController)
	rest.NewResource("/incomes", controllers.IncomeController)
	rest.NewResource("/expenses", controllers.ExpenseController)

	rest.Register(mux)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	server := http.Server{
		Handler: sessionMiddleware(middlewares.LoggingMiddleware(middlewares.Recovery(mux))),
		Addr:    ":" + config.Cfg.API_PORT(),
	}

	fmt.Printf("api started on port: %s\n", config.Cfg.API_PORT())
	utils.Must(server.ListenAndServe())
}
