package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/models"
)

type SignupForm struct {
	Email           string `form:"email" binding:"required" json:"email"`
	Password        string `form:"password" binding:"required" json:"password"`
	ConfirmPassword string `form:"confirm_password" binding:"required" json:"confirm_password"`
	UserType        string `form:"user_type" binding:"required" json:"user_type"`
	FirstName       string `form:"firstname" binding:"required" json:"first_name"`
	LastName        string `form:"lastname" binding:"required" json:"last_name"`
	Username        string `form:"username" binding:"required" json:"username"`
}

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

// Signup handles the signup form submission
func Signup(c *gin.Context) {
	// Using binding from struct above
	var form SignupForm
	if err := c.Bind(&form); err != nil {
		c.HTML(
			http.StatusBadRequest,
			"home/signup.html",
			gin.H{
				"title": "Signup",
				"error": err.Error(),
				"form":  form,
			})
		return
	}
	// Check if email and username already exists
	available, err := models.CheckEmailUsernameAvailable(form.Email, form.Username)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			"home/signup.html",
			gin.H{
				"title": "Signup",
				"alert": "An error occurred while checking availability",
			},
		)
		return
	}
	if !available {
		c.HTML(
			http.StatusIMUsed,
			"home/signup.html",
			gin.H{
				"title": "Signup",
				"alert": "Email or username already exists",
			},
		)
		return
	}
	// Check if passwords match
	if form.Password != form.ConfirmPassword {
		c.HTML(
			http.StatusNotAcceptable,
			"home/signup.html",
			gin.H{
				"alert": "Passwords do not match",
			},
		)
		return
	}
	// Create user
	user, err := models.UserCreate(
		form.Email,
		form.Password,
		form.FirstName,
		form.LastName,
		form.Username,
		models.UserType(form.UserType), // <- Cast string to models.UserType
	)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			"home/signup.html",
			gin.H{
				"alert": "Error creating user",
				"title": "Signup",
			},
		)
		return
	}

	if user.ID == 0 {
		c.HTML(
			http.StatusNotAcceptable,
			"home/signup.html",
			gin.H{
				"alert": "Error creating user",
				"title": "Signup",
			},
		)
	} else {
		helpers.SessionSet(c, uint64(user.ID))
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}
