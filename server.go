package tracker

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateServer() *gin.Engine {
	routes := gin.Default()

	routes.Use(Recovery())

	routes.GET("/status", func(context *gin.Context) {
		context.String(http.StatusOK, "I AM ALIVE")
	})

	return routes
}