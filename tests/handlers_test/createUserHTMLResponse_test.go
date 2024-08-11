package handlerstest

import (
	"bytes"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"../../static/ui/html/table.html",
		"../../static/ui/html/tableAllUsers.html",
		"../../static/ui/html/addUser.html",
		"../../static/ui/html/404.html",
		"../../static/ui/html/updateUserForms.html",
	}

	for _, page := range pages {
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[filepath.Base(page)] = ts
	}
	return cache, nil
}

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
	templateCache, err := newTemplateCache()
	if err != nil {
		panic(err)
	}

	app := &application.Application{
		Storage: data.Storage{
			User: mockStorage,
		},
		Logger:    logger,
		Templates: templateCache,
	}

	appCtx := application.AppContext{
		Context: c,
		App:     app,
	}

	c.Set("app", appCtx)

	err = api_handlers.CreateUser(appCtx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusSeeOther, rec.Code)
	expectedRedirectURL := "/users"
	assert.Equal(t, expectedRedirectURL, rec.Header().Get("Location"))
}
