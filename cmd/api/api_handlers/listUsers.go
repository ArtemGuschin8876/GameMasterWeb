package api_handlers

import (
	"strings"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo-contrib/session"
)

func ListUsers(c application.AppContext) error {
	app := c.App

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	flash := sess.Values["flash"]
	delete(sess.Values, "flash")
	sess.Save(c.Request(), c.Response())

	var flashMessage string

	if flash != nil {
		flashMessage, _ = flash.(string)
	}

	users, err := app.Storage.User.GetAll()
	if err != nil {
		return app.JsendError(c, "error getting the list of users")
	}

	data := data.TemplateData{
		Users: users,
		Flash: flashMessage,
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return app.JsendSuccess(c, users)
	} else {
		return app.RenderHTML(c, "tableAllUsers", data)
	}

}
