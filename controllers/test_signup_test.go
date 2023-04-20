package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
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

	// Test data
	testData := "email=test@example.com&password=testpassword&confirm_password=testpassword&user_type=Athlete&firstname=Test&lastname=User&username=testuser"

	// Create the request
	req, _ := http.NewRequest(http.MethodPost, "/signup", strings.NewReader(testData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response is a redirect (status code 301)
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// Clean up the user
	user, err := models.GetUserByEmail("test@example.com")
	if err == nil && user != nil {
		db.Database.Unscoped().Delete(user)
	}
}
