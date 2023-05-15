package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TestMarkersIndex is a test for the markers index page
func TestMarkersIndex(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a router with the test route
	funcMap := template.FuncMap{
		"mod":     func(i, j int) int { return i % j }, // used to order checkboxes for dives
		"shorten": h.Abbreviate,                        // used to abbreviate dive information. See helpers\abbreviate.go
		"seq":     h.Seq,                               // used to generate a sequence of numbers for the event show page
		"inc":     func(x int) int { return x + 1 },    // used to inc index on the event show page
	}

	// Create a gin router with the routes we need
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/markers", controllers.MarkersIndex)

	// Create two markers
	marker1, _ := models.MarkerCreate("Test Marker 1", "This is a small test description")
	marker2, _ := models.MarkerCreate("Test Marker 2", "This is a small test description")

	// Create a get request to /admin/markers
	req, err := http.NewRequest("GET", "/admin/markers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body for marker 1
	expected := marker1.Name
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	// Check the response body for marker 2
	expected = marker2.Name
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the response body for marker 1 description
	expected = marker1.Description
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	// Check the response body for marker 2 description
	expected = marker2.Description
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check if index page has the correct page text
	expected = "Markers"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&marker1)
	db.Database.Unscoped().Delete(&marker2)

}
