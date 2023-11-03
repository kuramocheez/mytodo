package controller

import (
	"fmt"
	"mytodo/helper"
	"mytodo/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryControllerInterface interface {
	AddCategory() echo.HandlerFunc
	GetCategories() echo.HandlerFunc
	GetCategory() echo.HandlerFunc
	UpdateCategory() echo.HandlerFunc
	DeleteCategory() echo.HandlerFunc
}

type CategoryController struct {
	model model.CategoryInterface
}

func NewCategoryControllerInterface(m model.CategoryInterface) CategoryControllerInterface {
	return &CategoryController{
		model: m,
	}
}

func (cc *CategoryController) AddCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(c.Get("user"))
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		data := model.Category{}
		if err := c.Bind(&data); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Bind Data", nil))
		}
		data.UserID = uint(id)
		res := cc.model.AddCategory(data)
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Create Category Failed", nil))
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse("Create Category Succesfull", nil))
	}
}

func (cc *CategoryController) GetCategories() echo.HandlerFunc {
	return func(c echo.Context) error {
		categories := []model.Category{}
		claims := helper.ExtractToken("user", c)
		id := claims["id"].(float64)
		pageString := c.QueryParam("page")
		page, err := strconv.Atoi(pageString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Get Page Value", nil))
		}
		perPageString := c.QueryParam("content")
		perPage, err := strconv.Atoi(perPageString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Get Content Value", nil))
		}
		categories = cc.model.GetCategories(page, perPage, uint(id))
		if categories == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Data Not Found", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Success Get Categories Data", &categories))
	}
}

func (cc *CategoryController) GetCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		idUser := claims["id"].(float64)
		idCategoryString := c.Param("id")
		idCategory, err := strconv.Atoi(idCategoryString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Id Category Format Wrong", nil))
		}
		res := cc.model.GetCategory(idCategory, uint(idUser))
		if res == nil {
			return c.JSON(http.StatusNotFound, helper.FormatResponse("Data Not Found", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Get Category Successfull", res))
	}
}

func (cc *CategoryController) UpdateCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		idUser := claims["id"].(float64)
		idCategoryString := c.Param("id")
		category := model.Category{}
		err := c.Bind(&category)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Error Bind Data", nil))
		}
		idCategory, err := strconv.Atoi(idCategoryString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Id Category Format Wrong", nil))
		}
		res := cc.model.UpdateCategory(category, idCategory, uint(idUser))
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Update Category Failed", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Update Category Successfull", nil))
	}
}

func (cc *CategoryController) DeleteCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := helper.ExtractToken("user", c)
		idUser := claims["id"].(float64)
		idCategoryString := c.Param("id")
		idCategory, err := strconv.Atoi(idCategoryString)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Id Category Format Wrong", nil))
		}
		res := cc.model.DeleteCategory(idCategory, uint(idUser))
		if !res {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Delete Category Failed", nil))
		}
		return c.JSON(http.StatusOK, helper.FormatResponse("Delete Category Successfull", nil))
	}
}
