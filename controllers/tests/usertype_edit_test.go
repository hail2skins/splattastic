package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeEdit renders the user type edit page
func TestUserTypeEdit(t *testing.T) {
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

	r.LoadHTMLGlob("../../templates/**/**")

	admin := r.Group("/admin")
	{
		admin.GET("/", controllers.AdminDashboard)

		// User types
		admin.GET("/usertypes/edit/:id", controllers.UserTypeEdit)

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

	req, err := http.NewRequest(http.MethodGet, "/admin/usertypes/edit/"+helpers.UintToString(testUserType.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assertions
	expectedText := "Editing User Type"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedText)

	// Cleanup
	db.Database.Unscoped().Delete(testUser)
	db.Database.Unscoped().Delete(&testUserType)

}
