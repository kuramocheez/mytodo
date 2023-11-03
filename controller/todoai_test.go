package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"

	"mytodo/config"
	"mytodo/model"
	"mytodo/model/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCategoryController_TodoAI(t *testing.T) {
	mockRequest := model.TodoAI{
		Todo: "Olahraga",
		Time: time.Date(2023, 11, 03, 14, 0, 0, 0, time.Local),
	}
	test := []struct {
		name             string
		mock             func(*mocks.TodoAIInterface)
		expectedHttpCode int
		in               any
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoAIInterface) {
				m.On("GetResponseAPI", mock.Anything, mock.Anything, mock.Anything).Return(openai.ChatCompletionResponse{}, nil)
			},
			expectedHttpCode: 201,
			in:               mockRequest,
		},
		{
			name: "Should be error, because unexpected return from TodoAI model",
			mock: func(m *mocks.TodoAIInterface) {
				m.On("GetResponseAPI", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Something error"))
			},
			expectedHttpCode: 500,
			in:               mockRequest,
		},
		{
			name: "Should be error, because invalid parse body",
			mock: func(m *mocks.TodoAIInterface) {
				m.On("GetResponseAPI", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in: map[string]interface{}{
				"todo": 1234,
			},
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoAiMockModel := new(mocks.TodoAIInterface)
			config := config.ProgramConfig{}

			tc.mock(todoAiMockModel)

			TodoAIController := NewTodoAIControllerInterface(todoAiMockModel, config)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/todoai", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)

			err = TodoAIController.TodoAI()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
