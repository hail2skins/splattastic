package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/setup"
)

func main() {
	setup.LoadEnv()
	//	setup.LoadDatabase()
	serveApplication()
}

func serveApplication() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/**/**")

	r.GET("/", controllers.Home)
	r.GET("/about", controllers.About)

	r.Static("/css", "./static/css")
	r.Static("/img", "./static/img")
	r.Static("/scss", "./static/scss")
	r.Static("/vendor", "./static/vendor")
	r.Static("/js", "./static/js")
	r.StaticFile("/favicon.ico", "./img/favicon.ico")

	log.Println("Server started")
	r.Run(":8080") // listen and serve on localhost:8080
}
