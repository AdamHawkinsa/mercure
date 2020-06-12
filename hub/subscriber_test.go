package hub

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDispatch(t *testing.T) {
	s := NewSubscriber("1", NewTopicSelectorStore())
	s.Topics = []string{"http://example.com"}
	go s.start()
	defer s.Disconnect()

	// Dispatch must be non-blocking
	// Messages coming from the history can be sent after live messages, but must be received first
	s.Dispatch(&Update{Topics: s.Topics, Event: Event{ID: "3"}}, false)
	s.Dispatch(&Update{Topics: s.Topics, Event: Event{ID: "1"}}, true)
	s.Dispatch(&Update{Topics: s.Topics, Event: Event{ID: "4"}}, false)
	s.Dispatch(&Update{Topics: s.Topics, Event: Event{ID: "2"}}, true)
	s.HistoryDispatched("")

	for i := 1; i <= 4; i++ {
		u := <-s.Receive()
		assert.Equal(t, strconv.Itoa(i), u.ID)
	}
}

func TestDisconnect(t *testing.T) {
	s := NewSubscriber("", NewTopicSelectorStore())
	s.Disconnect()
	// can be called two times without crashing
	s.Disconnect()

	assert.False(t, s.Dispatch(&Update{}, false))
}
