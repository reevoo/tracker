package logger

import (
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"testing"
)

var LastWrite string

type TestWriter struct{}

func (writer TestWriter) Write(p []byte) (n int, err error) {
	LastWrite = string(p[:])
	return len(LastWrite), nil
}

func TestIoJsonLogger(t *testing.T) {

	var message = map[string][]string{"foo": []string{"foo1", "foo2"}, "bar": []string{"baz"}}
	var logger = &IoJsonLogger{writer: TestWriter{}}

	Convey("the Log function ", t, func() {
		Convey("Logs JSON to the IO", func() {
			logger.Log(message)
			So(LastWrite, ShouldEqual, "{\"bar\":[\"baz\"],\"foo\":[\"foo1\",\"foo2\"]}\n")
		})
	})

}
