package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserTypeUpdate(t *testing.T) {
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
		admin.GET("/usertypes", controllers.UserTypeIndex)
		admin.GET("/usertypes/edit/:id", controllers.UserTypeEdit)
		admin.POST("/usertypes/:id", controllers.UserTypeUpdate)
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

	// Prepare a request to update the user type's name
	updatedName := "UpdatedTestType"
	req, err := http.NewRequest(http.MethodPost, "/admin/usertypes/"+helpers.UintToString(testUserType.ID), strings.NewReader("name="+url.QueryEscape(updatedName)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check if the request was successful
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// Retrieve the updated user type from the database
	var updatedUserType models.UserType
	db.Database.First(&updatedUserType, testUserType.ID)

	// Check if the name was updated correctly
	assert.Equal(t, updatedName, updatedUserType.Name)

	// Cleanup
	db.Database.Unscoped().Delete(testUser)
	db.Database.Unscoped().Delete(&testUserType)
}
