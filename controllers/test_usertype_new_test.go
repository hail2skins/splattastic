package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUserTypeNew function to test the new user type page
func TestUserTypeNew(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../templates/**/**")
	usertypes := r.Group("/usertypes")
	{
		usertypes.GET("/new", UserTypeNew)
	}
	// Assertions
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/usertypes/new", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New User Type")
}
