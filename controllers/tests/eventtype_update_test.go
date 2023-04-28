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
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestEventTypeUpdate is a test for the EventTypeUpdate controller
func TestEventTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "")

	// Create a new event type
	eventType, err := models.EventTypeCreate("TestEventType")
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/eventtypes", controllers.EventTypesIndex)
	r.POST("/admin/eventtypes/:id", controllers.EventTypeUpdate)

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "UpdatedTestEventType")
	req, err := http.NewRequest("POST", "/admin/eventtypes/"+helpers.UintToString(eventType.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the updated event type
	updatedEventType, err := models.EventTypeShow(uint64(eventType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	expected := "UpdatedTestEventType"
	if updatedEventType.Name != expected {
		t.Errorf("expected updated name %v, got %v", expected, updatedEventType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedEventType)
}
