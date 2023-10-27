package controller

import (
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
	model model.UsersInterface
}

func NewUsersControllerInterface(m model.UsersInterface) UsersControllerInterface {
	return &UsersController{
		model: m,
	}
}

func (uc *UsersController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := model.Users{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Register Error, Something Error When Bind Data", nil))
		}

		res := uc.model.Register(data)
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Register Error, Something Error When Create Data", nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse("Register Successfull", nil))
	}
}

func (uc *UsersController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := model.Login{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Login Error, Error When Bind Data", nil))
		}
		res := uc.model.Login(data)
		if res == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Login Error, Wrong Username Or Password", nil))
		}
		token, _ := helper.CreateToken(res.ID)
		if token == "" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Login Error, Token Error", nil))
		}
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Login Successfull",
			"token":   token,
		})
	}
}
