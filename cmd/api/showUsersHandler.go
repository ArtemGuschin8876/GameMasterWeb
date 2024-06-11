package main

import (
	"errors"
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
func (app *application) showUsersHandler(c echo.Context) error {
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

	// files := []string{
	// 	"./static/ui/html/table.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "error in parse files")
	// }

	ts, ok := app.template["./static/ui/html/table.html"]
	if !ok {
		return c.String(http.StatusBadRequest, "template doesnt exist in cache ")
	}

	err = ts.Execute(c.Response().Writer, user)
	if err != nil {
		return jsendError(c, "error execute template files")
	}

	return nil
}
