package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// TestEventMeetScore tests the EventMeetScore function
func TestEventMeetScore(t *testing.T) {
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
	// Route Get
	router.GET("/user/:id/event/:event_id/meet_score", func(c *gin.Context) {
		userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		c.Set("user_id", uint(userID))
		controllers.EventMeetScore(c)
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
	event, err := models.EventCreate(name, location, date, against, uint64(userID), uint64(eventtypeID), []uint64{uint64(dive.ID)})
	if err != nil {
		t.Fatalf("Error creating event with 1 dive: %v", err)
	}

	defer func() {
		// Get all scores for the dive
		scores := []models.Score{}
		db.Database.Unscoped().Where("dive_id = ?", dive.ID).Find(&scores)

		// Delete all scores for the dive
		for _, score := range scores {
			db.Database.Unscoped().Delete(&score)
		}
		// Get all scores for the other dive
		scores1 := []models.Score{}
		db.Database.Unscoped().Where("dive_id = ?", dive2.ID).Find(&scores1)

		// Delete all scores for the other dive
		for _, score1 := range scores1 {
			db.Database.Unscoped().Delete(&score1)
		}

		// Delete User Event Dive
		models.DeleteUserEventDivesByEventID(event.ID)

		// Clean up event2
		db.Database.Unscoped().Delete(&event)

		// Delete user
		db.Database.Unscoped().Delete(&user)

		// Delete user type
		db.Database.Unscoped().Delete(&ut1)

		// Delete event type
		db.Database.Unscoped().Delete(&et1)

		// Delete dives
		db.Database.Unscoped().Delete(&dive)
		db.Database.Unscoped().Delete(&dive2)

		// Remove test data
		helpers.CleanTestData(dg1, dg2, dt1, dt2, bt1, bt2, bh1, bh2)

	}()

	type testCase struct {
		numScores       int
		expectErr       bool
		expectScoreD1   float64
		expectScoreD1D2 float64
	}

	testCases := []testCase{
		{1, true, 0, 0},
		{2, true, 0, 0},
		{3, false, float64(8.5 * 3 * dive.Difficulty), float64(8.5 * 3 * (dive.Difficulty + dive2.Difficulty))},
		{4, true, 0, 0},
		{5, false, float64(8.5 * 3 * dive.Difficulty), float64(8.5 * 3 * (dive.Difficulty + dive2.Difficulty))},
		{6, true, 0, 0},
		{7, false, float64(8.5 * 3 * dive.Difficulty), float64(8.5 * 3 * (dive.Difficulty + dive2.Difficulty))},
	}

	for _, tc := range testCases {
		// Create valid scores
		validScores := make([]float64, tc.numScores)
		for i := range validScores {
			validScores[i] = 8.5 // or any valid score
		}

		// Loop through the valid scores and create them for dive1
		for i, scoreValue := range validScores {
			_, err = models.ScoreUpsert(uint64(userID), uint64(event.ID), uint64(dive.ID), i+1, scoreValue)
			if err != nil {
				t.Fatalf("Error creating valid score for dive: %v", err)
			}
		}

		// CalculateMeetScore by making a request to the endpoint
		req, err := http.NewRequest("GET", "/user/"+strconv.FormatUint(uint64(user.ID), 10)+"/event/"+strconv.FormatUint(uint64(event.ID), 10)+"/meet_score", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		// Make the request
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Check the status code
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		// Parse the response
		var meetScoreResponse map[string]float64
		err = json.Unmarshal(resp.Body.Bytes(), &meetScoreResponse)
		if err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		// Check that the total score is what we expect for dive1
		if meetScoreResponse["score"] != tc.expectScoreD1 {
			t.Fatalf("Evaluated score for Dive1. Expected meet score of %v, but got %v with %d scores", tc.expectScoreD1, meetScoreResponse["score"], tc.numScores)
		}

		// Loop through the valid scores and create them for dive2
		for i, scoreValue := range validScores {
			_, err = models.ScoreUpsert(uint64(userID), uint64(event.ID), uint64(dive2.ID), i+1, scoreValue)
			if err != nil {
				t.Fatalf("Error creating valid score for dive2: %v", err)
			}
		}

		// CalculateMeetScore by making a request to the endpoint
		req, err = http.NewRequest("GET", "/user/"+strconv.FormatUint(uint64(user.ID), 10)+"/event/"+strconv.FormatUint(uint64(event.ID), 10)+"/meet_score", nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		// Make the request
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Check the status code
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		// Parse the response
		err = json.Unmarshal(resp.Body.Bytes(), &meetScoreResponse)
		if err != nil {
			t.Fatalf("Error unmarshalling response: %v", err)
		}

		// Check that the total score is what we expect for dive1 and dive2
		if meetScoreResponse["score"] != tc.expectScoreD1 && meetScoreResponse["score"] != tc.expectScoreD1D2 {
			t.Fatalf("Evaluated score for Dive1 and Dive2. Expected meet score of %v or %v, but got %v with %d scores", tc.expectScoreD1, tc.expectScoreD1D2, meetScoreResponse["score"], tc.numScores)
		}

	}

}
