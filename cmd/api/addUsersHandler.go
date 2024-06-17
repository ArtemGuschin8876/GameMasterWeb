package main

import (
	"log"
	"net/http"

	"gamemasterweb.net/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type TemplateData struct {
	Errors     []string
	FormErrors map[string]string
	User       data.Users
}

func (app *application) addUsersHandler(c echo.Context) error {

	var user data.Users

	data := TemplateData{
		FormErrors: make(map[string]string),
		User:       user,
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusOK, jsendError(c, "invalid request payload"))
	}

	if err := user.ValidateUsers(); err != nil {
		if val, ok := err.(validation.Errors); ok {
			for field, valerr := range val {
				switch field {
				case "Name":
					data.FormErrors["name"] = valerr.Error()
				case "Nickname":
					data.FormErrors["nickname"] = valerr.Error()
				case "Email":
					data.FormErrors["email"] = valerr.Error()
				case "City":
					data.FormErrors["city"] = valerr.Error()
				case "About":
					data.FormErrors["about"] = valerr.Error()
				}
			}
		}
		return app.renderHTML(c, "addUser", data)
	}

	err := app.storage.Users.Add(&user)
	if err != nil {
		log.Println(err)
		return jsendError(c, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/users/successfully")
}
