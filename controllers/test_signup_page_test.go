package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignupPage(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.Default()

	// Load the templates
	r.LoadHTMLGlob("../templates/**/**")

	// Register the LoginPage route
	r.GET("/signup", SignupPage)

	// Create a request to the LoginPage route
	req, err := http.NewRequest(http.MethodGet, "/signup", nil)
	assert.NoError(t, err)

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response has the correct status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the response contains the correct title
	assert.Contains(t, w.Body.String(), "<title>Signup</title>")
}
