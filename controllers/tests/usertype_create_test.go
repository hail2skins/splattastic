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
	h "github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserTypeCreate(t *testing.T) {
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
		admin.GET("/usertypes/new", controllers.UserTypeNew)
		admin.POST("/usertypes", controllers.UserTypeCreate)

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

	// Create a request with POST method and form data
	form := url.Values{}
	form.Add("name", "Test UserType")
	req, err := http.NewRequest(http.MethodPost, "/admin/usertypes", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the UserTypeCreate function directly with the request and response recorder
	controllers.UserTypeCreate(h.GinContext(req, w))

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the user type was created in the database
	var userType models.UserType
	db.Database.Where("name = ?", "Test UserType").First(&userType)
	assert.Equal(t, "Test UserType", userType.Name)

	// Cleanup
	db.Database.Unscoped().Delete(&userType)
	db.Database.Unscoped().Delete(testUser)
	db.Database.Unscoped().Delete(&testUserType)
}
