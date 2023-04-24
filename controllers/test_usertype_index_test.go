package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestUserTypeIndex(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")

	admin := r.Group("/admin")
	{
		admin.GET("/", AdminDashboard)

		// User types
		admin.GET("/usertypes", UserTypeIndex)

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

	// Log in the test user by setting a session
	req, err := http.NewRequest(http.MethodGet, "/usertypes", nil)
	assert.NoError(t, err)

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the UserTypeIndex function with the request and response recorder
	r.ServeHTTP(w, req)

	// Check the response status code
	expectedType := "TestType"
	//assert.Equal(t, http.StatusOK, w1.Code)
	assert.Contains(t, w.Body.String(), expectedType)

	// cleanup
	db.Database.Unscoped().Delete(&testUserType)
	db.Database.Unscoped().Delete(testUser)
}
