package api_handlers

import (
	"log"
	"net/http"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type TemplateData struct {
	Errors      []string
	FormErrors  map[string]string
	User        data.Users
	Flash       string
	Users       []data.Users
	UserPointer *data.Users
}

func AddUsersHandler(c echo.Context) error {
	var user data.Users

	cc := c.(*application.CustomContext)
	app := cc.App

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusOK, app.JsendError(c, "invalid request payload"))
	}

	if err := user.ValidateUsers(); err != nil {
		tmplData := TemplateData{
			FormErrors: make(map[string]string),
			User:       user,
		}

		if val, ok := err.(validation.Errors); ok {
			for field, valerr := range val {
				switch field {
				case "Name":
					tmplData.FormErrors["name"] = valerr.Error()
				case "Nickname":
					tmplData.FormErrors["nickname"] = valerr.Error()
				case "Email":
					tmplData.FormErrors["email"] = valerr.Error()
				case "City":
					tmplData.FormErrors["city"] = valerr.Error()
				case "About":
					tmplData.FormErrors["about"] = valerr.Error()
				}
			}
		}
		return app.RenderHTML(c, "addUser", tmplData)
	}

	err := app.Storage.Users.Add(&user)
	if err != nil {
		log.Println(err)
		return app.JsendError(c, err.Error())
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("session creation error")
		return app.JsendError(c, "session creation error")
	}

	sess.Values["flash"] = "User" + user.Nickname + " successfully created!"
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/users")
}
