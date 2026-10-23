package routing

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/controllers"
	"github.com/linn221/expenses/services"
)

func loadControllers(c *app.Container, provider *services.Provider) controllers.Controllers {
	return controllers.LoadControllers(c, provider)
}
