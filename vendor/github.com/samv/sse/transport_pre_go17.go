// +build !go1.7

package sse

import (
	"net/http"
)

func (ssec *SSEClient) Close() {
	ssec.setReconnect(false)
	// TODO - this approach to canceling a request is deprecated, but
	// the new method is go 1.7+
	ssec.transport.CancelRequest(ssec.request)
}

type canceler struct {
	// facilities for canceling clients
	transport *http.Transport
	request   *http.Request
}

func (ssec *SSEClient) initTransport() {
	ssec.transport = &http.Transport{}
	ssec.Client.Transport = ssec.transport
}

// hook for go1.7+, pass-through before
func (ssec *SSEClient) wrapRequest(req *http.Request) *http.Request {
	ssec.request = req
	return req
}
