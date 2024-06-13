package main

import (
	"log"
	"net/http"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) addUsersHandler(c echo.Context) error {

	var user data.Users

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusOK, jsendError(c, "invalid request payload"))
	}

	if err := user.Validate(); err != nil {
		return jsendError(c, err.Error())
	}

	err := app.storage.Users.Add(&user)
	if err != nil {
		log.Println(err)
		return jsendError(c, err.Error())
	}

	return jsendSuccess(c, user)
}
