package main

import "github.com/labstack/echo/v4"

func (app *application) routes() *echo.Echo {
	e := echo.New()

	e.GET("/users", app.showUsersHandler)
	e.Static("/swagger/", "C:/Users/rxri2/go/projects/GameMasterWeb/static/dist")
	e.File("/docs/api/swagger.json", "C:/Users/rxri2/go/projects/GameMasterWeb/cmd/api/docs/swagger.json")

	return e
}
