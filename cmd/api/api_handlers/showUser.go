package api_handlers

import (
	"errors"
	"log"
	"strings"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func ShowUser(c echo.Context) error {

	cc := c.(*application.AppContext)
	app := cc.App

	id, err := app.ReadIDParam(c)
	if err != nil {
		return app.JsendError(c, "Id retrieval error")
	}

	user, err := app.Storage.User.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return app.JsendError(c, "the requested resource could not be found")
		default:
			return app.JsendError(c, "the server was unable to process your request")
		}
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return app.JsendSuccess(c, user)
	}

	err = app.RenderHTML(c, "table", user)
	if err != nil {
		log.Println("file rendering error")
		return app.JsendError(c, "file rendering error")
	}
	return nil
}
