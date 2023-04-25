package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/setup"
)

func main() {
	setup.LoadEnv()
	setup.LoadDatabase()
	serveApplication()
}

func serveApplication() {
	r := gin.Default()

	// Configure session middleware
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions("mysession", store))

	// Define basic authentication middle to protect signup functionality
	authMiddleware := gin.BasicAuth(gin.Accounts{
		os.Getenv("SIGNUP_USERNAME"): os.Getenv("SIGNUP_PASSWORD"),
	})

	r.Use(middlewares.AuthenticateUser())

	r.LoadHTMLGlob("templates/**/**")

	r.GET("/", controllers.Home)
	r.GET("/about", controllers.About)
	r.GET("/login", controllers.LoginPage)
	r.GET("/signup", authMiddleware, controllers.SignupPage) // Protect signup functionality
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)

	r.Static("/css", "./static/css")
	r.Static("/img", "./static/img")
	r.Static("/scss", "./static/scss")
	r.Static("/vendor", "./static/vendor")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./img/favicon.ico")

	admin := r.Group("/admin", middlewares.RequireAdmin())
	{
		admin.GET("/", controllers.AdminDashboard)

		// User types
		admin.GET("/usertypes", controllers.UserTypeIndex)
		admin.GET("/usertypes/new", controllers.UserTypeNew)
		admin.POST("/usertypes", controllers.UserTypeCreate)
		admin.GET("/usertypes/:id", controllers.UserTypeShow)
		admin.GET("/usertypes/edit/:id", controllers.UserTypeEdit)
		admin.POST("/usertypes/:id", controllers.UserTypeUpdate)
		admin.DELETE("/usertypes/:id", controllers.UserTypeDelete)

		// Board heights
		admin.GET("/boardheights/new", controllers.BoardHeightsNew)
	}

	log.Println("Server started")
	r.Run(":8080") // listen and serve on localhost:8080
}
