package tracker

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"os"
	"sync"
)

var env = os.Getenv("GO_ENV")

// Internal convenience struct for errors.
type trackerError struct {
	name string
	desc string
	meta map[string]string // Has to be string for Raven.
}

func newTrackerError(err interface{}, name string, meta map[string]string) trackerError {
	errStr := fmt.Sprint(err)

	return trackerError{
		name: name,
		desc: errStr,
		meta: meta,
	}
}

func (err trackerError) ToString() string {
	return fmt.Sprintf("%s: %s", err.name, err.desc)
}

func (err trackerError) ToMap() map[string]string {
	trackerMap := err.meta
	trackerMap["name"] = err.name
	trackerMap["desc"] = err.desc

	return trackerMap
}

// The error log is a thread-safe error logger.
type errorLog struct {
	errors []trackerError
	lock   sync.Mutex
}

func (log *errorLog) Push(err trackerError) {
	log.lock.Lock()
	log.errors = append(log.errors, err)
	log.lock.Unlock()
}

func (log errorLog) Count() int {
	log.lock.Lock()
	count := len(log.errors)
	log.lock.Unlock()

	return count
}

func (log *errorLog) Clear() {
	log.lock.Lock()
	log.errors = nil
	log.lock.Unlock()
}

// Naughty global variable for storing logged errors
var Errors = errorLog{}

// Track an error internally.
// This should be used instead of panicking!
func TrackError(err interface{}, name string, meta map[string]string) {
	tError := newTrackerError(err, name, meta)

	if env == "PROD" {
		// Post to Sentry if on production
		trackErrorWithSentry(tError)
	} else {
		// Add to Errors if not
		trackErrorLocally(tError)
	}
}

func trackErrorWithSentry(err trackerError) {
	packet := raven.NewPacket(err.name, raven.NewException(errors.New(err.name), raven.NewStacktrace(2, 3, nil)))
	raven.Capture(packet, err.ToMap())
}

func trackErrorLocally(err trackerError) {
	Errors.Push(err)
	fmt.Printf("[ERROR #%d] %s", Errors.Count(), err.ToString())
}
