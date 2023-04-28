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

// TestDiveTypeCreate tests the DiveTypeCreate function
func TestDiveTypeCreate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "TestDiveType")
	data.Set("letter", "Q")
	req, err := http.NewRequest("POST", "/admin/divetypes", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divetypes", controllers.DiveTypesIndex)
	r.POST("/admin/divetypes", controllers.DiveTypeCreate)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check the redirect location
	expected := "/admin/divetypes"
	if location := rr.Header().Get("Location"); location != expected {
		t.Errorf("handler returned unexpected redirect location: got %v want %v", location, expected)
	}

	// Check if the dive type is created in the database
	var createdDiveType models.DiveType
	if err := db.Database.Where("name = ?", "TestDiveType").First(&createdDiveType).Error; err != nil {
		t.Errorf("handler did not create the dive type in the database: %v", err)
	}
	if err := db.Database.Where("letter = ?", "Q").First(&createdDiveType).Error; err != nil {
		t.Errorf("handler did not create the dive type in the database: %v", err)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&createdDiveType)

}
