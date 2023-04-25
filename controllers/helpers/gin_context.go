package helpers

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func GinContext(req *http.Request, w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}
