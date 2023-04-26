package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

// DiveTypeNew renders the new dive type form
func DiveTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"divetypes/new.html",
		gin.H{
			"title":     "New Dive Type",
			"header":    "New Dive Type",
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}
