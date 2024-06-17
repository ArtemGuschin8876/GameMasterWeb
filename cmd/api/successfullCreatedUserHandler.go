package main

import "github.com/labstack/echo/v4"

func (app *application) successfullCreatedUserHandler(c echo.Context) error {
	return app.renderHTML(c, "successfullCreatedUser", nil)
}
