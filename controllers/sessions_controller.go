package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage renders the login page
func LoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/login.html",
		gin.H{
			"title": "Login",
		},
	)
}

// SignupPage renders the signup page
func SignupPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home/signup.html",
		gin.H{
			"title": "Signup",
		},
	)
}
