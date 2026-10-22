package routing

import (
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/controllers"
	"github.com/linn221/expenses/services"
)

type Controllers struct {
	categoryController *controllers.CategoryController
	itemController     *controllers.ItemController
	loginController    *controllers.LoginController
	balanceController  *controllers.BalanceController
}

func loadControllers(c *app.Container, provider *services.Provider) Controllers {
	myController := controllers.MyController{Container: c, Services: provider}
	return Controllers{
		categoryController: &controllers.CategoryController{
			MyController: myController,
		},
		itemController: &controllers.ItemController{
			MyController: myController,
		},
		loginController: &controllers.LoginController{
			Services: provider,
		},
		balanceController: &controllers.BalanceController{
			MyController: myController,
		},
	}
}
