package updateuserhandlertests

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

// Эти сессии ебучие это беда какая-то=)
func TestUpdateUserSessionAndRedirect(t *testing.T) {
	e := echo.New()

	mockStorage := &data.MockUserStorage{
		Users: map[string]*data.User{
			"1": {
				ID:       1,
				Name:     "Old Name",
				Nickname: "oldnickname",
				Email:    "oldemail@example.com",
				City:     "Old City",
				About:    "Old about",
			},
		},
	}

	body := `{"name":"New Name","nickname":"newnickname","email":"newemail@example.com","city":"New City","about":"New about","image":"newimage.jpg"}`

	req := httptest.NewRequest(http.MethodPost, "/users/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

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

	err := api_handlers.UpdateUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/users", rec.Header().Get("Location"))
}
