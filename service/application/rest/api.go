package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *restService) setViewRoutes(parentRouteGroup *gin.RouterGroup) {
	view := parentRouteGroup.Group("")

	view.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
