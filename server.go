package tracker

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// An Event is a structure holding information
// about something that has happened in one of our applications.
type Event struct {
	Name     string                 `json:"name" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Converts the Event to JSON format.
func (event Event) ToJson() string {
	jsonBytes, _ := json.Marshal(event)
	return string(jsonBytes[:])
}

// The Server is the Tracker API.
type Server struct {
	engine      *gin.Engine
	errorLogger ErrorLogger
	eventStore  EventStore
}

// Parameters passed to NewServer().
type ServerParams struct {
	ErrorLogger ErrorLogger
	EventStore  EventStore
}

// Create a new Server.
func NewServer(params ServerParams) Server {
	server := Server{
		engine:      gin.Default(),
		errorLogger: params.ErrorLogger,
		eventStore:  params.EventStore,
	}

	// Build the engine
	server.engine.Use(server.handleRecovery)
	server.engine.GET("/status", server.getStatus)
	server.engine.POST("/event", server.postEvent)

	return server
}

// Start the server and handle incoming requests.
func (server Server) Run(port string) {
	server.engine.Run(port)
}

func (server Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	server.engine.ServeHTTP(resp, req)
}

// Gin middleware that ensures a panic is recovered
// and the error is logged.
func (server Server) handleRecovery(context *gin.Context) {
	defer func() {
		meta := map[string]string{
			"endpoint": context.Request.URL.RequestURI(),
		}
		if err := recover(); err != nil {
			server.errorLogger.LogError(NewTrackerErrorFromError(err.(error), meta))
			context.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	context.Next()
}

// Report the status of the server.
func (server Server) getStatus(context *gin.Context) {
	context.String(http.StatusOK, "I AM ALIVE")
}

// Store an event.
func (server Server) postEvent(context *gin.Context) {
	// Ensure the JSON is valid before returning
	// context.BindJSON() sets an error status to the context on failure.
	var event Event
	err := context.BindJSON(&event)

	if err == nil {
		// We return the HTTP request quickly
		// and process the event in the background.
		go server.storeEvent(event)
		context.String(http.StatusOK, "")
	}
}

func (server Server) storeEvent(event Event) {
	err := server.eventStore.Store(event)

	if err != nil {
		server.errorLogger.LogError(TrackerError{
			Name: "CannotStoreEvent",
			// TODO: Get context
		})
	}
}
