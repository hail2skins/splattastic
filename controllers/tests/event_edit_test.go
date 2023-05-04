package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/models"
)

// TestEventEdit is a test function for the EventEdit controller
func TestEventEdit(t *testing.T) {
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
	router.GET("/user/:id/event/:event_id/edit", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		eventID, _ := strconv.ParseUint(c.Param("event_id"), 10, 64)
		c.Set("user_id", uint(userID))
		c.Set("event_id", uint(eventID))
		controllers.EventEdit(c)
	})

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two test dives (similar to what you did in the model test)
	dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Create a test event type
	et1, _ := models.EventTypeCreate("Test Event Type")

	// Create a test user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

	// Create a test event
	eventDate := time.Now()
	event, _ := models.EventCreate("Test Event", "Test Location", eventDate, "Test Against", uint64(user.ID), uint64(et1.ID), []uint64{})

	// Set userID
	userID := user.ID

	// Send the user to the session and through middleware
	helpers.SetupSession(router, uint(userID))
	router.Use(middlewares.AuthenticateUser())

	// Create a request to the test route
	req, err := http.NewRequest("GET", "/user/"+helpers.UintToString(userID)+"/event/"+helpers.UintToString(event.ID)+"/edit", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	resp := httptest.NewRecorder()

	// Serve the request to the recorder
	router.ServeHTTP(resp, req)

	// Check the status code
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if !strings.Contains(resp.Body.String(), "Test Event") {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), "Test Event")
	}
	// Check the response body for event name
	if !strings.Contains(resp.Body.String(), event.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), event.Name)
	}
	// Check the response body for event type
	if !strings.Contains(resp.Body.String(), et1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Body.String(), et1.Name)
	}

	// Delete the event
	db.Database.Unscoped().Delete(&event)

	// Delete the user
	db.Database.Unscoped().Delete(&user)

	// Delete the event type
	db.Database.Unscoped().Delete(&et1)

	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)

	// Delete the dives
	db.Database.Unscoped().Delete(&dive1)
	db.Database.Unscoped().Delete(&dive2)

	// Clean up the test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
