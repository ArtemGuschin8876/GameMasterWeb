package updateuserhandlertests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserHTMLRendering(t *testing.T) {
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

	formData := url.Values{
		"name":     {"New Name"},
		"nickname": {"newnickname"},
		"email":    {"newemail@example.com"},
		"city":     {"New City"},
		"about":    {"New about"},
		"image":    {"newimage.jpg"},
	}
	req := httptest.NewRequest(http.MethodPost, "/users/1", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set(echo.HeaderAccept, "text/html")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	templates, _ := application.ReadTemplateFromRootPath("../..")

	app := &application.Application{
		Storage: data.Storage{
			User: mockStorage,
		},
		Logger:    logger,
		Templates: templates,
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

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, "/users", rec.Header().Get("Location"))

	updatedUser, ok := mockStorage.Users["newnickname"]
	assert.True(t, ok)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "newnickname", updatedUser.Nickname)
	assert.Equal(t, "newemail@example.com", updatedUser.Email)
	assert.Equal(t, "New City", updatedUser.City)
	assert.Equal(t, "New about", updatedUser.About)
	assert.Equal(t, "newimage.jpg", updatedUser.Image)
}
