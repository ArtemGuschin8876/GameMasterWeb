package main

import (
	"errors"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) updateUsersHandler(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		jsendError(c, "the requested resource could not be found")
	}

	user, err := app.storage.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			jsendError(c, "the requested resource could not be found")
		default:
			jsendError(c, "the server was unable to process your request")
		}

	}

	var input struct {
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		City     string `json:"city"`
		About    string `json:"about"`
		Image    string `json:"image"`
	}

	if err := c.Bind(&input); err != nil {
		return jsendError(c, "database error")
	}

	user.Name = input.Name
	user.Nickname = input.Nickname
	user.Email = input.Email
	user.City = input.City
	user.About = input.About
	user.Image = input.Image

	err = app.storage.Users.Update(user)
	if err != nil {
		return jsendError(c, "error updating user")
	}

	return jsendSuccess(c, user)
}
