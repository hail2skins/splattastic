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

// User functions not within the sessions controller

// UserShow displays a user with associated UserType
func UserShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user
	user, err := models.UserShow(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the user's UserType
	userType, err := models.UserTypeShow(user.UserTypeID)
	if err != nil {
		log.Printf("Error getting user type: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the template
	c.HTML(
		http.StatusOK,
		"users/show.html",
		gin.H{
			"title":        "User Profile",
			"user":         user,
			"usertype":     userType,
			"logged_in":    h.IsUserLoggedIn(c),
			"test_run":     os.Getenv("TEST_RUN") == "true",
			"header":       "User Profile",
			"current_user": h.IsCurrentUser(c, uint64(user.ID)),
			"user_id":      c.GetUint("user_id"),
		},
	)
}

// UserEdit renders the user edit form
func UserEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the list of user types from the database
	userTypes, err := models.GetUserTypes()

	// Get the user
	user, err := models.UserShow(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the user's UserType
	userType, err := models.UserTypeShow(user.UserTypeID)
	if err != nil {
		log.Printf("Error getting user type: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the template
	c.HTML(
		http.StatusOK,
		"users/edit.html",
		gin.H{
			"title":        "Edit Profile",
			"user":         user,
			"usertype":     userType,
			"logged_in":    h.IsUserLoggedIn(c),
			"test_run":     os.Getenv("TEST_RUN") == "true",
			"header":       "Edit Profile",
			"current_user": h.IsCurrentUser(c, uint64(user.ID)),
			"user_id":      c.GetUint("user_id"),
			"userTypes":    userTypes,
		},
	)
}
