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
	"github.com/hail2skins/splattastic/models"
)

// TestUserShow tests the UserShow function
func TestUserShow(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a User Type
	ut1, err := models.CreateUserType("Test User Type")
	if err != nil {
		t.Fatal(err)
	}
	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

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
	router.GET("/user/:id", controllers.UserShow)

	t.Run("User with no events", func(t *testing.T) {
		// Create a new request
		req, err := http.NewRequest("GET", "/user/"+helpers.UintToString(user.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder so we can inspect the response
		rr := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// check username
		if !strings.Contains(rr.Body.String(), user.UserName) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.UserName)
		}
		// check firstname
		if !strings.Contains(rr.Body.String(), user.FirstName) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.FirstName)
		}
		// check lastname
		if !strings.Contains(rr.Body.String(), user.LastName) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), user.LastName)
		}
		// check user type
		if !strings.Contains(rr.Body.String(), ut1.Name) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), ut1.Name)
		}
		// check that there are no recent events
		if !strings.Contains(rr.Body.String(), "No Events Yet") {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "No Events Yet")
		}
		// check that there is a link to create a new event
		if !strings.Contains(rr.Body.String(), "/user/"+helpers.UintToString(user.ID)+"/event/new") {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "/user/"+helpers.UintToString(user.ID)+"/event/new")
		}
		// check that there is a link to events for the user
		if !strings.Contains(rr.Body.String(), "/user/"+helpers.UintToString(user.ID)+"/events") {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "/user/"+helpers.UintToString(user.ID)+"/events")
		}
		// ADD MORE HERE ONCE PAGE IS MORE READY
	})

	t.Run("User with 5 events", func(t *testing.T) {
		// create event type
		et1, _ := models.EventTypeCreate("Test Event Type")
		// Create 5 events
		for i := 0; i < 5; i++ {
			eventTime := time.Now().AddDate(0, 0, i)
			_, err := models.EventCreate("Event"+strconv.Itoa(i+1), "Test Location"+strconv.Itoa(i+1), eventTime, "Test Against"+strconv.Itoa(i+1), uint64(user.ID), uint64(et1.ID), []uint64{})
			if err != nil {
				t.Fatal(err)
			}
		}
		// Create a new request
		req, err := http.NewRequest("GET", "/user/"+helpers.UintToString(user.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder so we can inspect the response
		rr := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(rr, req)

		// check that there are 5 recent events
		for i := 0; i < 5; i++ {
			eventName := "Event" + strconv.Itoa(i+1)
			if !strings.Contains(rr.Body.String(), eventName) {
				t.Errorf("Event %d is missing: got %v want %v", i+1, rr.Body.String(), eventName)
			}
		}

		// Cleanup
		var events []models.Event
		db.Database.Where("user_id = ?", user.ID).Find(&events)
		for _, event := range events {
			db.Database.Unscoped().Delete(&event)
		}
		db.Database.Unscoped().Delete(&et1)
	})

	// Cleanup
	db.Database.Unscoped().Delete(&user)
	db.Database.Unscoped().Delete(&ut1)

}
