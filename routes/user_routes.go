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
		user.POST("/:id", middlewares.CheckCurrentUser(), controllers.UserUpdate)
		user.GET("/:id/event/new", middlewares.CheckCurrentUser(), controllers.EventNew)
		user.POST("/:id/event", middlewares.CheckCurrentUser(), controllers.EventCreate)
		user.GET("/:id/event/:event_id", controllers.EventShow)
		user.GET("/:id/events", controllers.GetUserEvents)
		user.GET("/:id/event/:event_id/edit", middlewares.CheckCurrentUser(), controllers.EventEdit)
		user.POST("/:id/event/:event_id", middlewares.CheckCurrentUser(), controllers.EventUpdate)
		user.DELETE("/:id/event/:event_id", middlewares.CheckCurrentUser(), controllers.EventDelete)
		user.POST("/:id/event/:event_id/scores", middlewares.CheckCurrentUser(), controllers.EventScoreUpsert)
		user.GET("/:id/event/:event_id/dive/:dive_id/scores", controllers.FetchScores)
		user.GET("/:id/event/:event_id/dive/:dive_id/total", controllers.EventDiveScoreTotal)
		user.GET("/:id/event/:event_id/meet_score", controllers.EventMeetScore)
	}
}
