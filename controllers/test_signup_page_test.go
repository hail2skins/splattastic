package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestSignupPage(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.Default()

	// Load the templates
	r.LoadHTMLGlob("../templates/**/**")

	// Register the SignupPage route
	r.GET("/signup", SignupPage)

	// Create user types
	userType1 := models.UserType{Name: "Test User Type"}
	userType2 := models.UserType{Name: "Another Test Type"}
	db.Database.Create(&userType1)
	db.Database.Create(&userType2)

	// Create a request to the SignupPage route
	req, err := http.NewRequest(http.MethodGet, "/signup", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response has the correct status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the response contains the correct title
	assert.Contains(t, w.Body.String(), "<title>Signup</title>")

	// Cleanup
	db.Database.Unscoped().Delete(&userType1)
	db.Database.Unscoped().Delete(&userType2)
}
