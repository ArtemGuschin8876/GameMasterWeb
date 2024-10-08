package api_handlers

import (
	"errors"
	"net/http"
	"strings"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"gamemasterweb.net/internal/logger"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo-contrib/session"
)

var zeroLog = logger.NewLogger()

func CreateUser(c application.AppContext) error {
	var user data.User
	app := c.App

	if err := c.Bind(&user); err != nil {
		zeroLog.Err(err).Msg("Problem with bind")
		return app.Respond(c, http.StatusBadRequest, app.JsendError(c, "invalid request payload"), "addUser", data.TemplateData{
			FormErrors: map[string]string{"error": "invalid request payload"},
			User:       &user,
		})
	}

	jsonData := data.JsonData{
		FormErrors: make(map[string]string),
	}

	if err := user.ValidateUser(); err != nil {
		zeroLog.Err(err).Msg("Problem with validation func")
		tmplData := data.TemplateData{
			FormErrors: make(map[string]string),
			User:       &user,
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
			return app.RenderHTML(c, "addUser", tmplData)
		}
	}

	err := app.Storage.User.Add(&user)
	if err != nil {
		zeroLog.Err(err).Msg("Error adding a user to the database")
		if errors.Is(err, data.ErrDuplicateEmail) {
			tmplData := data.TemplateData{
				FormErrors: map[string]string{"email": "Пользователь с таким email уже существует"},
			}
			return app.Respond(c, http.StatusConflict, tmplData, "addUser", tmplData)
		}

		if errors.Is(err, data.ErrDuplicateNickname) {
			tmplData := data.TemplateData{
				FormErrors: map[string]string{"nickname": "Пользователь с таким nickname уже существует"},
			}
			return app.Respond(c, http.StatusBadRequest, tmplData, "addUser", tmplData)
		}

		return app.Respond(c, http.StatusInternalServerError, app.JsendError(c, "internal server error"), "addUser", data.TemplateData{
			FormErrors: map[string]string{"error": "internal server error"},
		})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		zeroLog.Err(err).Msg("Error receiving session")
		return app.Respond(c, http.StatusInternalServerError, app.JsendError(c, "internal server error"), "addUser", data.TemplateData{
			FormErrors: map[string]string{"error": "internal server error"},
			User:       &user,
		})
	}

	sess.Values["flash"] = "User " + user.Nickname + " successfully created!"
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		zeroLog.Err(err).Msg("Error saving session")
		return app.Respond(c, http.StatusInternalServerError, app.JsendError(c, "internal server error"), "addUser", data.TemplateData{
			FormErrors: map[string]string{"error": "internal server error"},
			User:       &user,
		})
	}

	if c.Request().Header.Get("Accept") == "application/json" {
		return app.JsendSuccess(c, user)
	} else {
		return c.Redirect(http.StatusSeeOther, "/users")
	}
}
