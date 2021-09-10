package router

import (
	"go-echo-api-postgres/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router() *echo.Echo {
	e := echo.New()

	// Middlewares
	e.Use(middleware.BodyDump(handlers.BodyDumpHandler))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error}, message=${message}\n",
	}))
	e.Use(middleware.Recover())

	// Routers
	e.GET("/users", handlers.GetAllUsers)
	e.GET("/user/:id", handlers.GetUser)
	e.POST("/user", handlers.CreateUser)
	e.PUT("/user/:id", handlers.UpdateUser)
	e.DELETE("/user/:id", handlers.DeleteUser)

	return e
}
