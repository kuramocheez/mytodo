package controller

import (
	"bytes"
	"encoding/json"
	"io"

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

func TestCategoryController_AddCategory(t *testing.T) {
	mockRequest := model.Category{
		Category: "My Todo",
		Color:    "#FFFFFF",
		UserID:   1,
	}
	test := []struct {
		name             string
		mock             func(*mocks.CategoryInterface)
		expectedHttpCode int
		in               any
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.CategoryInterface) {
				m.On("AddCategory", mock.Anything).Return(true)
			},
			expectedHttpCode: 201,
			in:               mockRequest,
		},
		{
			name: "Should be error, because unexpected return from category model",
			mock: func(m *mocks.CategoryInterface) {
				m.On("AddCategory", mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
		},
		{
			name: "Should be error, because invalid parse body",
			mock: func(m *mocks.CategoryInterface) {
				m.On("AddCategory", mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in: map[string]interface{}{
				"category": 1234,
			},
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)

			err = categoryController.AddCategory()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)

			// request -> middleware(JWT-ECHO) -> jwt.token -> parsing id blba

			// request -> jwt.Token ? -> parsing nil
			// reequest -> token -> jwt.token ? -> parsing string -> interface{} -> jwt.Token

			// request -> jwt.Token <- manual -> -> parsing valid jwt
		})
	}

}
func TestCategoryController_GetCategories(t *testing.T) {
	mockRequest := model.Category{}
	test := []struct {
		name             string
		mock             func(*mocks.CategoryInterface)
		expectedHttpCode int
		in               any
		valuePage        string
		valueContent     string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategories", mock.Anything, mock.Anything, mock.Anything).Return([]model.Category{})
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "5",
		},
		{
			name: "Should be error, because unexpected return from category model",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategories", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 404,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "5",
		},
		{
			name: "Should be error, because page value format wrong",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategories", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			valuePage:        "abc",
			valueContent:     "5",
		},
		{
			name: "Should be error, because content value format wrong",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategories", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			valuePage:        "1",
			valueContent:     "abc",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/category", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")
			q := req.URL.Query()
			q.Add("page", tc.valuePage)
			q.Add("content", tc.valueContent)
			req.URL.RawQuery = q.Encode()
			res := httptest.NewRecorder()

			jwtMock := jwt.New(jwt.SigningMethodHS256)
			jwtMock.Claims = jwt.MapClaims{
				"id": float64(1),
			}

			var jwtMockItf interface{} = jwtMock

			ctx := e.NewContext(req, res)
			ctx.Set("user", jwtMockItf)

			err = categoryController.GetCategories()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}

func TestCategoryController_GetCategory(t *testing.T) {
	mockRequest := model.Category{}
	// expectedCategory := &model.Category{}
	test := []struct {
		name             string
		mock             func(*mocks.CategoryInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategory", mock.Anything, mock.Anything).Return(&model.Category{})
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from category model",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 404,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.CategoryInterface) {
				m.On("GetCategory", mock.Anything, mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/category/:id", strings.NewReader(buf.String()))
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

			err = categoryController.GetCategory()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
func TestCategoryController_UpdateCategory(t *testing.T) {
	mockRequest := model.Category{}
	test := []struct {
		name             string
		mock             func(*mocks.CategoryInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.CategoryInterface) {
				m.On("UpdateCategory", mock.Anything, mock.Anything, mock.Anything).Return(true)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from category model",
			mock: func(m *mocks.CategoryInterface) {
				m.On("UpdateCategory", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.CategoryInterface) {
				m.On("UpdateCategory", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
		{
			name: "Should be error, because bind data error",
			mock: func(m *mocks.CategoryInterface) {
				m.On("UpdateCategory", mock.Anything, mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in: map[string]any{
				"category": 123,
			},
			id: "1",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/category/:id", strings.NewReader(buf.String()))
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

			err = categoryController.UpdateCategory()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
func TestCategoryController_DeleteCategory(t *testing.T) {
	mockRequest := model.Category{}
	test := []struct {
		name             string
		mock             func(*mocks.CategoryInterface)
		expectedHttpCode int
		in               any
		id               string
	}{
		{
			name: "Should be Success",
			mock: func(m *mocks.CategoryInterface) {
				m.On("DeleteCategory", mock.Anything, mock.Anything).Return(true)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because unexpected return from category model",
			mock: func(m *mocks.CategoryInterface) {
				m.On("DeleteCategory", mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
			id:               "1",
		},
		{
			name: "Should be error, because id value format wrong",
			mock: func(m *mocks.CategoryInterface) {
				m.On("DeleteCategory", mock.Anything, mock.Anything).Return(false)
			},
			expectedHttpCode: 400,
			in:               mockRequest,
			id:               "!",
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(tt *testing.T) {
			e := echo.New()

			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/category/:id", strings.NewReader(buf.String()))
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

			err = categoryController.DeleteCategory()(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
