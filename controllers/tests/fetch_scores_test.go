package controllers

import (
	"encoding/json"
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
	"github.com/hail2skins/splattastic/models"
)

// TestFetchScores tests the FetchScores controller
func TestFetchScores(t *testing.T) {
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
	router.GET("/user/:id/event/:event_id/dive/:dive_id/scores", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.FetchScores(c)
	})

	// Insert test data and defer cleanup
	dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2 := helpers.CreateTestData()

	// Create two dives
	// Use the test data to create a new dive
	name := "Test Dive"
	number := 154
	difficulty := float32(2.5)
	divetypeID := dt1.ID
	divegroupID := dg1.ID
	boardtypeID := bt1.ID
	boardheightID := bh1.ID

	// Use this test data to create a second dive
	name2 := "Test Dive 2"
	number2 := 155
	difficulty2 := float32(3.5)
	divetypeID2 := dt2.ID
	divegroupID2 := dg2.ID
	boardtypeID2 := bt2.ID
	boardheightID2 := bh2.ID

	// Create the first dive
	dive, err := models.DiveCreate(name, number, difficulty, uint64(divetypeID), uint64(divegroupID), uint64(boardtypeID), uint64(boardheightID))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}
	// Create the second dive
	dive2, err := models.DiveCreate(name2, number2, difficulty2, uint64(divetypeID2), uint64(divegroupID2), uint64(boardtypeID2), uint64(boardheightID2))
	if err != nil {
		t.Fatalf("Error creating dive: %v", err)
	}

	// Create a UserType
	ut1, err := models.CreateUserType("Test User Type")
	if err != nil {
		t.Fatalf("Error creating user type: %v", err)
	}

	// Create a User
	user, _ := models.UserCreate("test@example.com", "testpassword", "Test", "User", "testuser", ut1.Name)

	// Create EventType
	et1, err := models.EventTypeCreate("Test Event Type")
	if err != nil {
		t.Fatalf("Error creating event type: %v", err)
	}

	// Use the test data to create a new event
	name = "Test Event"
	location := "Test Location"
	date := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	against := "Test Against"
	userID := user.ID
	eventtypeID := et1.ID

	// Test with 1 dive
	event1, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{uint64(dive.ID)})
	if err != nil {
		t.Fatalf("Error creating event with 1 dive: %v", err)
	}

	// Add 5 valid scores for a dive
	validScores := make([]float64, 5)
	for i := range validScores {
		validScores[i] = 8.5 // or any valid score
	}

	// Create scores for the dive
	for i, scoreValue := range validScores {
		_, err = models.ScoreCreate(uint64(userID), uint64(event1.ID), uint64(dive.ID), i+1, scoreValue)
		if err != nil {
			t.Fatalf("Error creating valid score: %v", err)
		}
	}

	// Perform an HTTP GET request to the "/user/:id/event/:event_id/dive/:dive_id/scores" route
	req, err := http.NewRequest("GET", "/user/"+strconv.FormatUint(uint64(user.ID), 10)+"/event/"+strconv.FormatUint(uint64(event1.ID), 10)+"/dive/"+strconv.FormatUint(uint64(dive.ID), 10)+"/scores", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v", http.StatusOK, w.Code)
	}

	// Parse the JSON response
	var response map[string][]models.Score
	errJson := json.Unmarshal(w.Body.Bytes(), &response)
	if errJson != nil {
		t.Fatalf("Error parsing JSON: %v", errJson)
	}

	// Check if the response contains the created scores
	if len(response["scores"]) != len(validScores) {
		t.Fatalf("Expected %v scores but got %v", len(validScores), len(response["scores"]))
	}
	for i, score := range response["scores"] {
		if score.Value != validScores[i] {
			t.Fatalf("Expected score value %v but got %v", validScores[i], score.Value)
		}
	}

	// Get all scores for the dive
	scores := []models.Score{}
	db.Database.Unscoped().Where("dive_id = ?", dive.ID).Find(&scores)

	// Delete all scores for the dive
	for _, score := range scores {
		db.Database.Unscoped().Delete(&score)
	}

	// Delete User Event Dive
	models.DeleteUserEventDivesByEventID(event1.ID)

	// Clean up event2
	db.Database.Unscoped().Delete(&event1)

	// Delete dives
	db.Database.Unscoped().Delete(&dive)
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
