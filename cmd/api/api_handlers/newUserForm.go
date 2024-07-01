package api_handlers

import (
	"gamemasterweb.net/internal/application"
)

func NewUserForm(c application.AppContext) error {
	app := c.App
	return app.RenderHTML(c, "addUser", nil)
}
