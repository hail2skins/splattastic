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
)

// TestDiveNew renders the new dive form
// Requires: DiveGroupGet, DiveTypeGet, BoardTypeGet, BoardHeightGet
func TestDiveNew(t *testing.T) {
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
	router.GET("/admin/dives/new", controllers.DiveNew)

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	req, _ := http.NewRequest("GET", "/admin/dives/new", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	body := resp.Body.String()
	//fmt.Println(body) // debugging

	if !strings.Contains(body, "New Dive") {
		t.Error("Expected body to contain 'New Dive'")
	}

	// Check if the response contains the test data values
	if !strings.Contains(body, dg1.Name) || !strings.Contains(body, dg2.Name) {
		t.Error("Expected body to contain dive group names")
	}

	if !strings.Contains(body, dt1.Name) || !strings.Contains(body, dt2.Name) {
		t.Error("Expected body to contain dive type names")
	}

	if !strings.Contains(body, bt1.Name) || !strings.Contains(body, bt2.Name) {
		t.Error("Expected body to contain board type names")
	}

	if !strings.Contains(body, fmt.Sprintf("%s Meter(s)", formatHeight(bh1.Height))) {
		t.Errorf("Expected body to contain board height %s Meter(s)", formatHeight(bh1.Height))
	}

	if !strings.Contains(body, fmt.Sprintf("%s Meter(s)", formatHeight(bh2.Height))) {
		t.Errorf("Expected body to contain board height %s Meter(s)", formatHeight(bh2.Height))
	}

	// Clean up test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}

func formatHeight(height float32) string {
	s := fmt.Sprintf("%.1f", height)
	if strings.HasSuffix(s, ".0") {
		return strings.TrimSuffix(s, ".0")
	}
	return s
}
