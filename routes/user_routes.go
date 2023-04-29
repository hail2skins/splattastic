package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
	"github.com/hail2skins/splattastic/middlewares"
)

// UserRoutes registers the user routes
// routes folder exists to keep main.go clean

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("")
	{
		user.GET("/:id", controllers.UserShow)
		user.GET("/edit/:id", middlewares.CheckCurrentUser(), controllers.UserEdit)
	}
}
