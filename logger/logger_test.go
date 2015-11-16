package logger_test

import (
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/reevoo/tracker/logger"
	"net"
	"os"
	"testing"
)

func TestIoJsonLogger(t *testing.T) {
	Convey("Setting Up a new logger", t, func() {
		Convey("With no config set", func() {
			l, _ := logger.New()
			So(l, ShouldHaveSameTypeAs, &logger.IoJsonLogger{})
		})

		Convey("With the fluent socket set in the environment", func() {
			socketFile := "/tmp/fluent.sock"
			listener, err := net.Listen("unix", socketFile)
			if err != nil {
				panic(err.Error())
			}
			defer listener.Close()
			os.Setenv("FLUENT_SOCKET", socketFile)
			l, _ := logger.New()
			So(l, ShouldHaveSameTypeAs, &logger.FluentLogger{})
		})
	})
}
