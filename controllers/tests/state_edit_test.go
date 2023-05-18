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
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TestStateEdit is a controller for editing a state
func TestStateEdit(t *testing.T) {
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
	r.GET("/admin/states/edit/:id", controllers.StateEdit)

	// Create a state
	state, _ := models.StateCreate("Test State", "TS")

	// Create a get request to /admin/states/edit/:id
	req, err := http.NewRequest("GET", "/admin/states/edit/"+helpers.UintToString(state.ID), nil)
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

	// Check if the state name is in the response
	if !strings.Contains(rr.Body.String(), "Test State") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Test State")
	}
	// Check if state code is in the response
	if !strings.Contains(rr.Body.String(), "TS") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "TS")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&state)

}
