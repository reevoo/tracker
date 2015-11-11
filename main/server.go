package main

import (
	. "github.com/reevoo/tracker"
	"os"
	"syscall"
)

var (
	dynamoUri           = os.Getenv("DYNAMODB_URI")
	env                 = os.Getenv("GO_ENV")
	Term      os.Signal = syscall.SIGTERM
)

func main() {
	// Release mode reduces the amount of logging.
	if env == "production" {
		SetServerMode("release")
	}

	server := NewServer(ServerParams{
		EventStore:  NewEventLogger(nil),
		ErrorLogger: SentryErrorLogger{},
	})

	server.Run(":3000")
}
