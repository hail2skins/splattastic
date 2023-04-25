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

func TestAdminDashboard(t *testing.T) {
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

	// Create a new HTTP request to the admin dashboard endpoint
	req, err := http.NewRequest(http.MethodGet, "/admin/", nil)
	assert.NoError(t, err)

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the AdminDashboard function with the request and response recorder
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the HTML response body for the expected title
	expectedTitle := "Admin Dashboard"
	assert.Contains(t, w.Body.String(), expectedTitle)

	// Cleanup
	db.Database.Unscoped().Delete(testUser)
	db.Database.Unscoped().Delete(&testUserType)
}
