package api_handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo-contrib/session"
)

func UpdateUser(c application.AppContext) error {
	app := c.App

	id, err := app.ReadIDParam(c)
	if err != nil {
		zeroLog.Err(err).Msg("error reading IDParam")
		app.JsendError(c, "the requested resource could not be found")
	}

	user, err := app.Storage.User.Get(id)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.JsendError(c, "the requested resource could not be found")
		default:
			app.JsendError(c, "the server was unable to process your request")
		}
		return err
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
		zeroLog.Err(err).Msg("error binding")
		return app.JsendError(c, "database error")
	}

	user.Name = input.Name
	user.Nickname = input.Nickname
	user.Email = input.Email
	user.City = input.City
	user.About = input.About
	user.Image = input.Image

	jsonData := data.JsonData{
		FormErrors: make(map[string]string),
	}

	if err := user.ValidateUser(); err != nil {
		zeroLog.Err(err).Msg("error validation user")
		tmplData := data.TemplateData{
			FormErrors: make(map[string]string),
			User:       user,
		}

		if valErrors, ok := err.(validation.Errors); ok {
			for field, valerr := range valErrors {
				mappedField := strings.ToLower(field)
				tmplData.FormErrors[mappedField] = valerr.Error()
				jsonData.FormErrors[mappedField] = valerr.Error()
			}
		}

		if c.Request().Header.Get("Accept") == "application/json" {
			return app.JsonError(c, jsonData)
		} else {
			return app.RenderHTML(c, "updateUserForms", tmplData)
		}
	}

	err = app.Storage.User.Update(user)
	if err != nil {
		zeroLog.Err(err).Msg("error updating user on database")
		return app.JsendError(c, "error updating user")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		zeroLog.Err(err).Msg("session creation error")
		panic("session creation error")
	}

	sess.Values["flash"] = user.Nickname + " successfully updated!"
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		zeroLog.Err(err).Msg("session saving error")
		return err
	}
	if c.Request().Header.Get("Accept") == "application/json" {
		return app.JsendSuccess(c, user)
	} else {
		return c.Redirect(http.StatusSeeOther, "/users")
	}
}
