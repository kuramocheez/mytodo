package controller

import (
	"fmt"
	"io"
	"testing"

	"bytes"
	"encoding/json"
	"mytodo/model"
	"mytodo/model/mocks"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserController_Register(t *testing.T) {

	mockRequest := model.Users{
		Name:     "bobi",
		Email:    "agus@gmail.com",
		Password: "Something",
	}
	mockUserResult := &model.Users{
		Name:     "Budi",
		Email:    "agus@gmail.com",
		Password: "Something",
	}

	tests := []struct {
		name             string
		mock             func(*mocks.UsersInterface)
		expectedHttpCode int
		in               any
	}{
		{
			name: "should be success",
			mock: func(m *mocks.UsersInterface) {
				m.On("Register", mock.Anything).Return(mockUserResult)
			},
			expectedHttpCode: 201,
			in:               mockRequest,
		},
		{
			name: "should be error, because unexpected return from register model",
			mock: func(m *mocks.UsersInterface) {
				m.On("Register", mock.Anything).Return(nil)
			},
			expectedHttpCode: 500,
			in:               mockRequest,
		},
		{
			name: "should be error, because invalid parse body",
			mock: func(m *mocks.UsersInterface) {
				m.On("Register", mock.Anything).Return(nil)
			},
			expectedHttpCode: 400,
			in: map[string]interface{}{
				"name": 1234,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			userMockModel := new(mocks.UsersInterface)
			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(userMockModel)

			userController := NewUsersControllerInterface(userMockModel, categoryMockModel)
			handlerFunc := userController.Register()

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")

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

func TestUsersController_Login(t *testing.T) {
	mockRequest := model.Login{
		Email:    "agus@gmail.com",
		Password: "Something",
	}

	mockUserResult := &model.Users{
		Name:     "Budi",
		Email:    "agus@gmail.com",
		Password: "Something",
	}

	tests := []struct {
		name             string
		mock             func(*mocks.UsersInterface)
		expectedHttpCode int
		in               any
	}{
		{
			name: "should be success",
			mock: func(m *mocks.UsersInterface) {
				m.On("Login", mock.Anything).Return(mockUserResult)
			},
			expectedHttpCode: 200,
			in:               mockRequest,
		},
		{
			name: "should be error, because invalid parse body",
			mock: func(m *mocks.UsersInterface) {
				// m.On("Login", mock.Anything).Return(&model.Users{
				// 	Name:     "Budi",
				// 	Email:    "agus@gmail.com",
				// 	Password: "Something",
				// })
			},
			expectedHttpCode: 400,
			in: map[string]interface{}{
				"email":    123,
				"password": 123,
			},
		},
		{
			name: "should be error, because email or password wrong",
			mock: func(m *mocks.UsersInterface) {
				m.On("Login", mock.Anything).Return(nil)
			},
			expectedHttpCode: 404,
			in: map[string]interface{}{
				"email":    "agus@gmail.com",
				"password": "Something",
			},
		},
		// {
		// 	name: "should be error, because token error",
		// 	mock: func(m *mocks.UsersInterface) {
		// 		m.On("Login", mock.Anything).Return("")
		// 	},
		// 	expectedHttpCode: 400,
		// 	in:               mockRequest,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userMockModel := new(mocks.UsersInterface)
			categoryMockModel := new(mocks.CategoryInterface)

			tc.mock(userMockModel)

			UsersController := NewUsersControllerInterface(userMockModel, categoryMockModel)
			handlerFunc := UsersController.Login()

			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(tc.in)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(buf.String()))
			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()

			e := echo.New()
			ctx := e.NewContext(req, res)

			err = handlerFunc(ctx)
			require.NoError(t, err)

			w := res.Result()
			body, err := io.ReadAll(w.Body)
			fmt.Println(string(body))
			require.NoError(t, err)

			require.Equal(t, tc.expectedHttpCode, w.StatusCode)
		})
	}
}
