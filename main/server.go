package main

import (
	. "github.com/reevoo/tracker"
	"os"
	"os/signal"
	"syscall"
)

var (
	dynamoUri           = os.Getenv("DYNAMODB_URI")
	env                 = os.Getenv("GO_ENV")
	Term      os.Signal = syscall.SIGTERM
)

func main() {
	go exitOnInterrupt()

	// Release mode reduces the amount of logging.
	if env == "production" {
		SetServerMode("release")
	}

	server := NewServer(ServerParams{
		EventStore:  DynamoDBEventStore{},
		ErrorLogger: SentryErrorLogger{},
	})

	server.Run(":3000")
}

func exitOnInterrupt() {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt, os.Kill, Term)
	<-sigchan
	os.Exit(0)
}
