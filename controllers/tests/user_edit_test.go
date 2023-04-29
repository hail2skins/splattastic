package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserEdit is a test for the UserEdit controller
func TestUserEdit(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a router with the test route
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/**/**")
	router.GET("/user/edit/:id", controllers.UserEdit)

	// Create two user types
	ut1, _ := models.CreateUserType("Test User Type 1")

	// Create a test user
	// Create a test user with hashed password
	user, err := models.UserCreate(
		"test@example.com",
		"testpassword",
		"testuser",
		"John",
		"Doe",
		ut1.Name,
	)
	if err != nil {
		t.Errorf("Error creating test user: %v", err)
	}

	// Make a http request to the test route
	req, err := http.NewRequest("GET", fmt.Sprintf("/user/edit/%d", user.ID), nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request to the recorder
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, user.Email) ||
		!strings.Contains(responseBody, user.UserName) ||
		!strings.Contains(responseBody, user.FirstName) ||
		!strings.Contains(responseBody, user.LastName) ||
		!strings.Contains(responseBody, ut1.Name) {
		t.Errorf("UserEdit did not return the correct user data")
	}

	// Cleanup user
	db.Database.Unscoped().Delete(&user)

	// Cleanup user type
	db.Database.Unscoped().Delete(&ut1)

}
