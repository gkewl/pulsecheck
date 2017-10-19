package sse

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	jsonNil = []byte("null")
)

// AnyEventFeed represents an event feed which does not return events
// which have their own byte marshall method.  Instead, encoding/json
// is used to marshal the events.
type AnyEventFeed interface {
	GetEventChan(clientCloseChan <-chan struct{}) <-chan interface{}
}

type jsonEncoderFeed struct {
	anyFeed   AnyEventFeed
	closeChan <-chan struct{}
}

type jsonEvent struct {
	data interface{}
}

func (je jsonEvent) GetData() ([]byte, error) {
	data, err := json.Marshal(je.data)
	if err == nil {
		if bytes.Equal(data, jsonNil) {
			return nil, nil
		}
	}
	return data, err
}

// NewJSONEncoderFeed converts an 'AnyEventFeed' to an EventFeed by
// marshalling each of the events returned via the EventFeed via
// encoding/json.
func NewJSONEncoderFeed(anyEventFeed AnyEventFeed) EventFeed {
	return &jsonEncoderFeed{
		anyFeed: anyEventFeed,
	}
}

// GetEventChan returns an SSE-compatible event feed, satisfying the EventFeed interface.
func (jea *jsonEncoderFeed) GetEventChan(clientCloseChan <-chan struct{}) <-chan SinkEvent {
	anyFeedClosedChan := make(chan struct{})
	jsonEventChan := jea.anyFeed.GetEventChan(anyFeedClosedChan)
	eventChan := make(chan SinkEvent)
	go func() {
	anyEventLoop:
		for {
			select {
			case _, ok := <-clientCloseChan:
				if !ok {
					close(anyFeedClosedChan)
					clientCloseChan = nil
				}
			case event, ok := <-jsonEventChan:
				if !ok {
					if eventChan != nil {
						close(eventChan)
						eventChan = nil
					}
					break anyEventLoop
				}
				eventChan <- &jsonEvent{event}
			}
		}
	}()
	return eventChan
}

// SinkJSONEvents is a wrapped-up handler for responding with anything
// that encoding/json can marshal, that you can happily `return` to in
// a net/http handler.
func SinkJSONEvents(w http.ResponseWriter, code int, feed AnyEventFeed) error {
	return SinkEvents(w, code, NewJSONEncoderFeed(feed))
}
