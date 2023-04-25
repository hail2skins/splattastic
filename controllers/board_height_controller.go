package controllers

import (
	"fmt"
	"log"
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

// BoardHeightShow renders the board height show page
func BoardHeightShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing id: %v\n", err)
	}

	boardheight, err := models.BoardHeightShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardheights/show.html",
		gin.H{
			"title":       "Board Height",
			"logged_in":   h.IsUserLoggedIn(c),
			"header":      "Board Height",
			"boardheight": boardheight,
			"test_run":    os.Getenv("TEST_RUN") == "true",
		},
	)
}

// BoardHeightEdit renders the board height edit page
func BoardHeightEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing id: %v\n", err)
	}

	boardheight, err := models.BoardHeightShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardheights/edit.html",
		gin.H{
			"title":       "Edit Board Height",
			"logged_in":   h.IsUserLoggedIn(c),
			"header":      "Edit Board Height",
			"boardheight": boardheight,
			"test_run":    os.Getenv("TEST_RUN") == "true",
		},
	)

}

// BoardHeightUpdate updates a board height
func BoardHeightUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	boardheight, err := models.BoardHeightShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	heightStr := c.PostForm("height")
	if heightStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Height cannot be blank"})
		return
	}

	// Convert the height to a float32
	height, err := strconv.ParseFloat(heightStr, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = boardheight.Update(float32(height))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/boardheights")
}
