package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveTypesIndex tests the DiveTypesIndex function
func TestDiveTypesIndex(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create two dive types
	diveType1, _ := models.DiveTypeCreate("TestDiveType1", "Q")
	diveType2, _ := models.DiveTypeCreate("TestDiveType2", "R")

	// Create a get request to /admin/divetypes
	req, err := http.NewRequest("GET", "/admin/divetypes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divetypes", controllers.DiveTypesIndex)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the dive types are in the response
	if !strings.Contains(rr.Body.String(), diveType1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveType1.Name)
	}
	if !strings.Contains(rr.Body.String(), diveType2.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveType2.Name)
	}
	if !strings.Contains(rr.Body.String(), diveType1.Letter) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveType1.Letter)
	}
	if !strings.Contains(rr.Body.String(), diveType2.Letter) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveType2.Letter)
	}

	// Check if index page has the correct page text
	if !strings.Contains(rr.Body.String(), "Dive Types") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Dive Types")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveType1)
	db.Database.Unscoped().Delete(&diveType2)

}
