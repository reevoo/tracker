package tracker

import (
  "github.com/reevoo/tracker/event"
	"github.com/reevoo/tracker/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"net/http"
)

// The Server is the Tracker API.
type Server struct {
	engine      *gin.Engine
	errorLogger ErrorLogger
	eventLogger  EventLogger
}

// Parameters passed to NewServer().
type ServerParams struct {
	ErrorLogger ErrorLogger
	EventLogger  EventLogger
}

func SetServerMode(mode string) {
	gin.SetMode(mode)
}

// Creates a new Server that does no logging.
func NewServer(params ServerParams) Server {
	engine := gin.New()
	return initServer(engine, params)
}

func initServer(engine *gin.Engine, params ServerParams) Server {
	server := Server{
		engine:      engine,
		errorLogger: params.ErrorLogger,
		eventLogger:  params.EventLogger,
	}

	// Build the engine
	server.engine.Use(server.handleRecovery)
	server.engine.GET("/status", server.getStatus)
	server.engine.GET("/event", server.trackEvent)

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

// Track an event.
func (server Server) trackEvent(context *gin.Context) {
	// Ensure the JSON is valid before returning
	// context.BindJSON() sets an error status to the context on failure.
	var e = event.New(
		context.Request.URL.Query(),
	)

	if e.Empty() {
		context.String(http.StatusBadRequest, "No event params given.")
	} else {
		// We return the HTTP request quickly
		// and process the event in the background.
		go server.storeEvent(e)
		context.String(http.StatusOK, "")
	}
}

func (server Server) storeEvent(e event.Event) {
	err := server.eventLogger.Log(e)

	if err != nil {
		server.errorLogger.LogError(TrackerError{
			Name: "CannotStoreEvent",
			// TODO: Get context
		})
	}
}
