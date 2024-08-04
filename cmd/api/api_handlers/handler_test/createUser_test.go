package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	userJSON = `{
  "name": "maaad",
  "nickname": "sswssdaww",
  "email": "dawdasrgff@gmail.com",
  "city": "Asdr",
  "about": "asd dsad sadsadasdasdw wdxsdsadsa awd a s",
  "image": ""
	}`
)

func TestCreateUserJSONResponse(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockStorage := &data.MockUserStorage{
		Users: make(map[string]*data.User),
	}

	app := &application.Application{
		Storage: data.Storage{
			UserMock: mockStorage,
		},
	}

	appCtx := application.AppContext{
		Context: c,
		App:     app,
	}

	c.Set("app", appCtx)
	fmt.Println("appCtx.App.Storage:", appCtx.App.Storage)
	fmt.Println("appCtx.App.Storage.UserMock:", appCtx.App.Storage.UserMock)

	if assert.NoError(t, api_handlers.CreateUser(appCtx)) {
		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Contains(t, rec.Header().Get("Location"), "/users")
		assert.JSONEq(t, userJSON, rec.Body.String())
	}
}
