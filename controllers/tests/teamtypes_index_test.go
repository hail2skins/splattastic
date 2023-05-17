package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypesIndex is a test for the team_types index page
func TestTeamTypesIndex(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env variable to true for views requiring logged_in
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env variable

	// Create a router with the test route
	funcMap := template.FuncMap{
		"mod":     func(i, j int) int { return i % j }, // used to order checkboxes for dives
		"shorten": h.Abbreviate,                        // used to abbreviate dive information. See helpers\abbreviate.go
		"seq":     h.Seq,                               // used to generate a sequence of numbers for the event show page
		"inc":     func(x int) int { return x + 1 },    // used to inc index on the event show page
	}

	// Create a gin router with the routes we need
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.SetFuncMap(funcMap)
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/teamtypes", controllers.TeamTypesIndex)

	// Create two team_types
	tt1, _ := models.TeamTypeCreate("Test Team Type 1")
	tt2, _ := models.TeamTypeCreate("Test Team Type 2")

	// Create a get request to /admin/teamtypes
	req, err := http.NewRequest("GET", "/admin/teamtypes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body for team_type 1
	expected := tt1.Name
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	// Check response body for team_type 2
	expected = tt2.Name
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check if index page has correct page text
	expected = "Team Types"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&tt1)
	db.Database.Unscoped().Delete(&tt2)

}
