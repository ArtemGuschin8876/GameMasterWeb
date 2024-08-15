package createUserHandler_tests

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
	userJSON = `{
  "name": "User",
  "nickname": "UserNickname",
  "email": "UserTest@gmail.com",
  "city": "Usersk",
  "about": "testing json data for test handler",
  "image": ""
}`

	expectedJSON = `{
	"status": "success",
	"data": {
		"id": 0,
		"name": "User",
		"nickname": "UserNickname",
		"email": "UserTest@gmail.com",
		"city": "Usersk",
		"about": "testing json data for test handler",
		"image": ""
	}
}`
)

func TestCreateUserJSONResponse(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
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

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}
