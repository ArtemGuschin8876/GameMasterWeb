package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) showAllUsersHandler(c echo.Context) error {

	users, err := app.storage.Users.GetAll()
	if err != nil {
		return jsendError(c, "error getting the list of users")
	}

	return c.JSON(http.StatusOK, users)
}
