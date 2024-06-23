package api_handlers

import (
	"gamemasterweb.net/internal/application"
	"github.com/labstack/echo/v4"
)

func NewUserForm(c echo.Context) error {
	cc := c.(*application.CustomContext)
	app := cc.App

	return app.RenderHTML(c, "addUser", nil)
}
