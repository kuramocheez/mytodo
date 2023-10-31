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
