package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

// TestTeamCreate tests the TeamCreate function
func TestTeamCreate(t *testing.T) {
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
	router.POST("/user/:id/team", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.TeamCreate(c)
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

	// use form to create team
	data := url.Values{}
	data.Set("team-name", "Test Team")
	data.Set("street", "123 Test St")
	data.Set("street1", "Suite 1")
	data.Set("city", "Test City")
	data.Set("zip", "12345")
	data.Set("abbreviation", "TT")
	data.Set("team_type_id", strconv.Itoa(int(tt1.ID)))
	data.Set("state-id", strconv.Itoa(int(state.ID)))
	data.Set("associate_user", "on")

	// Create a request to send to the above route
	req, err := http.NewRequest("POST", "/user/"+helpers.UintToString(userID)+"/team", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request to the recorder
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the created team
	var team models.Team
	result := db.Database.Where("name = ?", "Test Team").First(&team)
	if result.Error != nil {
		t.Errorf("handler returned unexpected error: got %v", result.Error)
	}

	// Clean up UserTeam association
	var userTeam models.UserTeam
	db.Database.Unscoped().Where("user_id = ?", userID).Delete(&userTeam)

	// Clean up UserMarkers in case they were created
	var userMarker models.UserMarker
	db.Database.Unscoped().Where("user_id = ?", userID).Delete(&userMarker)

	// Clean up Team
	db.Database.Unscoped().Delete(&team)

	// Clean up Team Type
	db.Database.Unscoped().Delete(&tt1)

	// Clean up State
	db.Database.Unscoped().Delete(&state)

	// Clean up User
	db.Database.Unscoped().Delete(&user)

	// Clean up User Type
	db.Database.Unscoped().Delete(&ut1)

}
