package controller

import (
	"mytodo/config"
	"mytodo/helper"
	"mytodo/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersControllerInterface interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
}

type UsersController struct {
	cfg   config.ProgramConfig
	model model.UsersInterface
}

func NewUsersControllerInterface(m model.UsersInterface, cf config.ProgramConfig) UsersControllerInterface {
	return &UsersController{
		model: m,
		cfg:   cf,
	}
}

func (uc *UsersController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := model.Users{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Register Failed, Error Bind Data", nil))
		}

		res := uc.model.Register(data)
		if res == nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Register Failed", nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse("Register Successfull", res))
	}
}

func (uc *UsersController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := model.Login{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Login Failed, Error Bind Data", nil))
		}
		res := uc.model.Login(data)
		if res == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Login Failed, Username or Password Wrong", nil))
		}
		token := helper.GenerateJWT(uc.cfg.Secret, res.ID)
		if token == nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Login Failed, Error Generate JWT", nil))
		}
		token["info"] = res
		c.Set("user", token["access_token"])
		return c.JSON(http.StatusOK, helper.FormatResponse("Login Successfull", token))
	}
}
