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

// TestEventTypesIndex is the controller for the event types index page
func TestEventTypesIndex(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create two event types
	eventType1, _ := models.EventTypeCreate("TestEventType1")
	eventType2, _ := models.EventTypeCreate("TestEventType2")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/eventtypes", controllers.EventTypesIndex)

	// Create a get request to /admin/eventtypes
	req, err := http.NewRequest("GET", "/admin/eventtypes", nil)
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

	// Check if the event types are in the response
	if !strings.Contains(rr.Body.String(), eventType1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), eventType1.Name)
	}
	if !strings.Contains(rr.Body.String(), eventType2.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), eventType2.Name)
	}

	// Check if index page has the correct page text
	if !strings.Contains(rr.Body.String(), "Event Types") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Event Types")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&eventType1)
	db.Database.Unscoped().Delete(&eventType2)

}
