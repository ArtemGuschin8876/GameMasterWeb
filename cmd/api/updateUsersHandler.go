package main

import (
	"errors"
	"log"
	"net/http"

	"gamemasterweb.net/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *application) updateUsersHandler(c echo.Context) error {

	method := c.FormValue("_method")
	if method != "" && method != "PUT" {
		return jsendError(c, "invalid method")
	}

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
		Name     string `json:"name" form:"name"`
		Nickname string `json:"nickname" form:"nickname"`
		Email    string `json:"email" form:"email"`
		City     string `json:"city" form:"city"`
		About    string `json:"about" form:"about"`
		Image    string `json:"image" form:"image"`
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

	tmplData := TemplateData{
		U:          user,
		FormErrors: make(map[string]string),
	}

	if err := user.ValidateUsers(); err != nil {
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
		return app.renderHTML(c, "updateUserForms", tmplData)
	}

	err = app.storage.Users.Update(user)
	if err != nil {
		return jsendError(c, "error updating user")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("session creation error")
		panic("session creation error")
	}

	sess.Values["flash"] = user.Nickname + " successfully updated!"
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/users")
}
