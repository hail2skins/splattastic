package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

func Home(c *gin.Context) {
	session := sessions.Default(c)
	alert := session.Flashes("alert")
	session.Save()

	c.HTML(
		http.StatusOK,
		"home/index.html",
		gin.H{
			"title":     "Splattastic",
			"logged_in": h.IsUserLoggedIn(c),
			"alert":     alert,
		})
}

func About(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/about.html",
		gin.H{
			"title":     "About",
			"logged_in": h.IsUserLoggedIn(c),
		})
}
