package controller

import (
	"mytodo/helper"
	"mytodo/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoControllerInterface interface {
	AddTodo() echo.HandlerFunc
	GetTodos() echo.HandlerFunc
	GetTodo() echo.HandlerFunc
	UpdateTodo() echo.HandlerFunc
	UpdateTodoStatus() echo.HandlerFunc
	DeleteTodo() echo.HandlerFunc
}

type TodoController struct {
	model model.TodoInterface
}

func NewTodoControllerInterface(m model.TodoInterface) TodoControllerInterface {
	return &TodoController{
		model: m,
	}
}

func (tc *TodoController) AddTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		data := model.Todo{}
		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Bind Data", nil))
		}
		data.Status = "OnGoing"
		data.UserID = uint(id)
		res := tc.model.AddTodo(data)
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Create Todo Failed", nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse("Create Todo Successfull", nil))
	}
}

func (tc *TodoController) GetTodos() echo.HandlerFunc {
	return func(c echo.Context) error {
		// todo := []model.Todo{}
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		pageString := c.QueryParam("page")
		page, err := strconv.Atoi(pageString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Get Page Value", nil))
		}
		perPageString := c.QueryParam("content")
		content, err := strconv.Atoi(perPageString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Get Content Value", nil))
		}
		status := c.QueryParam("status")
		date := c.QueryParam("date")
		todo := tc.model.GetTodos(page, content, uint(id), status, date)
		if todo == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Data Not Found", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Get Todo Successfull", todo))
	}
}

func (tc *TodoController) GetTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		idTodoString := c.Param("id")
		idTodo, err := strconv.Atoi(idTodoString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Id Category Format Wrong", nil))
		}
		res := tc.model.GetTodo(idTodo, uint(id))
		if res == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Data Not Found", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Get Todo Successfull", res))
	}
}

func (tc *TodoController) UpdateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		idTodoString := c.Param("id")
		idTodo, err := strconv.Atoi(idTodoString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Format Id", nil))
		}
		todo := model.Todo{}
		if err := c.Bind(&todo); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Bind Data", nil))
		}
		res := tc.model.UpdateTodo(idTodo, uint(id), todo)
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Update Todo Failed", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Update Todo Successfull", nil))
	}
}
func (tc *TodoController) UpdateTodoStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		idTodoString := c.Param("id")
		idTodo, err := strconv.Atoi(idTodoString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Format Id Wrong", nil))
		}
		res := tc.model.UpdateTodoStatus(idTodo, uint(id), "Done")
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Update Todo Status Failed", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Update Todo Status Successfull", nil))
	}
}

func (tc *TodoController) DeleteTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		idTodoString := c.Param("id")
		idTodo, err := strconv.Atoi(idTodoString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Format Id Wrong", nil))
		}
		res := tc.model.DeleteTodo(idTodo, uint(id))
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Delete Todo Failed", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Delete Todo Successfull", nil))
	}
}
