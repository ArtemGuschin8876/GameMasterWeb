package main

import (
	"log"
	"net/http"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) addUsersHandler(c echo.Context) error {

	if c.Request().Method == http.MethodGet {
		app.renderHTML(c, "addUser", nil)
		return nil
	}

	var user data.Users

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusOK, jsendError(c, "invalid request payload"))
	}

	if err := user.ValidateUsers(); err != nil {
		return jsendError(c, err.Error())
	}

	err := app.storage.Users.Add(&user)
	if err != nil {
		log.Println(err)
		return jsendError(c, err.Error())
	}

	err = app.renderHTML(c, "successfullCreatedUser", nil)
	if err != nil {
		log.Println("file rendering error")
		return jsendError(c, "file rendering error")
	}

	return nil
}
