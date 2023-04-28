package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightShow is a test function for the BoardHeightShow controller
func TestBoardHeightShow(t *testing.T) {
	// Setup code
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Set up the test server
	r := gin.Default()

	r.LoadHTMLGlob("../../templates/**/**")

	admin := r.Group("/admin")
	{
		admin.GET("/", controllers.AdminDashboard)

		// User types
		admin.GET("/boardheights/:id", controllers.BoardHeightShow)

	}

	// Create a test board height
	testBoardHeight := models.BoardHeight{Height: 1.5}
	db.Database.Create(&testBoardHeight)

	assert.NotNil(t, testBoardHeight)

	// request the board height show page
	req, err := http.NewRequest(http.MethodGet, "/admin/boardheights/"+helpers.UintToString(testBoardHeight.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	expectedText := fmt.Sprintf("Board Height: %.1fM", 1.5)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), expectedText)

	// Delete the test board height
	db.Database.Unscoped().Delete(&testBoardHeight)

}
