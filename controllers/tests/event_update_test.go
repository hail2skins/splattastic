package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

// TestEventUpdate is a test for the EventUpdate controller
func TestEventUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Create a router with the test route
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
	router.POST("/user/:id/event/:event_id", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.EventUpdate(c)
	})
	router.GET("/user/:id/event/:event_id", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		// You may need to implement the corresponding function in the controller
		controllers.EventShow(c)
	})

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two test dives (similar to what you did in the model test)
	dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Create an event type
	et1, _ := models.EventTypeCreate("Test Event Type")

	// Create a user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)
	// Set userID
	userID := user.ID

	helpers.SetupSession(router, uint(userID))
	router.Use(middlewares.AuthenticateUser())

	// Create an event with the specified dives
	eventDate := time.Now()
	event, _ := models.EventCreate("Test Event", "Test Location", eventDate, "Test Against", uint64(user.ID), uint64(et1.ID), []uint64{})

	type testCase struct {
		description   string
		existingDives []uint64
		newDives      []uint64
	}

	testCases := []testCase{
		{
			description:   "Add one dive",
			existingDives: []uint64{},
			newDives:      []uint64{uint64(dive1.ID)},
		},
		{
			description:   "Remove the first dive and add another one",
			existingDives: []uint64{uint64(dive1.ID)},
			newDives:      []uint64{uint64(dive2.ID)},
		},
		{
			description:   "Remove one dive and add two dives",
			existingDives: []uint64{uint64(dive1.ID)},
			newDives:      []uint64{uint64(dive1.ID), uint64(dive2.ID)},
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// Update the event with the existing dives
			event, _ = event.Update("Test Event", "Test Location", eventDate, "Test Against", uint64(et1.ID), test.existingDives)

			// Create a request with the new dives
			data := url.Values{}
			data.Set("name", "Test Event1")
			data.Set("location", "Test Location1")
			data.Set("event_date", "2019-01-03")
			data.Set("against", "Test Against1")
			data.Set("event_type_id", strconv.FormatUint(uint64(et1.ID), 10))

			for _, diveID := range test.newDives {
				data.Add("dive_id", strconv.FormatUint(diveID, 10))
			}
			req, err := http.NewRequest("POST", "/user/"+helpers.UintToString(user.ID)+"/event/"+helpers.UintToString(event.ID), strings.NewReader(data.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Execute the request
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Verify the status code
			if resp.Code != http.StatusFound {
				t.Errorf("Expected status code %d, got %d", http.StatusFound, resp.Code)
			}

			// Verify the event's dives
			currentDives, _ := models.GetDivesForEvent(uint64(event.ID))
			if len(currentDives) != len(test.newDives) {
				t.Errorf("Expected %d dives, got %d", len(test.newDives), len(currentDives))
			}
			for _, dive := range currentDives {
				if !contains(test.newDives, uint64(dive.ID)) {
					t.Errorf("Dive with ID %d should not be associated with the event", dive.ID)
				}
			}

			// Request the redirected page to ensure updated info
			req1, err := http.NewRequest("GET", "/user/"+helpers.UintToString(user.ID)+"/event/"+helpers.UintToString(event.ID), nil)
			if err != nil {
				t.Fatal(err)
			}
			resp1 := httptest.NewRecorder()
			router.ServeHTTP(resp1, req1)

			// Verify event name updated in response body
			if !strings.Contains(resp1.Body.String(), "Test Event1") {
				t.Errorf("Expected response body to contain %s, got %s", "Test Event1", resp1.Body.String())
			}
			// Verify event location updated in response body
			if !strings.Contains(resp1.Body.String(), "Test Location1") {
				t.Errorf("Expected response body to contain %s, got %s", "Test Location1", resp1.Body.String())
			}
			// Verify event against updated in response body
			if !strings.Contains(resp1.Body.String(), "Test Against1") {
				t.Errorf("Expected response body to contain %s, got %s", "Test Against1", resp1.Body.String())
			}
			// not testing date or event type here but probably should

			// Delete UserEventDives created during event.Update
			if err := models.DeleteUserEventDivesByEventID(event.ID); err != nil {
				t.Errorf("Error deleting user event dives: %v", err)
			}

			// Delete the event
			db.Database.Unscoped().Delete(&event)
		})
	}

	// Delete user event dives
	helpers.DeleteUserEventDives()

	// Delete event
	db.Database.Unscoped().Delete(&event)

	// Delete event type
	db.Database.Unscoped().Delete(&et1)

	// Delete dives
	db.Database.Unscoped().Delete(&dive1)
	db.Database.Unscoped().Delete(&dive2)

	// Delete user
	db.Database.Unscoped().Delete(&user)

	// Delete user type
	db.Database.Unscoped().Delete(&ut1)

	// Remove test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}

// Helper function to check if a uint64 slice contains a value
func contains(slice []uint64, value uint64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
