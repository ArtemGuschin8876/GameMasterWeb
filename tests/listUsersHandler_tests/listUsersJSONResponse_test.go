package listusershandlertests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListUsersJSONResponse(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUsers := []*data.User{
		{
			Name:     "User1",
			Nickname: "user1nick",
			Email:    "user1@example.com",
			City:     "City1",
			About:    "About User1About User1About User1About User1",
		},
		{
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

	for _, user := range mockUsers {
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
			"id": 0,
			"name": "User1",
			"nickname": "user1nick",
			"email": "user1@example.com",
			"city": "City1",
			"about": "About User1About User1About User1About User1",
			"image": ""
		  },
		  {
			"id": 1,
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
