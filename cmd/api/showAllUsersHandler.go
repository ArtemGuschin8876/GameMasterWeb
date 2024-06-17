package main

import (
	"log"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *application) showAllUsersHandler(c echo.Context) error {

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

	users, err := app.storage.Users.GetAll()
	if err != nil {
		return jsendError(c, "error getting the list of users")
	}

	acceptHeader := c.Request().Header.Get("Accept")

	if strings.Contains(acceptHeader, "application/json") {
		return jsendSuccess(c, users)
	}

	data := TemplateData{
		Users: users,
		Flash: flashMessage,
	}

	err = app.renderHTML(c, "tableAllUsers", data)
	if err != nil {
		log.Println("file rendering error")
		return jsendError(c, "file rendering error")
	}

	return nil
}
