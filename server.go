package tracker

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Event struct {
	Name     string                 `json:"name" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (event Event) ToJson() string {
	jsonBytes, _ := json.Marshal(event)
	return string(jsonBytes[:])
}

func postEventToDynamoDB(event Event) {

	// FIXME: Used in testing...
	if event.Name == "fail" {
		TrackError(nil, "PostFailed", map[string]string{
			"event": event.ToJson(),
		})
	}

}

func NewTrackerEngine() *gin.Engine {
	routes := gin.Default()
	routes.Use(Recovery())

	addRoutes(routes)

	return routes
}

func TestTrackerEngine() *gin.Engine {
	routes := gin.New()
	addRoutes(routes)

	return routes
}

func addRoutes(routes *gin.Engine) {
	routes.GET("/status", func(context *gin.Context) {
		context.String(http.StatusOK, "I AM ALIVE")
	})

	routes.POST("/event", func(context *gin.Context) {
		// Ensure the JSON is valid before returning
		var event Event
		err := context.BindJSON(&event)
		if err == nil {
			go postEventToDynamoDB(event)
			context.String(http.StatusOK, "")
		}
	})
}
