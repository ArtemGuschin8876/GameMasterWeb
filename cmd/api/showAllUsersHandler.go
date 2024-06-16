package main

import (
	"log"
	"strings"

	"github.com/labstack/echo/v4"
)

func (app *application) showAllUsersHandler(c echo.Context) error {

	users, err := app.storage.Users.GetAll()
	if err != nil {
		return jsendError(c, "error getting the list of users")
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return jsendSuccess(c, users)
	}

	err = app.renderHTML(c, "tableAllUsers", users)
	if err != nil {
		log.Println("file rendering error")
		return jsendError(c, "file rendering error")
	}

	return nil
}
