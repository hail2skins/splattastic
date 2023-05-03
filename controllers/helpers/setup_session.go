package helpers

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func SetupSession(router *gin.Engine, userID uint) {
	// Configure session middleware
	store := memstore.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	router.Use(sessions.Sessions("mysession", store))

	// Set the user ID in the session
	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userID", fmt.Sprintf("%d", userID))
		_ = session.Save()
		c.Next()
	})
}
