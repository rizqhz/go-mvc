package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/helpers"
	m "github.com/rizghz/api/models"
)

type BookController struct {
	model m.IBookModel
}

type IBookController interface {
	Index() echo.HandlerFunc
	Observe() echo.HandlerFunc
	Store() echo.HandlerFunc
	Edit() echo.HandlerFunc
	Destroy() echo.HandlerFunc
}

func NewBookController(model m.IBookModel) IBookController {
	return &BookController{
		model: model,
	}
}

func (c *BookController) Index() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data := c.model.Get()
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *BookController) Observe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid book id", nil))
		}
		data := c.model.Find(&id)
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *BookController) Store() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		book := m.Book{}
		if err := ctx.Bind(&book); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid book data", nil))
		}
		if _, data := c.model.Create(&book); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *BookController) Edit() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid book id", nil))
		}
		book := m.Book{}
		book.ID = uint(id)
		if err := ctx.Bind(&book); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid book data", nil))
		}
		if _, data := c.model.Update(&book); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *BookController) Destroy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid book id", nil))
		}
		if res := c.model.Delete(&id); !res {
			return ctx.JSON(http.StatusInternalServerError,
				helpers.FormatResponse("server error", nil))
		}
		return ctx.JSON(http.StatusNoContent, nil)
	}
}
