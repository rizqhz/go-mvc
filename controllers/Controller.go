package controllers

import (
	"github.com/labstack/echo/v4"
)

type Controller interface {
	index() echo.HandlerFunc
	show() echo.HandlerFunc
	create() echo.HandlerFunc
	edit() echo.HandlerFunc
	delete() echo.HandlerFunc
}
