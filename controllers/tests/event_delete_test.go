package controllers

import (
	"net/http"
	"net/http/httptest"
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

// TestEventDelete tests the event delete controller
func TestEventDelete(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

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
	router.DELETE("/user/:id/event/:event_id", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		// You may need to implement the corresponding function in the controller
		controllers.EventDelete(c)
	})

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two test dives (similar to what you did in the model test)
	dive1, _ := models.DiveCreate("Test Dive 1", 154, 2.5, uint64(dt1.ID), uint64(dg1.ID), uint64(bt1.ID), uint64(bh1.ID))
	dive2, _ := models.DiveCreate("Test Dive 2", 155, 3.5, uint64(dt2.ID), uint64(dg2.ID), uint64(bt2.ID), uint64(bh2.ID))

	// Create an event type
	et1, _ := models.EventTypeCreate("Test Event Type")

	// create a test user type
	ut1, _ := models.CreateUserType("Test User Type")

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)
	// Set userID
	userID := user.ID

	// Create an event with the specified dives
	eventDate := time.Now()
	event, _ := models.EventCreate("Test Event", "Test Location", eventDate, "Test Against", uint64(user.ID), uint64(et1.ID), []uint64{uint64(dive1.ID), uint64(dive2.ID)})

	helpers.SetupSession(router, uint(userID))
	router.Use(middlewares.AuthenticateUser())

	// Perform an HTTP DELETE request on the test route
	req, err := http.NewRequest("DELETE", "/user/"+helpers.UintToString(userID)+"/event/"+helpers.UintToString(event.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Verify that the event was deleted
	_, err = models.EventShow(uint64(event.ID))
	if err == nil {
		t.Errorf("event was not deleted")
	}

	// Verify the User Event Dive is soft deleted
	userEventDives, _ := models.GetDivesForEvent(uint64(event.ID))

	for _, userEventDive := range userEventDives {
		if userEventDive.DeletedAt.Valid {
			t.Logf("DeletedAt: %v", userEventDive.DeletedAt.Time)
		} else {
			t.Fatalf("User Event Dive not deleted: %v", userEventDive)
		}
	}

	// Delete the User Event Dives
	_ = models.DeleteUserEventDivesByEventID(event.ID)

	// Delete the event
	db.Database.Unscoped().Delete(&event)
	// Delete the event type
	db.Database.Unscoped().Delete(&et1)
	// Delete the dives
	db.Database.Unscoped().Delete(&dive1)
	db.Database.Unscoped().Delete(&dive2)
	// Delete the user
	db.Database.Unscoped().Delete(&user)
	// Delete the user type
	db.Database.Unscoped().Delete(&ut1)
	// Delete the test data
	helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

}
