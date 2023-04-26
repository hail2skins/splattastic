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
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveGroupCreate tests the DiveTypeCreate function
func TestDiveGroupCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "TestDiveGroup")
	req, err := http.NewRequest("POST", "/admin/divegroups", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divegroups", controllers.DiveGroupsIndex)
	r.POST("/admin/divegroups", controllers.DiveGroupCreate)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check the redirect location
	expected := "/admin/divegroups"
	if location := rr.Header().Get("Location"); location != expected {
		t.Errorf("handler returned unexpected redirect location: got %v want %v", location, expected)
	}

	// Check if the dive type is created in the database
	var createdDiveGroup models.DiveGroup
	if err := db.Database.Where("name = ?", "TestDiveGroup").First(&createdDiveGroup).Error; err != nil {
		t.Errorf("handler did not create dive group in database: %v", err)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&createdDiveGroup)

}
