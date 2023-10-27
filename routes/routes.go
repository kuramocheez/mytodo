package routes

import (
	"mytodo/controller"

	"github.com/labstack/echo/v4"
	// mid "github.com/labstack/echo/v4/middleware"
)

func RouteUsers(e *echo.Echo, uc controller.UsersControllerInterface){
	// e.Group)
	e.POST("/signup", uc.Register())
	e.POST("/auth", uc.Login())
}