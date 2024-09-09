package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEditUserForm(t *testing.T) {
	e := echo.New()

	t.Run("HTML Response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/edit/5", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

		templates, _ := application.ReadTemplateFromRootPath("../..")

		app := &application.Application{
			Logger:    logger,
			Templates: templates,
		}

		appCtx := application.AppContext{
			Context: c,
			App:     app,
		}

		c.Set("app", appCtx)

		err := api_handlers.EditUserForm(appCtx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("ID doesn't exist, JSON Response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/edit/999", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		mockStorage := &data.MockUserStorage{
			Users: make(map[string]*data.User),
		}

		app := &application.Application{
			Storage: data.Storage{
				User: mockStorage,
			},
		}

		appCtx := application.AppContext{
			Context: c,
			App:     app,
		}

		c.Set("app", appCtx)

		err := api_handlers.EditUserForm(appCtx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "the requested resource could not be found")
	})

	t.Run("ID doesn't exist HTML Response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/edit/5", nil)
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

		err := api_handlers.EditUserForm(appCtx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
