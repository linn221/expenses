package routing

import (
	"fmt"
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
	"github.com/linn221/expenses/controllers"
	"github.com/linn221/expenses/middlewares"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

func Serve(c *app.Container) {
	addr := ":" + config.Cfg.API_PORT()
	fmt.Printf("api started on port: %s\n", config.Cfg.API_PORT())
	utils.Must(http.ListenAndServe(addr, Handler(c)))
}

func Handler(c *app.Container) http.Handler {
	mux := http.NewServeMux()
	services := services.New(c)
	controllers := controllers.New(c, services)
	sessionMiddleware := middlewares.Session(c.Cache)
	var rateLimiter = middlewares.NewRateLimiter(c.Counter, config.Cfg.RateLimitWindow(), config.Cfg.RateLimit(), "login", func(r *http.Request) (string, error) {
		return utils.GetClientIP(r), nil
	})

	r := NewNormalRouter()
	r.NewInstance("POST /login", controllers.LoginController.Login(), rateLimiter)
	r.NewInstance("POST /register", controllers.LoginController.Register(), rateLimiter)

	authR := NewNormalRouter(middlewares.Auth())
	authR.NewInstance("POST /logout", controllers.LoginController.Logout())
	authR.NewInstance("GET /profile", controllers.ProfileController.Get())
	authR.NewInstance("GET /incomes/list/{walletId}", controllers.IncomeController.ListByWalletId())
	authR.NewInstance("GET /incomes/paginate/{walletId}", controllers.IncomeController.PaginateByWalletId())
	authR.NewInstance("GET /expenses/list/{walletId}", controllers.ExpenseController.ListByWalletId())
	authR.NewInstance("GET /expenses/paginate/{walletId}", controllers.ExpenseController.PaginateByWalletId())

	rest := NewRestRouter(middlewares.Auth())
	rest.NewResource("/categories", controllers.CategoryController)
	rest.NewResource("/wallets", controllers.WalletController)
	rest.NewResource("/wallet-transfers", controllers.WalletTransferController)
	rest.NewResource("/items", controllers.ItemController)
	rest.NewResource("/incomes", controllers.IncomeController)
	rest.NewResource("/expenses", controllers.ExpenseController)
	rest.NewResource("/cost-perspectives", controllers.CostPerspectiveController)

	r.Register(mux)
	authR.Register(mux)
	rest.Register(mux)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	return middlewares.CORS(sessionMiddleware(middlewares.LoggingMiddleware(middlewares.Recovery(mux))))
}
