package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestUserShow tests the UserShow function
func TestUserShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a user type
	userType, err := models.CreateUserType("TestUserType")

	// Create a user
	user := models.User{
		Email:      "test@example.com",
		UserName:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
		UserTypeID: uint64(userType.ID),
	}
	err = db.Database.Create(&user).Error
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/users/:id", controllers.UserShow)

	// Create a new request
	req, err := http.NewRequest("GET", "/users/"+helpers.UintToString(user.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// check username
	if !strings.Contains(rr.Body.String(), user.UserName) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.UserName)
	}
	// check firstname
	if !strings.Contains(rr.Body.String(), user.FirstName) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.FirstName)
	}
	// check lastname
	if !strings.Contains(rr.Body.String(), user.LastName) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.LastName)
	}
	// check user type
	if !strings.Contains(rr.Body.String(), userType.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), userType.Name)
	}

	// ADD MORE HERE ONCE PAGE IS MORE READY

	// Cleanup
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&userType)

}
