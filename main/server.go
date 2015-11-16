package main

import (
	. "github.com/reevoo/tracker"
	"github.com/reevoo/tracker/logger"
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

	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	server := NewServer(ServerParams{
		EventLogger: logger,
		ErrorLogger: SentryErrorLogger{},
	})

	server.Run(":3000")
}
