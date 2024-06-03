package main

import "github.com/labstack/echo/v4"

func (app *application) routes() *echo.Echo {
	e := echo.New()

	e.GET("/users", app.showUsersHandler)

	return e
}
