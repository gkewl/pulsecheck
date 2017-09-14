package sse

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type ReadyState int32
type WantFlag int32

const (
	// bitwise flags for selecting which events to pass through
	WantErrors    WantFlag = 1
	WantOpenClose WantFlag = 2
	WantMessages  WantFlag = 4

	// ReadyState is an enum representing the current status of the connection
	Connecting ReadyState = 0
	Open       ReadyState = 1
	Closed     ReadyState = 2
)

// SSEClient is a wrapper for an http Client which can also hold one
// SSE session open, and implements EventSource.  You can use this API
// directly, but the EventSource API will steer you towards message
// protocols and behavior which also work from browsers.
type SSEClient struct {
	http.Client

	// all state changes are guarded by this
	sync.Mutex

	// for clean shutdown
	wg sync.WaitGroup

	// info of the request.  The SSE spec only allows GET requests and no body!
	url    *url.URL
	origin string

	// current state of the client
	readyState int32

	// current SSE response, if connected
	response *http.Response

	// reader that decodes the response
	reader      *eventStreamReader
	readerError error
	eventStream chan *Event

	// whether to reconnect and after what time
	reconnect     bool
	reconnectTime time.Duration

	canceler

	// reconnect time after losing connection
	// what messages are being sent through
	wantStd int32

	// channels for these "standard" messages
	messageChan chan *Event // standard message channel
	openChan    chan bool   // open/close notification channel
	errorChan   chan error  // error channel
}

// NewSSEClient creates a new client which can make a single SSE call.
func NewSSEClient(options ...ClientOption) *SSEClient {
	ssec := &SSEClient{
		readyState:    int32(Connecting),
		reconnectTime: time.Second,
		reconnect:     true,
	}
	for _, option := range options {
		option.Apply(ssec)
	}
	ssec.initTransport()
	if ssec.wantStd == 0 {
		ssec.wantStd = int32(WantMessages)
	}
	if ssec.messageChan == nil {
		ssec.messageChan = make(chan *Event)
	}
	if ssec.openChan == nil {
		ssec.openChan = make(chan bool)
	}
	if ssec.errorChan == nil {
		ssec.errorChan = make(chan error)
	}
	return ssec
}

// GetStream makes a GET request and returns a channel for *all* events read
func (ssec *SSEClient) GetStream(uri string) error {
	ssec.Lock()
	defer ssec.Unlock()
	var err error
	if ssec.url, err = url.Parse(uri); err != nil {
		return errors.Wrap(err, "error parsing URL")
	}
	ssec.wg.Add(1)
	go ssec.process()
	return err
}

func (ssec *SSEClient) makeRequest() (*http.Request, error) {
	ssec.Lock()
	defer ssec.Unlock()
	request, err := http.NewRequest("GET", ssec.url.String(), nil)
	if err != nil {
		return nil, err
	}
	request = ssec.wrapRequest(request)

	request.Header.Set("Accept", "text/event-stream")
	return request, nil
}

func (ssec *SSEClient) connect() error {
	var err error
	var req *http.Request
	if req, err = ssec.makeRequest(); err != nil {
		return errors.Wrap(err, "error making request")
	}

	ssec.response, err = ssec.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error fetching URL")
	}
	ssec.origin = makeOrigin(ssec.response.Request.URL)
	switch ssec.response.StatusCode {
	case http.StatusOK:
		var typeOk bool
		var contentType string
		// FIXME - this is not RFC-compliant
		for _, mimeType := range ssec.response.Header["Content-Type"] {
			contentType = mimeType
			if strings.Index(mimeType, "text/event-stream") >= 0 {
				typeOk = true
				break
			}
		}
		if typeOk {
			return nil
		} else {
			// HTTP 200 OK responses that have a Content-Type other
			// than text/event-stream (or some other supported type)
			// must cause the user agent to fail the connection.
			err = errors.Errorf("Content type not text/event-stream: %s", contentType)
		}
	case http.StatusNoContent, http.StatusResetContent:
		// HTTP 204 No Content, and 205 Reset Content responses are
		// equivalent to 200 OK responses with the right MIME type but
		// no content, and thus must reset the connection.
		return nil
	case http.StatusMovedPermanently, statusPermanentRedirect:
		// TODO - update the URL, origin
	case http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		// TODO - retry different URL but don't update URL.  Update origin?
	case http.StatusUseProxy:
		// TODO
	case http.StatusUnauthorized, http.StatusProxyAuthRequired:
		// TODO
	default:
		err = errors.Errorf("Bad response: %s", ssec.response.Status)
	}

	// TODO: reset the connection here if error

	return err
}

