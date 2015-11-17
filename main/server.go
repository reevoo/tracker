package main

import (
	. "github.com/reevoo/tracker"
	"github.com/reevoo/tracker/logger"
	"os"
)

var (
	dynamoUri = os.Getenv("DYNAMODB_URI")
)

func main() {
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
