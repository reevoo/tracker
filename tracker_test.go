package tracker_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ExampleMetadata = Metadata{
	"meta_str":   "Hello, World!",
	"meta_int":   123,
	"meta_float": 1.23,
	"meta_bool":  true,
	"meta_nil":   nil,
	"meta_map": Metadata{
		"really_meta_str": "Goodbye, cruel world!",
	},
	"meta_list": []interface{}{"one", 2, 3.14, true, nil, Metadata{}},
}

func init() {
	// Disable verbose logging on Test.
	fmt.Print("NOTE: Verbose logging is disabled on test.\nTo change this, go to `tracker_test.go`.\n\n")
	SetServerMode("release")
}

// Performs a GET request.
func get(server *Server, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

// Performs a POST request.
func post(server *Server, url string, body string) *httptest.ResponseRecorder {
	bodyReader := bytes.NewBufferString(body)
	req, _ := http.NewRequest("POST", url, bodyReader)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

// Hook the Go testing framework into Ginkgo and Gomega.
func TestTracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracker Suite")
}
