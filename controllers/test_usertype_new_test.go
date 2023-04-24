package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeNew function to test the new user type page
func TestUserTypeNew(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	r.LoadHTMLGlob("../templates/**/**")

	admin := r.Group("/admin")
	{
		admin.GET("/", AdminDashboard)

		// User types
		admin.GET("/usertypes/new", UserTypeNew)
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

	// Log in the test user by setting a session
	req, err := http.NewRequest(http.MethodGet, "/admin/usertypes/new", nil)
	assert.NoError(t, err)

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the UserTypeIndex function with the request and response recorder
	r.ServeHTTP(w, req)

	// Check the response status code
	expectedText := "New User Type"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedText)

	// Cleanup
	db.Database.Unscoped().Delete(adminUser)
	db.Database.Unscoped().Delete(&adminUserType)
}
