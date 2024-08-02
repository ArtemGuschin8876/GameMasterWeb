package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
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

func TestCreateUser(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

	c.Set("app", app)

	// if assert.NoError(t, api_handlers.CreateUser()) {
	// 	assert.Equal(t, http.StatusSeeOther, rec.Code)
	// 	assert.Equal(t, userJSON, rec.Body.String())
	// }
}
