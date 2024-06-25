package api_handlers

import (
	"errors"
	"log"
	"net/http"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	var user data.User

	cc := c.(*application.AppContext)
	app := cc.App

	if err := c.Bind(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, app.JsendError(c, "invalid request payload"))
	}

	if err := user.ValidateUser(); err != nil {

		tmplData := data.TemplateData{
			FormErrors: make(map[string]string),
			User:       &user,
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

	err := app.Storage.User.Add(&user)
	if err != nil {
		log.Println(err)
		if errors.Is(err, data.ErrDuplicateEmail) {
			tmplData := data.TemplateData{
				FormErrors: map[string]string{"email": "Пользователь с таким email уже существует"},
				User:       &user,
			}
			return app.RenderHTML(c, "addUser", tmplData)
		}

		if errors.Is(err, data.ErrDuplicateNickname) {
			tmplData := data.TemplateData{
				FormErrors: map[string]string{"nickname": "Пользователь с таким nickname уже существует"},
				User:       &user,
			}
			return app.RenderHTML(c, "addUser", tmplData)

		}
		log.Println(err)
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
