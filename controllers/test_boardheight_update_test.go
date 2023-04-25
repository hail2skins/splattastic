package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightUpdate is a controller for testing the update of a board height
func TestBoardHeightUpdate(t *testing.T) {
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

	r.LoadHTMLGlob("../templates/**/**")

	admin := r.Group("/admin")
	{
		admin.GET("/", AdminDashboard)

		// User types
		admin.GET("/boardheights", BoardHeightsIndex)
		admin.GET("/boardheights/edit/:id", BoardHeightEdit)
		admin.POST("/boardheights/:id", BoardHeightUpdate)
	}

	// Create a test board height
	testBoardHeight := models.BoardHeight{Height: 20}
	db.Database.Create(&testBoardHeight)

	assert.NotNil(t, testBoardHeight)

	// Test the controller
	updatedHeight := float32(30.5)
	form := url.Values{}
	form.Add("height", fmt.Sprintf("%f", updatedHeight))
	req, err := http.NewRequest(http.MethodPost, "/admin/boardheights/"+helpers.UintToString(testBoardHeight.ID), strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the board height was updated
	assert.Equal(t, http.StatusFound, w.Code)

	// Check to see if the board height was updated
	var updatedBoardHeight models.BoardHeight
	db.Database.First(&updatedBoardHeight, testBoardHeight.ID)

	// Check to see if the board height was updated
	assert.Equal(t, updatedHeight, updatedBoardHeight.Height)

	// Teardown code
	db.Database.Unscoped().Delete(&testBoardHeight)

}
