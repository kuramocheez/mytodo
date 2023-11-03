package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"mytodo/model"
	"mytodo/model/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTodoController_AddTodo(t *testing.T) {
	mockRequest := model.Todo{
		Memo:       "Kelas Live Session Golang",
		DateTime:   time.Date(2023, 11, 03, 16, 0, 0, 0, time.Local),
		CategoryID: 1,
	}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("AddTodo", mock.Anything).Return(true)
			},
			expectedHttpCode: 201,
			in:               mockRequest,
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("AddTodo", mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
		},
		{
			name: "Should be error, because invalid parse body",
			mock: func(m *mocks.TodoInterface) {
				m.On("AddTodo", mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in: map[string]interface{}{
				"memo": 1234,
			},
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			todoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)

			err = todoController.AddTodo()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}

func TestTodoController_GetTodos(t *testing.T) {
	mockRequest := model.Todo{}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
		valuePage        string
		valueContent     string
		status           string
		date             string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodos", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.Todo{})
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "5",
			status:           "",
			date:             "",
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodos", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 404,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "5",
			status:           "",
			date:             "",
		},
		{
			name: "Should be error, because page value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodos", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			valuePage:        "abc",
			valueContent:     "5",
			status:           "",
			date:             "",
		},
		{
			name: "Should be error, because content value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodos", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "abc",
			status:           "",
			date:             "",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			TodoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/todo", nil)
			req.Header.Set("Content-Type", "application/json")
			q := req.URL.Query()
			q.Add("page", tc.valuePage)
			q.Add("content", tc.valueContent)
			q.Add("status", tc.status)
			q.Add("date", tc.date)
			req.URL.RawQuery = q.Encode()
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)

			err = TodoController.GetTodos()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}

func TestTodoController_GetTodo(t *testing.T) {
	mockRequest := model.Todo{}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodo", mock.Anything, mock.Anything).Return(&model.Todo{})
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 404,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("GetTodo", mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			todoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/todo/:id", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.id)

			err = todoController.GetTodo()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}

func TestTodoController_UpdateTodo(t *testing.T) {
	mockRequest := model.Todo{}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(true)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
		{
			name: "Should be error, because bind data error",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodo", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in: map[string]any{
				"memo": 123,
			},
			id: "1",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			todoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/todo/:id", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.id)

			err = todoController.UpdateTodo()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
func TestTodoController_UpdateTodoStatus(t *testing.T) {
	mockRequest := model.Todo{}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodoStatus", mock.Anything, mock.Anything, mock.Anything).Return(true)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodoStatus", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("UpdateTodoStatus", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			todoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/todo/status/:id", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.id)

			err = todoController.UpdateTodoStatus()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}

func TestTodoController_DeleteTodo(t *testing.T) {
	mockRequest := model.Todo{}
	test := []struct {
		name             string
		mock             func(*mocks.TodoInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.TodoInterface) {
				m.On("DeleteTodo", mock.Anything, mock.Anything).Return(true)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from todo model",
			mock: func(m *mocks.TodoInterface) {
				m.On("DeleteTodo", mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.TodoInterface) {
				m.On("DeleteTodo", mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			todoMockModel := new(mocks.TodoInterface)

			tc.mock(todoMockModel)

			todoController := NewTodoControllerInterface(todoMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/todo/:id", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.id)

			err = todoController.DeleteTodo()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
