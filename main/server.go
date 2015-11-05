package main

import (
	. "github.com/reevoo/tracker"
	"os"
	"os/signal"
	"syscall"
)

var (
	dynamoUri           = os.Getenv("DYNAMODB_URI")
	Term      os.Signal = syscall.SIGTERM
)

func main() {
	go exitOnInterrupt()

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
