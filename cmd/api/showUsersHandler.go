package main

import (
	"errors"
	"fmt"
	"text/template"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func (app *application) showUsersHandler(c echo.Context) error {
	fmt.Println(c)
	id, err := app.readIDParam(c)
	fmt.Println("id", id)
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

	// data := &pagesData{Page: user}

	files := []string{
		"./static/ui/html/table.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		return jsendError(c, "error read template files")
	}

	err = ts.Execute(c.Response().Writer, user)
	if err != nil {
		return jsendError(c, "error execute template files")
	}

	return nil
}
