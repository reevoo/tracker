package main

import (
	. "github.com/reevoo/tracker"
	"github.com/getsentry/raven-go"
	"os"
	"os/signal"
	"syscall"
)

var (
	dynamoUri            = os.Getenv("DYNAMODB_URI")
	Term       os.Signal = syscall.SIGTERM
)

func init() {
  raven.SetDSN(os.Getenv("SENTRY_DSN"))
  raven.SetRelease(os.Getenv("SENTRY_RELEASE"))
}

func main() {

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt, os.Kill, Term)
		<-sigchan
		os.Exit(0)
	}()

	routes := NewTrackerEngine()
	routes.Run(":3000")
}
