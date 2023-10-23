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

type ValidBookMockModel struct{}

func (mock *ValidBookMockModel) Get() []models.Book {
	return []models.Book{
		{Title: "Buku A", Author: "Author A", Publisher: "Publisher A"},
		{Title: "Buku B", Author: "Author B", Publisher: "Publisher B"},
		{Title: "Buku C", Author: "Author C", Publisher: "Publisher C"},
	}
}

func (mock *ValidBookMockModel) Find(key *int) *models.Book {
	return &models.Book{
		Title:     "Buku A",
		Author:    "Author A",
		Publisher: "Publisher A",
	}
}

func (mock *ValidBookMockModel) Create(book *models.Book) (bool, *models.Book) {
	return true, book
}

func (mock *ValidBookMockModel) Update(book *models.Book) (bool, *models.Book) {
	return true, book
}

func (mock *ValidBookMockModel) Delete(key *int) bool {
	return true
}

type InvalidBookMockModel struct{}

func (mock *InvalidBookMockModel) Get() []models.Book {
	return nil
}

func (mock *InvalidBookMockModel) Find(key *int) *models.Book {
	return nil
}

func (mock *InvalidBookMockModel) Create(book *models.Book) (bool, *models.Book) {
	return false, nil
}

func (mock *InvalidBookMockModel) Update(book *models.Book) (bool, *models.Book) {
	return false, nil
}

func (mock *InvalidBookMockModel) Delete(key *int) bool {
	return false
}

type BookResponseA struct {
	Data    []models.Book `json:"data"`
	Message string        `json:"message"`
}

type BookResponseB struct {
	Data    models.Book `json:"data"`
	Message string      `json:"message"`
}

func TestBookIndex(t *testing.T) {
	e := echo.New()
	req, res := httptest.NewRequest(http.MethodGet, "/books", nil), BookResponseA{}

	t.Run("Valid Book Index", func(t *testing.T) {
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/books", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Book Index (empty)", func(t *testing.T) {
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/books", controller.Index())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Index()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})
}

func TestBookObserve(t *testing.T) {
	e := echo.New()

	t.Run("Valid Book Observe", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/books/1", nil), BookResponseB{}
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/books/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Book Observe (empty)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/books/1", nil), BookResponseB{}
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/books/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.Empty(t, res.Data)
		}
	})

	t.Run("Invalid Book Observe (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodGet, "/books/satu", nil), BookResponseB{}
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.GET("/books/:id", controller.Observe())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid book id", res.Message)
		}
	})
}

func TestBookStore(t *testing.T) {
	e := echo.New()
	data := []byte(`{"title":"Buku Baru", "author":"Author Baru", "publisher":"Publisher Baru"}`)

	t.Run("Valid Book Store", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(data)), BookResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/books", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Book Store (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(data)), BookResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/books", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})

	t.Run("Invalid Book Store (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(data)), BookResponseB{}
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.POST("/books", controller.Store())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Observe()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid book data", res.Message)
		}
	})
}

func TestBookEdit(t *testing.T) {
	e := echo.New()
	data := []byte(`{"title":"Buku Baru", "author":"Author Baru", "publisher":"Publisher Baru"}`)

	t.Run("Valid Book Edit", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(data)), BookResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/books/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "success", res.Message)
			assert.NotEmpty(t, res.Data)
		}
	})

	t.Run("Invalid Book Edit (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/books/satu", bytes.NewReader(data)), BookResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/books/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid book id", res.Message)
		}
	})

	t.Run("Invalid Book Edit (payload)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(data)), BookResponseB{}
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/books/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid book data", res.Message)
		}
	})

	t.Run("Invalid Book Edit (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(data)), BookResponseB{}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.PUT("/books/:id", controller.Edit())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Edit()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}

func TestBookDestroy(t *testing.T) {
	e := echo.New()

	t.Run("Valid Book Destroy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/books/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Invalid Book Destroy (id)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/books/satu", nil), BookResponseB{}
		controller := NewBookController(&ValidBookMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/books/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid book id", res.Message)
		}
	})

	t.Run("Invalid Book Destroy (server)", func(t *testing.T) {
		req, res := httptest.NewRequest(http.MethodDelete, "/books/1", nil), BookResponseB{}
		controller := NewBookController(&InvalidBookMockModel{})
		rec := httptest.NewRecorder()
		e.DELETE("/books/:id", controller.Destroy())
		e.ServeHTTP(rec, req)
		json.Unmarshal(rec.Body.Bytes(), &res)
		if assert.NoError(t, nil, controller.Destroy()) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "server error", res.Message)
		}
	})
}
