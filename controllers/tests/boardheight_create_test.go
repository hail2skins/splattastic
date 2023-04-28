package controllers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	h "github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
	"github.com/stretchr/testify/assert"
)

// TestBoardHeightCreate is a test function that tests the BoardHeightsNew controller
func TestBoardHeightCreate(t *testing.T) {
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
		admin.POST("/boardheights", controllers.BoardHeightCreate)
	}

	// Create a request with POST method and form data
	form := url.Values{}
	form.Add("height", "9.5")
	req, err := http.NewRequest(http.MethodPost, "/admin/boardheight", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new response recorder
	w := httptest.NewRecorder()

	// Call the BoardHeightCreate function
	controllers.BoardHeightCreate(h.GinContext(req, w))

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the user type was created in the database
	var boardHeight models.BoardHeight
	db.Database.Where("height = ?", 9.5).First(&boardHeight)
	assert.Equal(t, float32(9.5), boardHeight.Height)

	// Delete the board height from the database
	db.Database.Unscoped().Delete(&boardHeight)
}
