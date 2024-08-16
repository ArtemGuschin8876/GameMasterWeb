package listusershandlertests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListUsersJSONResponse(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	mockUsers := []*data.User{
		{
			ID:       1,
			Name:     "User1",
			Nickname: "user1nick",
			Email:    "user1@example.com",
			City:     "City1",
			About:    "About User1About User1About User1About User1",
		},
		{
			ID:       2,
			Name:     "User2",
			Nickname: "user2nick",
			Email:    "user2@example.com",
			City:     "City2",
			About:    "About User2About User2About User2About User2",
		},
	}

	mockStorage := &data.MockUserStorage{
		Users: make(map[string]*data.User),
	}

	for i, user := range mockUsers {
		user.ID = int64(i + 1)
		mockStorage.Users[user.Nickname] = user
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

	err := api_handlers.ListUsers(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	expectedJSON := `{
		"status": "success",
		"data": [
		  {
			"id": 1,
			"name": "User1",
			"nickname": "user1nick",
			"email": "user1@example.com",
			"city": "City1",
			"about": "About User1About User1About User1About User1",
			"image": ""
		  },
		  {
			"id": 2,
			"name": "User2",
			"nickname": "user2nick",
			"email": "user2@example.com",
			"city": "City2",
			"about": "About User2About User2About User2About User2",
			"image": ""
		  }
		]
	  }`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}
