package api_handlers

import (
	"errors"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
)

type envelope map[string]interface{}

func DeleteUser(a application.AppContext) error {

	app := a.App
	c := a.Context

	id, err := app.ReadIDParam(c)
	if err != nil {
		return app.JsendError(c, "Id retrieval error")
	}

	err = app.Storage.User.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.JsendError(c, "the requested resource could not be found")
		default:
			return app.JsendError(c, "the server was unable to process your request")
		}
	}

	return app.JsendSuccess(c, envelope{"message": "user successfully deleted"})
}
