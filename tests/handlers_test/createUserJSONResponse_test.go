package handlerstest

import (
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
  "name": "dsadsaw",
  "nickname": "asweadasdqd",
  "email": "yreuwdasisadsadadeirot@gmail.com",
  "city": "Asdr",
  "about": "asd dsad sadsadasdasdw wdxsdsadsa awd a s",
  "image": ""
}`
)

func TestCreateUserJSONResponse(t *testing.T) {

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
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

	err := api_handlers.CreateUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

}
