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
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func TestBoardTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "") // Reset the TEST_RUN env var

	// Create a new board type
	boardType, err := models.BoardTypeCreate("TestBoardType")
	if err != nil {
		t.Fatal(err)
	}

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "TestBoardTypeUpdated")
	req, err := http.NewRequest("POST", "/admin/boardtypes/"+helpers.UintToString(boardType.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/boardtypes", controllers.BoardTypesIndex)
	r.POST("/admin/boardtypes/:id", controllers.BoardTypeUpdate)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the updated board type
	updatedBoardType, err := models.BoardTypeShow(uint64(boardType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check if the name was updated
	expected := "TestBoardTypeUpdated"
	if updatedBoardType.Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedBoardType.Name, expected)
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedBoardType)
}
