package api_handlers

import (
	"errors"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
)

func EditUserForm(c application.AppContext) error {

	app := c.App

	id, err := app.ReadIDParam(c)
	if err != nil {
		zeroLog.Err(err).Msg("error reading id")
		return app.JsendError(c, "the requested resource could not be found")
	}

	user, err := app.Storage.User.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			tmplData := data.TemplateData{
				User:       user,
				FormErrors: make(map[string]string),
			}

			zeroLog.Err(err).Msg("request not found")

			if c.Request().Header.Get("Accept") == "application/json" {
				return app.JsendError(c, "user id not found")
			} else {
				return app.RenderHTML(c, "404", tmplData)
			}

		default:
			zeroLog.Err(err).Msg("the server did not process the request")
			return app.JsendError(c, "the server was unable to process your request")
		}
	}

	return nil
}
