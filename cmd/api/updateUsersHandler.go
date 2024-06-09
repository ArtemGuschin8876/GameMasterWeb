package main

import (
	"errors"
	"net/http"

	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) updateUsersHandler(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		c.JSON(http.StatusNotFound, jsendError(c, "the requested resource could not be found"))
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, jsendError(c, "the requested resource could not be found"))
		default:
			c.JSON(http.StatusInternalServerError, jsendError(c, "the server was unable to process your request"))
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
		return c.JSON(http.StatusBadRequest, jsendError(c, "bad response to an enquiry"))
	}

	user.Name = input.Name
	user.Nickname = input.Nickname
	user.Email = input.Email
	user.City = input.City
	user.About = input.About
	user.Image = input.Image

	err = app.models.Users.Update(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, jsendError(c, "error updating user"))
	}

	return c.JSON(http.StatusOK, jsendSuccess(c, user))
}
