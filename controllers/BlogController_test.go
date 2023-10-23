package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/models"
	"github.com/stretchr/testify/assert"
)

type ValidBlogMockModel struct{}

func (mock *ValidBlogMockModel) Get() []models.Blog {
	return []models.Blog{
		{Title: "Title A", Content: "Content A", UserID: 1},
		{Title: "Title B", Content: "Content B", UserID: 2},
		{Title: "Title C", Content: "Content C", UserID: 3},
	}
}

func (mock *ValidBlogMockModel) Find(key *int) *models.Blog {
	return &models.Blog{
		Title:   "Title A",
		Content: "Content A",
		UserID:  1,
	}
}

func (mock *ValidBlogMockModel) Create(blog *models.Blog) (bool, *models.Blog) {
	return true, blog
}

func (mock *ValidBlogMockModel) Update(blog *models.Blog) (bool, *models.Blog) {
	return true, blog
}

func (mock *ValidBlogMockModel) Delete(key *int) bool {
	return true
}

type InvalidBlogMockModel struct{}

func (mock *InvalidBlogMockModel) Get() []models.Blog {
	return nil
}

func (mock *InvalidBlogMockModel) Find(key *int) *models.Blog {
	return nil
}

func (mock *InvalidBlogMockModel) Create(blog *models.Blog) (bool, *models.Blog) {
	return false, nil
}

func (mock *InvalidBlogMockModel) Update(blog *models.Blog) (bool, *models.Blog) {
	return false, nil
}

func (mock *InvalidBlogMockModel) Delete(key *int) bool {
	return false
}

type BlogResponseA struct {
	Data    []models.Blog `json:"data"`
	Message string        `json:"message"`
}

type BlogResponseB struct {
	Data    models.Blog `json:"data"`
	Message string      `json:"message"`
}

func TestBlogIndex(t *testing.T) {
	e := echo.New()
	req, res := httptest.NewRequest(http.MethodGet, "/blogs", nil), BlogResponseA{}

	t.Run("Valid Blog Index", func(t *testing.T) {
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/blogs", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Blog Index (empty)", func(t *testing.T) {
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/blogs", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})
}

func TestBlogObserve(t *testing.T) {
	e := echo.New()

	t.Run("Valid Blog Observe", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/blogs/1", nil), BlogResponseB{}
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/blogs/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Blog Observe (empty)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/blogs/1", nil), BlogResponseB{}
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/blogs/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})

	t.Run("Invalid Blog Observe (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/blogs/satu", nil), BlogResponseB{}
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/blogs/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid blog id", res.Message)
		}
	})
}

func TestBlogStore(t *testing.T) {
	e := echo.New()
	data := []byte(`{"title": "Title Baru", "content": "Content Baru", "user_id": 1}`)

	t.Run("Valid Blog Store", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/blogs", bytes.NewReader(data)), BlogResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/blogs", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Blog Store (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/blogs", bytes.NewReader(data)), BlogResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/blogs", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})

	t.Run("Invalid Blog Store (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/blogs", bytes.NewReader(data)), BlogResponseB{}
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/blogs", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid blog data", res.Message)
		}
	})
}

func TestBlogEdit(t *testing.T) {
	e := echo.New()
	data := []byte(`{"title": "Title Baru", "content": "Content Baru", "user_id": 1}`)

	t.Run("Valid Blog Edit", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/blogs/1", bytes.NewReader(data)), BlogResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/blogs/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Blog Edit (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/blogs/satu", bytes.NewReader(data)), BlogResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/blogs/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid blog id", res.Message)
		}
	})

	t.Run("Invalid Blog Edit (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/blogs/1", bytes.NewReader(data)), BlogResponseB{}
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/blogs/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid blog data", res.Message)
		}
	})

	t.Run("Invalid Blog Edit (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/blogs/1", bytes.NewReader(data)), BlogResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/blogs/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}

func TestBlogDestroy(t *testing.T) {
	e := echo.New()

	t.Run("Valid Blog Destroy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/blogs/1", nil)
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/blogs/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Invalid Blog Destroy (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/blogs/satu", nil), BlogResponseB{}
		controller := NewBlogController(&ValidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/blogs/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid blog id", res.Message)
		}
	})

	t.Run("Invalid Blog Destroy (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/blogs/1", nil), BlogResponseB{}
		controller := NewBlogController(&InvalidBlogMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/blogs/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}