func (ssec *SSEClient) setReconnect(should bool) (changed bool) {
	ssec.Lock()
	changed = (should != ssec.reconnect)
	ssec.reconnect = false
	ssec.Unlock()
	return
}

// Reopen allows a connection which was closed to be re-opened again.
func (ssec *SSEClient) Reopen() {
	if ssec.setReconnect(true) {
		ssec.wg.Add(1)
		go ssec.process()
	}
}

func (ssec *SSEClient) shouldReconnect() bool {
	ssec.Lock()
	should := ssec.reconnect
	ssec.Unlock()
	return should
}

func (ssec *SSEClient) wants(what WantFlag) bool {
	return (atomic.LoadInt32(&ssec.wantStd) & int32(what)) != 0
}

func (ssec *SSEClient) demand(what WantFlag) {
	if !ssec.wants(what) {
		// seems silly to use both atomic load and mutex, but this is
		// a read, modify, update, and the fastpath is lockless atomic read
		ssec.Lock()
		wants := atomic.LoadInt32(&ssec.wantStd) | int32(what)
		atomic.StoreInt32(&ssec.wantStd, wants)
		ssec.Unlock()
	}
}

// Messages returns a channel from which events can be read
func (ssec *SSEClient) Messages() <-chan *Event {
	ssec.demand(WantMessages)
	return ssec.messageChan
}

func (ssec *SSEClient) emit(event *Event) {
	// only send "message" events down the Messages() channel
	switch event.Type {
	case MessageType:
		if ssec.wants(WantMessages) {
			ssec.messageChan <- event
		}
	case ErrorType:
	default:
	}
}

func (ssec *SSEClient) Opens() <-chan bool {
	ssec.demand(WantOpenClose)
	return ssec.openChan
}

func (ssec *SSEClient) emitOpenClose(which bool) {
	if ssec.wants(WantOpenClose) {
		ssec.openChan <- which
	}
}

// Errors returns a channel from which errors will be returned.  As an
// error indicates that the SSE channel is closed, there will only be
// one error returned before you reset the client (via Reopen())
func (ssec *SSEClient) Errors() <-chan error {
	ssec.demand(WantErrors)
	return ssec.errorChan
}

func (ssec *SSEClient) emitError(err error) {
	if ssec.wants(WantErrors) {
		ssec.errorChan <- err
	}
	packagedError := &Event{
		Origin: ssec.origin, // not entirely true - might be the wrong thing
		Error:  err,
		Type:   ErrorType,
	}
	ssec.emit(packagedError)
}

// URL returns the configured URL of the client
func (ssec *SSEClient) URL() *url.URL {
	ssec.Lock()
	retUrl := ssec.url
	ssec.Unlock()
	return retUrl
}

func (ssec *SSEClient) readStream() {
	ssec.eventStream = make(chan *Event)
	ssec.reader = newEventStreamReader(ssec.response.Body, ssec.origin)
	go ssec.reader.decode(ssec.eventStream)
}

func (ssec *SSEClient) process() {
processLoop:
	for {
		if atomic.LoadInt32(&ssec.readyState) == int32(Connecting) {
			Logger.Printf("connecting to %s", ssec.url)
			err := ssec.connect()
			if err != nil {
				Logger.Printf("error; state=closed: %s", err)
				atomic.StoreInt32(&ssec.readyState, int32(Closed))
				ssec.emitError(err)
				break processLoop
			} else {
				atomic.StoreInt32(&ssec.readyState, int32(Open))
				Logger.Printf("connected; state=open: %s", err)
				ssec.emitOpenClose(true)
				ssec.readStream()
			}
		}
		// Logger.Printf("ready for an event - reading from %v", ssec.eventStream)
		select {
		case ev, ok := <-ssec.eventStream:
			// Logger.Printf("ev = %v, ok = %v", ev, ok)
			if ok {
				Logger.Printf("event; state=open: %v", ev)
				ssec.emit(ev)
			} else {
				// TODO - do we need to clean this up? ssec.response.Body.Close()
				ssec.emitOpenClose(false)
				ssec.response = nil
				if ssec.shouldReconnect() {
					Logger.Printf("detected close; reconnecting in %v", ssec.reconnectTime)
					atomic.StoreInt32(&ssec.readyState, int32(Connecting))
					if ssec.reconnectTime > time.Duration(0) {
						time.Sleep(ssec.reconnectTime)
					}
				} else {
					break processLoop
				}
			}
		}
	}
	Logger.Print("all done in process")
	close(ssec.messageChan)
	close(ssec.openChan)
	close(ssec.errorChan)
	ssec.wg.Done()
}
