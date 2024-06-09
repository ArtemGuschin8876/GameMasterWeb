package main

import (
	"errors"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func (app *application) showUsersHandler(c echo.Context) error {

	id, err := app.readIDParam(c)
	if err != nil {
		return jsendError(c, "Id retrieval error")
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			jsendError(c, "the requested resource could not be found")
		default:
			jsendError(c, "the server was unable to process your request")
		}

	}

	return jsendSuccess(c, user)
}
