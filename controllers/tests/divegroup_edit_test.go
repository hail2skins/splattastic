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

// TestDiveGroupEdit is the test for the DiveGroupEdit controller
func TestDiveGroupEdit(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a new dive group
	diveGroup, _ := models.DiveGroupCreate("TestDiveGroup")

	// Set up the router
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divegroups/edit/:id", controllers.DiveGroupEdit)

	// Set up the request
	req, err := http.NewRequest("GET", "/admin/divegroups/edit/"+helpers.UintToString(diveGroup.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set up the response recorder
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Edit Dive Group"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveGroup)

}
