package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestSignup tests the signup controller success path
func TestSignupSuccess(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/signup", Signup)

	// Create user types
	userType1 := models.UserType{Name: "Test User Type"}
	userType2 := models.UserType{Name: "Another Test Type"}
	db.Database.Create(&userType1)
	db.Database.Create(&userType2)

	// Test data
	userTypeID := strconv.Itoa(int(userType1.ID))
	testData := "email=test@example.com&password=testpassword&confirm_password=testpassword&user_type=" + userTypeID + "&firstname=Test&lastname=User&username=testuser"

	// Create the request
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(testData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response is a redirect (status code 301)
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// Clean up
	var user models.User
	db.Database.Where("email = ?", "test@example.com").First(&user)
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&userType1)
	db.Database.Unscoped().Delete(&userType2)
}

// TestSignupEmailOrUsernameExists tests a signup controller failure path
func TestSignupEmailOrUsernameExists(t *testing.T) {
	// ... setup code from the previous test ...
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/signup", Signup)

	// Create user types
	userType1 := models.UserType{Name: "Test User Type"}
	userType2 := models.UserType{Name: "Another Test Type"}
	db.Database.Create(&userType1)
	db.Database.Create(&userType2)

	// Create a test user
	user := models.User{
		Email:     "test@example.com",
		UserName:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		UserType:  userType1,
	}
	err := db.Database.Create(&user).Error
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Test data
	testData := "email=test@example.com&password=testpassword&confirm_password=testpassword&user_type=Test User Type&firstname=Test&lastname=User&username=testuser"

	// Create the request
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(testData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response has status code 226 (IM Used)
	assert.Equal(t, http.StatusIMUsed, w.Code)
	assert.Contains(t, w.Body.String(), "Email or username already exists")

	// Clean up the test setup
	db.Database.Where("email = ?", "test@example.com").First(&user)
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&userType1)
	db.Database.Unscoped().Delete(&userType2)

}

// TestSignupPasswordsDoNotMatch tests a signup controller failure path
func TestSignupPasswordsDoNotMatch(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	// Sessions init
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("../templates/**/**")
	r.POST("/signup", Signup)

	// Create user types
	userType1 := models.UserType{Name: "Test User Type"}
	userType2 := models.UserType{Name: "Another Test Type"}
	db.Database.Create(&userType1)
	db.Database.Create(&userType2)

	// Test data with non-matching passwords
	testData := "email=test@example.com&password=testpassword&confirm_password=differentpassword&user_type=Test User Type&firstname=Test&lastname=User&username=testuser"

	// Create the request
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(testData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response has status code 406 (Not Acceptable)
	assert.Equal(t, http.StatusNotAcceptable, w.Code)
	assert.Contains(t, w.Body.String(), "Passwords do not match")

	// Clean up the test setup
	db.Database.Unscoped().Delete(&userType1)
	db.Database.Unscoped().Delete(&userType2)
}
