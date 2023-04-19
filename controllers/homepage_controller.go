package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/index.html",
		gin.H{
			"title": "Splattastic",
		})
}

func About(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/about.html",
		gin.H{
			"title": "About",
		})
}
