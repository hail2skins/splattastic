package controllers

import (
	"fmt"
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

// TestDiveUpdate tests the dive update functionality
func TestDiveUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a router with the test route
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/**/**")
	router.POST("/admin/dives/:id", controllers.DiveUpdate) // Add the DiveUpdate route
	router.GET("/admin/dives", controllers.DivesIndex)      // Add the DiveIndex route

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// create a test dive
	dive, _ := models.DiveCreate("Test Dive", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))

	// use form to create updated dive
	data := url.Values{}
	data.Set("name", "Test Dive Updated")
	data.Set("number", "155")
	data.Set("difficulty", "2.6")
	data.Set("divegroup_id", fmt.Sprintf("%d", dg2.ID))
	data.Set("divetype_id", fmt.Sprintf("%d", dt2.ID))
	data.Set("boardtype_id", fmt.Sprintf("%d", bt2.ID))
	data.Set("boardheight_id", fmt.Sprintf("%d", bh2.ID))

	// Create a request and submit it to the router
	req, err := http.NewRequest(http.MethodPost, "/admin/dives/"+helpers.UintToString(dive.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create a new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Submit the request
	router.ServeHTTP(w, req)

	// Check the status code
	if status := w.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check the updated dive
	updatedDive, _ := models.DiveShow(uint64(dive.ID))
	if updatedDive.Name != "Test Dive Updated" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.Name, "Test Dive Updated")
	}
	if updatedDive.Number != 155 {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.Number, 155)
	}
	if updatedDive.Difficulty != 2.6 {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.Difficulty, 2.6)
	}
	if updatedDive.DiveTypeID != uint64(dt2.ID) {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.DiveTypeID, dt2.ID)
	}

	if updatedDive.DiveGroupID != uint64(dg2.ID) {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.DiveGroupID, dg2.ID)
	}
	if updatedDive.BoardTypeID != uint64(bt2.ID) {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.BoardTypeID, bt2.ID)
	}
	if updatedDive.BoardHeightID != uint64(bh2.ID) {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDive.BoardHeightID, bh2.ID)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedDive)

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}
