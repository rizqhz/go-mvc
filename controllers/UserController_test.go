package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/models"
	"github.com/stretchr/testify/assert"
)

type ValidUserMockModel struct{}

func (mock *ValidUserMockModel) Get() []models.User {
	return []models.User{
		{Name: "User A", Email: "a@mail.com", Password: "A123"},
		{Name: "User B", Email: "b@mail.com", Password: "B123"},
		{Name: "User C", Email: "c@mail.com", Password: "C123"},
	}
}

func (mock *ValidUserMockModel) Find(key *int) *models.User {
	return &models.User{
		Name:     "User A",
		Email:    "a@mail.com",
		Password: "A123",
	}
}

func (mock *ValidUserMockModel) Create(user *models.User) (bool, *models.User) {
	return true, user
}

func (mock *ValidUserMockModel) Update(user *models.User) (bool, *models.User) {
	return true, user
}

func (mock *ValidUserMockModel) Delete(key *int) bool {
	return true
}

func (mock *ValidUserMockModel) Check(user *models.User) (*models.User, error) {
	user.Token = "falsidfk2j3r123klflasjf"
	return user, nil
}

type InvalidUserMockModel struct{}

func (mock *InvalidUserMockModel) Get() []models.User {
	return nil
}

func (mock *InvalidUserMockModel) Find(key *int) *models.User {
	return nil
}

func (mock *InvalidUserMockModel) Create(user *models.User) (bool, *models.User) {
	return false, nil
}

func (mock *InvalidUserMockModel) Update(user *models.User) (bool, *models.User) {
	return false, nil
}

func (mock *InvalidUserMockModel) Delete(key *int) bool {
	return false
}

func (mock *InvalidUserMockModel) Check(user *models.User) (*models.User, error) {
	return nil, errors.New("Invalid")
}

type UserResponseA struct {
	Data    []models.User `json:"data"`
	Message string        `json:"message"`
}

type UserResponseB struct {
	Data    models.User `json:"data"`
	Message string      `json:"message"`
}

func TestUserIndex(t *testing.T) {
	e := echo.New()
	req, res := httptest.NewRequest(http.MethodGet, "/users", nil), UserResponseA{}

	t.Run("Valid User Index", func(t *testing.T) {
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/users", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid User Index (empty)", func(t *testing.T) {
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/users", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})
}

func TestUserObserve(t *testing.T) {
	e := echo.New()

	t.Run("Valid User Observe", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/users/1", nil), UserResponseB{}
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/users/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid User Observe (empty)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/users/1", nil), UserResponseB{}
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/users/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})

	t.Run("Invalid User Observe (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/users/satu", nil), UserResponseB{}
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/users/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user id", res.Message)
		}
	})
}

func TestUserStore(t *testing.T) {
	e := echo.New()
	data := []byte(`{"name":"User Baru", "email":"baru@mail.com", "password":"Baru321"}`)

	t.Run("Valid User Store", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/users", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid User Store (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/users", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})

	t.Run("Invalid User Store (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(data)), UserResponseB{}
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/users", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user data", res.Message)
		}
	})
}

func TestUserEdit(t *testing.T) {
	e := echo.New()
	data := []byte(`{"name":"User Baru", "email":"baru@mail.com", "password":"Baru321"}`)

	t.Run("Valid User Edit", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/users/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid User Edit (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/users/satu", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/users/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user id", res.Message)
		}
	})

	t.Run("Invalid User Edit (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(data)), UserResponseB{}
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/users/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user data", res.Message)
		}
	})

	t.Run("Invalid User Edit (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/users/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}

func TestUserDestroy(t *testing.T) {
	e := echo.New()

	t.Run("Valid User Destroy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/users/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Invalid User Destroy (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/users/satu", nil), UserResponseB{}
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/users/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user id", res.Message)
		}
	})

	t.Run("Invalid User Destroy (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/users/1", nil), UserResponseB{}
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/users/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}

func TestUserLogin(t *testing.T) {
	e := echo.New()
	data := []byte(`{"email":"a@mail.com", "password":"A123"}`)

	t.Run("Valid User Login", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&ValidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/login", controller.Login())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Login()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data.Token)
		}
	})

	t.Run("Invalid User Login (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(data)), UserResponseB{}
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/login", controller.Login())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Login()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid user data", res.Message)
		}
	})

	t.Run("Invalid User Login (empty)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/login", bytes.NewReader(data)), UserResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewUserController(&InvalidUserMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/login", controller.Login())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Login()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Invalid", res.Message)
		}
	})
}
