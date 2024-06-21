package main

import (
	"errors"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) editUserFormHandler(c echo.Context) error {

	id, err := app.readIDParam(c)
	if err != nil {
		return jsendError(c, "the requested resource could not be found")
	}

	user, err := app.storage.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			jsendError(c, "the requested resource could not be found")
		default:
			jsendError(c, "the server was unable to process your request")
		}
	}

	tmplData := TemplateData{
		U:          user,
		FormErrors: make(map[string]string),
	}

	return app.renderHTML(c, "updateUserForms", tmplData)
}
