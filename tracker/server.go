package main

import (
	"github.com/getsentry/raven-go"
	. "github.com/reevoo/tracker"
	"os"
	"os/signal"
	"syscall"
)

var (
	dynamoUri           = os.Getenv("DYNAMODB_URI")
	Term      os.Signal = syscall.SIGTERM
)

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.SetRelease(os.Getenv("SENTRY_RELEASE"))
}

func main() {
	go exitOnInterrupt()

	server := NewServer(
		SentryErrorLogger{},
	)

	server.Run(":3000")
}

func exitOnInterrupt() {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt, os.Kill, Term)
	<-sigchan
	os.Exit(0)
}
