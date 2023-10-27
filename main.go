package main

import (
	"mytodo/config"
	"mytodo/controller"
	"mytodo/model"
	"mytodo/routes"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main(){
	e := echo.New()
	config := config.InitConfig()

	db := model.InitModel(*config)
	model.Migrate(db)

	usersModel := model.NewUsersModel(db)

	usersController := controller.NewUsersControllerInterface(usersModel)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, latency_human=${latency_human}\n",
		}))
	routes.RouteUsers(e, usersController)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}