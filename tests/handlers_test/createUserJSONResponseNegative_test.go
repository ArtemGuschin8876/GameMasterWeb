package handlerstest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	WrongUserJSON = `{
		"name": "",
		"nickname": "",
		"email": "",
		"city": "",
		"about": "",
		"image": ""
	  }`

	expectedErrorJSON = `{
  	"status": "error",
  	"data": {
    	"form_errors": {
      "about": "Это поле является обязательным",
      "city": "Это поле является обязательным",
      "email": "Это поле является обязательным",
      "name": "Это поле является обязательным",
      "nickname": "Это поле является обязательным"
    }
  }
}`
)

func TestCreateUserJSONResponseNegative(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(WrongUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

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

	err := api_handlers.CreateUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, expectedErrorJSON, rec.Body.String())
}
