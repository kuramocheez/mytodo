package main

import (
	"fmt"
	"mytodo/config"
	"mytodo/controller"
	"mytodo/model"
	"mytodo/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	config := config.InitConfig()

	db := model.InitModel(*config)
	model.Migrate(db)

	usersModel := model.NewUsersModel(db)
	categoryModel := model.NewCategoryModel(db)
	todoModel := model.NewTodoModel(db)
	todoAIModel := model.NewTodoAIModel(db)

	usersController := controller.NewUsersControllerInterface(usersModel, *config)
	categoryController := controller.NewCategoryControllerInterface(categoryModel)
	todoController := controller.NewTodoControllerInterface(todoModel)
	todoAIController := controller.NewTodoAIControllerInterface(todoAIModel, *config)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, latency_human=${latency_human}\n",
		}))
	routes.RouteUsers(e, usersController)
	routes.RouteCategory(e, categoryController, *config)
	routes.RouteTodo(e, todoController, *config)
	routes.RouteTodoAI(e, todoAIController, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
