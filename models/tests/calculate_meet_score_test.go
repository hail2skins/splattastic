package models

import (
	"testing"
	"time"

	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestCalculateMeetScore tests the CalculateMeetScore function
func TestCalculateMeetScore(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

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
			_, err = models.ScoreUpsert(uint64(userID), uint64(event1.ID), uint64(dive.ID), i+1, scoreValue)
			if err != nil {
				t.Fatalf("Error creating valid score for dive: %v", err)
			}
		}

		// Test CalculateMeetScore
		meetScore, err := models.CalculateMeetScore(uint64(event1.ID))
		if err != nil && err != models.ErrInvalidScores {
			t.Fatalf("Unexpected error calculating meet score with %d scores: %v", tc.numScores, err)
		}

		// Check that the total score is what we expect for dive1
		if err == nil && meetScore != tc.expectScoreD1 {
			t.Fatalf("Expected meet score of %v, but got %v with %d scores", tc.expectScoreD1, meetScore, tc.numScores)
		}

		// Loop through the valid scores and create them for dive2
		for i, scoreValue := range validScores {
			_, err = models.ScoreUpsert(uint64(userID), uint64(event1.ID), uint64(dive2.ID), i+1, scoreValue)
			if err != nil {
				t.Fatalf("Error creating valid score for dive2: %v", err)
			}
		}

		// Test CalculateMeetScore again after scoring dive2
		meetScore, err = models.CalculateMeetScore(uint64(event1.ID))
		if err != nil && err != models.ErrInvalidScores {
			t.Fatalf("Unexpected error calculating meet score with %d scores: %v", tc.numScores, err)
		}

		// Check that the total score is what we expect for dive1 and dive2
		if err == nil {
			if meetScore != tc.expectScoreD1 && meetScore != tc.expectScoreD1D2 {
				t.Fatalf("Expected meet score of %v or %v, but got %v with %d scores", tc.expectScoreD1, tc.expectScoreD1D2, meetScore, tc.numScores)
			}
		}
	}

	// Get all scores for the dives
	scores := []models.Score{}
	db.Database.Unscoped().Where("dive_id = ?", dive.ID).Find(&scores)
	scores1 := []models.Score{}
	db.Database.Unscoped().Where("dive_id = ?", dive2.ID).Find(&scores1)

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
