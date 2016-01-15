package logger

import (
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"net"
	"os"
	"testing"
)

func TestFluentLogger(t *testing.T) {
	Convey("Seting Up a new Fluent Logger", t, func() {
		socketFile := "/tmp/fluent.sock"
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			panic(err.Error())
		}
		defer l.Close()

		fluentLogger, _ := newFluentLogger(socketFile)

		So(fluentLogger.logger.Config.FluentNetwork, ShouldEqual, "unix")
		So(fluentLogger.logger.Config.FluentSocketPath, ShouldEqual, socketFile)

		os.Remove(socketFile)

	})

	Convey("Setting up when the fluentd server is not running", t, func() {
		socketFile := "/tmp/not-fluent.sock"
		_, err := newFluentLogger(socketFile)
		So(err.Error(), ShouldEqual, "dial unix /tmp/not-fluent.sock: connect: no such file or directory")
	})
}
