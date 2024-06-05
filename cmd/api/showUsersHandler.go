package main

import (
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

	user := data.Users{
		ID:       id,
		Name:     "Oleg",
		Nickname: "Parlis",
		Email:    "OlegSuka@gmail.com",
		City:     "Saratov",
		About:    "I am pidor",
	}

	return c.JSON(200, jsendSuccess(c, user))

}
