package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

// TestBoardTypesIndex is a test for the BoardTypesIndex controller
func TestBoardTypesIndex(t *testing.T) {
	LoadEnv()
	db.Connect()

	// Set TestRun to true and defer to false when done
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "")

	// Create two board types
	boardType1, _ := models.BoardTypeCreate("TestBoardType1")
	boardType2, _ := models.BoardTypeCreate("TestBoardType2")

	// Create a get request to /admin/boardtypes
	req, err := http.NewRequest("GET", "/admin/boardtypes", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/boardtypes", controllers.BoardTypesIndex)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the board types are in the response
	if !strings.Contains(rr.Body.String(), boardType1.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), boardType1.Name)
	}
	if !strings.Contains(rr.Body.String(), boardType2.Name) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), boardType2.Name)
	}

	// Check if index page has the correct page text
	if !strings.Contains(rr.Body.String(), "Board Types") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Board Types")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&boardType1)
	db.Database.Unscoped().Delete(&boardType2)

}
