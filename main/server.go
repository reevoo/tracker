package main

import (
	. "github.com/reevoo/tracker"
	"github.com/reevoo/tracker/logger"
	"os"
)

var (
	dynamoUri = os.Getenv("DYNAMODB_URI")
	env       = os.Getenv("GO_ENV")
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	server := NewServer(ServerParams{
		EventLogger: logger,
		ErrorLogger: SentryErrorLogger{},
		Environment: env,
	})

	server.Run(":3000")
}
