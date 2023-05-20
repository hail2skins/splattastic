package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/models"
)

// TestTeamNew tests the team new page
func TestTeamNew(t *testing.T) {
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
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.SetFuncMap(funcMap)

	router.LoadHTMLGlob("../../templates/**/**")
	router.GET("/user/:id/team/new", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.TeamNew(c)
	})

	// Create a user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)
	// Set userID
	userID := user.ID

	// Create a team type
	tt1, _ := models.TeamTypeCreate("Test Team Type")

	// Create a state
	state, _ := models.StateCreate("Test State", "TS")

	// Send the user to the session and through middleware
	helpers.SetupSession(router, uint(userID))
	router.Use(middlewares.AuthenticateUser())

	// Create a request to the route /user/:id/team/new
	req, err := http.NewRequest("GET", "/user/"+helpers.UintToString(userID)+"/team/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response body has the team type name
	if !strings.Contains(rr.Body.String(), tt1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt1.Name)
	}

	// Check if the response body has the state name
	if !strings.Contains(rr.Body.String(), state.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), state.Name)
	}

	// Check if the response body has the expected content
	if !strings.Contains(rr.Body.String(), "New Team") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "New Team")
	}

	// Cleanup
	// Delete the UserMarker
	var userMarker models.UserMarker
	db.Database.Unscoped().Where("user_id = ?", user.ID).Delete(&userMarker)
	// Delete the user
	db.Database.Unscoped().Delete(&user)
	// Delete the team type
	db.Database.Unscoped().Delete(&tt1)
	// Delete the state
	db.Database.Unscoped().Delete(&state)
	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)

}
