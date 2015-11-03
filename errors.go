package tracker

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"os"
)

// Naughty global variable for storing logged errors
var Errors = make(map[string]interface{})
var env    = os.Getenv("GO_ENV")

// Track an error internally.
// This should be used instead of panicking!
func TrackError(error interface{}, name string, meta map[string]string) {
	if(env == "PROD") {
		// Post to Sentry if on production
		trackErrorWithSentry(error, name, meta)
	} else {
		// Add to Errors if not
		trackErrorLocally(error, name, meta)
	}
}


func trackErrorWithSentry(err interface{}, name string, meta map[string]string) {
	meta["name"] = name

	errStr := fmt.Sprint(err)
	packet := raven.NewPacket(errStr, raven.NewException(errors.New(errStr), raven.NewStacktrace(2, 3, nil)))
	raven.Capture(packet, meta)
}

func trackErrorLocally(err interface{}, name string, meta map[string]string) {
	Errors[name] = err
	fmt.Printf("[ERROR] %s: %s", name, err)
}