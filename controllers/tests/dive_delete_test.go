package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestDiveDelete tests the deletion of a dive
func TestDiveDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create a dive
	dive, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))

	// Create a router with the test route
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.LoadHTMLGlob("../../templates/**/**")
	router.DELETE("/admin/dives/:id", controllers.DiveDelete)

	// Perform an HTTP DELETE request to the "/admin/dives/:id" route
	req, err := http.NewRequest("DELETE", "/admin/dives/"+helpers.UintToString(dive.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Verify the dive was deleted
	_, err = models.DiveShow(uint64(dive.ID))
	if err == nil {
		t.Errorf("Dive with ID %d not deleted", dive.ID)
	}

	// Teardown
	db.Database.Unscoped().Delete(&dive)

	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
