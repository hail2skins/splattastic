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

// TestDiveGroupUpdate is a test for the DiveGroupUpdate function
func TestDiveGroupUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "")

	// Create a new dive group
	diveGroup, err := models.DiveGroupCreate("TestDiveGroup")
	if err != nil {
		t.Fatal(err)
	}

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "UpdatedTestDiveGroup")
	req, err := http.NewRequest("POST", "/admin/divegroups/update/"+helpers.UintToString(diveGroup.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divegroups", controllers.DiveGroupsIndex)
	r.POST("/admin/divegroups/update/:id", controllers.DiveGroupUpdate)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the updated dive group
	updatedDiveGroup, err := models.DiveGroupShow(uint64(diveGroup.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	expected := "UpdatedTestDiveGroup"
	if updatedDiveGroup.Name != expected {
		t.Errorf("expected updated name %v, got %v", expected, updatedDiveGroup.Name)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedDiveGroup)
}
