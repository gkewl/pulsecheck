package sse

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"time"
)

// protocol header constants
var (
	IDHeader    = []byte("id")
	EventHeader = []byte("event")
	DataHeader  = []byte("data")
	RetryHeader = []byte("retry")
)

// message types/names
const (
	MessageType = "message"
	ErrorType   = "error"
)

type eventStreamReader struct {
	reader io.Reader
	Origin string
	Retry  time.Duration
}

func newEventStreamReader(reader io.Reader, origin string) *eventStreamReader {
	return &eventStreamReader{
		reader: reader,
		Origin: origin,
	}
}

func (decoder *eventStreamReader) decode(events chan<- *Event) error {
	scanner := bufio.NewScanner(decoder.reader)
	scanner.Split(SplitFunc())

	var lastEventID string

	// When a stream is parsed, a _data_ buffer and an _event name_
	// buffer must be associated with it. They must be initialized to
	// the empty string
	var data bytes.Buffer
	var eventName string

	dispatch := func() {
		// If the _data_ buffer is an empty string, set the _data_
		// buffer and the _event name_ buffer to the empty string
		// and abort these steps.
		if data.Len() != 0 {
			// If the data buffer's last character is a U+000A
			// LINE FEED (LF) character, then remove the last
			// character from the data buffer.
			data := data.Bytes()
			dataLen := len(data)
			if data[dataLen-1] == '\n' {
				dataLen--
			}
			eventData := make([]byte, dataLen)
			copy(eventData, data[:dataLen])
			// Otherwise, create an event that uses the
			// `MessageEvent` interface, with the event name
			// `message`...
			if eventName == "" {
				eventName = "message"
			}
			event := &Event{
				EventID: lastEventID,
				Origin:  decoder.Origin,
				Data:    eventData,
				Type:    eventName,
			}
			events <- event
		}
		eventName = ""
		data.Reset()
	}

	// Lines must be processed, in the order they are received, as
	// follows:
	for scanner.Scan() {
		token := scanner.Bytes()
		switch {
		case bytes.Equal(token, endOfLine):
			// If the line is empty (a blank line)
			// Dispatch the event
			dispatch()
		case bytes.Equal(token, commentMarker):
			// If the line starts with a U+003A COLON character (:)
			// Ignore the line.
		default:
			// Collect the characters on the line before the first
			// U+003A COLON character (:), and let _field_ be that
			// string.
			field := token

			// Collect the characters on the line after the first
			// U+003A COLON character (:), and let _value_ be that
			// string.
			var value []byte
			var hasValue bool
			if scanner.Scan() && bytes.Equal(scanner.Bytes(), fieldDelim) {
				hasValue = true
				if scanner.Scan() {
					value = scanner.Bytes()
					if bytes.Equal(value, endOfLine) {
						value = []byte{}
					} else {
						// swallow the end of line
						scanner.Scan()
					}
				}
			}

			// Field names must be compared literally, with no case
			// folding performed.
			switch {
			case bytes.Equal(field, EventHeader):
				// If the field name is "event"
				// Set the _event name_ buffer to field value.
				eventName = string(value)
			case bytes.Equal(field, DataHeader):
				// If the field name is "data"
				// Append the field value to the _data_ buffer, then
				// append a single U+000A LINE FEED (LF) character to
				// the _data_ buffer.
				if hasValue {
					data.Write(value)
					data.WriteRune('\n')
				}
			case bytes.Equal(field, IDHeader):
				// If the field name is "id"
				// Set the event stream's _last event ID_ to the field value.
				lastEventID = string(value)
			case bytes.Equal(field, RetryHeader):
				// If the field name is "retry"
				//
				// If the field value consists of only characters in
				// the range U+0030 DIGIT ZERO (0) to U+0039 DIGIT
				// NINE (9), then interpret the field value as an
				// integer in base ten, and set the event stream's
				// reconnection time to that integer. Otherwise,
				// ignore the field.
				retryMillis, err := strconv.Atoi(string(value))
				if err == nil && retryMillis >= 0 {
					decoder.Retry = time.Millisecond * time.Duration(retryMillis)
				}
			default:
				// Otherwise
				// The field is ignored.
			}
		}
	}
	dispatch()
	close(events)

	// the SSE spec defines all garbage input as valid; eventually we'll
	// implement saner garbage detection.
	return nil
}

func (decoder *eventStreamReader) decodeChan() <-chan *Event {
	eventChan := make(chan *Event)
	go decoder.decode(eventChan)
	return eventChan
}
