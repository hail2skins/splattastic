package main

import (
	"log"
	"os"
	"text/template"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/helpers"
	"github.com/hail2skins/splattastic/middlewares"
	"github.com/hail2skins/splattastic/routes"
	"github.com/hail2skins/splattastic/setup"
)

func main() {
	setup.LoadEnv()
	setup.LoadDatabase()
	serveApplication()

}

// serveApplication serves the application
func serveApplication() {
	// Define custom functions for templates
	// This allows a "mod" function in our edit_event_new.html page to order checkboxes for dives in such a way
	// as they form colums.
	funcMap := template.FuncMap{
		"mod":     func(i, j int) int { return i % j }, // used to order checkboxes for dives
		"shorten": helpers.Abbreviate,                  // used to abbreviate dive information. See helpers\abbreviate.go
		"seq":     helpers.Seq,                         // used to generate a sequence of numbers for the event show page
		"inc":     func(x int) int { return x + 1 },    // used to inc index on the event show page
	}

	r := gin.Default()
	r.SetFuncMap(funcMap)

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

	adminRoutes := r.Group("/admin", middlewares.RequireAdmin())
	routes.AdminRoutes(adminRoutes)

	userRoutes := r.Group("/user")
	routes.UserRoutes(userRoutes)

	log.Println("Server started")
	r.Run(":8080") // listen and serve on localhost:8080
}
