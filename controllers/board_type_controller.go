package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// BoardTypeNew renders the new board type form
func BoardTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"boardtypes/new.html",
		gin.H{
			"title":     "New Board Type",
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "New Board Type",
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// BoardTypeCreate creates a new board type
func BoardTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.BoardTypeCreate(name)
	c.Redirect(http.StatusFound, "/admin/boardtypes")
}

// BoardTypesIndex renders the board types index page
func BoardTypesIndex(c *gin.Context) {
	boardtypes, err := models.BoardTypesGet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardtypes/index.html",
		gin.H{
			"title":      "Board Types",
			"logged_in":  h.IsUserLoggedIn(c),
			"header":     "Board Types",
			"boardtypes": boardtypes,
			"test_run":   os.Getenv("TEST_RUN") == "true",
			"user_id":    c.GetUint("user_id"),
		},
	)
}

// BoardTypeShow renders the board type show page
func BoardTypeShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing board type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	boardtype, err := models.BoardTypeShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardtypes/show.html",
		gin.H{
			"title":     "Board Type",
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "Board Type",
			"boardtype": boardtype,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// BoardTypeEdit renders the board type edit page
func BoardTypeEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing board type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	boardtype, err := models.BoardTypeShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"boardtypes/edit.html",
		gin.H{
			"title":     "Edit Board Type",
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "Edit Board Type",
			"boardtype": boardtype,
			"test_run":  os.Getenv("TEST_RUN") == "true",
			"user_id":   c.GetUint("user_id"),
		},
	)
}

// BoardTypeUpdate updates a board type
func BoardTypeUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing board type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	boardtype, err := models.BoardTypeShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = boardtype.Update(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/boardtypes")
}

// BoardTypeDelete deletes a board type
func BoardTypeDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing board type id: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = models.BoardTypeDelete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/boardtypes")
}
