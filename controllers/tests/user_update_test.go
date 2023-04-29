package controllers

import (
	"fmt"
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
)

// TestUserUpdate is a test function for the UserUpdate controller
func TestUserUpdate(t *testing.T) {
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
	router.POST("/users/:id", controllers.UserUpdate) // Add the UserUpdate route

	// Create two test usertypes
	ut1, _ := models.CreateUserType("Test User Type 1")
	ut2, _ := models.CreateUserType("Test User Type 2")

	// Create a test user
	email := "test@example.com"
	password := "testpassword"
	firstName := "John"
	lastName := "Doe"
	userName := "testuser"
	usertypeName := ut1.Name

	user, err := models.UserCreate(email, password, firstName, lastName, userName, usertypeName)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	// use form to create updated user
	data := url.Values{}
	data.Set("email", "test1@example.com")
	data.Set("firstname", "Jane")
	data.Set("lastname", "Dime")
	data.Set("username", "testuser1")
	data.Set("user_type_id", fmt.Sprintf("%d", ut2.ID))

	// Create an http request and submit it to the router
	req, err := http.NewRequest(http.MethodPost, "/users/"+helpers.UintToString(user.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create a new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Submit the request to the router
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusFound {
		t.Errorf("expected status %v but got %v", http.StatusFound, w.Code)
	}

	// Check the updated user
	updatedUser, _ := models.UserShow(uint64(user.ID))
	if updatedUser.Email != "test1@example.com" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedUser.Email, "test1@example.com")
	}
	if updatedUser.FirstName != "Jane" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedUser.FirstName, "Jane")
	}
	if updatedUser.LastName != "Dime" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedUser.LastName, "Dime")
	}
	if updatedUser.UserName != "testuser1" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedUser.UserName, "testuser1")
	}
	if updatedUser.UserType.Name != "Test User Type 2" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedUser.UserType.Name, "Test User Type 2")
	}

	// Delete the test user
	db.Database.Unscoped().Delete(&user)

	// Delete the test usertypes
	db.Database.Unscoped().Delete(&ut1)
	db.Database.Unscoped().Delete(&ut2)

}
