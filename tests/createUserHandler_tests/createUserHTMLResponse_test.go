package createUserHandler_tests

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHTMLResponse(t *testing.T) {
	e := echo.New()

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	writer.WriteField("name", "User")
	writer.WriteField("nickname", "UserNickname")
	writer.WriteField("email", "UserTest@gmail.com")
	writer.WriteField("city", "New York")
	writer.WriteField("about", "testing formdata for create_user handler")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

	mockStorage := &data.MockUserStorage{
		Users: make(map[string]*data.User),
	}

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

	err := api_handlers.CreateUser(appCtx)
	assert.NoError(t, err)

	expectedRedirectURL := "/users"

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	assert.Equal(t, expectedRedirectURL, rec.Header().Get("Location"))
}
