package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

// BoardHeightsNew renders the new board height form
func BoardHeightsNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"boardheights/new.html",
		gin.H{
			"title":     "New Board Height",
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "New Board Height",
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}
