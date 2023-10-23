package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rizghz/api/configs"
	"github.com/rizghz/api/controllers"
	"github.com/rizghz/api/models"
	"github.com/rizghz/api/routes"
)

func main() {
	env, err := configs.NewDatabaseEnv()
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	conf := configs.NewDatabaseConfig(env)
	db, err := models.Init(conf)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	err = models.Migrate(db)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	mUser := models.NewUserModel(db)
	cUser := controllers.NewUserController(mUser)

	mBook := models.NewBookModel(db)
	cBook := controllers.NewBookController(mBook)

	mBlog := models.NewBlogModel(db)
	cBlog := controllers.NewBlogController(mBlog)

	e := echo.New()

	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} time=${time_rfc3339}\n",
	}))

	routes.UserRoute(e, cUser)
	routes.BookRoute(e, cBook)
	routes.BlogRoute(e, cBlog)

	e.GET("/coba", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]any{
			"message": "success",
		})
	})

	e.Logger.Fatal(e.Start(":8008").Error())
}
