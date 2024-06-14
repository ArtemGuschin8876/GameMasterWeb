package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func (app *application) showOneUserHandler(c echo.Context) error {

	id, err := app.readIDParam(c)
	if err != nil {
		return jsendError(c, "Id retrieval error")
	}

	user, err := app.storage.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return jsendError(c, "the requested resource could not be found")
		default:
			return jsendError(c, "the server was unable to process your request")
		}
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return jsendSuccess(c, user)
	}

	ts, ok := app.templates["table.html"]
	if !ok {
		return c.String(http.StatusBadRequest, "template doesn't exist in cache")
	}

	err = ts.Execute(c.Response().Writer, user)

	if err != nil {
		log.Println(err)
		return jsendError(c, "error execute template files")
	}

	return nil
}
