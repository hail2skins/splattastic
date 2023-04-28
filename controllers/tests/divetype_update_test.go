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

// TestDiveTypeUpdate is a test for the DiveTypeUpdate controller
func TestDiveTypeUpdate(t *testing.T) {
	// Setup
	LoadEnv()
	db.Connect()

	// Sets the TEST_RUN env var to true for views requiring logged in user but tests that don't require a logged in user
	os.Setenv("TEST_RUN", "true")
	defer os.Setenv("TEST_RUN", "")

	// Create a new dive type
	diveType, err := models.DiveTypeCreate("TestDiveType", "Q")
	if err != nil {
		t.Fatal(err)
	}

	// Create a post request with form data
	data := url.Values{}
	data.Set("name", "UpdatedTestDiveType")
	data.Set("letter", "R")
	req, err := http.NewRequest("POST", "/admin/divetypes/"+helpers.UintToString(diveType.ID), strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a gin router with the routes we need
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/**/**")
	r.GET("/admin/divetypes", controllers.DiveTypesIndex)
	r.POST("/admin/divetypes/:id", controllers.DiveTypeUpdate)

	// Create a response recorder so we can inspect the response
	rr := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Get the updated dive type
	updatedDiveType, err := models.DiveTypeShow(uint64(diveType.ID))
	if err != nil {
		t.Fatal(err)
	}

	// Check the name of the updated dive type
	expected := "UpdatedTestDiveType"
	if updatedDiveType.Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDiveType.Name, expected)
	}
	if updatedDiveType.Letter != "R" {
		t.Errorf("handler returned unexpected body: got %v want %v", updatedDiveType.Letter, "R")
	}

	// Cleanup
	db.Database.Unscoped().Delete(&updatedDiveType)
}
