package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hail2skins/splattastic/helpers"
)

func CheckCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/?alert=Invalid%20user%20ID")
			c.Abort()
			return
		}

		if !helpers.IsCurrentUser(c, id) {
			c.Redirect(http.StatusSeeOther, "/?alert=You%20are%20not%20authorized%20to%20access%20this%20resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
