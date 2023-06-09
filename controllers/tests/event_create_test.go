package controllers

import (
	"fmt"
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

// TestEventCreate tests the EventCreate function
func TestEventCreate(t *testing.T) {
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
	router.POST("/user/:id/event", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.EventCreate(c)
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

	// use form to create event
	data := url.Values{}
	data.Set("name", "Test Event")
	data.Set("location", "Test Location")
	data.Set("event_date", "2019-01-02")
	data.Set("against", "Test Against")
	data.Set("user_id", fmt.Sprintf("%d", user.ID))
	data.Set("event_type_id", fmt.Sprintf("%d", et1.ID))
	testCases := []struct {
		name  string
		dives []uint64
	}{
		{
			name:  "zero dives",
			dives: []uint64{},
		},
		{
			name:  "one dive",
			dives: []uint64{uint64(dive1.ID)},
		},
		{
			name:  "two dives",
			dives: []uint64{uint64(dive1.ID), uint64(dive2.ID)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the dive_ids in the data only if there are dives in the test case
			data.Del("dive_id") // Remove any previously set "dive_id" fields
			for _, diveID := range tc.dives {
				data.Add("dive_id", strconv.FormatUint(diveID, 10))
			}

			// Create a request to send to the above route
			req, err := http.NewRequest("POST", "/user/"+helpers.UintToString(user.ID)+"/event", strings.NewReader(data.Encode()))
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check to see if the response was what you expected
			if w.Code != http.StatusFound {
				t.Errorf("Expected status code %v but got %v", http.StatusFound, w.Code)
			}

			// Get the created event
			var createdEvent models.Event
			result := db.Database.Where("name = ?", "Test Event").First(&createdEvent)

			if result.Error != nil {
				t.Errorf("Error retrieving the created event: %v", result.Error)
			}

			// Check to see if the response redirected to the event show page
			if w.Header().Get("Location") != "/user/"+helpers.UintToString(user.ID)+"/event/"+helpers.UintToString(createdEvent.ID) {
				t.Errorf("Expected redirect to /user/%d/event/1 but got %s", user.ID, w.Header().Get("Location"))
			}

			// Check if the number of created UserEventDives is equal to the number of dives specified in the test case
			var userEventDives []models.UserEventDive
			db.Database.Where("event_id = ?", createdEvent.ID).Find(&userEventDives)
			if len(userEventDives) != len(tc.dives) {
				t.Errorf("Expected %d UserEventDives, but got %d", len(tc.dives), len(userEventDives))
			}

			// Clean up the UserEventDives associated with the event
			for _, userEventDive := range userEventDives {
				db.Database.Unscoped().Delete(&userEventDive)
			}

			// Clean up the created event
			db.Database.Unscoped().Delete(&createdEvent)
		})
	}

	// Delete user event dives
	helpers.DeleteUserEventDives()

	// Delete dives
	db.Database.Unscoped().Delete(&dive1)
	db.Database.Unscoped().Delete(&dive2)

	// Delete user
	db.Database.Unscoped().Delete(&user)

	// Delete user type
	db.Database.Unscoped().Delete(&ut1)

	// Delete event type
	db.Database.Unscoped().Delete(&et1)

	// Remove test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
