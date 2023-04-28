package controllers

import (
	"fmt"
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

// TestDiveEdit is a test for the DiveEdit controller
func TestDiveEdit(t *testing.T) {
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
	router.GET("/admin/dives/edit/:id", controllers.DiveEdit)

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create a test dive
	dive, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))

	// Make a http request to the test route
	req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/dives/edit/%d", dive.ID), nil)

	// Create a response recorder
	resp := httptest.NewRecorder()

	// Serve the request to the recorder
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	responseBody := resp.Body.String()
	if !strings.Contains(responseBody, "Test Dive") ||
		!strings.Contains(responseBody, dg1.Name) ||
		!strings.Contains(responseBody, dt1.Name) ||
		!strings.Contains(responseBody, bt1.Name) {
		t.Errorf("DiveEdit did not return the correct dive data")
	}

	// Clean up test dive
	db.Database.Unscoped().Delete(dive)

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}
