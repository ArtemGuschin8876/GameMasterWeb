package main

import (
	"errors"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

type envelope map[string]interface{}

func (app *application) deleteUsersHandler(c echo.Context) error {

	id, err := app.readIDParam(c)
	if err != nil {
		return jsendError(c, "Id retrieval error")
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return jsendError(c, "the requested resource could not be found")
		default:
			return jsendError(c, "the server was unable to process your request")
		}
	}

	return jsendSuccess(c, envelope{"message": "user successfully deleted"})
}
