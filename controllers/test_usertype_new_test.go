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
	"github.com/hail2skins/splattastic/middlewares"
	m "github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeNew function to test the new user type page
func TestUserTypeNew(t *testing.T) {
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

	r.Use(m.AuthenticateUser())

	r.LoadHTMLGlob("../templates/**/**")

	r.POST("/login", Login)

	admin := r.Group("/admin", middlewares.RequireAdmin())
	{
		usertypes := admin.Group("/usertypes")
		{
			usertypes.GET("/new", UserTypeNew)
		}
	}

	// Create a test user type
	adminUserType := models.UserType{Name: "Admin"}
	db.Database.Create(&adminUserType)

	// Create a test user with hashed password and set its user type to admin
	adminUser, err := models.UserCreate(
		"admin@example.com",
		"adminpassword",
		"adminuser",
		"John",
		"Doe",
		"Admin",
	)
	assert.NoError(t, err)
	assert.NotNil(t, adminUser)
	// Set the admin flag to true
	adminUser.Admin = true
	db.Database.Save(&adminUser)

	// Login as the admin user
	form := url.Values{}
	form.Add("email", adminUser.Email)
	form.Add("password", "adminpassword")

	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Navigate to /admin/usertypes/new
	req, err = http.NewRequest(http.MethodGet, "/admin/usertypes/new", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the session cookie to the one from the login
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Cleanup
	db.Database.Unscoped().Delete(adminUser)
	db.Database.Unscoped().Delete(&adminUserType)
}
