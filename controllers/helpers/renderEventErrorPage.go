package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	h "github.com/hail2skins/splattastic/helpers"
)

// RenderEventErrorPage renders the event error page
// Specifically within the EventNew function
func RenderEventErrorPage(c *gin.Context, message string) {
	c.HTML(
		http.StatusInternalServerError,
		"home/index.html",
		gin.H{
			"title":     "Splattastic",
			"alert":     message,
			"logged_in": h.IsUserLoggedIn(c),
			"header":    "Splattastic",
			"user_id":   c.GetUint("user_id"),
		},
	)
}
