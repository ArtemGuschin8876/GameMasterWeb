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
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestShowUser(t *testing.T) {

	e := echo.New()

	t.Run("JSON Response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

		mockStorage := &data.MockUserStorage{
			Users: map[string]*data.User{
				"1": {
					ID:       1,
					Name:     "Test User",
					Nickname: "testuser",
					Email:    "testuser@example.com",
					City:     "Test City",
					About:    "This is a test user.",
				},
			},
		}

		c.SetParamNames("id")
		c.SetParamValues("1")

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

		err := api_handlers.ShowUser(appCtx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		expected := `{
		"status": "success",
		"data": {
			"id": 1,
			"name": "Test User",
			"nickname": "testuser",
			"email": "testuser@example.com",
			"city": "Test City",
			"about": "This is a test user.",
			"image": ""
		}
	}`
		assert.JSONEq(t, expected, rec.Body.String())
	})

	t.Run("HTML Resposne", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		mockStorage := &data.MockUserStorage{
			Users: make(map[string]*data.User),
		}

		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

		templates, _ := application.ReadTemplateFromRootPath("..")

		app := &application.Application{
			Storage: data.Storage{
				User: mockStorage,
			},
			Logger:    logger,
			Templates: templates,
		}

		appCtx := application.AppContext{
			Context: c,
			App:     app,
		}

		c.Set("app", appCtx)

		err := api_handlers.ShowUser(appCtx)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
