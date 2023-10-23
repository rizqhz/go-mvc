package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/helpers"
	m "github.com/rizghz/api/models"
)

type UserController struct {
	model m.IUserModel
}

type IUserController interface {
	Index() echo.HandlerFunc
	Observe() echo.HandlerFunc
	Store() echo.HandlerFunc
	Edit() echo.HandlerFunc
	Destroy() echo.HandlerFunc
	Login() echo.HandlerFunc
}

func NewUserController(model m.IUserModel) IUserController {
	return &UserController{
		model: model,
	}
}

func (c *UserController) Index() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data := c.model.Get()
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *UserController) Observe() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user id", nil))
		}
		data := c.model.Find(&id)
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", data))
	}
}

func (c *UserController) Store() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := m.User{}
		if err := ctx.Bind(&user); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user data", nil))
		}
		if _, data := c.model.Create(&user); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *UserController) Edit() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user id", nil))
		}
		user := m.User{}
		user.ID = uint(id)
		if err := ctx.Bind(&user); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user data", nil))
		}
		if _, data := c.model.Update(&user); data != nil {
			return ctx.JSON(http.StatusCreated,
				helpers.FormatResponse("success", data))
		}
		return ctx.JSON(http.StatusInternalServerError,
			helpers.FormatResponse("server error", nil))
	}
}

func (c *UserController) Destroy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user id", nil))
		}
		if res := c.model.Delete(&id); !res {
			return ctx.JSON(http.StatusInternalServerError,
				helpers.FormatResponse("server error", nil))
		}
		return ctx.JSON(http.StatusNoContent, nil)
	}
}

func (c *UserController) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := m.User{}
		if err := ctx.Bind(&user); err != nil {
			return ctx.JSON(http.StatusBadRequest,
				helpers.FormatResponse("invalid user data", nil))
		}
		res, err := c.model.Check(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusOK,
			helpers.FormatResponse("success", res))
	}
}
