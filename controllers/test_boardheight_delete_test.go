package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightDelete is a test function for the BoardHeightDelete controller
func TestBoardHeightDelete(t *testing.T) {
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
		admin.DELETE("/boardheights/:id", BoardHeightDelete)
	}

	// Create a test board height
	testBoardHeight := models.BoardHeight{Height: 11.0}
	db.Database.Create(&testBoardHeight)

	// Create a request to send to the above route
	req, err := http.NewRequest(http.MethodDelete, "/admin/boardheights/"+helpers.UintToString(testBoardHeight.ID), nil)
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}

	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	assert.Equal(t, http.StatusFound, w.Code)

	// Check to see if the board height was deleted
	var softDeletedBoardHeight models.BoardHeight
	db.Database.Unscoped().Where("id = ?", testBoardHeight.ID).First(&softDeletedBoardHeight)
	assert.NotNil(t, softDeletedBoardHeight.DeletedAt)

	// Clean up
	db.Database.Unscoped().Delete(&softDeletedBoardHeight)

}
