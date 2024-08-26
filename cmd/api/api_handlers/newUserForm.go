package api_handlers

import (
	"gamemasterweb.net/internal/application"
)

func NewUserForm(c application.AppContext) error {

	app := c.App

	err := app.RenderHTML(c, "addUser", nil)
	if err != nil {
		zeroLog.Err(err).Msg("error render html in newUserForm")
	}

	return nil
}
