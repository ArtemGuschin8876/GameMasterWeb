package tests

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gamemasterweb.net/cmd/api/api_handlers"
	"gamemasterweb.net/internal/application"
	"gamemasterweb.net/internal/data"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()

	t.Run("JSON Response", func(t *testing.T) {

		userJSON := `{
			"name": "User",
			"nickname": "UserNickname",
			"email": "UserTest@gmail.com",
			"city": "Usersk",
			"about": "testing json data for test handler",
			"image": ""
		  }`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

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

		expectedJSON := `{
			"status": "success",
			"data": {
				"id": 0,
				"name": "User",
				"nickname": "UserNickname",
				"email": "UserTest@gmail.com",
				"city": "Usersk",
				"about": "testing json data for test handler",
				"image": ""
			}
		}`

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	})

	t.Run("HTML Response", func(t *testing.T) {
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

		templates, _ := application.ReadTemplateFromRootPath("..")

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
	})

	t.Run("Negative JSON", func(t *testing.T) {

		WrongUserJSON := `{
			"name": "",
			"nickname": "",
			"email": "",
			"city": "",
			"about": "",
			"image": ""
		  }`

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(WrongUserJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

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

		expectedErrorJSON := `{
			"status": "error",
			"data": {
			  "form_errors": {
			"about": "Это поле является обязательным",
			"city": "Это поле является обязательным",
			"email": "Это поле является обязательным",
			"name": "Это поле является обязательным",
			"nickname": "Это поле является обязательным"
		  }
		}
	  }`

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, expectedErrorJSON, rec.Body.String())
	})
}
