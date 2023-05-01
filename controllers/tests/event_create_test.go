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

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
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
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	h.InitRouterWithFuncMap(router)
	router.LoadHTMLGlob("../../templates/**/**")
	router.POST("/user/:id/event", controllers.EventCreate)

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

	// use form to create event
	data := url.Values{}
	data.Set("name", "Test Event")
	data.Set("location", "Test Location")
	data.Set("date", "01/01/2019")
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
			// Set the dive_ids in the data
			divesStr := make([]string, len(tc.dives))
			for i, diveID := range tc.dives {
				divesStr[i] = strconv.FormatUint(diveID, 10)
			}
			data.Set("dive_ids", strings.Join(divesStr, ","))

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
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %v but got %v", http.StatusOK, w.Code)
			}
			// Get the created event
			var createdEvent models.Event
			db.Database.Where("name = ?", "Test Event").First(&createdEvent)

			// Clean up the created event
			db.Database.Unscoped().Delete(&createdEvent)

			// Clean up the UserEventDives associated with the event
			var userEventDives []models.UserEventDive
			db.Database.Where("event_id = ?", createdEvent.ID).Find(&userEventDives)
			for _, userEventDive := range userEventDives {
				db.Database.Unscoped().Delete(&userEventDive)
			}
		})
	}

	// Delete user event dives
	deleteUserEventDives()

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

func deleteUserEventDives() error {
	var userEventDives []models.UserEventDive
	result := db.Database.Find(&userEventDives)
	if result.Error != nil {
		return result.Error
	}

	//fmt.Println("UserEventDives before deletion:", userEventDives)

	for _, userEventDive := range userEventDives {
		db.Database.Unscoped().Delete(&userEventDive)
	}

	// Check if the user event dives are deleted
	var userEventDivesAfterDeletion []models.UserEventDive
	result = db.Database.Find(&userEventDivesAfterDeletion)
	if result.Error != nil {
		return result.Error
	}
	//fmt.Println("UserEventDives after deletion:", userEventDivesAfterDeletion)

	return nil
}
