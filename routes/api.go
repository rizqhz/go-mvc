package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rizghz/api/configs"
	. "github.com/rizghz/api/controllers"
)

var (
	key, _ = configs.NewJwtEnv()
)

func UserRoute(e *echo.Echo, c IUserController) {
	users, _ := e.Group("/users"), echojwt.JWT([]byte(key["SECRET_KEY"].(string)))
	users.GET("/login", c.Login())
	// users.GET("", c.Index(), jwt)
	users.GET("", c.Index())
	// users.GET("/:id", c.Observe(), jwt)
	users.GET("/:id", c.Observe())
	users.POST("", c.Store())
	// users.PUT("/:id", c.Edit(), jwt)
	users.PUT("/:id", c.Edit())
	// users.DELETE("/:id", c.Destroy(), jwt)
	users.DELETE("/:id", c.Destroy())
}

func BookRoute(e *echo.Echo, c IBookController) {
	books := e.Group("/books")
	// books.Use(echojwt.JWT(echojwt.JWT([]byte(key["SECRET_KEY"].(string)))))
	books.GET("", c.Index())
	books.GET("/:id", c.Observe())
	books.POST("", c.Store())
	books.PUT("/:id", c.Edit())
	books.DELETE("/:id", c.Destroy())
}

func BlogRoute(e *echo.Echo, c IBlogController) {
	blogs := e.Group("/blogs")
	blogs.GET("", c.Index())
	blogs.GET("/:id", c.Observe())
	blogs.POST("", c.Store())
	blogs.PUT("/:id", c.Edit())
	blogs.DELETE("/:id", c.Destroy())
}
