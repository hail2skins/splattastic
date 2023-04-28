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

// TestDiveShow is a test for the DiveShow controller
func TestDiveShow(t *testing.T) {
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
	router.GET("/admin/dives/:id", controllers.DiveShow)

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Perform an HTTP GET request to the "/admin/dives/:id" route
	req, _ := http.NewRequest("GET", "/admin/dives/"+helpers.UintToString(dive1.ID), nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check if the response contains the created dive
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "Test Dive 1") {
		t.Errorf("DiveShow did not return the created dive")
	}
	if !strings.Contains(responseBody, "154") {
		t.Errorf("DiveShow did not return the created dive")
	}
	if !strings.Contains(responseBody, "2.5") {
		t.Errorf("DiveShow did not return the created dive")
	}
	if !strings.Contains(responseBody, "Dive Group 1") {
		t.Errorf("DiveShow did not return the created dive")
	}
	if !strings.Contains(responseBody, "Dive Type 1") {
		t.Errorf("DiveShow did not return the created dive")
	}
	if !strings.Contains(responseBody, "Board Type 1") {
		t.Errorf("DiveShow did not return the created dive")
	}

	// Clean up test dives
	db.Database.Unscoped().Delete(dive1)
	db.Database.Unscoped().Delete(dive2)

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
