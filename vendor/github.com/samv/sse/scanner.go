package sse

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var (
	// standard forms of tokens used in this implementation
	endOfLine  = []byte{'\n'}
	fieldDelim = []byte{':', ' '}

	// Comments are used for keep-alives
	commentMarker = []byte{':'}

	// the spec specifies this as a replacement character for errors
	replacementChar = []byte("\uFFFD")
)

// parser state.  Use a pointer to an integer as a very simple state
// machine.
type parseState int

const (
	startOfStream parseState = iota
	startOfLine
	comment
	readField
	readDelim
	midEOL
	atEndOfLine
	endOfStream
)

func newSplitter() *parseState {
	state := startOfStream
	return &state
}

// SplitFunc returns a bufio.SplitFunc for the text/event-stream MIME
// type.  This is a stateful scanner which tokenizes the event stream.
func SplitFunc() func([]byte, bool) (int, []byte, error) {
	dfa := newSplitter()
	return func(data []byte, atEOF bool) (int, []byte, error) {
		return dfa.scan(data, atEOF)
	}
}

// scan is the function called by Scanner.Split() to find the next token.
func (state *parseState) scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		if *state == atEndOfLine || *state == midEOL {
			*state = endOfStream
			return 0, endOfLine, nil
		}
		return 0, nil, nil
	}
	switch *state {
	case startOfStream:
		advance, token, err = state.scanBOM(data, atEOF)
		if advance != -1 {
			return advance, token, err
		}
		fallthrough
	case startOfLine:
		return state.scanStartOfLine(data, atEOF)
	case midEOL:
		if data[0] == '\n' {
			return 1, nil, nil
		}
		*state = startOfLine
		return state.scanStartOfLine(data, atEOF)
	case comment, readDelim:
		return state.scanAnyChar(data, atEOF)
	case readField:
		return state.scanDelim(data, atEOF)
	case atEndOfLine:
		return state.scanEndOfLine(data, atEOF)
	default:
		panic(fmt.Sprintf("bad state %v", state))
	}
}

// scanBOM strips off a useless byte-order mark
// > eventsource ABNF rule: stream        = [ bom ] *event
func (state *parseState) scanBOM(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var maybe, no bool
	if data[0] != '\xEF' {
		no = true
	} else if len(data) == 1 {
		maybe = true
	} else if data[1] != '\xBB' {
		no = true
	} else if len(data) == 2 {
		maybe = true
	} else if data[2] != '\xBF' {
		no = true
	}
	if maybe {
		return 0, nil, nil
	}
	*state = startOfLine
	if no {
		return -1, nil, nil
	}
	return 3, nil, nil
}

// scanStartOfLine scans most of this rule:
// >   event         = *( comment / field ) end-of-line
func (state *parseState) scanStartOfLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch data[0] {
	case '\r', '\n':
		// an event is complete
		*state = atEndOfLine
		return state.scanEndOfLine(data, atEOF)
	case ':':
		*state = comment
		return 1, commentMarker, nil
	default:
		// field name character
		idx := bytes.IndexAny(data, "\n\r:")
		if idx == -1 {
			if !atEOF {
				return 0, nil, nil
			}
			idx = len(data)
		}
		*state = readField
		return idx, validUTF8(data[:idx]), nil
	}
}

// scanAnyChar scans:
// > any-char      = %x0000-0009 / %x000B-000C / %x000E-10FFFF
// >                ; a Unicode character other than U+000A LINE FEED (LF) or U+000D CARRIAGE RETURN (CR)
func (state *parseState) scanAnyChar(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch data[0] {
	case '\r', '\n':
		// no data.  can't emit an empty token, so...
		*state = atEndOfLine
		return state.scanEndOfLine(data, atEOF)
	default:
		idx := bytes.IndexAny(data, "\n\r")
		if idx == -1 && !atEOF {
			return 0, nil, nil
		}
		if idx == -1 {
			idx = len(data)
		}
		*state = atEndOfLine
		return idx, validUTF8(data[:idx]), nil
	}
}

// scanDelim scans the colon and space from:
// > field         = 1*name-char [ colon [ space ] *any-char ] end-of-line
func (state *parseState) scanDelim(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch data[0] {
	case '\r', '\n':
		// no data.  can't emit an empty token, so...
		*state = atEndOfLine
		return state.scanEndOfLine(data, atEOF)
	case ':':
		if len(data) == 1 {
			if atEOF {
				*state = readDelim
				return 1, fieldDelim, nil
			}
			return 0, nil, nil
		}
		*state = readDelim
		advance = 1
		if data[1] == ' ' {
			advance = 2
		}
		return advance, fieldDelim, nil
	default:
		// should never happen; can't get into readField state unless
		// next char is '\r', '\n' or ':'
		return 0, nil, errors.Errorf("expected: field delimiter or end of line, saw: %v", data[0])
	}
}

// scanEndOfLine scans a line ending.  Lines may be terminated by \r,
// \n *or* \r\n (or the end of the stream), complicating this somewhat.
// rule:
// > end-of-line   = ( cr lf / cr / lf / eof )
func (state *parseState) scanEndOfLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch data[0] {
	case '\r':
		if len(data) == 1 {
			if !atEOF {
				*state = midEOL
			}
			return 1, endOfLine, nil
		} else if data[1] == '\n' {
			*state = startOfLine
			return 2, endOfLine, nil
		} else {
			*state = startOfLine
			return 1, endOfLine, nil
		}
	case '\n':
		*state = startOfLine
		return 1, endOfLine, nil
	default:
		// should never happen; can't get into readField state unless
		// next char is '\r' or '\n'
		return 0, nil, errors.Errorf("Expected end-of-line, found: %v", data[0])
	}
}

// validUTF8 is a function which ensures this part of the spec:
// > Bytes or sequences of bytes that are not valid UTF-8 sequences must
// > be interpreted as the U+FFFD REPLACEMENT CHARACTER.
func validUTF8(network []byte) []byte {
	if utf8.Valid(network) {
		return network
	}
	// FIXME - can I just return []byte(string([]rune(network)))?
	valid := make([]byte, 0, len(network)+32)
	for len(network) > 0 {
		rune, size := utf8.DecodeRune(network)
		if rune == utf8.RuneError {
			// "or sequences of bytes" - skip to the next rune start
			network = network[1:]
			for len(network) > 0 && !utf8.RuneStart(network[0]) {
				network = network[1:]
			}
			valid = append(valid, []byte(string(rune))...)
		} else {
			valid = append(valid, network[:size]...)
			network = network[size:]
		}
	}
	return valid
}
