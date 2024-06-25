package api_handlers

import (
	"errors"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
)

func EditUserForm(a application.AppContext) error {

	app := a.App
	c := a.Context

	id, err := app.ReadIDParam(c)
	if err != nil {
		return app.JsendError(c, "the requested resource could not be found")
	}

	user, err := app.Storage.User.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.JsendError(c, "the requested resource could not be found")
		default:
			app.JsendError(c, "the server was unable to process your request")
		}
	}

	tmplData := data.TemplateData{
		User:       user,
		FormErrors: make(map[string]string),
	}

	return app.RenderHTML(c, "updateUserForms", tmplData)
}
