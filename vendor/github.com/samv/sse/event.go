package sse

import (
	"bytes"
	"io"
)

// Event is a structure holding SSE-compliant events
type Event struct {
	// EventID is the ID of the event, or a previous event
	EventID string

	// Type is variously called "event type" and "event name" in the
	// TR.  Defaults to "message".  You must listen for specific named
	// event types to receive them.
	Type string

	// Error contains a go error if the error came from the client
	// (eg, connection problems)
	Error error

	// Data is the body of the event, and always terminated with a
	// line feed.  "Simple" events have this empty.  Returned as a
	// []byte as go decoders generally use that, but can't be binary!
	Data []byte

	// the RFC "Origin" field
	Origin string
}

// Reader returns a reader for convenient passing to decoder
// functions.
func (ev *Event) Reader() io.ReadSeeker {
	return bytes.NewReader(ev.Data)
}

// GetData implements the SinkEvent interface, so that sse.Event can
// be used for sinks and sources.
func (ev *Event) GetData() ([]byte, error) {
	return ev.Data, nil
}
