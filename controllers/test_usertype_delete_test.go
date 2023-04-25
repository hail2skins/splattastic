package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeDelete function to test soft delete of user type.
// TestUserTypeDelete function to test soft delete of user type.
func TestUserTypeDelete(t *testing.T) {
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
		admin.GET("/usertypes", UserTypeIndex)
		admin.DELETE("/usertypes/:id", UserTypeDelete)
	}

	// Create a test user type
	testUserType := models.UserType{Name: "TestType"}
	db.Database.Create(&testUserType)

	req, err := http.NewRequest(http.MethodDelete, "/admin/usertypes/"+helpers.UintToString(testUserType.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// Check if the user type was soft-deleted correctly
	var softDeletedUserType models.UserType
	db.Database.Unscoped().Where("id = ?", testUserType.ID).First(&softDeletedUserType)
	assert.NotNil(t, softDeletedUserType.DeletedAt)

	// Cleanup
	db.Database.Unscoped().Delete(&softDeletedUserType)
}
