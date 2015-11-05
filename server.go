package tracker

import (
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"net/http"
)

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

func SetServerMode(mode string) {
	gin.SetMode(mode)
}

// Create a new Server.
func NewServer(params ServerParams) Server {
	return initServer(gin.Default(), params)
}

// Creates a new Server that does no logging.
// Handy in testing.
func NewSilentServer(params ServerParams) Server {
	engine := gin.New()
	return initServer(engine, params)
}

func initServer(engine *gin.Engine, params ServerParams) Server {
	server := Server{
		engine:      engine,
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
		// Set a new ID
		id, _ := uuid.NewV4()
		event.Id = *id

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
