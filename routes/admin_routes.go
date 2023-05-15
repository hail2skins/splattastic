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
		admin.GET("/boardtypes/edit/:id", controllers.BoardTypeEdit)
		admin.POST("/boardtypes/:id", controllers.BoardTypeUpdate)
		admin.DELETE("/boardtypes/:id", controllers.BoardTypeDelete)

		// Dive types
		admin.GET("/divetypes", controllers.DiveTypesIndex)
		admin.GET("/divetypes/new", controllers.DiveTypeNew)
		admin.POST("/divetypes", controllers.DiveTypeCreate)
		admin.GET("/divetypes/:id", controllers.DiveTypeShow)
		admin.GET("/divetypes/edit/:id", controllers.DiveTypeEdit)
		admin.POST("/divetypes/:id", controllers.DiveTypeUpdate)
		admin.DELETE("/divetypes/:id", controllers.DiveTypeDelete)

		// Dive groups
		admin.GET("/divegroups", controllers.DiveGroupsIndex)
		admin.GET("/divegroups/new", controllers.DiveGroupNew)
		admin.POST("/divegroups", controllers.DiveGroupCreate)
		admin.GET("/divegroups/:id", controllers.DiveGroupShow)
		admin.GET("/divegroups/edit/:id", controllers.DiveGroupEdit)
		admin.POST("/divegroups/:id", controllers.DiveGroupUpdate)
		admin.DELETE("/divegroups/:id", controllers.DiveGroupDelete)

		// Dives
		admin.GET("/dives", controllers.DivesIndex)
		admin.GET("/dives/new", controllers.DiveNew)
		admin.POST("/dives", controllers.DiveCreate)
		admin.GET("/dives/:id", controllers.DiveShow)
		admin.GET("/dives/edit/:id", controllers.DiveEdit)
		admin.POST("/dives/:id", controllers.DiveUpdate)
		admin.DELETE("/dives/:id", controllers.DiveDelete)

		// Event Types
		admin.GET("/eventtypes", controllers.EventTypesIndex)
		admin.GET("/eventtypes/new", controllers.EventTypeNew)
		admin.POST("/eventtypes", controllers.EventTypeCreate)
		admin.GET("/eventtypes/:id", controllers.EventTypeShow)
		admin.GET("/eventtypes/edit/:id", controllers.EventTypeEdit)
		admin.POST("/eventtypes/:id", controllers.EventTypeUpdate)
		admin.DELETE("/eventtypes/:id", controllers.EventTypeDelete)

		// Markers
		admin.GET("/markers", controllers.MarkersIndex)
		admin.GET("/markers/new", controllers.MarkerNew)
		admin.POST("/markers", controllers.MarkerCreate)
		admin.GET("/markers/:id", controllers.MarkerShow)

	}

}
