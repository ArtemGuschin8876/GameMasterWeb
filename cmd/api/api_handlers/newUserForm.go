package api_handlers

import (
	"gamemasterweb.net/internal/application"
)

func NewUserForm(a application.AppContext) error {
	app := a.App
	c := a.Context

	return app.RenderHTML(c, "addUser", nil)
}
