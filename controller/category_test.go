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
			expectedHttpCode: 400,
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
			categoryMockModel := new(mocks.CategoryInterface)
			// config := config.InitConfig()

			tc.mock(categoryMockModel)

			categoryController := NewCategoryControllerInterface(categoryMockModel)
			handlerFunc := categoryController.AddCategory()

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(buf.String()))

			res := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(req, res)

			err = handlerFunc(ctx)
			require.NoError(t, err)

			w := res.Result()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}

}
