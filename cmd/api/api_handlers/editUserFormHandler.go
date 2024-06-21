package api_handlers

import (
	"errors"

	"gamemasterweb.net/internal/app"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func EditUserFormHandler(c echo.Context, app *app.Application) error {

	id, err := app.ReadIDParam(c)
	if err != nil {
		return app.JsendError(c, "the requested resource could not be found")
	}

	user, err := app.Storage.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.JsendError(c, "the requested resource could not be found")
		default:
			app.JsendError(c, "the server was unable to process your request")
		}
	}

	tmplData := TemplateData{
		U:          user,
		FormErrors: make(map[string]string),
	}

	return app.RenderHTML(c, "updateUserForms", tmplData)
}
