package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/controllers/helpers"
	db "github.com/hail2skins/splattastic/database"
	"github.com/hail2skins/splattastic/models"
)

func TestBoardTypeDelete(t *testing.T) {
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

	// Create a delete request
	req, err := http.NewRequest(http.MethodDelete, "/admin/boardtypes/"+helpers.UintToString(boardType.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.DELETE("/admin/boardtypes/:id", controllers.BoardTypeDelete)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Verify the board type was deleted
	_, err = models.BoardTypeShow(uint64(boardType.ID))
	if err == nil {
		t.Errorf("Board type with ID %d not deleted", boardType.ID)
	}

	// Teardown
	db.Database.Unscoped().Delete(&boardType)
}
