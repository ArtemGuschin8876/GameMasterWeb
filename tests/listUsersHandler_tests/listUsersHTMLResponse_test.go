package listusershandlertests

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

func TestListUsersHTMLResponse(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	mockStorage := &data.MockUserStorage{
		Users: make(map[string]*data.User),
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	templates, _ := application.ReadTemplateFromRootPath("../..")

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

	err := api_handlers.ListUsers(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
}
