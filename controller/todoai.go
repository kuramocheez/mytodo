package controller

import (
	"mytodo/config"
	"mytodo/helper"
	"mytodo/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type TodoAIControllerInterface interface {
	TodoAI() echo.HandlerFunc
}

type TodoAIController struct {
	model model.TodoAIInterface
	cfg   config.ProgramConfig
}

func NewTodoAIControllerInterface(m model.TodoAIInterface, cf config.ProgramConfig) TodoAIControllerInterface {
	return &TodoAIController{
		model: m,
		cfg:   cf,
	}
}

func (tc *TodoAIController) TodoAI() echo.HandlerFunc {
	return func(c echo.Context) error {
		todoai := model.TodoAI{}
		err := c.Bind(&todoai)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse("Get Response Error Because Bind Data Error", nil))
		}
		key := tc.cfg.ApiKey
		res, err := tc.model.GetResponseAPI(c, key, todoai)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse("Get Recomendation Todo Failed", nil))
		}
		resp := openai.ChatCompletionMessage{
			Content: res.Choices[0].Message.Content,
		}
		return c.JSON(http.StatusCreated, helper.FormatResponse("Get Recomendation Todo Successfull", resp))
	}
}
