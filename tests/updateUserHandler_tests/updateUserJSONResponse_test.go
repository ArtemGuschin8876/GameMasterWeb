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
				Image:    "",
			},
		},
	}

	body := `{"name":"New Name","nickname":"newnickname","email":"newemail@example.com","city":"New City","about":"New aboutaboutaboutaboutaboutaboutabout","image":"newimage.jpg"}`
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

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	expectedResponse := `{"status":"success","data":{"id":1,"name":"New Name","nickname":"newnickname","email":"newemail@example.com","city":"New City","about":"New aboutaboutaboutaboutaboutaboutabout","image":"newimage.jpg"}}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	updatedUser, ok := mockStorage.Users["newnickname"]
	assert.True(t, ok)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "newnickname", updatedUser.Nickname)
	assert.Equal(t, "newemail@example.com", updatedUser.Email)
	assert.Equal(t, "New City", updatedUser.City)
	assert.Equal(t, "New aboutaboutaboutaboutaboutaboutabout", updatedUser.About)
	assert.Equal(t, "newimage.jpg", updatedUser.Image)

}
