package controllers

import (
	"net/http"

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
