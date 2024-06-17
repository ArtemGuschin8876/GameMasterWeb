package main

import "github.com/labstack/echo/v4"

func (app *application) showUserForm(c echo.Context) error {
	return app.renderHTML(c, "addUser", nil)
}
