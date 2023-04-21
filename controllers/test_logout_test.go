package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	th "github.com/hail2skins/splattastic/controllers/testhelpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	// ... setup code from the previous test ...
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/logout", Logout)
	r.POST("/login", Login)

	// Create a test user
	user := th.CreateTestUser()

	// Create form data for login
	form := url.Values{}
	form.Add("email", user.Email)
	form.Add("password", "testpassword")

	// Create a new HTTP request with form data
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Serve the request to the router
	r.ServeHTTP(w, req)

	// Create a new HTTP request to logout
	reql, err := http.NewRequest(http.MethodPost, "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the session cookie from the previous response to the request
	cookies := w.Header()["Set-Cookie"]
	if len(cookies) == 0 {
		t.Fatal("No cookies in the response")
	}
	reql.Header.Set("Cookie", cookies[0])

	// Send the request to the server
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reql)

	// Check that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Successfully logged out")

	// Cleanup
	db.Database.Unscoped().Delete(user)

}
