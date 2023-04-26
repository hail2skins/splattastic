package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/controllers"
)

// AdminRoutes registers the admin routes
// routes folder exists to keep main.go clean
func AdminRoutes(r *gin.RouterGroup) {
	admin := r.Group("")
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
		admin.GET("/boardheights", controllers.BoardHeightsIndex)
		admin.GET("/boardheights/new", controllers.BoardHeightNew)
		admin.POST("/boardheights", controllers.BoardHeightCreate)
		admin.GET("/boardheights/:id", controllers.BoardHeightShow)
		admin.GET("/boardheights/edit/:id", controllers.BoardHeightEdit)
		admin.POST("/boardheights/:id", controllers.BoardHeightUpdate)
		admin.DELETE("/boardheights/:id", controllers.BoardHeightDelete)

		// Board types
		admin.GET("/boardtypes", controllers.BoardTypesIndex)
		admin.GET("/boardtypes/new", controllers.BoardTypeNew)
		admin.POST("/boardtypes", controllers.BoardTypeCreate)
		admin.GET("/boardtypes/:id", controllers.BoardTypeShow)
	}

}
