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

// TestDiveCreate tests the creation of a new dive
// Requires: DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func TestDiveCreate(t *testing.T) {
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
	router.POST("/admin/dives", controllers.DiveCreate)

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// use form to create dive
	// Prepare form data for creating a new dive
	data := url.Values{}
	data.Set("name", "Test Dive")
	data.Set("number", "101")
	data.Set("difficulty", "2.5")
	data.Set("divetype_id", fmt.Sprintf("%d", dt1.ID))
	data.Set("divegroup_id", fmt.Sprintf("%d", dg1.ID))
	data.Set("boardtype_id", fmt.Sprintf("%d", bt1.ID))
	data.Set("boardheight_id", fmt.Sprintf("%d", bh1.ID))

	// Create a request and submit it to the router
	req, err := http.NewRequest(http.MethodPost, "/admin/dives", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("Failed to create a new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check if the response status code is 302 (redirect)
	if w.Code != http.StatusFound {
		t.Errorf("Expected status code 302, got %d", w.Code)
	}

	// Get the newly created dive and hard-delete it to avoid FK constraint errors
	newDiveID := getLastInsertedDiveID()
	if newDiveID != 0 {
		err := db.Database.Unscoped().Delete(&models.Dive{}, newDiveID).Error
		if err != nil {
			t.Fatalf("Error hard-deleting dive: %v", err)
		}
	}

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}

func getLastInsertedDiveID() uint {
	var dive models.Dive
	err := db.Database.Order("id desc").First(&dive).Error
	if err != nil {
		return 0
	}
	return dive.ID
}
