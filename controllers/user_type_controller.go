package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

// UserTypeNew function to render the new user type page with title and logged_in
func UserTypeNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"usertypes/new.html",
		gin.H{
			"title":     "New User Type",
			"logged_in": h.IsUserLoggedIn(c),
		},
	)
}

// UserTypeCreate function to create a new user type
func UserTypeCreate(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	models.CreateUserType(name)
	c.Redirect(http.StatusMovedPermanently, "/")
}

// UserTypeIndex function to render the user type index page with title and logged_in
func UserTypeIndex(c *gin.Context) {
	usertypes, err := models.GetUserTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"usertypes/index.html",
		gin.H{
			"title":     "User Types",
			"logged_in": h.IsUserLoggedIn(c),
			"usertypes": usertypes,
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}

// UserTypeShow function to render the user type show page with title and logged_in
func UserTypeShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing User Type id: %v\n", err)
	}

	usertype, err := models.UserTypeShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(
		http.StatusOK,
		"usertypes/show.html",
		gin.H{
			"usertype":  usertype,
			"title":     usertype.Name,
			"logged_in": h.IsUserLoggedIn(c),
			"test_run":  os.Getenv("TEST_RUN") == "true",
		},
	)
}
