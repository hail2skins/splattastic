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

// TestDiveGroupsIndex is a test for DiveGroupsIndex
func TestDiveGroupsIndex(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create two dive groups
	diveGroup1, _ := models.DiveGroupCreate("TestDiveGroup1")
	diveGroup2, _ := models.DiveGroupCreate("TestDiveGroup2")

	// Create a get request to /admin/divegroups
	req, err := http.NewRequest("GET", "/admin/divegroups", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divegroups", controllers.DiveGroupsIndex)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the dive groups are in the response
	if !strings.Contains(rr.Body.String(), diveGroup1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveGroup1.Name)
	}
	if !strings.Contains(rr.Body.String(), diveGroup2.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), diveGroup2.Name)
	}

	// Check if index page has the correct page text
	if !strings.Contains(rr.Body.String(), "Dive Groups") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Dive Groups")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&diveGroup1)
	db.Database.Unscoped().Delete(&diveGroup2)

}
