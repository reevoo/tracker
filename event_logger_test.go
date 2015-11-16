package tracker

import (
  "github.com/reevoo/tracker/event"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
	"net"
	"os"
)

var LastWrite string

type TestWriter struct{}

func (writer TestWriter) Write(p []byte) (n int, err error) {
	LastWrite = string(p[:])
	return len(LastWrite), nil
}

var _ = Describe("NewEventLogger", func() {
	var (
		eventLogger EventLogger
		e       event.Event
	)

	BeforeEach(func() {
		e = event.New(map[string][]string{"param1": []string{"val1"}})
	})

	Describe("IoEventLogger", func() {
		BeforeEach(func() {
			eventLogger, _ = NewEventLogger(TestWriter{})
		})

		Describe("Store()", func() {

			It("Writes JSON to writer", func() {
				eventLogger.Log(e)
				Expect(LastWrite).To(ContainSubstring(e.ToJson()))
			})

			It("Writes a new line", func() {
				eventLogger.Log(e)
				Expect(LastWrite).To(HaveSuffix("\n"))
			})

		})
	})

	Describe("FluentEventLogger", func() {

		var socketFile = "/tmp/fluent.sock"

		BeforeEach(func() {
			os.Setenv("FLUENT_SOCKET", socketFile)

			l, err := net.Listen("unix", socketFile)
			if err != nil {
				Fail(err.Error())
			}
			defer l.Close()

			eventLogger, _ = NewEventLogger(nil)
		})

		AfterEach(func() {
			os.Remove(socketFile)
			os.Setenv("FLUENT_SOCKET", "")
		})

		It("Sets up the fluent logger correctly", func() {
			fluentLogger := eventLogger.(*FluentEventLogger)
			Expect(fluentLogger.logger.Config.FluentNetwork).To(Equal("unix"))
			Expect(fluentLogger.logger.Config.FluentSocketPath).To(Equal(socketFile))
		})

	})

})
