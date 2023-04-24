package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

// AdminDashboard is the controller for the admin dashboard
func AdminDashboard(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"admin/dashboard.html",
		gin.H{
			"title":     "Admin Dashboard",
			"logged_in": h.IsUserLoggedIn(c),
		})
}
