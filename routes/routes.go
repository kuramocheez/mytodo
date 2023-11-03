package routes

import (
	"mytodo/config"
	"mytodo/controller"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func RouteUsers(e *echo.Echo, uc controller.UsersControllerInterface) {
	e.POST("/signup", uc.Register())
	e.POST("/auth", uc.Login())
}

func RouteCategory(e *echo.Echo, cc controller.CategoryControllerInterface, cfg config.ProgramConfig) {
	auth := e.Group("/category")
	auth.Use(mid.JWT([]byte(cfg.Secret)))
	auth.GET("", cc.GetCategories())
	auth.GET("/:id", cc.GetCategory())
	auth.POST("", cc.AddCategory())
	auth.PUT("/:id", cc.UpdateCategory())
	auth.DELETE("/:id", cc.DeleteCategory())
}

func RouteTodo(e *echo.Echo, tc controller.TodoControllerInterface, cfg config.ProgramConfig) {
	auth := e.Group("/todo")
	auth.Use(mid.JWT([]byte(cfg.Secret)))
	auth.GET("", tc.GetTodos())
	auth.GET("/:id", tc.GetTodo())
	auth.POST("", tc.AddTodo())
	auth.PUT("/:id", tc.UpdateTodo())
	auth.PUT("/status/:id", tc.UpdateTodoStatus())
	auth.DELETE("/:id", tc.DeleteTodo())
}

func RouteTodoAI(e *echo.Echo, tc controller.TodoAIControllerInterface, cfg config.ProgramConfig) {
	auth := e.Group("/todoai")
	auth.Use(mid.JWT([]byte(cfg.Secret)))
	auth.POST("", tc.TodoAI())
}
