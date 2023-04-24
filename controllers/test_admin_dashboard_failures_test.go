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
	db "github.com/hail2skins/splattastic/database"
	m "github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestAdminDashboardWithoutLogin tests the admin dashboard without logging in
func TestAdminDashboardWithoutLogin(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(m.AuthenticateUser())

	r.LoadHTMLGlob("../templates/**/**")

	r.GET("/", Home)
	r.GET("/login", LoginPage)
	r.POST("/login", Login)

	admin := r.Group("/admin", m.RequireAdmin())
	{
		admin.GET("/", AdminDashboard)

	}

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Create a test user with hashed password and set its user type to admin
	testUser, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"TestType",
	)
	assert.NoError(t, err)
	assert.NotNil(t, testUser)

	// Create a new HTTP request with an empty session cookie
	req, err := http.NewRequest(http.MethodGet, "/admin", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	// Send the request and get the response
	r.ServeHTTP(w, req)

	// Check the response status code
	if w.Code == http.StatusMovedPermanently {
		req, err = http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Send the request and get the response
		r.ServeHTTP(w, req)
	}

	// Check for specific alert
	expectedAlert := "Login"
	assert.Contains(t, w.Body.String(), expectedAlert)

	// cleanup
	db.Database.Unscoped().Delete(&testUserType)
	db.Database.Unscoped().Delete(testUser)
}

// Expected to fail due to admin protection on route
func TestAdminDashboardWithLogin(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(m.AuthenticateUser())

	r.LoadHTMLGlob("../templates/**/**")

	r.GET("/", Home)
	r.GET("/login", LoginPage)
	r.POST("/login", Login)

	admin := r.Group("/admin", m.RequireAdmin())
	{
		admin.GET("/", AdminDashboard)

	}

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	// Create a test user with hashed password and set its user type to admin
	testUser, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"TestType",
	)
	assert.NoError(t, err)
	assert.NotNil(t, testUser)

	// Create a new HTTP request to the login endpoint
	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("password", "testpassword")
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Get the session cookie from the response
	sessionCookie := w.Header().Get("Set-Cookie")

	// Create a new HTTP request with the session cookie
	req, err = http.NewRequest(http.MethodGet, "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", sessionCookie)
	w = httptest.NewRecorder()

	// Send the request and get the response
	r.ServeHTTP(w, req)

	// Check the response status code
	if w.Code == http.StatusMovedPermanently {
		req, err = http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Cookie", sessionCookie)
		w = httptest.NewRecorder()

		// Send the request and get the response
		r.ServeHTTP(w, req)
	}

	// Check for specific alert
	expectedAlert := "Logout"
	assert.Contains(t, w.Body.String(), expectedAlert)

	// cleanup
	db.Database.Unscoped().Delete(&testUserType)
	db.Database.Unscoped().Delete(testUser)

}
