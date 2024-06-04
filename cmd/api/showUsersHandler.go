package main

import (
	"github.com/labstack/echo/v4"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func (app *application) showUsersHandler(c echo.Context) error {

	err := app.writeJSON(c)
	if err != nil {
		return c.JSON(200, response{
			Status: "error",
			//			data:    nil,
			Message: "problema in Oleg",
		})
	}

	return nil
}
