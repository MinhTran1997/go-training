package router

import (
	"go-echo-api-mongodb/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Router() *echo.Echo {
	e := echo.New()

	//middleware
	e.Use(middleware.BodyDump(handlers.BodyDumpHandler))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "host=${host}, method=${method}, uri=${uri}, status=${status}, error=${error}, message=${message}\n",
	}))
	e.Use(middleware.Recover())

	//routers
	e.GET("/employees", handlers.GetAllEmployees)
	e.GET("/employee/:id", handlers.GetEmployeeByID)
	e.POST("/employee", handlers.AddEmployee)
	e.PUT("/employee/:id", handlers.UpdateEmployeeById)
	e.DELETE("/employee/:id", handlers.DeleteEmployeeByID)

	return e
}
