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

// TestDivesIndex is a test function for DiveIndex
func TestDivesIndex(t *testing.T) {
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
	router.GET("/admin/dives", controllers.DivesIndex)

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two test dives (similar to what you did in the model test)
	dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Perform an HTTP GET request to the "/admin/dives" route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/dives", nil)
	router.ServeHTTP(w, req)

	// Check if the response contains the created dives
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "Test Dive 1") || !strings.Contains(responseBody, "Test Dive 2") {
		t.Errorf("DivesIndex did not return the created dives")
	}
	if !strings.Contains(responseBody, "Test Dive Group 1") || !strings.Contains(responseBody, "Test Dive Group 2") {
		t.Errorf("DivesIndex did not return the created dive groups")
	}
	if !strings.Contains(responseBody, "Test Dive Type 1") || !strings.Contains(responseBody, "Test Dive Type 2") {
		t.Errorf("DivesIndex did not return the created dive types")
	}

	// Clean up test dives
	db.Database.Unscoped().Delete(dive1)
	db.Database.Unscoped().Delete(dive2)

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}
