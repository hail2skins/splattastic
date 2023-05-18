package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamTypeUpdate is a controller for testing the update of a team_type
func TestTeamTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

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
	r.POST("/admin/teamtypes/:id", controllers.TeamTypeUpdate)

	// Create a team_type
	tt1, _ := models.TeamTypeCreate("TestTeamType")

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "UpdatedTestTeamType")
	req, err := http.NewRequest("POST", "/admin/teamtypes/"+helpers.UintToString(tt1.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the updated team_type
	updatedTeamType, _ := models.TeamTypeShow(uint64(tt1.ID))

	// Check the updated team_type name
	if updatedTeamType.Name != "UpdatedTestTeamType" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedTeamType.Name, "UpdatedTestTeamType")
	}

	// Delete the team_type
	db.Database.Unscoped().Delete(&updatedTeamType)
}
