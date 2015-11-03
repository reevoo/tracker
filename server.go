package tracker

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
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
	TrackError(nil, "PostFailed", map[string]string{
		"event": event.ToJson(),
	})
}

func CreateServer() *gin.Engine {
    routes := gin.Default()

    routes.Use(Recovery())

	routes.GET("/status", func(context *gin.Context) {
		context.String(http.StatusOK, "I AM ALIVE")
	})

	routes.POST("/event", func(context *gin.Context) {
		// Ensure the JSON is valid before returning
		var event Event
		err := context.BindJSON(&event)
		if(err == nil) {
			go postEventToDynamoDB(event)
			context.String(http.StatusOK, "")
		}		
	})

	return routes
}