package updateuserhandlertests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	WrongUserFormData = url.Values{
		"name":     {""},
		"nickname": {""},
		"email":    {""},
		"city":     {""},
		"about":    {""},
		"image":    {""},
	}

	expectedErrorJSON = `{
        "status": "error",
        "data": {
            "form_errors": {
                "name": "Это поле является обязательным",
                "nickname": "Это поле является обязательным",
                "email": "Это поле является обязательным",
                "city": "Это поле является обязательным",
                "about": "Это поле является обязательным"
            }
        }
    }`
)

func TestUpdateUserJSONResponseNegative(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/users/1", strings.NewReader(WrongUserFormData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	mockStorage := &data.MockUserStorage{
		Users: map[string]*data.User{
			"1": {
				ID:       1,
				Name:     "Old Name",
				Nickname: "oldnickname",
				Email:    "oldemail@example.com",
				City:     "Old City",
				About:    "Old about",
				Image:    "oldimage.jpg",
			},
		},
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
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := api_handlers.UpdateUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, expectedErrorJSON, rec.Body.String())
}
