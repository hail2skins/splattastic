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

// TestDiveTypeShow is a test function for the DiveTypeShow controller function
func TestDiveTypeShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a dive type
	diveType, err := models.DiveTypeCreate("TestDiveType")
	if err != nil {
		t.Fatal(err)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/admin/divetypes/"+helpers.UintToString(diveType.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divetypes/:id", controllers.DiveTypeShow)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the dive type is in the response
	if !strings.Contains(rr.Body.String(), diveType.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveType.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveType)
}
