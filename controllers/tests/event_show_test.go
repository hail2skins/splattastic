package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

func TestEventShow(t *testing.T) {
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

	// Create a user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)
	// Set userID
	userID := user.ID

	helpers.SetupSession(router, uint(userID))
	router.Use(middlewares.AuthenticateUser())

	testCases := []struct {
		name          string
		diveIDs       []uint64
		expectedDives int
	}{
		{
			name:          "0 dives",
			diveIDs:       []uint64{},
			expectedDives: 0,
		},
		{
			name:          "1 dive",
			diveIDs:       []uint64{uint64(dive1.ID)},
			expectedDives: 1,
		},
		{
			name:          "multiple dives",
			diveIDs:       []uint64{uint64(dive1.ID), uint64(dive2.ID)},
			expectedDives: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an event type
			et1, _ := models.EventTypeCreate("Test Event Type")
			// Create an event with the specified dives
			eventDate := time.Now()
			event, _ := models.EventCreate("Test Event", "Test Location", eventDate, "Test Against", uint64(user.ID), uint64(et1.ID), tc.diveIDs)

			// Create a request to the test route
			req, err := http.NewRequest("GET", "/user/"+helpers.UintToString(user.ID)+"/event/"+helpers.UintToString(event.ID), nil)
			if err != nil {
				t.Errorf("Error creating request: %v", err)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request to the recorder
			router.ServeHTTP(w, req)

			// Check to see if the response was what you expected
			if w.Code != http.StatusOK {
				t.Errorf("Expected response code %v, got %v", http.StatusOK, w.Code)
			}

			// Check if the number of dives is as expected
			var eventData map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &eventData)
			if err != nil {
				t.Errorf("Error unmarshalling response JSON: %v", err)
			}
			dives := eventData["dives"].([]interface{})
			if len(dives) != tc.expectedDives {
				t.Errorf("Expected %d dives, got %d", tc.expectedDives, len(dives))
			}

			// Cleanup User Event Dives
			helpers.DeleteUserEventDives()

			// Cleanup
			db.Database.Unscoped().Delete(&event)

			// Delete event type
			db.Database.Unscoped().Delete(&et1)
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

	// Remove test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)
}
