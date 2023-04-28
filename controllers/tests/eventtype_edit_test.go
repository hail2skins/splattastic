package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeEdit is the controller for the event type edit page
func TestEventTypeEdit(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create an event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/eventtypes/edit/:id", controllers.EventTypeEdit)

	// Create a new request
	req, err := http.NewRequest("GET", "/admin/eventtypes/edit/"+helpers.UintToString(eventType.ID), nil)
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

	// Check if the event type is in the response
	if !strings.Contains(rr.Body.String(), eventType.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), eventType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&eventType)
}
