// +build go1.7

package sse

import "context"

import (
	"net/http"
)

type canceler struct {
	cancel func()
	ctx    context.Context
}

func (ssec *SSEClient) Close() {
	ssec.setReconnect(false)
	ssec.cancel()
}

func (ssec *SSEClient) wrapRequest(req *http.Request) *http.Request {
	// go 1.7+ presumably
	ssec.ctx, ssec.cancel = context.WithCancel(context.Background())
	return req.WithContext(ssec.ctx)
}

// hook for go1.6-, no-op after
func (ssec *SSEClient) initTransport() {
}

// SetContext allows the context to be specified - this affects
// cancelation and timeouts.  Affects active client on reconnection only.
func (ssec *SSEClient) SetContext(ctx context.Context) {
	ssec.Lock()
	ssec.ctx = ctx
	ssec.Unlock()
}
