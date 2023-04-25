package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// BoardHeightsNew renders the new board height form
func BoardHeightNew(c *gin.Context) {
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

// BoardHeightCreate creates a new board height
func BoardHeightCreate(c *gin.Context) {
	heightStr := c.PostForm("height")
	if heightStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Convert the height to a float32
	height, err := strconv.ParseFloat(heightStr, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, err = models.CreateBoardHeight(float32(height))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/admin")
}

// BoardHeightsIndex renders the board heights index page
func BoardHeightsIndex(c *gin.Context) {
	boardheights, err := models.GetBoardHeights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardheights/index.html",
		gin.H{
			"title":        "Board Heights",
			"logged_in":    h.IsUserLoggedIn(c),
			"header":       "Board Heights",
			"boardheights": boardheights,
			"test_run":     os.Getenv("TEST_RUN") == "true",
		},
	)
}
