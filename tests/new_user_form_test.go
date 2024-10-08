package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewUserForm(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users/new", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	templates, _ := application.ReadTemplateFromRootPath("..")

	app := &application.Application{
		Logger:    logger,
		Templates: templates,
	}

	appCtx := application.AppContext{
		Context: c,
		App:     app,
	}

	c.Set("app", appCtx)

	err := api_handlers.NewUserForm(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

}
