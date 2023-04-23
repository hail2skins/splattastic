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
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
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
	r.POST("/login", Login)

	// Create a test user
	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Create a test user with hashed password
	user, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"Test User Type",
	)
	assert.NoError(t, err)
	assert.NotNil(t, user)

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

	// Check the response status code
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// Cleanup
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(usertype)
}

// TestLoginIncorrectEmail tests the login controller with an incorrect email which is same as an invalid or empty user that doesn't exist in the database
func TestLoginIncorrectEmail(t *testing.T) {
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

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/login", Login)

	// Create a test user
	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Create a test user with hashed password
	user, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"Test User Type",
	)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Create form data for login
	form := url.Values{}
	form.Add("email", "purple@pretty.com")
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

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid email or password")

	// Cleanup
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(usertype)
}

// TestLoginIncorrectPassword tests the login controller with an incorrect password
func TestLoginIncorrectPassword(t *testing.T) {
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

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/login", Login)

	// Create a test user
	usertype := models.UserType{Name: "Test User Type"}
	db.Database.Create(&usertype)

	// Create a test user with hashed password
	user, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		"Test User Type",
	)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Create form data for login
	form := url.Values{}
	form.Add("email", user.Email)
	form.Add("password", "funnywrongpassword")

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

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Password does not match")

	// Cleanup
	db.Database.Unscoped().Delete(user)
	db.Database.Unscoped().Delete(usertype)
}
