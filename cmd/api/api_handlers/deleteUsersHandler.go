package api_handlers

import (
	"errors"

	"gamemasterweb.net/internal/app"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

type envelope map[string]interface{}

func DeleteUsersHandler(c echo.Context, app *app.Application) error {

	id, err := app.ReadIDParam(c)
	if err != nil {
		return app.JsendError(c, "Id retrieval error")
	}

	err = app.Storage.Users.Delete(id)
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
