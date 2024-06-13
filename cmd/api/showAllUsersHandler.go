package main

import (
	"net/http"
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

	ts, ok := app.templates["tableAllUsers.html"]
	if !ok {
		return c.String(http.StatusBadRequest, "template doesn't exist in cache")
	}

	// data := map[string]interface{}{
	// 	"IsList": true,
	// 	"User":   users,
	// }

	err = ts.Execute(c.Response().Writer, users)
	if err != nil {
		return jsendError(c, "error execute template files")
	}

	return nil
}
