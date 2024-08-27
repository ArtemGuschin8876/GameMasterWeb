package api_handlers

import (
	"errors"
	"strings"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
)

// @Summary Получить всех пользователей
// @Description Получить информацию о пользователе по его ID
// @Produce json
// @Success 200 {object} User
// @Failure 404 {string} string
func ShowUser(c application.AppContext) error {

	app := c.App

	id, err := app.ReadIDParam(c)
	if err != nil {
		zeroLog.Err(err).Msg("error reading id")
		return app.JsendError(c, "Id retrieval error")
	}

	user, err := app.Storage.User.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			zeroLog.Err(err).Msg("not found")
			return app.JsendError(c, "the requested resource could not be found")
		default:
			zeroLog.Err(err).Msg("incorrect request")
			return app.JsendError(c, "the server was unable to process your request")
		}
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return app.JsendSuccess(c, user)
	}

	err = app.RenderHTML(c, "table", user)
	if err != nil {
		zeroLog.Err(err).Msg("file rendering error")
		return app.JsendError(c, "file rendering error")
	}
	return nil
}
