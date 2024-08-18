package deleteuserhandlertests

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

func TestDeleteUserJSONResponse(t *testing.T) {
	e := echo.New()

	mockStorage := &data.MockUserStorage{
		Users: map[string]*data.User{
			"1": {
				ID:       1,
				Name:     "Test User",
				Nickname: "testnickname",
				Email:    "testemail@example.com",
				City:     "Test City",
				About:    "Test about",
				Image:    "testimage.jpg",
			},
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

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

	err := api_handlers.DeleteUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	expectedJson := `{"status": "success","data": {"message": "user successfully deleted"}}`
	assert.JSONEq(t, expectedJson, rec.Body.String())

	_, ok := mockStorage.Users["1"]
	assert.False(t, ok)
}
