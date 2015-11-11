package tracker_test

import (
	. "github.com/reevoo/tracker"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
	"os"
)

var LastWrite string

type TestStore struct{}

func (store TestStore) Write(p []byte) (n int, err error) {
	LastWrite = string(p[:])
	return len(LastWrite), nil
}

var _ = Describe("NewEventLogger", func() {

	It("Writes to STDOUT by default", func() {
		eventLogger := NewEventLogger(nil)
		Expect(eventLogger.Writer).To(Equal(os.Stdout))
	})

	Describe("Store()", func() {
		var (
			eventLogger EventLogger
			event Event
		)

		BeforeEach(func () {
			eventLogger = NewEventLogger(TestStore{})
			event = NewEvent(map[string][]string{"param1": []string{"val1"}})
		})

		It("Writes JSON to writer", func() {
			eventLogger.Store(event)
			Expect(LastWrite).To(ContainSubstring(event.ToJson()))
		})

		It("Writes a new line", func() {
			eventLogger.Store(event)
			Expect(LastWrite).To(HaveSuffix("\n"))
		})

	})
})
