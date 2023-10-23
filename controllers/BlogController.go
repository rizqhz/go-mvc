package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/helpers"
	m "github.com/rizghz/api/models"
)

type BlogController struct {
	model m.IBlogModel
}

type IBlogController interface {
	Index() echo.HandlerFunc
	Observe() echo.HandlerFunc
	Store() echo.HandlerFunc
	Edit() echo.HandlerFunc
	Destroy() echo.HandlerFunc
}

func NewBlogController(model m.IBlogModel) IBlogController {
	return &BlogController{
		model: model,
	}
}

func (c *BlogController) Index() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data := c.model.Get()
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *BlogController) Observe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid blog id", nil))
		}
		data := c.model.Find(&id)
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *BlogController) Store() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		blog := m.Blog{}
		if err := ctx.Bind(&blog); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid blog data", blog))
		}
		if _, data := c.model.Create(&blog); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *BlogController) Edit() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid blog id", nil))
		}
		blog := m.Blog{}
		blog.ID = uint(id)
		if err := ctx.Bind(&blog); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid blog data", nil))
		}
		if _, data := c.model.Update(&blog); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *BlogController) Destroy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid blog id", nil))
		}
		if res := c.model.Delete(&id); !res {
			return ctx.JSON(http.StatusInternalServerError,
				helpers.FormatResponse("server error", nil))
		}
		return ctx.JSON(http.StatusNoContent, nil)
	}
}
