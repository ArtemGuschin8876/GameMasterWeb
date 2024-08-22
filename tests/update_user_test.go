package tests

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
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	e := echo.New()

	t.Run("JSON Response", func(t *testing.T) {
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
	})
	t.Run("Negative JSON", func(t *testing.T) {

		WrongUserFormData := url.Values{
			"name":     {""},
			"nickname": {""},
			"email":    {""},
			"city":     {""},
			"about":    {""},
			"image":    {""},
		}

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

		expectedErrorJSON := `{
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

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, expectedErrorJSON, rec.Body.String())
	})
	t.Run("HTML Response", func(t *testing.T) {
		mockStorage := &data.MockUserStorage{
			Users: map[string]*data.User{
				"1": {
					ID:       1,
					Name:     "Old Name",
					Nickname: "oldnickname",
					Email:    "oldemail@example.com",
					City:     "Old City",
					About:    "Old adsaadsadssadsadabout",
					Image:    "",
				},
			},
		}

		formData := url.Values{
			"name":     {"NewName"},
			"nickname": {"newnickname"},
			"email":    {"newemail@example.com"},
			"city":     {"NewCity"},
			"about":    {"New about with enough length"},
			"image":    {"newimage.jpg"},
		}

		req := httptest.NewRequest(http.MethodPost, "/users/1", strings.NewReader(formData.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		req.Header.Set(echo.HeaderAccept, echo.MIMETextHTML)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.Set("_session_store", sessions.NewCookieStore([]byte("secret")))

		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

		templates, err := application.ReadTemplateFromRootPath("..")
		if err != nil {
			t.Fatalf("Error loading templates: %v", err)
		}

		app := &application.Application{
			Storage:   data.Storage{User: mockStorage},
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

		err = api_handlers.UpdateUser(appCtx)
		if err != nil {
			t.Fatal("Handler error:", err)
		}

		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/users", rec.Header().Get("Location"))

		updatedUser, ok := mockStorage.Users["newnickname"]
		if !ok {
			t.Fatal("User should exist in mock storage with new nickname")
		}
		assert.Equal(t, "NewName", updatedUser.Name)
		assert.Equal(t, "newnickname", updatedUser.Nickname)
		assert.Equal(t, "newemail@example.com", updatedUser.Email)
		assert.Equal(t, "NewCity", updatedUser.City)
		assert.Equal(t, "New about with enough length", updatedUser.About)
		assert.Equal(t, "newimage.jpg", updatedUser.Image)
	})
}
